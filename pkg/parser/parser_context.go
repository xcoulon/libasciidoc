package parser

import (
	"github.com/bytesparadise/libasciidoc/pkg/configuration"
	"github.com/bytesparadise/libasciidoc/pkg/types"

	log "github.com/sirupsen/logrus"
)

type parserContext struct {
	config        configuration.Configuration
	attributes    map[string]string
	substitutions map[SubstitutionKind]bool
	levelOffsets  []levelOffset
	blockLevels   *stack
}

func newParserContext(config configuration.Configuration) *parserContext {
	return &parserContext{
		config:       config,
		attributes:   map[string]string{},
		levelOffsets: []levelOffset{},
		blockLevels:  newStack(),
		substitutions: map[SubstitutionKind]bool{
			// default substitutions
			AttributesSubstitution: true,
		},
	}
}

func (ctx *parserContext) clone() *parserContext {
	return &parserContext{
		config:        ctx.config,
		attributes:    ctx.attributes, // TODO: should we clone this too? ie, can an attribute declared in a child doc be used in rest of the parent doc?
		levelOffsets:  append([]levelOffset{}, ctx.levelOffsets...),
		blockLevels:   ctx.blockLevels,
		substitutions: ctx.substitutions,
	}
}

func (ctx *parserContext) isSectionRuleEnabled() bool {
	return ctx.blockLevels.empty()
}

func (ctx *parserContext) onBlockDelimiter(d types.BlockDelimiter) {
	currentLevel := ctx.blockLevels.get()
	if currentLevel == d.Kind {
		ctx.blockLevels.pop() // discard current level, assuming we've just parsed the ending delimiter of the current block
		return
	}
	ctx.blockLevels.push(d.Kind) // push current delimiter kind, assuming we've just parsed the starting delimiter of a new block
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("current substitution context is now: %v", ctx.blockLevels.get())
	}
}

// called when block attributes are initialized.
// if the `subs` attribute is set, it is retained in the parser context
// and used to enable/disable subsequent parser rules
// TODO: take into account blank lines and block delimiters to reset the substitutions to apply
func (ctx *parserContext) onBlockAttributes(attrs types.Attributes) {
	subs, exists := attrs[types.AttrSubstitutions]
	if !exists {
		return
	}
	switch subs {
	// 		case "normal":
	// 			subs = subs.append(
	// 				"specialcharacters",
	// 				"quotes",
	// 				"attributes",
	// 				"replacements",
	// 				"macros",
	// 				"post_replacements",
	// 			)
	// 		case "inline_passthrough", "callouts", "specialcharacters", "specialchars", "quotes", "attributes", "macros", "replacements", "post_replacements", "none":
	// 			subs = subs.append(s)
	// 		case "+callouts", "+specialcharacters", "+specialchars", "+quotes", "+attributes", "+macros", "+replacements", "+post_replacements", "+none":
	// 			if len(subs) == 0 {
	// 				subs = subs.append(block.DefaultSubstitutions()...)
	// 			}
	// 			subs = subs.append(strings.ReplaceAll(s, "+", ""))
	// 		case "callouts+", "specialcharacters+", "specialchars+", "quotes+", "attributes+", "macros+", "replacements+", "post_replacements+", "none+":
	// 			if len(subs) == 0 {
	// 				subs = subs.append(block.DefaultSubstitutions()...)
	// 			}
	// 			subs = subs.prepend(strings.ReplaceAll(s, "+", ""))
	// 		case "-callouts", "-specialcharacters", "-specialchars", "-quotes", "-attributes", "-macros", "-replacements", "-post_replacements", "-none":
	// 			if len(subs) == 0 {
	// 				subs = subs.append(block.DefaultSubstitutions()...)
	// 			}
	// 			subs = subs.remove(strings.ReplaceAll(s, "-", ""))
	// 		default:
	// 			return nil, fmt.Errorf("unsupported substitution: '%s", s)
	// 		}
	}
}

type SubstitutionKind string

const (
	// substitutions the key in which substitutions are stored
	substitutionsKey string = "substitution"

	// InlinePassthroughsSubstitution the "inline_passthrough" substitution
	InlinePassthroughSubstitution SubstitutionKind = "inline_passthrough"
	// AttributesSubstitution the "attributes" substitution
	AttributesSubstitution SubstitutionKind = "attributes"
	// SpecialCharactersSubstitution the "specialcharacters" substitution
	SpecialCharactersSubstitution SubstitutionKind = "specialcharacters"
	// CalloutsSubstitution the "callouts" substitution
	CalloutsSubstitution SubstitutionKind = "callouts"
	// QuotesSubstitution the "quotes" substitution
	QuotesSubstitution SubstitutionKind = "quotes"
	// ReplacementsSubstitution the "replacements" substitution
	ReplacementsSubstitution SubstitutionKind = "replacements"
	// MacrosSubstitution the "macros" substitution
	MacrosSubstitution SubstitutionKind = "macros"
	// PortReplacementsSubstitution the "post_replacements" substitution
	PortReplacementsSubstitution SubstitutionKind = "post_replacements"
	// NoneSubstitution the "none" substitution
	NoneSubstitution SubstitutionKind = "none"
)

func (ctx *parserContext) substitutionEnabled(kind SubstitutionKind) bool {
	if value, exists := ctx.substitutions[kind]; exists {
		return value
	}
	return false
}
