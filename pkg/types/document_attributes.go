package types

// DocumentAttributes the document attributes
type DocumentAttributes map[string]interface{}

const (
	// AttrDocType the "doctype" attribute
	AttrDocType string = "doctype"
	// AttrSyntaxHighlighter the attribute to define the syntax highlighter on code source blocks
	AttrSyntaxHighlighter string = "source-highlighter"
	// AttrIDPrefix the key to retrieve the ID Prefix
	AttrIDPrefix string = "idprefix"
	// DefaultIDPrefix the default ID Prefix
	DefaultIDPrefix string = "_"
	// AttrTableOfContents the `toc` attribute at document level
	AttrTableOfContents string = "toc"
	// AttrTableOfContentsLevels the document attribute which specifies the number of levels to display in the ToC
	AttrTableOfContentsLevels string = "toclevels"
	// AttrNoHeader attribute to disable the rendering of document footer
	AttrNoHeader string = "noheader"
	// AttrNoFooter attribute to disable the rendering of document footer
	AttrNoFooter string = "nofooter"
)

// Has returns the true if an entry with the given key exists
func (a DocumentAttributes) Has(key string) bool {
	_, ok := a[key]
	return ok
}

// AddAll adds the given attributes
func (a DocumentAttributes) AddAll(attrs map[string]interface{}) {
	for k, v := range attrs {
		a.Add(k, v)
	}
}

// Add adds the given attribute if its value is non-nil
// TODO: raise a warning if there was already a name/value
func (a DocumentAttributes) Add(key string, value interface{}) {
	a[key] = value
}

// AddNonEmpty adds the given attribute if its value is non-nil and non-empty
// TODO: raise a warning if there was already a name/value
func (a DocumentAttributes) AddNonEmpty(key string, value interface{}) {
	// do not add nil or empty values
	if value == "" {
		return
	}
	a.Add(key, value)
}

// AddDeclaration adds the given attribute
// TODO: raise a warning if there was already a name/value
func (a DocumentAttributes) AddDeclaration(attr DocumentAttributeDeclaration) {
	a.Add(attr.Name, attr.Value)
}

// Delete deletes the given attribute
func (a DocumentAttributes) Delete(attr DocumentAttributeReset) {
	delete(a, attr.Name)
}

// GetAsString gets the string value for the given key (+ `true`),
// or empty string (+ `false`) if none was found
// TODO: raise a warning if there was no entry found
func (a DocumentAttributes) GetAsString(key string) (string, bool) {
	// check in predefined attributes
	if value, found := Predefined[key]; found {
		return value, true
	}

	if value, found := a[key]; found {
		if value, ok := value.(string); ok {
			return value, true
		}
	}
	return "", false
}

// GetAsStringWithDefault gets the string value for the given key,
// or returns the given default value
// TODO: raise a warning if there was no entry found
func (a DocumentAttributes) GetAsStringWithDefault(key, defaultValue string) string {
	// check in predefined attributes
	if value, found := Predefined[key]; found {
		return value
	}
	if value, found := a[key]; found {
		if value, ok := value.(string); ok {
			return value
		}
	}
	return defaultValue
}

// GetAuthors returns the authors or empty slice if none was set in the document
func (a DocumentAttributes) GetAuthors() []DocumentAuthor {
	if authors, ok := a[AttrAuthors].([]DocumentAuthor); ok {
		return authors
	}
	return []DocumentAuthor{}
}

// DocumentAttributesWithOverrides the document attributes with some overrides provided by the CLI (for example)
type DocumentAttributesWithOverrides struct {
	Content   map[string]interface{}
	Overrides map[string]string
}

// All returns all attributes
func (a DocumentAttributesWithOverrides) All() DocumentAttributes {
	result := DocumentAttributes{}
	for k, v := range a.Content {
		result[k] = v
	}
	for k, v := range a.Overrides {
		result[k] = v
	}
	return result
}

// Add add the given attribute
func (a DocumentAttributesWithOverrides) Add(key string, value interface{}) {
	a.Content[key] = value
}

// AddAll adds the given attributes
func (a DocumentAttributesWithOverrides) AddAll(attrs map[string]interface{}) {
	for k, v := range attrs {
		a.Content[k] = v
	}
}

// Delete deletes the given attribute
func (a DocumentAttributesWithOverrides) Delete(key string) {
	delete(a.Content, key)
}

// GetAsString gets the string value for the given key (+ `true`),
// or empty string (+ `false`) if none was found
func (a DocumentAttributesWithOverrides) GetAsString(key string) (string, bool) {
	// if value is overridden
	if value, found := a.Overrides[key]; found {
		return value, true
	}
	// check in predefined attributes
	if value, found := Predefined[key]; found {
		return value, true
	}
	// if value is reset
	if _, found := a.Overrides["!"+key]; found {
		return "", false
	}
	if value, found := a.Content[key].(string); found {
		return value, true
	}
	// TODO: raise a warning if there was no entry found
	return "", false
}

// GetAsStringWithDefault gets the string value for the given key,
// or returns the given default value
func (a DocumentAttributesWithOverrides) GetAsStringWithDefault(key, defaultValue string) string {
	if value, found := a.Overrides[key]; found {
		return value
	}
	// check in predefined attributes
	if value, found := Predefined[key]; found {
		return value
	}
	if value, found := a.Content[key].(string); found {
		return value
	}
	// TODO: raise a warning if there was no entry found
	return defaultValue
}
