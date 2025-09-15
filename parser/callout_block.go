package parser

import (
	"bytes"

	"github.com/gomarkdown/markdown/ast"
)

// returns callout prefix length
func (p *Parser) calloutPrefix(data []byte) int {
	i := 0
	n := len(data)
	for i < 3 && i < n && data[i] == ' ' {
		i++
	}
	if i < n && data[i] == '(' {
		if i+1 < n && data[i+1] == ' ' {
			return i + 2
		}
		return i + 1
	}
	return 0
}

// callout ends with at least one blank line
// followed by something without a callout prefix
func (p *Parser) terminateCallout(data []byte, beg, end int) bool {
	if IsEmpty(data[beg:]) <= 0 {
		return false
	}
	if end >= len(data) {
		return true
	}
	return p.calloutPrefix(data[end:]) == 0 && IsEmpty(data[end:]) == 0
}

// parse a callout fragment
func (p *Parser) callout(data []byte) int {
	var raw bytes.Buffer
	beg, end := 0, 0
	for beg < len(data) {
		end = beg
		// Step over whole lines, collecting them. While doing that, check for
		// fenced code and if one's found, incorporate it altogether,
		// irregardless of any contents inside it
		for end < len(data) && data[end] != '\n' {
			if p.extensions&FencedCode != 0 {
				if i := p.fencedCodeBlock(data[end:], false); i > 0 {
					// -1 to compensate for the extra end++ after the loop:
					end += i - 1
					break
				}
			}
			end++
		}
		end = skipCharN(data, end, '\n', 1)
		if pre := p.calloutPrefix(data[beg:]); pre > 0 {
			// skip the prefix
			beg += pre
		} else if p.terminateCallout(data, beg, end) {
			break
		}
		// this line is part of the callout
		raw.Write(data[beg:end])
		beg = end
	}

	block := p.AddBlock(&ast.CalloutBlock{})
	p.Block(raw.Bytes())
	p.Finalize(block)
	return end
}
