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
func ApplySubstitutions(ctx *processContext, done <-chan interface{}, fragmentStream <-chan types.DocumentFragment) chan types.DocumentFragment {
	processedFragmentStream := make(chan types.DocumentFragment)
	go func() {
		defer close(processedFragmentStream)
		for f := range fragmentStream {
			select {
			case <-done:
				log.WithField("pipeline_stage", "apply_substitutions").Debug("received 'done' signal")
				return
			case processedFragmentStream <- applySubstitutions(ctx, f):
			}
		}
		log.WithField("pipeline_stage", "apply_substitutions").Debug("done processing upstream content")
	}()
	return processedFragmentStream
}

func applySubstitutions(ctx *processContext, f types.DocumentFragment) types.DocumentFragment {
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.WithField("pipeline_stage", "fragment_processing").Debugf("incoming fragment:\n%s", spew.Sdump(f))
	// }
	// if the fragment already contains an error, then send it as-is downstream
	if err := f.Error; err != nil {
		log.Debugf("skipping substitutions: %v", f.Error)
		return f
	}
	elements, err := applySubstitutionsOnElements(ctx, f.Elements)
	if err != nil {
		return types.NewErrorFragment(f.LineOffset, err)
	}
	f.Elements = elements
	return f
}

func applySubstitutionsOnElements(ctx *processContext, elements []interface{}) ([]interface{}, error) {
	result := make([]interface{}, len(elements))
	for i, element := range elements {
		var err error
		if result[i], err = applySubstitutionsOnElement(ctx, element); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func applySubstitutionsOnElement(ctx *processContext, element interface{}) (interface{}, error) {
	switch e := element.(type) {
	case *types.AttributeDeclaration:
		ctx.addAttribute(e.Name, e.Value)
		return e, nil
	case *types.AttributeReset:
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
func replaceAttributeRefsInAttributes(ctx *processContext, b types.WithAttributes) error {
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debugf("applying attribute substitutions on attributes\n%s", spew.Sdump(b.GetAttributes()))
	// }
	for key, value := range b.GetAttributes() {
		switch value := value.(type) {
		case []interface{}: // multi-value attributes
			value, err := replaceAttributeRefsInElements(ctx, value)
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
			value, err := replaceAttributeRefsInElement(ctx, value)
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
	log.Debugf("processing block with elements of type '%T'", block)
	if err := processAttributes(ctx, block); err != nil {
		return err
	}
	// log.Debugf("applying substitutions on elements of block of type '%T'", block)
	subs, err := newSubstitutions(block)
	if err != nil {
		return err
	}
	elements, err := subs.processElements(ctx, block.GetElements())
	if err != nil {
		return err
	}
	if err := block.SetElements(elements); err != nil {
		return err
	}
	// also, process terms of labeled list elements
	if l, ok := block.(*types.LabeledListElement); ok {
		var err error
		if l.Term, err = subs.processElements(ctx, l.Term); err != nil {
			return err
		}
	}
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debugf("processed block with elements: %s", spew.Sdump(block))
	// }
	return nil
}

// func processLabeledListElement(ctx *processContext, block *types.LabeledListElement) error {
// 	log.Debugf("processing LabeledListElement")
// 	if err := processAttributes(ctx, block); err != nil {
// 		return err
// 	}
// 	// log.Debugf("applying substitutions on elements of block of type '%T'", block)
// 	plan, err := newSubstitutions(block)
// 	if err != nil {
// 		return err
// 	}
// 	term, err := plan.processElements(ctx, block.Term)
// 	if err != nil {
// 		return err
// 	}
// 	block.Term = term
// 	elements, err := plan.processElements(ctx, block.GetElements())
// 	if err != nil {
// 		return err
// 	}
// 	if err := block.SetElements(elements); err != nil {
// 		return errors.Wrapf(err, "failed to process substitutions on block of type '%T'", block)
// 	}
// 	return nil
// }

func processBlockWithLocation(ctx *processContext, block types.WithLocation) error {
	log.Debugf("processing block with attributes and location")
	if err := processAttributes(ctx, block); err != nil {
		return err
	}
	log.Debugf("applying substitutions on `location` of block of type '%T'", block)
	elements := block.GetLocation().Path
	elements, err := replaceAttributeRefsInElements(ctx, elements)
	if err != nil {
		return err
	}
	block.GetLocation().SetPath(elements)
	imagesdir := ctx.attributes.GetAsStringWithDefault("imagesdir", "")
	block.GetLocation().SetPathPrefix(imagesdir)
	return nil
}

func processAttributes(ctx *processContext, block types.WithAttributes) error {
	if err := replaceAttributeRefsInAttributes(ctx, block); err != nil {
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

type substitutions []*substitution

func newSubstitutions(b types.WithAttributes) (substitutions, error) {
	// TODO: introduce a `types.BlockWithSubstitution` interface?
	// note: would also be helpful for paragraphs with `[listing]` style.
	defaultSub, err := defaultSubstitution(b)
	if err != nil {
		return nil, errors.Wrap(err, "unable to determine substitutions")
	}
	subs := strings.Split(b.GetAttributes().GetAsStringWithDefault(types.AttrSubstitutions, defaultSub), ",")
	allIncremental, err := ValidateSubstitutions(subs)
	if err != nil {
		return nil, err
	}
	// when dealing with incremental substitutions
	if allIncremental {
		d, err := newSubstitution(defaultSub)
		if err != nil {
			return nil, err
		}
		result := substitutions{d}
		for _, sub := range subs {
			switch {
			case strings.HasSuffix(sub, "+"): // prepend
				s, err := newSubstitution(strings.TrimSuffix(sub, "+"))
				if err != nil {
					return nil, err
				}
				result = append(substitutions{s}, result...)
			case strings.HasPrefix(sub, "+"): // append
				s, err := newSubstitution(strings.TrimPrefix(sub, "+"))
				if err != nil {
					return nil, err
				}
				result = append(result, s)
			case strings.HasPrefix(sub, "-"): // remove from all substitutions
				for _, s := range result {
					s.disable(substitutionKind(strings.TrimPrefix(sub, "-")))
				}
			}
		}
		return result, nil
	}

	result := make([]*substitution, len(subs))
	for i, sub := range subs {
		if result[i], err = newSubstitution(sub); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// can't mix incremental additions/removals with sets. eg:
// "attributes+,+replacements,-callouts" // OK
// "attributes+,normal" // NOT OK
func ValidateSubstitutions(subs []string) (bool, error) {
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("checking substitutions: %v", subs)
	}
	if len(subs) == 0 { // should not happen, there is always a default sub
		return false, nil
	}
	// init with first substitution
	allIncremental := isIncrementalSubstitution(subs[0])
	// check others
	for i, sub := range subs {
		if i == 0 {
			// skip, no need to check this one
			continue
		}
		if isIncrementalSubstitution(sub) != allIncremental {
			return false, fmt.Errorf("cannot mix incremental and non-incremental substitutions: '%s'", strings.Join(subs, ", "))
		}
	}
	return allIncremental, nil
}

func isIncrementalSubstitution(sub string) bool {
	return strings.HasPrefix(sub, "+") ||
		strings.HasPrefix(sub, "-") ||
		strings.HasSuffix(sub, "+")
}
func (s substitutions) processElements(ctx *processContext, elements []interface{}) ([]interface{}, error) {
	// skip if there's nothing to do
	// (and no need to return an empty slice, btw)
	if len(elements) == 0 {
		return nil, nil
	}
	for _, substitution := range s {
		var err error
		elements, err = substitution.processElements(ctx, elements)
		if err != nil {
			return nil, err
		}
	}
	return elements, nil
}

func defaultSubstitution(b interface{}) (string, error) {
	switch b := b.(type) {
	case *types.DelimitedBlock:
		switch b.Kind {
		case types.Listing, types.Fenced, types.Literal:
			return "verbatim", nil
		case types.Example, types.Quote:
			return "normal", nil
		default:
			return "", fmt.Errorf("unsupported kind of delimited block: '%v'", b.Kind)
		}
	case *types.Paragraph, *types.GenericList, types.ListElement, *types.QuotedText:
		return "normal", nil
	case *types.Section:
		return "header", nil
	default:
		return "", fmt.Errorf("unsupported kind of element: '%T'", b)
	}
}

type substitutionRule string

const (
	AttributesGroup        substitutionRule = "AttributesGroup"
	ElementAttributesGroup substitutionRule = "ElementAttributesGroup"
	HeaderGroup            substitutionRule = "HeaderGroup"
	MacrosGroup            substitutionRule = "MacrosGroup"
	NoneGroup              substitutionRule = "NoneGroup"
	NormalGroup            substitutionRule = "NormalGroup"
	QuotesGroup            substitutionRule = "QuotesGroup"
	ReplacementsGroup      substitutionRule = "ReplacementsGroup"
	PostReplacementsGroup  substitutionRule = "PostReplacementsGroup"
	SpecialcharactersGroup substitutionRule = "SpecialCharactersGroup"
	VerbatimGroup          substitutionRule = "VerbatimGroup"
)

type substitution struct {
	rule                      substitutionRule
	enablements               map[substitutionKind]bool
	hasAttributeSubstitutions bool // TODO: replace with a key in the parser's store
}

type SubstitutionGroup struct {
}

func newSubstitution(kind string) (*substitution, error) {
	switch kind {
	case "attributes":
		return &substitution{
			rule: AttributesGroup,
			enablements: map[substitutionKind]bool{
				InlinePassthroughs: true,
				Attributes:         true,
			},
		}, nil
	case "element_attributes":
		return &substitution{
			rule: ElementAttributesGroup,
			enablements: map[substitutionKind]bool{
				InlinePassthroughs: true,
				Attributes:         true,
				Quotes:             true,
				SpecialCharacters:  true, // TODO: is it needed?
			},
		}, nil
	case "header":
		return &substitution{
			rule: HeaderGroup,
			enablements: map[substitutionKind]bool{
				InlinePassthroughs: true,
				SpecialCharacters:  true,
				Attributes:         true,
			},
		}, nil
	case "macros":
		return &substitution{
			rule: MacrosGroup,
			enablements: map[substitutionKind]bool{
				Macros: true,
			},
		}, nil
	case "normal":
		return &substitution{
			rule: NormalGroup,
			enablements: map[substitutionKind]bool{
				InlinePassthroughs: true,
				SpecialCharacters:  true,
				Attributes:         true,
				Quotes:             true,
				Replacements:       true,
				Macros:             true,
				PostReplacements:   true,
			},
		}, nil
	case "none":
		return &substitution{
			rule:        NoneGroup,
			enablements: map[substitutionKind]bool{},
		}, nil
	case "quotes":
		return &substitution{
			rule: QuotesGroup,
			enablements: map[substitutionKind]bool{
				Quotes: true,
			},
		}, nil
	case "replacements":
		return &substitution{
			rule: ReplacementsGroup,
			enablements: map[substitutionKind]bool{
				Replacements: true,
			},
		}, nil
	case "post_replacements":
		return &substitution{
			rule: PostReplacementsGroup,
			enablements: map[substitutionKind]bool{
				PostReplacements: true,
			},
		}, nil
	case "specialchars":
		return &substitution{
			rule: SpecialcharactersGroup,
			enablements: map[substitutionKind]bool{
				SpecialCharacters: true,
			},
		}, nil
	case "verbatim":
		return &substitution{
			rule: VerbatimGroup,
			enablements: map[substitutionKind]bool{
				SpecialCharacters: true,
				Callouts:          true,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported kind of substitution: '%v'", kind)
	}
}

func (s *substitution) disable(kinds ...substitutionKind) {
	for _, k := range kinds {
		switch k {
		case "specialchars":
			delete(s.enablements, SpecialCharacters)
		default:
			delete(s.enablements, k)
		}
	}
}

func (s *substitution) hasEnablements() bool {
	return len(s.enablements) > 0
}

// // disables the "inline_passthroughs", "special_characters" and "attributes" substitution
// // return `true` if there are more enablements to apply, `false` otherwise (ie, no substitution would be applied if the content was parsed again)
// func (s *substitution) reduce() bool { // TODO: rename this func
// 	s.enablements = s.defaults() // reset
// 	for sub := range s.enablements {
// 		switch sub {
// 		case InlinePassthroughs, SpecialCharacters, Attributes:
// 			delete(s.enablements, sub)
// 		}
// 	}
// 	// if log.IsLevelEnabled(log.DebugLevel) {
// 	// 	log.Debugf("new enablements for '%s': %s", s.group, spew.Sdump(s.enablements))
// 	// }
// 	return len(s.enablements) > 0
// }

func (s *substitution) processElements(ctx *processContext, elements []interface{}) ([]interface{}, error) {
	// log.Debugf("applying step '%s'", step.group)
	elements, err := parseElements(ctx, elements, s.rule, GlobalStore(substitutionsKey, s))
	if err != nil {
		return nil, err
	}
	if s.hasAttributeSubstitutions {
		// log.WithField("group", step.group).Debug("attribute substitutions detected during parsing")
		// apply substitutions on elements
		elements, err = replaceAttributeRefsInElements(ctx, elements)
		if err != nil {
			return nil, err
		}
		elements = types.Merge(elements)
		// re-run the parser, skipping the `inline_passthrough` and `attribute` rules this time
		s.disable(InlinePassthroughs, Attributes)
		if s.hasEnablements() {
			if elements, err = parseElements(ctx, elements, s.rule, GlobalStore(substitutionsKey, s)); err != nil {
				return nil, err
			}
		}
	}
	return elements, nil
}

// sets the new substitution plan in the golbal store, overriding any exist–ing one
// NOTE: will override any existing substitution context
// TODO: is there any case where a stack would be needed, so we don't override an existing context?
func (c *current) setCurrentSubstitution(kind string) error {
	p, err := newSubstitution(kind)
	if err != nil {
		return err
	}
	c.globalStore[substitutionsKey] = p
	return nil
}

func (c *current) unsetCurrentSubstitution() {
	delete(c.globalStore, substitutionsKey)
}

func (c *current) lookupCurrentSubstitution() (*substitution, error) {
	ctx, ok := c.globalStore[substitutionsKey].(*substitution)
	if !ok {
		return nil, fmt.Errorf("unable to look-up the substitution context in the parser's global store")
	}
	return ctx, nil
}

// called when an attribute substitution occurred
// TODO: find a better name for this method
func (c *current) hasAttributeSubstitutions() error {
	phase, err := c.lookupCurrentSubstitution()
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

func (c *current) isSubstitutionEnabled(k substitutionKind) (bool, error) {
	s, err := c.lookupCurrentSubstitution()
	if err != nil {
		return false, err
	}
	enabled, found := s.enablements[k]
	if !found {
		return false, nil
	}
	return enabled, nil
}

type substitutionKind string

const (
	// substitutionsKey the key in which substitutions are stored in the parser's GlobalStore
	substitutionsKey string = "substitutions"

	// Attributes the "attributes" substitution
	Attributes substitutionKind = "attributes"
	// Callouts the "callouts" substitution
	Callouts substitutionKind = "callouts"
	// InlinePassthroughs the "inline_passthrough" substitution
	InlinePassthroughs substitutionKind = "inline_passthrough"
	// Macros the "macros" substitution
	Macros substitutionKind = "macros"
	// None the "none" substitution
	None substitutionKind = "none"
	// PostReplacements the "post_replacements" substitution
	PostReplacements substitutionKind = "post_replacements"
	// Quotes the "quotes" substitution
	Quotes substitutionKind = "quotes"
	// Replacements the "replacements" substitution
	Replacements substitutionKind = "replacements"
	// SpecialCharacters the "specialchars" substitution
	SpecialCharacters substitutionKind = "specialchars"
)

func parseElements(ctx *processContext, elements []interface{}, rule substitutionRule, opts ...Option) ([]interface{}, error) {
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.WithField("group", group).Debug("parsing elements")
	// }
	serialized, placeholders := serialize(elements)
	if len(serialized) == 0 {
		return nil, nil
	}
	opts = append(opts)
	result, err := parseContent(serialized, append(opts, Entrypoint(string(rule)))...)
	if err != nil {
		return nil, err
	}
	// also, apply the substitutions on the placeholders, case by case
	for _, element := range placeholders.elements {
		log.Debugf("processing placeholder of type '%T'", element)
		// if e, ok := element.(types.WithAttributes); ok {
		// 	attrs, err := parseAttributes(e.GetAttributes(), rule, opts...)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	e.SetAttributes(attrs)
		// }
		// if e, ok := element.(*types.LabeledListElement); ok {
		// 	term, err := parseElements(e.Term, rule, opts...)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	e.Term = term
		// }
		if w, ok := element.(types.WithElements); ok {
			if err := processBlockWithElements(ctx, w); err != nil {
				return nil, err
			}
		} else if elements, ok := element.([]interface{}); ok {
			for _, e := range elements {
				if w, ok := e.(types.WithElements); ok {
					if err := processBlockWithElements(ctx, w); err != nil {
						return nil, err
					}
				}
			}
		} else {
			log.Debugf("skipping substitutions on block of type '%T'", element)
		}
		// 	elements, err := parseElements(e, rule, opts...)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	placeholders.elements[key] = elements
		// }
		// if e, ok := element.(types.WithLocation); ok {
		// 	if err := parseElementWithLocation(e, group, opts...); err != nil {
		// 		return nil, err
		// 	}
		// }
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debugf("processed placeholder: %s", spew.Sdump(element))
		}
	}
	result = placeholders.restoreElements(result)
	return result, nil
}

func parseAttributes(attributes types.Attributes, rule substitutionRule, opts ...Option) (types.Attributes, error) {
	if !(rule == AttributesGroup || rule == QuotesGroup) { // TODO: include special_characters?
		// log.Debugf("no need to parse attributes for group '%s'", rule)
		return attributes, nil
	}
	for name, value := range attributes {
		switch value := value.(type) {
		case []interface{}:
			serialized, placeholders := serialize(value)
			if len(serialized) == 0 {
				continue
			}
			elements, err := parseContent(serialized, append(opts, Entrypoint(string(rule)))...)
			if err != nil {
				return nil, err
			}
			elements = placeholders.restoreElements(elements)
			attributes[name] = elements
		case string:
			elements, err := parseContent([]byte(value), append(opts, Entrypoint(string(rule)))...)
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
func replaceAttributeRefsInElement(ctx *processContext, element interface{}) (interface{}, error) {
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("replacing attribute references in element of type '%T'", element)
	}
	switch e := element.(type) {
	case []interface{}:
		return replaceAttributeRefsInElements(ctx, e)
	case types.AttributeSubstitution:
		return &types.StringElement{
			Content: ctx.attributes.GetAsStringWithDefault(e.Name, "{"+e.Name+"}"),
		}, nil
	case types.CounterSubstitution:
		return applyCounterSubstitution(ctx, e)
	case types.WithElements:
		// replace AttributeSubstitutions on attributes
		if err := replaceAttributeRefsInAttributes(ctx, e); err != nil {
			return nil, err
		}
		// replace AttributeSubstitutions on nested elements
		elements, err := replaceAttributeRefsInElements(ctx, e.GetElements())
		if err != nil {
			return nil, err
		}
		// elements = types.Merge(elements)
		if err := e.SetElements(elements); err != nil {
			return nil, errors.Wrapf(err, "failed to replace attribute references in block of type '%T'", e)
		}
		return e, nil
	case types.WithLocation:
		// replace AttributeSubstitutions on attributes
		if err := replaceAttributeRefsInAttributes(ctx, e); err != nil {
			return nil, err
		}
		// replace AttributeSubstitutions on embedded location
		if err := replaceAttributeRefsInLocation(ctx, e); err != nil {
			return nil, err
		}
		return e, nil
	case types.WithAttributes:
		// replace AttributeSubstitutions on attributes
		if err := replaceAttributeRefsInAttributes(ctx, e); err != nil {
			return nil, err
		}
		return e, nil
	default:
		return e, nil
	}
}

func replaceAttributeRefsInLocation(ctx *processContext, b types.WithLocation) error {
	path, err := replaceAttributeRefsInElements(ctx, b.GetLocation().Path)
	if err != nil {
		return err
	}
	b.GetLocation().Path = path
	return nil
}

func replaceAttributeRefsInElements(ctx *processContext, elements []interface{}) ([]interface{}, error) {
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("replacing attribute refs in elements:\n%s", spew.Sdump(elements))
	}
	result := make([]interface{}, len(elements)) // maximum capacity should exceed initial input
	for i, element := range elements {
		element, err := replaceAttributeRefsInElement(ctx, element)
		if err != nil {
			return nil, err
		}
		result[i] = element
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("replaced attribute refs in elements:\n%s", spew.Sdump(result))
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
			return &types.StringElement{Content: ""}, nil
		}
		return &types.StringElement{
			Content: strconv.Itoa(counter),
		}, nil
	case rune:
		if increment {
			counter++
		}
		ctx.counters[c.Name] = counter
		if c.Hidden {
			// return empty string facilitates merging
			return &types.StringElement{Content: ""}, nil
		}
		return &types.StringElement{
			Content: string(counter),
		}, nil
	default:
		return &types.StringElement{}, fmt.Errorf("unexpected type of counter value: '%T'", counter)
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
			case *types.StringElement:
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
