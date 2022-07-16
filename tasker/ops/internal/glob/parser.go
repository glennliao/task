package glob

import (
	"regexp"
	"strings"
)

func cleanPattern(input string) string {
	// windows path segments to posix characters
	input = strings.ReplaceAll(input, "\\", "/")
	input = strings.Replace(input, "***", "*", 1)
	input = strings.Replace(input, "**/**", "**", 1)
	input = strings.Replace(input, "**/**/**", "**", 1)

	return input
}

type ParsedPattern struct {
	// regExp is the regular expression as string that falls out of the parsedPattern
	stringPattern string
	isBaseSet     bool

	// input is the original glob pattern
	Input      string
	RegExp     *regexp.Regexp
	IsGlobstar bool
	// base is the base folder that can be used for matching a glob.
	// For example if a glob starts with `src/**/*.ts` we don't need to crawl all
	// folders in the current working directory as we see the `src` as base folder
	Base string
}

func (p *ParsedPattern) String() string {
	return p.RegExp.String()
}

func (p *ParsedPattern) setBase(str string) {
	if !p.isBaseSet {
		p.Base = str
		p.isBaseSet = true
	}
}

func (p *ParsedPattern) Compile() (*ParsedPattern, error) {
	re, err := regexp.Compile(`^` + p.stringPattern + `$`)
	if err != nil {
		return nil, err
	}
	p.RegExp = re
	return p, nil
}

func (p *ParsedPattern) add(char string) {
	p.stringPattern += char
}

func Parse(input string) (*ParsedPattern, error) {

	input = cleanPattern(input)

	parsed := ParsedPattern{
		Input:      input,
		IsGlobstar: false,
	}

	for i := 0; i < len(input); i++ {
		cur := string(input[i])

		var nextChar string
		if i < len(input)-1 {
			nextChar = string(input[i+1])
		}

		switch cur {
		case "/":
			fallthrough
		case "$":
			fallthrough
		case "^":
			fallthrough
		case "+":
			fallthrough
		case ".":
			fallthrough
		case "(":
			fallthrough
		case ")":
			fallthrough
		case "=":
			fallthrough
		case "!":
			fallthrough
		case "|":

			parsed.add(`\` + cur)

		case "?":
			parsed.add(`.`)

		case "*":

			parsed.setBase(input[0:i])

			starCount := 1

			if nextChar == "*" {
				starCount++
				i++
			}

			isGlobstar := starCount > 1

			if !parsed.IsGlobstar && isGlobstar {
				parsed.IsGlobstar = true
			}

			if isGlobstar {

				parsed.add(GLOBSTER)
				i++
			} else {

				parsed.add(STAR)
			}

		default:
			parsed.add(cur)
		}

	}

	return parsed.Compile()
}
