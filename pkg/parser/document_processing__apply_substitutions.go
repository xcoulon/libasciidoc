package parser

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// TODO: convert `ctx *processContext` as a local variable instead of a func param
func ProcessSubstitutions(ctx *processContext, done <-chan interface{}, fragmentStream <-chan types.DocumentFragment) chan types.DocumentFragment {
	processedFragmentStream := make(chan types.DocumentFragment)
	go func() {
		defer close(processedFragmentStream)
		for f := range fragmentStream {
			select {
			case <-done:
				log.WithField("pipeline_stage", "apply_substitutions").Debug("received 'done' signal")
				return
			case processedFragmentStream <- processSubstitutions(ctx, f):
			}
		}
		log.WithField("pipeline_stage", "apply_substitutions").Debug("done processing upstream content")
	}()
	return processedFragmentStream
}

func processSubstitutions(ctx *processContext, f types.DocumentFragment) types.DocumentFragment {
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.WithField("pipeline_stage", "fragment_processing").Debugf("incoming fragment:\n%s", spew.Sdump(f))
	// }
	// if the fragment already contains an error, then send it as-is downstream
	if err := f.Error; err != nil {
		return f
	}
	element, err := processSubstitutionsOnElement(ctx, f.Content)
	if err != nil {
		return types.NewErrorFragment(f.LineOffset, err)
	}
	f.Content = element
	return f
}

func processSubstitutionsOnElement(ctx *processContext, element interface{}) (interface{}, error) {
	switch e := element.(type) {
	case *types.AttributeDeclaration:
		ctx.addAttribute(e.Name, e.Value)
		return e, nil
	case types.AttributeReset:
		delete(ctx.attributes, e.Name)
		return e, nil
	case types.WithElements:
		if err := processBlockWithElements(ctx, e); err != nil {
			return nil, err
		}
		return e, nil
	case types.WithLocation:
		if err := processBlockWithLocation(ctx, e); err != nil {
			return nil, err
		}
		return e, nil
	default:
		log.WithField("pipeline_stage", "fragment_processing").Debugf("forwarding fragment content of type '%T' as-is", e)
		return element, nil
	}
}

// replaces the AttributeSubstitution by their actual values.
// TODO: returns `true` if at least one AttributeSubstitution was found (whatever its replacement)?
func replaceAttributeSubstitutionsInAttributes(ctx *processContext, b types.WithAttributes) error {
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debugf("applying attribute substitutions on attributes\n%s", spew.Sdump(b.GetAttributes()))
	// }
	for key, value := range b.GetAttributes() {
		switch value := value.(type) {
		case []interface{}: // multi-value attributes
			value, err := replaceAttributeSubstitutionsInElements(ctx, value)
			if err != nil {
				return err
			}
			// (a bit hack-ish): do not merge values when the attribute is `roles` or `options`
			switch key {
			case types.AttrRoles, types.AttrOptions:
				b.GetAttributes()[key] = value
			default:
				b.GetAttributes()[key] = types.Reduce(value)
			}
		default: // single-value attributes
			value, err := replaceAttributeSubstitutionsInElement(ctx, value)
			if err != nil {
				return err
			}
			b.GetAttributes()[key] = types.Reduce(value)
		}
	}
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debugf("applied attribute substitutions:\n%s", spew.Sdump(b.GetAttributes()))
	// }
	return nil
}

func processBlockWithElements(ctx *processContext, block types.WithElements) error {
	log.Debugf("processing block of type '%T' with attributes and nested elements", block)
	if err := processAttributes(ctx, block); err != nil {
		return err
	}
	// log.Debugf("applying substitutions on elements of block of type '%T'", block)
	plan, err := newSubstitutionPlan(block)
	if err != nil {
		return err
	}
	elements, err := processElements(ctx, block.GetElements(), plan)
	if err != nil {
		return err
	}
	if err := block.SetElements(elements); err != nil {
		return errors.Wrapf(err, "failed to process substitutions on block of type '%T'", block)
	}
	return nil
}

