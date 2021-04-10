package parser

import (
	"fmt"

	"github.com/bytesparadise/libasciidoc/pkg/configuration"
	"github.com/bytesparadise/libasciidoc/pkg/types"
)

// extra methods on the generated parser's `storeDict` type

const attributesKey = "attributes"

const parserContextKey = "parserContext"

const usermacrosKey = "user_macros"

func (c storeDict) pushAttributes(value interface{}) {
	if s, ok := c[attributesKey].(*stack); ok {
		s.push(value)
		return
	}
	s := newStack()
	s.push(value)
	c[attributesKey] = s
}

func (c storeDict) discardAttributes() {
	if s, ok := c[attributesKey].(*stack); ok {
		s.pop()
	}
}

func (c storeDict) getAttributes() interface{} {
	if s, ok := c[attributesKey].(*stack); ok {
		return s.get()
	}
	return nil
}

func (c storeDict) hasUserMacro(name string) bool {
	if macros, exists := c[usermacrosKey].(map[string]configuration.MacroTemplate); exists {
		_, exists := macros[name]
		return exists
	}
	return false
}

func (c *current) parserContext() (*parserContext, error) {
	if ctx, ok := c.globalStore[parserContextKey].(*parserContext); ok {
		return ctx, nil
	}
	return nil, fmt.Errorf("unable to look-up the parse context in the parser's global store")
}

func (c *current) isCommentBlockContentEnabled(kind types.DelimiterKind) (bool, error) {
	if ctx, ok := c.globalStore[parserContextKey].(*parserContext); ok {
		return ctx.isCommentBlockContentEnabled(kind), nil
	}
	return false, fmt.Errorf("unable to look-up the parse context in the parser's global store")

}

// func (c storeDict) pushSubsitutionContext(ctx substitutionContext) {
// 	if s, ok := c[substitutionContextKey].(*stack); ok {
// 		s.push(value)
// 		return
// 	}
// 	s := newStack()
// 	s.push(value)
// 	c[substitutionContextKey] = s
// }

// func (c storeDict) discardSubstitutionContext() {
// 	if s, ok := c[substitutionContextKey].(*stack); ok {
// 		s.pop()
// 	}
// }

// func (c storeDict) getSubstitutionContext() (*substitutionContext, error) {
// 	if ctx, ok := c[substitutionContextKey].(*substitutionContext); ok {
// 		return ctx, nil
// 	}
// 	return nil, fmt.Errorf("unable to look-up the substitution context in the parser's global store")
// }

// type substitutionContext struct {
// 	// within which kind of block are we currently parsing some content?
// 	// empty by default, ie, currently parsing blocks at the document level (sections, paragraphs, delimited blocks, etc.)
// 	// but we could have nested delimited blocks, too. Hence a stack.
// 	blockLevels *stack
// }

// func newSubstitutionContext() *substitutionContext {
// 	return &substitutionContext{
// 		blockLevels: newStack(),
// 	}
// }
// func (ctx *substitutionContext) isSectionRuleEnabled() bool {
// 	return ctx.blockLevels.empty()
// }

// func (ctx *substitutionContext) OnBlockDelimiter(d types.BlockDelimiter) {
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

// const attributeDeclatationsKey = "attributeDeclatations"

// func (c storeDict) AddAttributeDeclaration(d types.AttributeDeclaration) {
// 	if attrDeclarations, ok := c[attributeDeclatationsKey].(map[string]interface{}); ok {
// 		attrDeclarations[d.Name] = d.Value
// 		return
// 	}
// 	c[attributeDeclatationsKey] = map[string]interface{}{
// 		d.Name: d.Value,
// 	}
// }

// func (c storeDict) AttributeDeclarations() map[string]interface{} {
// 	if attrDeclarations, ok := c[attributeDeclatationsKey].(map[string]interface{}); ok {
// 		return attrDeclarations
// 	}
// 	return nil
// }
