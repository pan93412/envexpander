package envexpander

import (
	"strings"
)

func Expand(originalVariables map[string]string) map[string]string {
	const (
		Unvisited = iota
		Visiting
		Visited
	)

	variablePositions := make(map[string][]VariablePos, len(originalVariables))
	for key := range originalVariables {
		variablePositions[key] = Refs(originalVariables[key])
	}

	ordered := make([]string, 0, len(originalVariables))
	visited := make(map[string]int, len(originalVariables))
	circular := make(map[string]bool, len(originalVariables))

	var visit func(string) bool
	visit = func(node string) bool {
		if visited[node] == Visited { // Circular dependency found, ignore it
			circular[node] = true
			return false
		}
		if visited[node] == Visiting { // Already processed this node
			return true
		}
		visited[node] = Visiting

		for _, pos := range variablePositions[node] {
			variable := pos.Variable(originalVariables[node])

			if variable == node { // Ignore self-references
				continue
			}
			if _, exists := variablePositions[variable]; !exists { // Ignore unknown references
				continue
			}
			if !visit(variable) {
				continue
			}
		}

		visited[node] = Visited
		ordered = append(ordered, node)
		return true
	}

	for key := range originalVariables {
		if !visit(key) {
			continue
		}
	}

	// Expand variables.
	mapping := make(map[string]string, len(originalVariables))
	for _, key := range ordered {
		mapping[key] = ExpandOne(originalVariables[key], mapping, variablePositions[key])
	}

	return mapping
}

// ExpandOne expands a single variable in the mapping.
//
// "refs" is a slice of references in this s. If there
// is none, we call [Refs] to get them.
func ExpandOne(s string, mapping map[string]string, positions ...[]VariablePos) string {
	result := strings.Builder{}
	prevBegin := 0

	var resolvedPositions []VariablePos
	if len(positions) == 0 {
		resolvedPositions = Refs(s)
	} else {
		resolvedPositions = positions[0]
	}

	for _, pos := range resolvedPositions {
		variable := pos.Variable(s)
		value, ok := mapping[variable]
		if ok || value == variable {
			result.WriteString(s[prevBegin:pos.Begin])
			result.WriteString(value)
		} else {
			result.WriteString(s[prevBegin:pos.End])
		}
		prevBegin = pos.End
	}
	result.WriteString(s[prevBegin:])

	return result.String()
}

type Token struct {
	Type    string
	Content string
}

// Refs gets the variable referenced in value.
func Refs(value string) []VariablePos {
	refs := []VariablePos{}
	ptr := 0

	var currentRef VariablePos

	for ptr < len(value) {
		if value[ptr] == '$' && ptr+1 < len(value) && value[ptr+1] == '{' {
			// if this $ is escaped, we skip it.
			if ptr > 0 && value[ptr-1] == '\\' {
				ptr += 3
				continue
			}

			// make sure if we are not in the middle of a ref
			if currentRef.End == 0 {
				currentRef = VariablePos{
					Begin: ptr,
					End:   ptr + 2,
				}
				ptr += 2
				continue
			}
		}

		if value[ptr] == '}' {
			if currentRef.End != 0 {
				currentRef.End = ptr + 1
				refs = append(refs, currentRef)
			}

			currentRef.End = 0
			ptr++
			continue
		}

		if currentRef.End != 0 {
			// append the current character to the last token
			currentRef.End++
		}

		ptr++
	}

	return refs
}

// It is designed for `${VAR}`-like interpolation.
type VariablePos struct {
	Begin int // $
	End   int // }; != 0 unless not initialized
}

// Variable returns the variable name based on the position.
func (p *VariablePos) Variable(r string) string {
	return r[p.Begin+2 : p.End-1]
}