func processLabeledListElement(ctx *processContext, block *types.LabeledListElement) error {
	log.Debugf("processing LabeledListElement")
	if err := processAttributes(ctx, block); err != nil {
		return err
	}
	// log.Debugf("applying substitutions on elements of block of type '%T'", block)
	plan, err := newSubstitutionPlan(block)
	if err != nil {
		return err
	}
	term, err := processElements(ctx, block.Term, plan)
	if err != nil {
		return err
	}
	block.Term = term
	elements, err := processElements(ctx, block.GetElements(), plan)
	if err != nil {
		return err
	}
	if err := block.SetElements(elements); err != nil {
		return errors.Wrapf(err, "failed to process substitutions on block of type '%T'", block)
	}
	return nil
}

func processBlockWithLocation(ctx *processContext, block types.WithLocation) error {
	log.Debugf("processing block with attributes and location")
	if err := processAttributes(ctx, block); err != nil {
		return err
	}
	log.Debugf("applying substitutions on `location` of block of type '%T'", block)
	elements := block.GetLocation().Path
	elements, err := replaceAttributeSubstitutionsInElements(ctx, elements)
	if err != nil {
		return err
	}
	block.GetLocation().SetPath(elements)
	imagesdir := ctx.attributes.GetAsStringWithDefault("imagesdir", "")
	block.GetLocation().SetPathPrefix(imagesdir)
	return nil
}

func processElements(ctx *processContext, elements []interface{}, plan *substitutionPlan) ([]interface{}, error) {
	// skip if there's nothing to do
	if len(elements) == 0 {
		return elements, nil
	}
	var err error
	for _, step := range plan.steps {
		// log.Debugf("applying step '%s'", step.group)
		elements, err = parseElements(elements, step.group, GlobalStore(substitutionPhaseKey, step))
		if err != nil {
			return nil, err
		}
		if step.hasAttributeSubstitutions {
			// log.WithField("group", step.group).Debug("attribute substitutions detected during parsing")
			// apply substitutions on elements
			elements, err = replaceAttributeSubstitutionsInElements(ctx, elements)
			if err != nil {
				return nil, err
			}
			elements = types.Merge(elements)
			// re-run the Parser, skipping attribute substitutions and earlier rules this time
			if step.reduce() {
				if elements, err = parseElements(elements, step.group, GlobalStore(substitutionPhaseKey, step)); err != nil {
					return nil, err
				}
			}
		}
	}
	return elements, nil
}

func processAttributes(ctx *processContext, block types.WithAttributes) error {
	if err := replaceAttributeSubstitutionsInAttributes(ctx, block); err != nil {
		return err
	}
	// TODO: only parse element attributes if an attribute substitution occurred?
	attrs, err := parseAttributes(block.GetAttributes(), ElementAttributesGroup)
	if err != nil {
		return err
	}
	block.SetAttributes(attrs)
	return nil
}

type substitutionPlan struct {
	steps []*substitutionStep
}

func newSubstitutionPlan(b types.WithAttributes) (*substitutionPlan, error) {
	// TODO: introduce a `types.BlockWithSubstitution` interface?
	// note: would also be helpful for paragraphs with `[listing]` style.
	s, err := defaultSubstitution(b)
	if err != nil {
		return nil, errors.Wrap(err, "unable to determine substitution plan")
	}
	subs := b.GetAttributes().GetAsStringWithDefault(types.AttrSubstitutions, s)
	steps := strings.Split(subs, ",")
	plan := &substitutionPlan{
		steps: make([]*substitutionStep, len(steps)),
	}
	for i, step := range steps {
		s, err := newSubstitutionStep(step)
		if err != nil {
			return nil, err
		}
		plan.steps[i] = s
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		steps := make([]substitutionGroup, len(plan.steps))
		for i, p := range plan.steps {
			steps[i] = p.group
		}
		// log.Debugf("applying steps: '%s'", steps)
	}
	return plan, nil
}

func defaultSubstitution(b interface{}) (string, error) {
	switch b := b.(type) {
	case *types.DelimitedBlock:
		switch b.Kind {
		case types.Listing, types.Fenced:
			return "verbatim", nil
		default:
			return "", fmt.Errorf("unsupported kind of delimited block: '%v'", b.Kind)
		}
	case *types.Paragraph, *types.GenericList, *types.Section:
		return "normal", nil
	default:
		return "", fmt.Errorf("unsupported kind of element: '%T'", b)
	}
}

type substitutionStep struct {
	group                     substitutionGroup
	enablements               map[SubstitutionKind]bool
	hasAttributeSubstitutions bool
}

type substitutionGroup string

