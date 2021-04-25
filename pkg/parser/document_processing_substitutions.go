package parser

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

func applySubstitutionsOnParagraph(ctx *processContext, e *types.RawParagraph) (types.Paragraph, error) {
	// apply the attribute substitutions on the paragraph attributes
	e.Attributes = applyAttributeSubstitutionsOnAttributes(ctx, e.Attributes)
	// parse the lines, using the default substitutions
	plan, err := newSubstitutionPlan(e.Attributes.GetAsStringWithDefault(types.AttrSubstitutions, ""))
	if err != nil {
		return types.Paragraph{}, err
	}
	lines := e.Lines
	for _, step := range plan.steps {
		log.Debugf("applying step '%s'", step.group)
		if lines, err = parseParagraphLines(lines, step.group, GlobalStore(substitutionPhaseKey, step)); err != nil {
			return types.Paragraph{}, err
		}
		if step.hasAttributeSubstitutions {
			log.Debugf("attribute substitutions detected during parsing of group '%s'", step.group)
			// apply substitutions on attributes
			lines = applyAttributeSubstitutionsOnElements(ctx, lines)
			lines = types.Merge(lines)
			// re-run the Parser, ignoring attribute substitutions and earlier ones
			if step.reduce() {
				if lines, err = parseParagraphLines(lines, step.group, GlobalStore(substitutionPhaseKey, step)); err != nil {
					return types.Paragraph{}, err
				}
			}
		}
	}

	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("paragraph lines")
		spew.Fdump(log.StandardLogger().Out, lines)
	}
	p, err := types.NewParagraph(lines, e.Attributes)
	if err != nil {
		return types.Paragraph{}, err
	}
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.WithField("stage", "fragment_processing").Debugf("new paragraph: %s", spew.Sdump(p))
	// }
	return p, nil
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
	// case types.DocumentAuthor, types.DocumentRevision, types.AttributeDeclaration:
	// 	return &substitutionPlan{
	// 		group: "HeaderGroup",
	// 		substitutions: map[SubstitutionKind]bool{
	// 			// default substitutions
	// 			InlinePassthroughs: true,
	// 			SpecialCharacters:  true,
	// 			Attributes:         true,
	// 		},
	// 	}, nil
	// default:
	// 	return nil, fmt.Errorf("unsupported kind of block to process substitutions: '%v'", block)
	// }
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

// func (ctx *substitutionPlan) onAttributeDeclaration(d types.AttributeDeclaration) error {
// 	value := substituteAttributes(d.Value, ctx.attributes)
// 	switch value := value.(type) {
// 	case types.StringElement:
// 		ctx.attributes[d.Name] = value.Content
// 	case string:
// 		ctx.attributes[d.Name] = value
// 	case nil:
// 		ctx.attributes[d.Name] = ""
// 	default:
// 		return fmt.Errorf("unexpected type of value after substituing attributes: '%T'", value)
// 	}
// 	return nil
// }

// func (ctx *substitutionPlan) onBlockDelimiter(d types.BlockDelimiter) {
// 	currentLevel := ctx.blockLevels.get()
// 	if currentLevel == d.Kind {
// 		ctx.blockLevels.pop() // discard current level, assuming we've just parsed the ending delimiter of the current block
// 		return
// 	}
// 	ctx.blockLevels.push(d.Kind) // push current delimiter kind, assuming we've just parsed the starting delimiter of a new block
// 	if log.IsLevelEnabled(log.DebugLevel) {
// 		log.Debugf("current substitution context is now: %v", ctx.blockLevels.get())
// 	}
// }

