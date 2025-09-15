package markdown

import (
	"testing"

	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func TestCalloutBlock(t *testing.T) {
	tests := []string{
		"( Callout line 1\n( Callout line 2\n\nAfter\n",
		"<div class=\"callout\">\n<p>Callout line 1\nCallout line 2</p>\n</div>\n<p>After</p>\n",
	}
	params := TestParams{
		extensions: parser.CommonExtensions | parser.Callouts,
		Flags:      html.UseXHTML,
	}
	doTestsParam(t, tests, params)
}