const (
	AttributesGroup        substitutionGroup = "AttributesGroup"
	ElementAttributesGroup substitutionGroup = "ElementAttributesGroup"
	HeaderGroup            substitutionGroup = "HeaderGroup"
	MacrosGroup            substitutionGroup = "MacrosGroup"
	NoneGroup              substitutionGroup = "NoneGroup"
	NormalGroup            substitutionGroup = "NormalGroup"
	QuotesGroup            substitutionGroup = "QuotesGroup"
	ReplacementsGroup      substitutionGroup = "ReplacementsGroup"
	SpecialcharactersGroup substitutionGroup = "SpecialCharactersGroup"
	VerbatimGroup          substitutionGroup = "VerbatimGroup"
)

//TODO: simplify by using a single grammar rule and turn-off choices?
var substitutionGroups = map[string]substitutionGroup{
	"":                  NormalGroup,
	"attributes":        AttributesGroup,
	"header":            HeaderGroup,
	"macros":            MacrosGroup,
	"normal":            NormalGroup,
	"none":              NoneGroup,
	"quotes":            QuotesGroup,
	"replacements":      ReplacementsGroup,
	"specialchars":      SpecialcharactersGroup,
	"specialcharacters": SpecialcharactersGroup,
	"verbatim":          VerbatimGroup,
}

func newSubstitutionStep(kind string) (*substitutionStep, error) {
	group, found := substitutionGroups[kind]
	if !found {
		return nil, fmt.Errorf("unsupported kind of substitution: '%v'", kind)
	}
	s := &substitutionStep{
		group: group,
	}
	s.reset()
	return s, nil
}

func (s *substitutionStep) reset() {
	switch s.group {
	case AttributesGroup:
		s.enablements = map[SubstitutionKind]bool{
			InlinePassthroughs: true,
			Attributes:         true,
		}
	case ElementAttributesGroup:
		s.enablements = map[SubstitutionKind]bool{
			InlinePassthroughs: true,
			Attributes:         true,
			Quotes:             true,
			SpecialCharacters:  true, // TODO: is it needed?
		}
	case HeaderGroup:
		s.enablements = map[SubstitutionKind]bool{
			InlinePassthroughs: true,
			SpecialCharacters:  true,
			Attributes:         true,
		}
	case MacrosGroup:
		s.enablements = map[SubstitutionKind]bool{
			Macros: true,
		}
	case NoneGroup:
		s.enablements = map[SubstitutionKind]bool{}
	case NormalGroup:
		s.enablements = map[SubstitutionKind]bool{
			InlinePassthroughs: true,
			SpecialCharacters:  true,
			Attributes:         true,
			Quotes:             true,
			Replacements:       true,
			Macros:             true,
			PostReplacements:   true,
		}
	case QuotesGroup:
		s.enablements = map[SubstitutionKind]bool{
			Quotes: true,
		}
	case ReplacementsGroup:
		s.enablements = map[SubstitutionKind]bool{
			Replacements: true,
		}
	case SpecialcharactersGroup:
		s.enablements = map[SubstitutionKind]bool{
			SpecialCharacters: true,
		}
	case VerbatimGroup:
		s.enablements = map[SubstitutionKind]bool{
			SpecialCharacters: true,
			Callouts:          true,
		}
	}
}

// disables the "inline_passthroughs", "special_characters" and "attributes" substitution
// return `true` if there are more enablements to apply, `false` otherwise (ie, no substitution would be applied if the content was parsed again)
func (s *substitutionStep) reduce() bool { // TODO: rename this func
	s.reset()
	for sub := range s.enablements {
		switch sub {
		case InlinePassthroughs, SpecialCharacters, Attributes:
			delete(s.enablements, sub)
		}
	}
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debugf("new enablements for '%s': %s", s.group, spew.Sdump(s.enablements))
	// }
	return len(s.enablements) > 0
}

// sets the new substitution plan in the golbal store, overriding any exist–ing one
// NOTE: will override any existing substitution context
// TODO: is there any case where a stack would be needed, so we don't override an existing context?
func (c *current) setSubstitutionPhase(kind string) error {
	p, err := newSubstitutionStep(kind)
	if err != nil {
		return err
	}
	c.globalStore[substitutionPhaseKey] = p
	return nil
}

func (c *current) unsetSubstitutionPhase() {
	delete(c.globalStore, substitutionPhaseKey)
}