// // called when block attributes are initialized.
// // if the `subs` attribute is set, it is retained in the parser context
// // and used to enable/disable subsequent parser rules
// // TODO: take into account blank lines and block delimiters to reset the substitutions to apply
// func (ctx *substitutionPlan) onBlockAttributes(attrs types.Attributes) {
// 	subs, exists := attrs[types.AttrSubstitutions]
// 	if !exists {
// 		return
// 	}
// 	switch subs {
// 	// 		case "normal":
// 	// 			subs = subs.append(
// 	// 				"specialcharacters",
// 	// 				"quotes",
// 	// 				"attributes",
// 	// 				"replacements",
// 	// 				"macros",
// 	// 				"post_replacements",
// 	// 			)
// 	// 		case "inline_passthrough", "callouts", "specialcharacters", "specialchars", "quotes", "attributes", "macros", "replacements", "post_replacements", "none":
// 	// 			subs = subs.append(s)
// 	// 		case "+callouts", "+specialcharacters", "+specialchars", "+quotes", "+attributes", "+macros", "+replacements", "+post_replacements", "+none":
// 	// 			if len(subs) == 0 {
// 	// 				subs = subs.append(block.DefaultSubstitutions()...)
// 	// 			}
// 	// 			subs = subs.append(strings.ReplaceAll(s, "+", ""))
// 	// 		case "callouts+", "specialcharacters+", "specialchars+", "quotes+", "attributes+", "macros+", "replacements+", "post_replacements+", "none+":
// 	// 			if len(subs) == 0 {
// 	// 				subs = subs.append(block.DefaultSubstitutions()...)
// 	// 			}
// 	// 			subs = subs.prepend(strings.ReplaceAll(s, "+", ""))
// 	// 		case "-callouts", "-specialcharacters", "-specialchars", "-quotes", "-attributes", "-macros", "-replacements", "-post_replacements", "-none":
// 	// 			if len(subs) == 0 {
// 	// 				subs = subs.append(block.DefaultSubstitutions()...)
// 	// 			}
// 	// 			subs = subs.remove(strings.ReplaceAll(s, "-", ""))
// 	// 		default:
// 	// 			return nil, fmt.Errorf("unsupported substitution: '%s", s)
// 	// 		}
// 	}
// }

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

func parseParagraphLines(content []interface{}, group substitutionGroup, options ...Option) ([]interface{}, error) {
	log.Debug("parsing lines")
	serialized, placeholders := serialize(content)
	// options = append(options)
	lines, err := parseElements(serialized, append(options, Entrypoint(string(group)))...)
	if err != nil {
		return nil, err
	}
	// also, apply the same substitution group on the placeholders, case by case
	for _, p := range placeholders.elements {
		switch p := p.(type) {
		case types.InlineLink: // TODO: define interface for elements with Attributes
			if err := parseElementAttributes(p.Attributes, group, options...); err != nil {
				return nil, err
			}
		}
	}
	lines = placeholders.restoreElements(lines)
	return lines, nil
}

func parseElementAttributes(attributes types.Attributes, group substitutionGroup, options ...Option) error {
	if !(group == AttributesGroup || group == QuotesGroup) { // TODO: include special_characters?
		log.Debugf("no need to parse attributes for group '%s'", group)
		return nil
	}
	for name, value := range attributes {
		switch value := value.(type) {
		case []interface{}:
			serialized, placeholders := serialize(value)
			result, err := parseElements(serialized, append(options, Entrypoint(string(group)))...)
			if err != nil {
				return err
			}
			result = placeholders.restoreElements(result)
			attributes[name] = result
		case string:
			result, err := parseElements([]byte(value), append(options, Entrypoint(string(group)))...)
			if err != nil {
				return err
			}
			attributes[name] = types.Reduce(result)
		default:
			return fmt.Errorf("unexpected type of attribute value: '%T'", value)
		}
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("parsed attributes for group '%s': %s", group, spew.Sdump(attributes))
	}
	return nil
}

