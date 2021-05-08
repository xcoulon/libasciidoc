package parser

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/bytesparadise/libasciidoc/pkg/configuration"
	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

// ParseDocument parses the content of the reader identitied by the filename
func ParseDocument(r io.Reader, config configuration.Configuration, opts ...Option) (types.Document, error) {
	done := make(chan interface{})
	defer close(done)

	attributes := types.Attributes{}
	blocks := []interface{}{}
	inHeader := true
	ctx := newProcessContext()
	pipeline := processFragments(ctx, done,
		AssembleFragments(done,
			ParseDocumentFragmentGroups(r, done, opts...)))
	for f := range pipeline {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.WithField("_state", "document_parsing").Debugf("received document fragment: %s", spew.Sdump(f))
		}
		if err := f.Error; err != nil {
			log.WithField("_state", "document_parsing").WithError(err).Error("error occurred")
			return types.Document{}, err
		}
		switch b := f.Content.(type) {
		case types.Section:
			if b.Level != 0 && inHeader { // do not even allow 2ndary section with level 0 as headers
				inHeader = false
			}
		case types.AttributeDeclaration:
			if inHeader {
				attributes.Set(b.Name, b.Value)
			}
		case types.AttributeReset:
			delete(attributes, b.Name)
		default:
			// anything else and we're not in the header anynore
			inHeader = false
		}
		blocks = append(blocks, f.Content)
	}
	// fragments := make(chan types.DocumentFragment)
	// err := ParseDocumentFragmentGroups							(r, fragments, opts...)
	// if err != nil {
	// 	return types.Document{}, err
	// }

	// elements, err := assemble(fragments)

	// // draftDoc, err := ApplySubstitutions(rawDoc, config)
	// // if err != nil {
	// // 	return types.Document{}, err
	// // }
	// // if log.IsLevelEnabled(log.DebugLevel) {
	// // 	log.Debug("draft doc:")
	// // 	spew.Fdump(log.StandardLogger().Out, draftDoc)
	// // }

	// // now, merge list items into proper lists
	// blocks, err := rearrangeListItems(elements, false)
	// if err != nil {
	// 	return types.Document{}, err
	// }
	// // filter out blocks not needed in the final doc
	// blocks = filter(blocks, allMatchers...)

	// blocks, footnotes := processFootnotes(blocks)
	// // now, rearrange elements in a hierarchical manner
	// doc, err := rearrangeSections(blocks)
	// if err != nil {
	// 	return types.Document{}, err
	// }
	// // also, set the footnotes
	// doc.Footnotes = footnotes
	// // insert the preamble at the right location
	// doc = includePreamble(doc)
	// // doc.Attributes = doc.Attributes.SetAll(draftDoc.Attributes)
	// // also insert the table of contents
	// doc = includeTableOfContentsPlaceHolder(doc)
	// // finally
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debug("final document:")
	// 	spew.Fdump(log.StandardLogger().Out, doc)
	// }
	// return doc, nil
	return types.Document{
		Attributes: attributes.NilIfEmpty(),
		Elements:   blocks,
	}, nil
}

type processContext struct {
	attributes types.Attributes
	counters   map[string]interface{}
}

func newProcessContext() *processContext {
	return &processContext{
		attributes: types.Attributes{},
		counters:   map[string]interface{}{},
	}
}
func (ctx *processContext) addAttribute(name string, value interface{}) {
	ctx.attributes[name] = value
}

// ContextKey a non-built-in type for keys in the context
type ContextKey string

// LevelOffset the key for the level offset of the file to include
const LevelOffset ContextKey = "leveloffset"

func processFragments(ctx *processContext, done <-chan interface{}, fragmentStream <-chan types.DocumentFragment) chan types.DocumentFragment {
	processedFragmentStream := make(chan types.DocumentFragment)
	go func() {
		defer close(processedFragmentStream)
		for f := range fragmentStream {
			select {
			case <-done:
				log.WithField("stage", "fragment_processing").Debug("received 'done' signal")
				return
			case processedFragmentStream <- processFragment(ctx, f):
			}
		}
		log.Debug("end of fragment processing")
	}()
	return processedFragmentStream
}