func (c *current) lookupSubstitutionPhase() (*substitutionStep, error) {
	ctx, ok := c.globalStore[substitutionPhaseKey].(*substitutionStep)
	if !ok {
		return nil, fmt.Errorf("unable to look-up the substitution context in the parser's global store")
	}
	return ctx, nil
}

// called when an attribute substitution occurred
// TODO: find a better name for this method
func (c *current) hasAttributeSubstitutions() error {
	phase, err := c.lookupSubstitutionPhase()
	if err != nil {
		return err
	}
	phase.hasAttributeSubstitutions = true
	// // also, disable all subsitutions post `attributes` (`macros`, etc.)
	// for s := range phase.enablements {
	// 	switch s {
	// 	case Quotes, Replacements, Macros, PostReplacements: // TODO: avoid hard-coded entries
	// 		phase.enablements[s] = false // disabled
	// 	}
	// }
	return nil
}

func (c *current) isSubstitutionEnabled(k SubstitutionKind) (bool, error) {
	phase, err := c.lookupSubstitutionPhase()
	if err != nil {
		return false, err
	}
	enabled, found := phase.enablements[k]
	if !found {
		return false, nil
	}
	return enabled, nil
}

type SubstitutionKind string

const (
	// substitutionPhaseKey the key in which substitutions contexts are stored
	substitutionPhaseKey string = "substitution_contexts"

	// Attributes the "attributes" substitution
	Attributes SubstitutionKind = "attributes"
	// Callouts the "callouts" substitution
	Callouts SubstitutionKind = "callouts"
	// InlinePassthroughs the "inline_passthrough" substitution
	InlinePassthroughs SubstitutionKind = "inline_passthrough"
	// Macros the "macros" substitution
	Macros SubstitutionKind = "macros"
	// None the "none" substitution
	None SubstitutionKind = "none"
	// PostReplacements the "post_replacements" substitution
	PostReplacements SubstitutionKind = "post_replacements"
	// Quotes the "quotes" substitution
	Quotes SubstitutionKind = "quotes"
	// Replacements the "replacements" substitution
	Replacements SubstitutionKind = "replacements"
	// SpecialCharacters the "specialcharacters" substitution
	SpecialCharacters SubstitutionKind = "specialcharacters"
)

func parseElements(elements []interface{}, group substitutionGroup, opts ...Option) ([]interface{}, error) {
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.WithField("group", group).Debug("parsing elements")
	// }
	serialized, placeholders := serialize(elements)
	if len(serialized) == 0 {
		return nil, nil
	}
	opts = append(opts)
	result, err := parseContent(serialized, append(opts, Entrypoint(string(group)))...)
	if err != nil {
		return nil, err
	}
	// also, apply the same substitution group on the placeholders, case by case
	for key, element := range placeholders.elements {
		log.Debugf("processing placeholder of type '%T'", element)
		if e, ok := element.(types.WithAttributes); ok {
			attrs, err := parseAttributes(e.GetAttributes(), group, opts...)
			if err != nil {
				return nil, err
			}
			e.SetAttributes(attrs)
		}
		if e, ok := element.(*types.LabeledListElement); ok {
			term, err := parseElements(e.Term, group, opts...)
			if err != nil {
				return nil, err
			}
			e.Term = term
		}
		if e, ok := element.(types.WithElements); ok {
			elements, err := parseElements(e.GetElements(), group, opts...)
			if err != nil {
				return nil, err
			}
			if err := e.SetElements(elements); err != nil {
				return nil, errors.Wrapf(err, "failed to parse elements of block of type '%T'", e)
			}
		}
		if e, ok := element.([]interface{}); ok {
			elements, err := parseElements(e, group, opts...)
			if err != nil {
				return nil, err
			}
			placeholders.elements[key] = elements
		}
		// if e, ok := element.(types.WithLocation); ok {
		// 	if err := parseElementWithLocation(e, group, opts...); err != nil {
		// 		return nil, err
		// 	}
		// }
	}
	result = placeholders.restoreElements(result)
	return result, nil
}

