package envexpander

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

// MarkVariablePositions constructs a state machine and
// marks the variables positions from the value without
// any backtracing.
func MarkVariablePositions(value string) []VariablePos {
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