func processFragment(ctx *processContext, f types.DocumentFragment) types.DocumentFragment {
	if log.IsLevelEnabled(log.DebugLevel) {
		log.WithField("stage", "fragment_processing").Debugf("incoming fragment:\n%s", spew.Sdump(f))
	}
	// if the fragment contains an error, then send it as-is downstream
	if err := f.Error; err != nil {
		return f
	}
	switch e := f.Content.(type) {
	case types.AttributeDeclaration:
		ctx.addAttribute(e.Name, e.Value)
		return types.NewDocumentFragment(f.LineOffset, e)
	case types.AttributeReset:
		delete(ctx.attributes, e.Name)
		return types.NewDocumentFragment(f.LineOffset, e)
	case types.BlockWithNestedElements:
		log.Debugf("processing block with attributes and nested elements")
		if err := applyAttributeSubstitutionsOnAttributes(ctx, e); err != nil {
			return types.DocumentFragment{
				LineOffset: f.LineOffset,
				Error:      err,
			}
		}
		// TODO: only parse element attributes if an attribute substitution occurred?
		if err := parseElementAttributes(e, ElementAttributesGroup); err != nil {
			return types.DocumentFragment{
				LineOffset: f.LineOffset,
				Error:      err,
			}
		}
		if err := processNestedElements(ctx, e); err != nil {
			return types.DocumentFragment{
				LineOffset: f.LineOffset,
				Error:      err,
			}
		}
		return types.DocumentFragment{
			LineOffset: f.LineOffset,
			Content:    e,
		}
	case types.BlockWithLocation:
		log.Debugf("processing block with attributes and location")
		if err := applyAttributeSubstitutionsOnAttributes(ctx, e); err != nil {
			return types.DocumentFragment{
				LineOffset: f.LineOffset,
				Error:      err,
			}
		}
		// TODO: only parse element attributes if an attribute substitution occurred?
		if err := parseElementAttributes(e, ElementAttributesGroup); err != nil {
			return types.DocumentFragment{
				LineOffset: f.LineOffset,
				Error:      err,
			}
		}
		if err := processLocation(ctx, e); err != nil {
			return types.DocumentFragment{
				LineOffset: f.LineOffset,
				Error:      err,
			}
		}
		return types.DocumentFragment{
			LineOffset: f.LineOffset,
			Content:    e,
		}
	default:
		log.WithField("stage", "fragment_processing").Debugf("forwarding fragment content of type '%T' as-is", e)
		return types.NewDocumentFragment(f.LineOffset, e)
	}
}

// replaces the AttributeSubstitution elements by their actual values.
// TODO: returns `true` if at least one AttributeSubstitution was found (whatever its replacement)?
func applyAttributeSubstitutionsOnAttributes(ctx *processContext, b types.BlockWithAttributes) error {
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("applying attribute substitutions on attributes\n%s", spew.Sdump(b.GetAttributes()))
	}
	for key, value := range b.GetAttributes() {
		switch value := value.(type) {
		case []interface{}: // multi-value attributes
			value, err := applyAttributeSubstitutionsOnElements(ctx, value)
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
			value, err := applyAttributeSubstitutionsOnElement(ctx, value)
			if err != nil {
				return err
			}
			b.GetAttributes()[key] = types.Reduce(value)
		}
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("applied attribute substitutions:\n%s", spew.Sdump(b.GetAttributes()))
	}
	return nil
}