func parseAttributes(attributes types.Attributes, group substitutionGroup, opts ...Option) (types.Attributes, error) {
	if !(group == AttributesGroup || group == QuotesGroup) { // TODO: include special_characters?
		log.Debugf("no need to parse attributes for group '%s'", group)
		return attributes, nil
	}
	for name, value := range attributes {
		switch value := value.(type) {
		case []interface{}:
			serialized, placeholders := serialize(value)
			if len(serialized) == 0 {
				continue
			}
			elements, err := parseContent(serialized, append(opts, Entrypoint(string(group)))...)
			if err != nil {
				return nil, err
			}
			elements = placeholders.restoreElements(elements)
			attributes[name] = elements
		case string:
			elements, err := parseContent([]byte(value), append(opts, Entrypoint(string(group)))...)
			if err != nil {
				return nil, err
			}
			attributes[name] = types.Reduce(elements)
		default:
			return nil, fmt.Errorf("unexpected type of attribute value: '%T'", value)
		}
	}
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debugf("parsed attributes for group '%s': %s", group, spew.Sdump(element.GetAttributes()))
	// }
	return attributes, nil
}

func parseContent(content []byte, opts ...Option) ([]interface{}, error) {
	result, err := Parse("", content, opts...)
	if err != nil {
		return nil, err
	}
	r, ok := result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type of content after parsing elements: '%T'", result)
	}
	return r, nil
}

// func applyAttributeSubstitutionsOnElements(ctx *processContext, elements []interface{}) ([]interface{}, error) {
// 	result := make([]interface{}, len(elements)) // maximum capacity should exceed initial input
// 	for i, element := range elements {
// 		element, err := applyAttributeSubstitutionsOnElement(ctx, element)
// 		if err != nil {
// 			return nil, err
// 		}
// 		result[i] = element
// 	}
// 	return result, nil
// }

// replaces the AttributeSubstitution or Counter substitution with its actual value, recursively if the given `element`
// is a slice
func replaceAttributeSubstitutionsInElement(ctx *processContext, element interface{}) (interface{}, error) {
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("applying attribute substitution on element of type '%T'", element)
	}
	switch e := element.(type) {
	case []interface{}:
		return replaceAttributeSubstitutionsInElements(ctx, e)
	case types.AttributeSubstitution:
		return types.StringElement{
			Content: ctx.attributes.GetAsStringWithDefault(e.Name, "{"+e.Name+"}"),
		}, nil
	case types.CounterSubstitution:
		return applyCounterSubstitution(ctx, e)
	case types.WithElements:
		// replace AttributeSubstitutions on attributes
		if err := replaceAttributeSubstitutionsInAttributes(ctx, e); err != nil {
			return nil, err
		}
		// replace AttributeSubstitutions on nested elements
		elements, err := replaceAttributeSubstitutionsInElements(ctx, e.GetElements())
		if err != nil {
			return nil, err
		}
		// elements = types.Merge(elements)
		if err := e.SetElements(elements); err != nil {
			return nil, errors.Wrapf(err, "failed to apply attribute substitutions on block of type '%T'", e)
		}
		return e, nil
	case types.WithLocation:
		// replace AttributeSubstitutions on attributes
		if err := replaceAttributeSubstitutionsInAttributes(ctx, e); err != nil {
			return nil, err
		}
		// replace AttributeSubstitutions on embedded location
		if err := applyAttributeSubstitutionsOnLocation(ctx, e); err != nil {
			return nil, err
		}
		return e, nil
	case types.WithAttributes:
		// replace AttributeSubstitutions on attributes
		if err := replaceAttributeSubstitutionsInAttributes(ctx, e); err != nil {
			return nil, err
		}
		return e, nil
	default:
		return e, nil
	}
}

func applyAttributeSubstitutionsOnLocation(ctx *processContext, b types.WithLocation) error {
	path, err := replaceAttributeSubstitutionsInElements(ctx, b.GetLocation().Path)
	if err != nil {
		return err
	}
	b.GetLocation().Path = path
	return nil
}

func replaceAttributeSubstitutionsInElements(ctx *processContext, elements []interface{}) ([]interface{}, error) {
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("applying attribute substitutions on elements:\n%s", spew.Sdump(elements))
	}
	result := make([]interface{}, len(elements)) // maximum capacity should exceed initial input
	for i, element := range elements {
		element, err := replaceAttributeSubstitutionsInElement(ctx, element)
		if err != nil {
			return nil, err
		}
		result[i] = element
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("applied attribute substitutions on elements:\n%s", spew.Sdump(result))
	}
	return result, nil
}

