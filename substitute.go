package envexpander

import "strings"

// Replacer replace the interpolation in the value with the
// actual variables value.
type Replacer struct {
	Value string

	// variables is a map of the variables that are referenced by this variable.
	Variables map[string]string
}

// Integrate replaces the variables interpolation with
// the variables value.
func (r *Replacer) Integrate() string {
	var builder strings.Builder
	varPos := MarkVariablePositions(r.Value)

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
			builder.WriteString(value)
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

// Integrate replaces the variables interpolation with
// the variables value.
func (r *Replacer) IntegrateWithCache(cvp CachedVariablePos) string {
	var builder strings.Builder
	varPos := cvp.MarkVariablePositions(r.Value)

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
			builder.WriteString(value)
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
