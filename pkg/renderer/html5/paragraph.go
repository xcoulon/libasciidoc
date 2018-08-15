package html5

import (
	"bytes"
	"strings"
	texttemplate "text/template"

	"github.com/bytesparadise/libasciidoc/pkg/renderer"
	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var paragraphTmpl texttemplate.Template
var admonitionParagraphTmpl texttemplate.Template
var listParagraphTmpl texttemplate.Template

// initializes the templates
func init() {
	paragraphTmpl = newTextTemplate("paragraph",
		`{{ $ctx := .Context }}{{ with .Data }}{{ $renderedLines := renderLines $ctx .Lines | printf "%s" }}{{ if ne $renderedLines "" }}<div {{ if ne .ID "" }}id="{{ .ID }}" {{ end }}class="paragraph">{{ if ne .Title "" }}
<div class="doctitle">{{ .Title }}</div>{{ end }}
<p>{{ $renderedLines }}</p>
</div>{{ end }}{{ end }}`,
		texttemplate.FuncMap{
			"renderLines": renderLines,
		})

	admonitionParagraphTmpl = newTextTemplate("admonition paragraph",
		`{{ $ctx := .Context }}{{ with .Data }}{{ $renderedLines := renderLines $ctx .Lines | printf "%s"  }}{{ if ne $renderedLines "" }}<div {{ if .ID }}id="{{ .ID }}" {{ end }}class="admonitionblock {{ .Class }}">
<table>
<tr>
<td class="icon">
<div class="title">{{ .Icon }}</div>
</td>
<td class="content">{{ if .Title }}
<div class="title">{{ .Title }}</div>{{ end }}
{{ $renderedLines }}
</td>
</tr>
</table>
</div>{{ end }}{{ end }}`,
		texttemplate.FuncMap{
			"renderLines": renderLines,
		})

	listParagraphTmpl = newTextTemplate("list paragraph",
		`{{ $ctx := .Context }}{{ with .Data }}<p>{{ renderLines $ctx .Lines | printf "%s" }}</p>{{ end }}`,
		texttemplate.FuncMap{
			"renderLines": renderLines,
		})
}

func renderParagraph(ctx *renderer.Context, p types.Paragraph) ([]byte, error) {
	if len(p.Lines) == 0 {
		return make([]byte, 0), nil
	}
	result := bytes.NewBuffer(nil)
	var id, title string
	if i, ok := p.Attributes[types.AttrID].(string); ok {
		id = i
	}
	if t, ok := p.Attributes[types.AttrTitle].(string); ok {
		title = strings.TrimSpace(t)
	}
	var err error
	if _, ok := p.Attributes[types.AttrAdmonitionKind]; ok {
		log.Debug("rendering admonition paragraph...")
		k, ok := p.Attributes[types.AttrAdmonitionKind].(types.AdmonitionKind)
		if !ok {
			return nil, errors.Errorf("failed to render admonition with unknown kind: %T", p.Attributes[types.AttrAdmonitionKind])
		}
		err = admonitionParagraphTmpl.Execute(result, ContextualPipeline{
			Context: ctx,
			Data: struct {
				ID    string
				Class string
				Icon  string
				Title string
				Lines []types.InlineElements
			}{
				ID:    id,
				Class: getClass(k),
				Icon:  getIcon(k),
				Title: getTitle(p.Attributes[types.AttrTitle]),
				Lines: p.Lines,
			},
		})
	} else if ctx.WithinDelimitedBlock() || ctx.WithinList() {
		log.Debugf("rendering paragraph with %d lines within a delimited block or a list", len(p.Lines))
		err = listParagraphTmpl.Execute(result, ContextualPipeline{
			Context: ctx,
			Data: struct {
				ID    string
				Title string
				Lines []types.InlineElements
			}{
				ID:    id,
				Title: title,
				Lines: p.Lines,
			},
		})
	} else {
		log.Debug("rendering a standalone paragraph")
		err = paragraphTmpl.Execute(result, ContextualPipeline{
			Context: ctx,
			Data: struct {
				ID    string
				Title string
				Lines []types.InlineElements
			}{
				ID:    id,
				Title: title,
				Lines: p.Lines,
			},
		})

	}
	if err != nil {
		return nil, errors.Wrapf(err, "unable to render paragraph")
	}
	log.Debugf("rendered paragraph: '%s'", result.String())
	return result.Bytes(), nil
}

func getClass(kind types.AdmonitionKind) string {
	switch kind {
	case types.Tip:
		return "tip"
	case types.Note:
		return "note"
	case types.Important:
		return "important"
	case types.Warning:
		return "warning"
	case types.Caution:
		return "caution"
	default:
		log.Error("unexpected kind of admonition: %v", kind)
		return ""
	}
}

func getIcon(kind types.AdmonitionKind) string {
	switch kind {
	case types.Tip:
		return "Tip"
	case types.Note:
		return "Note"
	case types.Important:
		return "Important"
	case types.Warning:
		return "Warning"
	case types.Caution:
		return "Caution"
	default:
		log.Error("unexpected kind of admonition: %v", kind)
		return ""
	}
}

func getTitle(value interface{}) string {
	if t, ok := value.(string); ok {
		return strings.TrimSpace(t)
	}
	return ""
}