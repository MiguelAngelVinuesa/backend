package rng_magic

import (
	"fmt"
	"regexp"
	"strings"

	magic "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/analysis/simulation/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
)

func parse(statement string, conditions map[string]*magic.Condition, symbols *slots.SymbolSet) (Function, error) {
	subs := rx1.FindAllString(statement, -1)
	if len(subs) == 0 {
		return nil, fmt.Errorf(msgInvalidStatement)
	}

	var level int
	var boolF, not bool
	var f Function
	var stack []Function

	for ix := range subs {
		sub := strings.Trim(subs[ix], " \t\n\r\f")
		switch sub {
		case "(":
			level++

			if f != nil {
				if !f.AND() && !f.OR() {
					return nil, fmt.Errorf(msgParenthesisAfterFunc)
				}

				stack = append(stack, f)
				f = nil
			}

		case ")":
			level--

			if f != nil {
				if o, ok := f.(*or); ok {
					o.closed = true
				}
				if a, ok := f.(*and); ok {
					a.closed = true
				}
			}

			if l := len(stack) - 1; l >= 0 {
				f2 := stack[l]
				stack = stack[:l]

				if f2 != nil {
					if !f2.AND() && !f2.OR() {
						return nil, fmt.Errorf(msgMultipleFunc)
					}

					f2.Append(f)
				}

				f = f2
			}

		case "AND", "and":
			switch {
			case f == nil:
				return nil, fmt.Errorf(msgBoolNoPredecessor)
			case boolF:
				return nil, fmt.Errorf(msgMultipleBool)
			case f.OR():
				return nil, fmt.Errorf(msgMixedBools)
			case f.AND():
			default:
				f = &and{list: []Function{f}}
			}

			boolF = true

		case "OR", "or":
			switch {
			case f == nil:
				return nil, fmt.Errorf(msgBoolNoPredecessor)
			case boolF:
				return nil, fmt.Errorf(msgMultipleBool)
			case f.AND():
				return nil, fmt.Errorf(msgMixedBools)
			case f.OR():
			default:
				f = &or{list: []Function{f}}
			}

			boolF = true

		case "NOT", "not":
			not = !not

		default:
			parts := rx2.FindAllString(sub, -1)
			if len(parts) == 0 {
				return nil, fmt.Errorf(msgInvalidFunction, sub)
			}

			name := parts[0]
			cond, ok := conditions[name]
			if !ok {
				return nil, fmt.Errorf(msgInvalidFunctionName, name)
			}

			f2 := &function{not: not, name: name}

			var err error
			if f2.params, err = checkParams(cond, symbols, parts[1:]); err != nil {
				return nil, fmt.Errorf(msgInvalidParameters, name, err)
			}

			switch {
			case f == nil:
				f = f2
			case f.OR() || f.AND():
				f.Append(f2)
			default:
				return nil, fmt.Errorf(msgMultipleFunc)
			}

			not = false
			boolF = false
		}
	}

	if level != 0 {
		return nil, fmt.Errorf(msgUnmatchedParenthesis, level)
	}
	if boolF {
		return nil, fmt.Errorf(msgBoolNoSuccessor)
	}
	if f == nil {
		return nil, fmt.Errorf(msgInvalidFunction, statement)
	}

	return f, nil
}

var (
	rx1 = regexp.MustCompile(`\s*(\(|\)|and|or|not|AND|OR|NOT|[a-zA-z0-9\-]+\([^)]*\))\s*`)
	rx2 = regexp.MustCompile(`([^,()]+)`)
)

const (
	msgInvalidStatement     = "not a valid statement (rx1 failed)"
	msgInvalidFunction      = "not a valid function [%s] (rx2 failed)"
	msgInvalidFunctionName  = "not a valid function name [%s]"
	msgMultipleBool         = "multiple boolean operators in a row"
	msgBoolNoPredecessor    = "no preceding function for boolean operator"
	msgBoolNoSuccessor      = "no successive function for boolean operator"
	msgMixedBools           = "mixed boolean operators not supported; please use brackets"
	msgMultipleFunc         = "multiple functions in a row"
	msgUnmatchedParenthesis = "unmatched parentheses; level = %d"
	msgParenthesisAfterFunc = "parenthesis after function"
	msgInvalidParameters    = "invalid function parameters for [%s]: %w"
)