// applyCounterSubstitutions is called by applyAttributeSubstitutionsOnElement.  Unless there is an error with
// the element (the counter is the wrong type, which should never occur), it will return a `StringElement, true`
// (because we always either find the element, or allocate one), and `nil`.  On an error it will return `nil, false`,
// and the error.  The extra boolean here is to fit the calling expectations of our caller.  This function was
// factored out of a case from applyAttributeSubstitutionsOnElement in order to reduce the complexity of that
// function, but otherwise it should have no callers.
func applyCounterSubstitution(ctx *processContext, c types.CounterSubstitution) (interface{}, error) {
	counter := ctx.counters[c.Name]
	if counter == nil {
		counter = 0
	}
	increment := true
	if c.Value != nil {
		ctx.counters[c.Name] = c.Value
		counter = c.Value
		increment = false
	}
	switch counter := counter.(type) {
	case int:
		if increment {
			counter++
		}
		ctx.counters[c.Name] = counter
		if c.Hidden {
			// return empty string facilitates merging
			return types.StringElement{Content: ""}, nil
		}
		return types.StringElement{
			Content: strconv.Itoa(counter),
		}, nil
	case rune:
		if increment {
			counter++
		}
		ctx.counters[c.Name] = counter
		if c.Hidden {
			// return empty string facilitates merging
			return types.StringElement{Content: ""}, nil
		}
		return types.StringElement{
			Content: string(counter),
		}, nil
	default:
		return types.StringElement{}, fmt.Errorf("unexpected type of counter value: '%T'", counter)
	}
}

type placeholders struct {
	seq      int
	elements map[string]interface{}
}

func newPlaceholders() *placeholders {
	return &placeholders{
		seq:      0,
		elements: map[string]interface{}{},
	}
}

func (p *placeholders) add(element interface{}) types.ElementPlaceHolder {
	p.seq++
	p.elements[strconv.Itoa(p.seq)] = element
	return types.ElementPlaceHolder{
		Ref: strconv.Itoa(p.seq),
	}

}

// replace the placeholders with their original element in the given elements
func (p *placeholders) restoreElements(elements []interface{}) []interface{} {
	// skip if there's nothing to restore
	if len(p.elements) == 0 {
		return elements
	}
	for i, e := range elements {
		//
		if e, ok := e.(types.ElementPlaceHolder); ok {
			elements[i] = p.elements[e.Ref]
		}
		// also check nested elements (eg, in QuotedText, etc.)
		// for each element, check *all* interfaces to see if there's a need to replace the placeholders
		if e, ok := e.(types.WithPlaceholdersInElements); ok {
			elements[i] = e.RestoreElements(p.elements)
		}
		if e, ok := e.(types.WithPlaceholdersInAttributes); ok {
			elements[i] = e.RestoreAttributes(p.elements)
		}
		if e, ok := e.(types.WithPlaceholdersInLocation); ok {
			elements[i] = e.RestoreLocation(p.elements)
		}
	}
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debugf("restored elements:\n%v", spew.Sdump(elements))
	// }
	return elements
}

func serialize(content interface{}) ([]byte, *placeholders) {
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debugf("serializing:\n%v", spew.Sdump(content))
	// }
	placeholders := newPlaceholders()
	result := bytes.NewBuffer(nil)
	switch content := content.(type) {
	case string: // for attributes with simple (string) values
		result.WriteString(content)
	case []interface{}: // for paragraph lines, attributes with complex values, etc.
		for i, element := range content {
			switch element := element.(type) {
			case types.RawLine:
				result.WriteString(string(element))
				// add `\n` unless the next element is a single-line comment
				if i < len(content)-1 {
					if _, ok := content[i+1].(*types.SingleLineComment); !ok {
						result.WriteString("\n")
					}
				}
			case *types.SingleLineComment:
				// replace with placeholder
				p := placeholders.add(element)
				result.WriteString(p.String())
				// add `\n` unless the next element is a single-line comment
				if i < len(content)-1 {
					if _, ok := content[i+1].(*types.SingleLineComment); !ok {
						result.WriteString("\n")
					}
				}
			case types.StringElement:
				result.WriteString(element.Content)
			// case types.SingleLineComment:
			// 	// replace with placeholder
			// 	p := placeholders.add(element)
			// 	result.WriteString(p.String())
			default:
				// replace with placeholder
				p := placeholders.add(element)
				result.WriteString(p.String())
			}
		}
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("serialized lines: '%s'\nplaceholders: %v", result.Bytes(), spew.Sdump(placeholders.elements))
	}
	return result.Bytes(), placeholders
}
