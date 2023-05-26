// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package lang_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/stretchr/testify/assert"
)

// Test builders

func TestBuildIntValue(t *testing.T) {
	wantVal := &lang.IntValue{Val: 1}

	gotVal := lang.BuildIntValue(1)

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildTrueValue(t *testing.T) {
	wantVal := lang.BuildBoolValue(true)

	gotVal := lang.BuildTrueValue()

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildFalseValue(t *testing.T) {
	wantVal := lang.BuildBoolValue(false)

	gotVal := lang.BuildFalseValue()

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildBoolValue(t *testing.T) {
	wantVal := &lang.BoolValue{Val: true}

	gotVal := lang.BuildBoolValue(true)

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildStringValue(t *testing.T) {
	str := "Lorem Ipsum"

	wantVal := &lang.StringValue{Val: str}

	gotVal := lang.BuildStringValue(str)

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildTimeValue(t *testing.T) {
	wantVal := &lang.TimeValue{Val: 1}

	gotVal := lang.BuildTimeValue(1)

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildArrayValue(t *testing.T) {
	elVals := []lang.Value{lang.BuildIntValue(1)}
	wantVal := &lang.ArrayValue{Vals: elVals}

	gotVal := lang.BuildArrayValue(elVals)

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildFunctionValue(t *testing.T) {
	fn := func(args []lang.Value) lang.Value {
		return &lang.IntValue{Val: 0}
	}

	wantVal := &lang.FunctionValue{fn}

	gotVal := lang.BuildFunctionValue(fn)

	assert.True(t, wantVal.Equals(gotVal))
}

// Test kind

func TestIntValueKind(t *testing.T) {
	wantVal := lang.INT_VALUE

	intVal := &lang.IntValue{Val: 1}
	gotVal := intVal.Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestBoolValueKind(t *testing.T) {
	wantVal := lang.BOOL_VALUE

	boolVal := &lang.BoolValue{Val: true}
	gotVal := boolVal.Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestStringValueKind(t *testing.T) {
	wantVal := lang.STRING_VALUE

	strVal := &lang.StringValue{Val: "Lorem Ipsum"}
	gotVal := strVal.Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestTimeValueKind(t *testing.T) {
	wantVal := lang.TIME_VALUE

	timeVal := &lang.TimeValue{Val: 1}
	gotVal := timeVal.Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestArrayValueKind(t *testing.T) {
	wantVal := lang.ARRAY_VALUE

	arrayVal := &lang.ArrayValue{Vals: []lang.Value{}}
	gotVal := arrayVal.Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestFunctionValueKind(t *testing.T) {
	wantVal := lang.FUNCTION_VALUE

	fnVal := &lang.FunctionValue{
		func(args []lang.Value) lang.Value {
			return &lang.IntValue{Val: 0}
		},
	}
	gotVal := fnVal.Kind()

	assert.Equal(t, wantVal, gotVal)
}

// Test has kind of

func TestIntValueHasKindOf(t *testing.T) {
	intVal := &lang.IntValue{Val: 0}

	assert.True(t, intVal.HasKindOf(lang.INT_VALUE))
}

func TestBoolValueHasKindOf(t *testing.T) {
	boolVal := &lang.BoolValue{Val: true}

	assert.True(t, boolVal.HasKindOf(lang.BOOL_VALUE))
}

func TestStringValueHasKindOf(t *testing.T) {
	strVal := &lang.StringValue{Val: "Lorem Ipsum"}

	assert.True(t, strVal.HasKindOf(lang.STRING_VALUE))
}

func TestTimeValueHasKindOf(t *testing.T) {
	timeVal := &lang.TimeValue{Val: 1}

	assert.True(t, timeVal.HasKindOf(lang.TIME_VALUE))
}

func TestArrayValueHasKindOf(t *testing.T) {
	arrayVal := &lang.ArrayValue{Vals: []lang.Value{}}

	assert.True(t, arrayVal.HasKindOf(lang.ARRAY_VALUE))
}

func TestFunctionValueHasKindOf(t *testing.T) {
	fnVal := &lang.FunctionValue{
		func(args []lang.Value) lang.Value {
			return &lang.IntValue{Val: 0}
		},
	}

	assert.True(t, fnVal.HasKindOf(lang.FUNCTION_VALUE))
}

// Test equals

func TestIntValueEquals_WhenDiffKinds(t *testing.T) {
	intVal := &lang.IntValue{Val: 0}
	otherVal := &lang.BoolValue{Val: true}

	assert.False(t, intVal.Equals(otherVal))
}

func TestIntValueEquals_WhenTrue(t *testing.T) {
	intVal := &lang.IntValue{Val: 0}
	otherVal := &lang.IntValue{Val: 0}

	assert.True(t, intVal.Equals(otherVal))
}

func TestIntValueEquals_WhenFalse(t *testing.T) {
	intVal := &lang.IntValue{Val: 0}
	otherVal := &lang.IntValue{Val: 1}

	assert.False(t, intVal.Equals(otherVal))
}

func TestBoolValueEquals_WhenDiffKinds(t *testing.T) {
	boolVal := &lang.BoolValue{Val: true}
	otherVal := &lang.IntValue{Val: 0}

	assert.False(t, boolVal.Equals(otherVal))
}

func TestBoolValueEquals_WhenTrue(t *testing.T) {
	boolVal := &lang.BoolValue{Val: true}
	otherVal := &lang.BoolValue{Val: true}

	assert.True(t, boolVal.Equals(otherVal))
}

func TestBoolValueEquals_WhenFalse(t *testing.T) {
	boolVal := &lang.BoolValue{Val: true}
	otherVal := &lang.BoolValue{Val: false}

	assert.False(t, boolVal.Equals(otherVal))
}

func TestStringValueEquals_WhenDiffKinds(t *testing.T) {
	strVal := &lang.StringValue{Val: "Lorem Ipsum"}
	otherVal := &lang.IntValue{Val: 0}

	assert.False(t, strVal.Equals(otherVal))
}

func TestStringValueEquals_WhenTrue(t *testing.T) {
	strVal := &lang.StringValue{Val: "Lorem Ipsum"}
	otherVal := &lang.StringValue{Val: "Lorem Ipsum"}

	assert.True(t, strVal.Equals(otherVal))
}

func TestStringValueEquals_WhenFalse(t *testing.T) {
	strVal := &lang.StringValue{Val: "Lorem Ipsum #1"}
	otherVal := &lang.StringValue{Val: "Lorem Ipsum #2"}

	assert.False(t, strVal.Equals(otherVal))
}

func TestTimeValueEquals_WhenDiffKinds(t *testing.T) {
	timeVal := &lang.TimeValue{Val: 1}
	otherVal := &lang.IntValue{Val: 0}

	assert.False(t, timeVal.Equals(otherVal))
}

func TestTimeValueEquals_WhenTrue(t *testing.T) {
	timeVal := &lang.TimeValue{Val: 1}
	otherVal := &lang.TimeValue{Val: 1}

	assert.True(t, timeVal.Equals(otherVal))
}

func TestTimeValueEquals_WhenFalse(t *testing.T) {
	timeVal := &lang.TimeValue{Val: 0}
	otherVal := &lang.TimeValue{Val: 1}

	assert.False(t, timeVal.Equals(otherVal))
}

func TestArrayValueEquals_WhenDiffKinds(t *testing.T) {
	arrayVal := &lang.ArrayValue{Vals: []lang.Value{}}
	otherVal := &lang.IntValue{Val: 0}

	assert.False(t, arrayVal.Equals(otherVal))
}

func TestArrayValueEquals_WhenDiffLength(t *testing.T) {
	arrayVal := &lang.ArrayValue{Vals: []lang.Value{}}
	otherVal := &lang.ArrayValue{Vals: []lang.Value{&lang.BoolValue{Val: true}}}

	assert.False(t, arrayVal.Equals(otherVal))
}

func TestArrayValueEquals_WhenDiffElems(t *testing.T) {
	arrayVal := &lang.ArrayValue{Vals: []lang.Value{&lang.BoolValue{Val: false}}}
	otherVal := &lang.ArrayValue{Vals: []lang.Value{&lang.BoolValue{Val: true}}}

	assert.False(t, arrayVal.Equals(otherVal))
}

func TestArrayValueEquals_WhenTrue(t *testing.T) {
	arrayVal := &lang.ArrayValue{Vals: []lang.Value{&lang.BoolValue{Val: true}}}
	otherVal := &lang.ArrayValue{Vals: []lang.Value{&lang.BoolValue{Val: true}}}

	assert.True(t, arrayVal.Equals(otherVal))
}

func TestArrayValueEquals_WhenFalse(t *testing.T) {
	arrayVal := &lang.ArrayValue{Vals: []lang.Value{&lang.BoolValue{Val: true}}}
	otherVal := &lang.ArrayValue{Vals: []lang.Value{&lang.BoolValue{Val: false}}}

	assert.False(t, arrayVal.Equals(otherVal))
}

func TestFunctionValueEquals_WhenTrue(t *testing.T) {
	fnVal := &lang.FunctionValue{
		func(args []lang.Value) lang.Value {
			return &lang.IntValue{Val: 0}
		},
	}

	otherVal := &lang.FunctionValue{
		func(args []lang.Value) lang.Value {
			return &lang.IntValue{Val: 0}
		},
	}

	assert.True(t, fnVal.Equals(otherVal))
}

func TestFunctionValueEquals_WhenFalse(t *testing.T) {
	fnVal := &lang.FunctionValue{
		func(args []lang.Value) lang.Value {
			return &lang.IntValue{Val: 0}
		},
	}

	otherVal := &lang.IntValue{Val: 0}

	assert.False(t, fnVal.Equals(otherVal))
}

func TestJSONValueEquals(t *testing.T) {
	tests := map[string]struct {
		firstVal  *lang.JSONValue
		secondVal *lang.JSONValue
		equals    bool
	}{
		"when types are the same but values differ": {
			firstVal:  lang.BuildJSONValue(1),
			secondVal: lang.BuildJSONValue(10),
			equals:    false,
		},
		"when types are the same and values are the same": {
			firstVal:  lang.BuildJSONValue(true),
			secondVal: lang.BuildJSONValue(true),
			equals:    true,
		},
		"when arrays have same values": {
			firstVal:  lang.BuildJSONValue([]interface{}{1, 2, 3, "text", true, false}),
			secondVal: lang.BuildJSONValue([]interface{}{1, 2, 3, "text", true, false}),
			equals:    true,
		},
		"when arrays have out of order values": {
			firstVal:  lang.BuildJSONValue([]interface{}{1, 2, 3}),
			secondVal: lang.BuildJSONValue([]interface{}{1, 3, 2}),
			equals:    false,
		},
		"when map has same keys and same values but out of order": {
			firstVal: lang.BuildJSONValue(map[string]interface{}{
				"a": "b",
				"c": 1,
				"d": true,
				"e": false,
				"g": []interface{}{1},
			}),
			secondVal: lang.BuildJSONValue(map[string]interface{}{
				"a": "b",
				"d": true,
				"c": 1,
				"g": []interface{}{1},
				"e": false,
			}),
			equals: true,
		},
		"when map has same keys and different values": {
			firstVal: lang.BuildJSONValue(map[string]interface{}{
				"a": "b",
				"c": 1,
				"d": true,
				"e": false,
				"g": []interface{}{1},
			}),
			secondVal: lang.BuildJSONValue(map[string]interface{}{
				"a": "b",
				"c": 1.0,
				"d": false,
				"e": true,
				"g": []interface{}{2},
			}),
			equals: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := test.firstVal.Equals(test.secondVal)

			assert.Equal(t, test.equals, env)
		})
	}
}