func processNestedElements(ctx *processContext, b types.BlockWithNestedElements) error {
	log.Debugf("applying substitutions on elements of block of type '%T'", b)
	plan, err := newSubstitutionPlan(b.GetAttributes().GetAsStringWithDefault(types.AttrSubstitutions, ""))
	if err != nil {
		return err
	}

	elements := b.GetElements()
	for _, step := range plan.steps {
		log.Debugf("applying step '%s'", step.group)
		if elements, err = parseElements(elements, step.group, GlobalStore(substitutionPhaseKey, step)); err != nil {
			return err
		}
		if step.hasAttributeSubstitutions {
			log.Debugf("attribute substitutions detected during parsing of group '%s'", step.group)
			// apply substitutions on elements
			elements, err = applyAttributeSubstitutionsOnElements(ctx, elements)
			if err != nil {
				return err
			}
			elements = types.Merge(elements)
			// re-run the Parser, skipping attribute substitutions and earlier rules this time
			if step.reduce() {
				if elements, err = parseElements(elements, step.group, GlobalStore(substitutionPhaseKey, step)); err != nil {
					return err
				}
			}
		}
	}
	b.SetElements(elements)
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("applied substitutions on block elements")
		spew.Fdump(log.StandardLogger().Out, b.GetElements()...)
	}
	return nil
}

func processLocation(ctx *processContext, b types.BlockWithLocation) error {
	log.Debugf("applying substitutions on `location` of block of type '%T'", b)
	elements := b.GetLocation().Path
	elements, err := applyAttributeSubstitutionsOnElements(ctx, elements)
	if err != nil {
		return err
	}
	b.GetLocation().SetPath(elements)
	imagesdir := ctx.attributes.GetAsStringWithDefault("imagesdir", "")
	b.GetLocation().SetPathPrefix(imagesdir)

	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("applied substitutions on block with location")
		spew.Fdump(log.StandardLogger().Out, b)
	}
	return nil
}

type substitutionPlan struct {
	steps []*substitutionStep
}

func newSubstitutionPlan(subs string) (*substitutionPlan, error) {
	phases := strings.Split(subs, ",")
	plan := &substitutionPlan{
		steps: make([]*substitutionStep, len(phases)),
	}
	for i, name := range phases {
		phase, err := newSubstitutionStep(name)
		if err != nil {
			return nil, err
		}
		plan.steps[i] = phase
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		phases := make([]substitutionGroup, len(plan.steps))
		for i, p := range plan.steps {
			phases[i] = p.group
		}
		log.Debugf("applying steps: '%s'", phases)
	}
	return plan, nil
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
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("new enablements for '%s': %s", s.group, spew.Sdump(s.enablements))
	}
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

	// InlinePassthroughs the "inline_passthrough" substitution
	InlinePassthroughs SubstitutionKind = "inline_passthrough"
	// Attributes the "attributes" substitution
	Attributes SubstitutionKind = "attributes"
	// SpecialCharacters the "specialcharacters" substitution
	SpecialCharacters SubstitutionKind = "specialcharacters"
	// Callouts the "callouts" substitution
	Callouts SubstitutionKind = "callouts"
	// Quotes the "quotes" substitution
	Quotes SubstitutionKind = "quotes"
	// Replacements the "replacements" substitution
	Replacements SubstitutionKind = "replacements"
	// Macros the "macros" substitution
	Macros SubstitutionKind = "macros"
	// PostReplacements the "post_replacements" substitution
	PostReplacements SubstitutionKind = "post_replacements"
	// None the "none" substitution
	None SubstitutionKind = "none"
)

