package parser

import (
	"fmt"
)

type substitutionContext struct {
	substitutions map[SubstitutionKind]bool
}

func newSubstitutionContext() *substitutionContext {
	return &substitutionContext{
		substitutions: map[SubstitutionKind]bool{
			// default substitutions
			InlinePassthroughs: true,
			SpecialCharacters:  true,
			Attributes:         true,
			Quotes:             true,
			Replacements:       true,
			Macros:             true,
			PostReplacements:   true,
		},
	}
}

func (c *current) isSubstitutionEnabled(k SubstitutionKind) (bool, error) {
	ctx, ok := c.globalStore[substitutionContextKey].(*substitutionContext)
	if !ok {
		return false, fmt.Errorf("unable to look-up the substitution context in the parser's global store")
	}
	enabled, found := ctx.substitutions[k]
	if !found {
		return false, nil
	}
	return enabled, nil
}

func (ctx *substitutionContext) isSubstitutionEnabled(k SubstitutionKind) bool {
	enabled, found := ctx.substitutions[k]
	if !found {
		return false
	}
	return enabled
}

// func (ctx *substitutionContext) onAttributeDeclaration(d types.AttributeDeclaration) error {
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

// func (ctx *substitutionContext) onBlockDelimiter(d types.BlockDelimiter) {
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
// func (ctx *substitutionContext) onBlockAttributes(attrs types.Attributes) {
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
	// substitutions the key in which substitutions are stored
	substitutionContextKey string = "substitutionContext"

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

func (ctx *substitutionContext) substitutionEnabled(kind SubstitutionKind) bool {
	if value, exists := ctx.substitutions[kind]; exists {
		return value
	}
	return false
}