func parseElements(content []byte, opts ...Option) ([]interface{}, error) {
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

func applyAttributeSubstitutionsOnElements(ctx *processContext, elements []interface{}) []interface{} {
	result := make([]interface{}, len(elements)) // maximum capacity should exceed initial input
	for i, element := range elements {
		result[i] = applyAttributeSubstitutionsOnElement(ctx, element)
	}
	return result
}

func applyAttributeSubstitutionsOnAttributes(ctx *processContext, attributes types.Attributes) types.Attributes {
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debugf("applying attribute substitutions on\n%s", spew.Sdump(attributes))
	// }
	for key, value := range attributes {
		switch value := value.(type) {
		case []interface{}: // multi-value attributes
			attributes[key] = applyAttributeSubstitutionsOnElements(ctx, value)
		default: // single-value attributes
			value = applyAttributeSubstitutionsOnElement(ctx, value)
			attributes[key] = types.Reduce(value)
		}
	}
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debugf("applied attribute substitutions:\n%s", spew.Sdump(attributes))
	// }
	return attributes
}

func applyAttributeSubstitutionsOnElement(ctx *processContext, element interface{}) interface{} {
	switch e := element.(type) {
	case []interface{}:
		return applyAttributeSubstitutionsOnElements(ctx, e)
	case types.AttributeDeclaration:
		ctx.attributes.Set(e.Name, e.Value)
		return e
	case types.AttributeReset:
		delete(ctx.attributes, e.Name)
		return e
	case types.AttributeSubstitution:
		return types.StringElement{
			Content: ctx.attributes.GetAsStringWithDefault(e.Name, "{"+e.Name+"}"),
		}
	case types.QuotedText:
		e.Elements = applyAttributeSubstitutionsOnElements(ctx, e.Elements)
		return e
	case types.InlineLink:
		e.Attributes = applyAttributeSubstitutionsOnAttributes(ctx, e.Attributes)
		e.Location.Path = applyAttributeSubstitutionsOnElement(ctx, e.Location.Path)
		return e
	case types.CounterSubstitution:
		return applyCounterSubstitution(ctx, e)
	default:
		return e
	}
	// case types.WithElementsToSubstitute:
	// 	elmts, err := applyAttributeSubstitutionsOnElements(ctx, e.ElementsToSubstitute())
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	element = e.ReplaceElements(types.Merge(elmts))
	// case types.WithLineSubstitution:
	// 	lines, err := applyAttributeSubstitutionsOnLines(ctx, e.LinesToSubstitute())
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	element = e.SubstituteLines(lines)
	// case types.ContinuedListItemElement:
	// 	if e.Element, err = applyAttributeSubstitutionsOnElement(ctx, e.Element); err != nil {
	// 		return nil, err
	// 	}
	// }
	// also, retain the attribute declaration value (if applicable)
}

// applyCounterSubstitutions is called by applyAttributeSubstitutionsOnElement.  Unless there is an error with
// the element (the counter is the wrong type, which should never occur), it will return a `StringElement, true`
// (because we always either find the element, or allocate one), and `nil`.  On an error it will return `nil, false`,
// and the error.  The extra boolean here is to fit the calling expectations of our caller.  This function was
// factored out of a case from applyAttributeSubstitutionsOnElement in order to reduce the complexity of that
// function, but otherwise it should have no callers.
func applyCounterSubstitution(ctx *processContext, c types.CounterSubstitution) interface{} {
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
			return types.StringElement{Content: ""}
		}
		return types.StringElement{
			Content: strconv.Itoa(counter),
		}
	case rune:
		if increment {
			counter++
		}
		ctx.counters[c.Name] = counter
		if c.Hidden {
			// return empty string facilitates merging
			return types.StringElement{Content: ""}
		}
		return types.StringElement{
			Content: string(counter),
		}
	}
	// TODO: make sure this case never happens by checking the counter value set in the context
	log.Errorf("unexpected type of counter value: '%T'", counter)
	return types.StringElement{Content: ""}
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
		log.Debugf("serialized lines:\n%s\nplaceholders: %v", result.Bytes(), spew.Sdump(placeholders))
	}
	return result.Bytes(), placeholders
}