func parseElements(elements []interface{}, group substitutionGroup, opts ...Option) ([]interface{}, error) {
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("parsing elements with group '%v'", group)
	}
	serialized, placeholders := serialize(elements)
	opts = append(opts)
	result, err := parseContent(serialized, append(opts, Entrypoint(string(group)))...)
	if err != nil {
		return nil, err
	}
	// also, apply the same substitution group on the placeholders, case by case
	for _, element := range placeholders.elements {
		if e, ok := element.(types.BlockWithAttributes); ok {
			if err := parseElementAttributes(e, group, opts...); err != nil {
				return nil, err
			}
		}
		if e, ok := element.(types.BlockWithNestedElements); ok {
			elements, err := parseElements(e.GetElements(), group, opts...)
			if err != nil {
				return nil, err
			}
			e.SetElements(elements)
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

func parseElementAttributes(element types.BlockWithAttributes, group substitutionGroup, opts ...Option) error {
	if !(group == AttributesGroup || group == QuotesGroup) { // TODO: include special_characters?
		log.Debugf("no need to parse attributes for group '%s'", group)
		return nil
	}
	for name, value := range element.GetAttributes() {
		switch value := value.(type) {
		case []interface{}:
			serialized, placeholders := serialize(value)
			elements, err := parseContent(serialized, append(opts, Entrypoint(string(group)))...)
			if err != nil {
				return err
			}
			elements = placeholders.restoreElements(elements)
			element.GetAttributes()[name] = elements
		case string:
			elements, err := parseContent([]byte(value), append(opts, Entrypoint(string(group)))...)
			if err != nil {
				return err
			}
			element.GetAttributes()[name] = types.Reduce(elements)
		default:
			return fmt.Errorf("unexpected type of attribute value: '%T'", value)
		}
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("parsed attributes for group '%s': %s", group, spew.Sdump(element.GetAttributes()))
	}
	return nil
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
func applyAttributeSubstitutionsOnElement(ctx *processContext, element interface{}) (interface{}, error) {
	switch e := element.(type) {
	case []interface{}:
		return applyAttributeSubstitutionsOnElements(ctx, e)
	case types.AttributeSubstitution:
		return types.StringElement{
			Content: ctx.attributes.GetAsStringWithDefault(e.Name, "{"+e.Name+"}"),
		}, nil
	case types.CounterSubstitution:
		return applyCounterSubstitution(ctx, e)
	case types.BlockWithNestedElements:
		// replace AttributeSubstitutions on attributes
		if err := applyAttributeSubstitutionsOnAttributes(ctx, e); err != nil {
			return nil, err
		}
		// replace AttributeSubstitutions on nested elements
		elements, err := applyAttributeSubstitutionsOnElements(ctx, e.GetElements())
		if err != nil {
			return nil, err
		}
		elements = types.Merge(elements)
		e.SetElements(elements)
		return e, nil
	case types.BlockWithLocation:
		// replace AttributeSubstitutions on attributes
		if err := applyAttributeSubstitutionsOnAttributes(ctx, e); err != nil {
			return nil, err
		}
		// replace AttributeSubstitutions on embedded location
		if err := applyAttributeSubstitutionsOnLocation(ctx, e); err != nil {
			return nil, err
		}
		return e, nil
	case types.BlockWithAttributes:
		// replace AttributeSubstitutions on attributes
		if err := applyAttributeSubstitutionsOnAttributes(ctx, e); err != nil {
			return nil, err
		}
		return e, nil
	default:
		return e, nil
	}
}

func applyAttributeSubstitutionsOnLocation(ctx *processContext, b types.BlockWithLocation) error {
	path, err := applyAttributeSubstitutionsOnElements(ctx, b.GetLocation().Path)
	if err != nil {
		return err
	}
	b.GetLocation().Path = path
	return nil
}

func applyAttributeSubstitutionsOnElements(ctx *processContext, elements []interface{}) ([]interface{}, error) {
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("applying attribute substitutions on elements:\n%s", spew.Sdump(elements))
	}
	result := make([]interface{}, len(elements)) // maximum capacity should exceed initial input
	for i, element := range elements {
		element, err := applyAttributeSubstitutionsOnElement(ctx, element)
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
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("restored elements:\n%v", spew.Sdump(elements))
	}
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
				if i < len(content)-1 {
					result.WriteString("\n")
				}
			case types.StringElement:
				result.WriteString(element.Content)
			case types.SingleLineComment:
				// replace with placeholder
				p := placeholders.add(element)
				result.WriteString(p.String())
			default:
				// replace with placeholder
				p := placeholders.add(element)
				result.WriteString(p.String())
			}
		}
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("serialized lines:\n%s\nplaceholders:\n%v", result.Bytes(), spew.Sdump(placeholders))
	}
	return result.Bytes(), placeholders
}
