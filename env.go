// Package envexpander provides a way to resolve environment variables statically.
package envexpander

import (
	"strings"
)

// Replacer is a struct that holds the value and the referenced variables,
// which is used to replace the interpolation with the value.
type Replacer struct {
	Value string

	// variables is a map of the variables that are referenced by this variable.
	//
	// nil for no referenced variables.
	Variables map[string]*string
}

// Integrate replaces the variables interpolation with
// the variables value.
func (r *Replacer) Integrate(cvp CachedVariablePos) string {
	var builder strings.Builder
	varPos := cvp.ExtractAllVariables(r.Value)

	previousBegin := 0

	for _, pos := range varPos {
		// Write the not variable region.
		if previousBegin != pos.Begin {
			builder.WriteString(r.Value[previousBegin:pos.Begin])
		}

		// Write the variable region.
		variable := pos.Variable(r.Value)
		value, ok := r.Variables[variable]

		if ok {
			if value != nil {
				// If the variable is defined, replace the interpolation with the value.
				builder.WriteString(*value)
			} else {
				// If the variable is `nil`, replace the interpolation with "".
				builder.WriteString("")
			}
		} else {
			// If the variable is not defined, replace the interpolation with the original text.
			builder.WriteString(r.Value[pos.Begin:pos.End])
		}

		// Set the previous begin position to the end of the variable region.
		previousBegin = pos.End
	}

	// Write the rest of the string.
	builder.WriteString(r.Value[previousBegin:])

	return builder.String()
}

// ResolveEnvVariable determines the env and expands
// environment variables.
func ResolveEnvVariable(env map[string]string) map[string]string {
	type replacerWithRefvar struct {
		Replacer
		Refvar map[string]struct{}
	}

	cvp := make(CachedVariablePos)
	resolved := make(map[string]string, len(env))
	dependents := make(map[string]replacerWithRefvar, len(env))

	// Extract all the referenced variables from the received environment variable list.
	for key, val := range env {
		refvar := ExtractReferencedVariable(cvp, val)

		dependents[key] = replacerWithRefvar{
			Replacer: Replacer{
				Value:     val,
				Variables: make(map[string]*string, len(refvar)),
			},
			Refvar: refvar,
		}
	}

	// Start resolving the variables.
	for len(resolved) != len(env) {
	deploop:
		for key, replacer := range dependents {
			// If the variable is already resolved, skip it.
			if _, ok := resolved[key]; ok {
				continue
			}

			// Check if all the referenced variables are resolved.
			for v := range replacer.Refvar {
				// If a variable has been resolved, skip it.
				if _, ok := replacer.Variables[v]; ok {
					continue
				}

				// If a variable is not defined in env, skip it.
				if _, ok := env[v]; !ok {
					continue
				}

				// If a variable is in circlular import, replace the interpolation with "".
				if _, ok := dependents[v].Refvar[key]; ok {
					replacer.Variables[v] = nil
					continue
				}

				// If a variable is resolved already, replace the interpolation
				// with the resolved value.
				if resolvedVal, ok := resolved[v]; ok {
					replacer.Variables[v] = &resolvedVal
					continue
				}

				// If there are still unreferenced variables,
				// we should resolve others.
				continue deploop
			}

			// If all the referenced variables are resolved, replace the interpolation
			// with the resolved value.
			resolved[key] = replacer.Integrate(cvp)
		}
	}

	// Return the resolved environment variable list.
	return resolved
}

// ExtractReferencedVariable extracts referenced variables (keys) from the value.
//
// It returns the set of keys referenced in the variable value.
func ExtractReferencedVariable(cvp CachedVariablePos, value string) map[string]struct{} {
	keys := make(map[string]struct{})
	extractedVars := cvp.ExtractAllVariables(value)

	for _, pos := range extractedVars {
		keys[pos.Variable(value)] = struct{}{}
	}

	return keys
}

// VariablePos records the position of the variable in the value.
//
// It is designed for `${VAR}`-like interpolation.
type VariablePos struct {
	Begin int
	End   int
}

// Variable returns the variable name based on the position.
func (p VariablePos) Variable(r string) string {
	return r[p.VarBegin():p.VarEnd()]
}

// VarBegin returns the begin index of the variable itself.
//
// For example, when Begin-End is `${VAR}`, the VarBegin-VarEnd is `VAR`.
func (p VariablePos) VarBegin() int {
	return p.Begin + 2
}

// VarEnd returns the end index of the variable itself.
//
// For example, when Begin-End is `${VAR}`, the VarBegin-VarEnd is `VAR`.
func (p VariablePos) VarEnd() int {
	return p.End - 1
}

// ExtractAllVariables constructs a state machine and
// extracts all the variables from the value in the O(n) way.
func ExtractAllVariables(value string) []VariablePos {
	type statusT int

	const statusWaitingForInterpolation statusT = 1 << 0
	const statusWaitingForVariable statusT = 1 << 1
	const statusWaitingForClosingBrace statusT = 1 << 2

	variables := make([]VariablePos, 0)

	begin := 0
	status := statusWaitingForInterpolation

	for index, char := range value {
		switch {
		case status&statusWaitingForInterpolation != 0 && char == '$':
			status = statusWaitingForVariable
			begin = index
		case status&statusWaitingForVariable != 0:
			if char == '{' {
				status = statusWaitingForClosingBrace
			} else {
				status = statusWaitingForInterpolation
			}
		case status&statusWaitingForClosingBrace != 0 && char == '}':
			status = statusWaitingForInterpolation

			variables = append(variables, VariablePos{
				Begin: begin,
				End:   index + 1,
			})
		}
	}

	return variables
}

// CachedVariablePos is a scoped cache for ExtractAllVariables.
type CachedVariablePos map[string][]VariablePos

// ExtractAllVariables returns the cached variable positions.
func (c CachedVariablePos) ExtractAllVariables(value string) []VariablePos {
	if cached, ok := c[value]; ok {
		return cached
	}

	c[value] = ExtractAllVariables(value)
	return c[value]
}
