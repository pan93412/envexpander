package envexpander_test

import (
	"testing"

	"github.com/pan93412/envexpander/v2"
	"github.com/stretchr/testify/assert"
)

func TestMarkVariablePositions(t *testing.T) {
	testmap := map[string][]envexpander.VariablePos{
		"A=${B}": {
			{2, 6},
		},
		"A=${B}${C}": {
			{2, 6},
			{6, 10},
		},
		"A=${B}${C}${D}": {
			{2, 6},
			{6, 10},
			{10, 14},
		},
		"A=${B}${C}${D}${D}": {
			{2, 6},
			{6, 10},
			{10, 14},
			{14, 18},
		},
		"A=${B}${C}${C}${D}${D}": {
			{2, 6},
			{6, 10},
			{10, 14},
			{14, 18},
			{18, 22},
		},
		"A=1234": {},
		"A=${B":  {},
		"A=$B${C}": {
			{4, 8},
		},
		"A=$B": {},
	}

	for k, v := range testmap {
		k := k
		v := v

		t.Run(k, func(t *testing.T) {
			t.Parallel()

			assert.ElementsMatch(t, envexpander.MarkVariablePositions(k), v)
		})
	}
}

// it covers varbegin and varend.
func TestVariablePos_Variable(t *testing.T) {
	testmap := []struct {
		Value     string
		Positions []envexpander.VariablePos
		Variables []string
	}{
		{
			Value: "A=${B}",
			Positions: []envexpander.VariablePos{
				{2, 6},
			},
			Variables: []string{"B"},
		},
		{
			Value: "A=${B}${C}",
			Positions: []envexpander.VariablePos{
				{2, 6},
				{6, 10},
			},
			Variables: []string{"B", "C"},
		},
		{
			Value: "A=${B}${C}${D}",
			Positions: []envexpander.VariablePos{
				{2, 6},
				{6, 10},
				{10, 14},
			},
			Variables: []string{"B", "C", "D"},
		},
		{
			Value: "A=${B}${C}${D}${D}",
			Positions: []envexpander.VariablePos{
				{2, 6},
				{6, 10},
				{10, 14},
				{14, 18},
			},
			Variables: []string{"B", "C", "D", "D"},
		},
	}

	for _, v := range testmap {
		v := v
		t.Run(v.Value, func(t *testing.T) {
			t.Parallel()

			for i, pos := range v.Positions {
				assert.Equal(t, v.Variables[i], pos.Variable(v.Value))
			}
		})
	}
}
