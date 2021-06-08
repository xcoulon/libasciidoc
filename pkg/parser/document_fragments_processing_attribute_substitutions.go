package parser

// TODO: rename this func
// func onInlineElements(ctx *substitutionContext, elements types.InlineElements) (types.InlineElements, error) {
// 	if log.IsLevelEnabled(log.DebugLevel) {
// 		log.Debug("inline elements:")
// 		spew.Fdump(log.StandardLogger().Out, elements)
// 	}
// 	// if attribute substitution is enabled...
// 	if !ctx.substitutionEnabled(AttributesSubstitution) {
// 		log.Debug("attributes substitution is not enabled")
// 		return elements, nil
// 	}
// 	// ... and if there are some attributes to substitute
// 	if !elements.HasAttributeSubstitutions() {
// 		log.Debug("no attributes to substitute")
// 		return elements, nil
// 	}
// 	elmts := substituteAttributes(elements, ctx.attributes)
// 	placeholders := newPlaceholders()
// 	line := serialize(elmts, placeholders)
// 	elmts, err := ParseReader("", strings.NewReader(line), Entrypoint("InlineElements"))
// 	if err != nil {
// 		return nil, err
// 	}
// 	result := restoreElements(elmts.(types.InlineElements), placeholders)
// 	return result, nil
// }

// func substituteAttributes(content interface{}, attributes map[string]string) interface{} {
// 	switch element := content.(type) {
// 	case types.InlineElements:
// 		for i, elmt := range element {
// 			element[i] = substituteAttributes(elmt, attributes)
// 		}
// 		return types.Merge(element...)
// 	case []interface{}:
// 		for i, elmt := range element {
// 			element[i] = substituteAttributes(elmt, attributes)
// 		}
// 		return types.Merge(element...)
// 	case types.AttributeSubstitution:
// 		if value, found := attributes[element.Name]; found {
// 			return types.StringElement{
// 				Content: value,
// 			}
// 		}
// 		log.Debugf("unable to substitute attribute '%s': no match found", element.Name)
// 		return types.StringElement{
// 			Content: "{" + element.Name + "}",
// 		}
// 	default:
// 		// do nothing, return as-is
// 		return content
// 	}
// }

// type placeholders struct {
// 	seq      int
// 	elements map[string]interface{}
// }

// func newPlaceholders() *placeholders {
// 	return &placeholders{
// 		seq:      0,
// 		elements: map[string]interface{}{},
// 	}
// }

// func (p *placeholders) add(element interface{}) types.ElementPlaceHolder {
// 	p.seq++
// 	p.elements[strconv.Itoa(p.seq)] = element
// 	return types.ElementPlaceHolder{
// 		Ref: strconv.Itoa(p.seq),
// 	}

// }

// func serialize(element interface{}, placeholders *placeholders) string {
// 	result := strings.Builder{}
// 	switch element := element.(type) {
// 	case types.InlineElements:
// 		for _, elmt := range element {
// 			result.WriteString(serialize(elmt, placeholders))
// 		}
// 	case []interface{}:
// 		for _, elmt := range element {
// 			result.WriteString(serialize(elmt, placeholders))
// 		}
// 	case types.StringElement:
// 		result.WriteString(element.Content)
// 	case types.SingleLineComment:
// 		// replace with placeholder
// 		p := placeholders.add(element)
// 		result.WriteString(p.String())
// 	default:
// 		// replace with placeholder
// 		p := placeholders.add(element)
// 		result.WriteString(p.String())
// 	}
// 	// if log.IsLevelEnabled(log.DebugLevel) {
// 	// 	log.Debug("serialized line:")
// 	// 	spew.Fdump(log.StandardLogger().Out, result.String())
// 	// }
// 	return result.String()
// }

// // replace the placeholders with their original element in the given elements
// func restoreElements(elements []interface{}, placeholders *placeholders) []interface{} {
// 	// skip if there's nothing to restore
// 	if len(placeholders.elements) == 0 {
// 		return elements
// 	}
// 	for i, e := range elements {
// 		//
// 		if e, ok := e.(types.ElementPlaceHolder); ok {
// 			elements[i] = placeholders.elements[e.Ref]
// 		}
// 		// // for each element, check *all* interfaces to see if there's a need to replace the placeholders
// 		// if e, ok := e.(types.WithPlaceholdersInElements); ok {
// 		// 	elements[i] = e.RestoreElements(placeholders.elements)
// 		// }
// 		// if e, ok := e.(types.WithPlaceholdersInAttributes); ok {
// 		// 	elements[i] = e.RestoreAttributes(placeholders.elements)
// 		// }
// 		// if e, ok := e.(types.WithPlaceholdersInLocation); ok {
// 		// 	elements[i] = e.RestoreLocation(placeholders.elements)
// 		// }
// 	}
// 	return elements
// }
