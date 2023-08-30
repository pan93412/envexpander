// Package envexpander provides a way to resolve environment variables statically.
package envexpander

// ResolveEnvVariable determines the env and expands
// environment variables.
func ResolveEnvVariable(env map[string]string) map[string]string {
	type replacerWithRefvar struct {
		Replacer
		Refvar map[string]struct{}
	}

	cvp := NewCachedVariablePos()
	resolved := make(map[string]string, len(env))
	dependents := make(map[string]replacerWithRefvar, len(env))

	// Extract all the referenced variables from the received environment variable list.
	for key, val := range env {
		refvar := FindVariableReferenceMap(cvp, val)

		dependents[key] = replacerWithRefvar{
			Replacer: Replacer{
				Value:     val,
				Variables: make(map[string]string, len(refvar)),
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
					replacer.Variables[v] = ""
					continue
				}

				// If a variable is resolved already, replace the interpolation
				// with the resolved value.
				if resolvedVal, ok := resolved[v]; ok {
					replacer.Variables[v] = resolvedVal
					continue
				}

				// If there are still unreferenced variables,
				// we should resolve others.
				continue deploop
			}

			// If all the referenced variables are resolved, replace the interpolation
			// with the resolved value.
			resolved[key] = replacer.IntegrateWithCache(cvp)
		}
	}

	// Return the resolved environment variable list.
	return resolved
}

// FindVariableReferenceMap extracts referenced variables (keys) from the value.
//
// It returns the set of keys referenced in the variable value.
func FindVariableReferenceMap(cvp CachedVariablePos, value string) map[string]struct{} {
	extractedVars := cvp.MarkVariablePositions(value)
	keys := make(map[string]struct{}, len(extractedVars))

	for _, pos := range extractedVars {
		keys[pos.Variable(value)] = struct{}{}
	}

	return keys
}
