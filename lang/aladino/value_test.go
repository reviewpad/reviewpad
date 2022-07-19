// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/stretchr/testify/assert"
)

// Test builders

func TestBuildIntValue(t *testing.T) {
	wantVal := &aladino.IntValue{Val: 1}

	gotVal := aladino.BuildIntValue(1)

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildTrueValue(t *testing.T) {
	wantVal := aladino.BuildBoolValue(true)

	gotVal := aladino.BuildTrueValue()

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildFalseValue(t *testing.T) {
	wantVal := aladino.BuildBoolValue(false)

	gotVal := aladino.BuildFalseValue()

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildBoolValue(t *testing.T) {
	wantVal := &aladino.BoolValue{Val: true}

	gotVal := aladino.BuildBoolValue(true)

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildStringValue(t *testing.T) {
	str := "Lorem Ipsum"

	wantVal := &aladino.StringValue{Val: str}

	gotVal := aladino.BuildStringValue(str)

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildTimeValue(t *testing.T) {
	wantVal := &aladino.TimeValue{Val: 1}

	gotVal := aladino.BuildTimeValue(1)

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildArrayValue(t *testing.T) {
	elVals := []aladino.Value{aladino.BuildIntValue(1)}
	wantVal := &aladino.ArrayValue{Vals: elVals}

	gotVal := aladino.BuildArrayValue(elVals)

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildFunctionValue(t *testing.T) {
	fn := func(args []aladino.Value) aladino.Value {
		return &aladino.IntValue{Val: 0}
	}

	wantVal := &aladino.FunctionValue{fn}

	gotVal := aladino.BuildFunctionValue(fn)

	assert.True(t, wantVal.Equals(gotVal))
}

// Test kind

func TestIntValueKind(t *testing.T) {
	wantVal := aladino.INT_VALUE

	intVal := &aladino.IntValue{Val: 1}
	gotVal := intVal.Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestBoolValueKind(t *testing.T) {
	wantVal := aladino.BOOL_VALUE

	boolVal := &aladino.BoolValue{Val: true}
	gotVal := boolVal.Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestStringValueKind(t *testing.T) {
	wantVal := aladino.STRING_VALUE

	strVal := &aladino.StringValue{Val: "Lorem Ipsum"}
	gotVal := strVal.Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestTimeValueKind(t *testing.T) {
	wantVal := aladino.TIME_VALUE

	timeVal := &aladino.TimeValue{Val: 1}
	gotVal := timeVal.Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestArrayValueKind(t *testing.T) {
	wantVal := aladino.ARRAY_VALUE

	arrayVal := &aladino.ArrayValue{Vals: []aladino.Value{}}
	gotVal := arrayVal.Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestFunctionValueKind(t *testing.T) {
	wantVal := aladino.FUNCTION_VALUE

	fnVal := &aladino.FunctionValue{
		func(args []aladino.Value) aladino.Value {
			return &aladino.IntValue{Val: 0}
		},
	}
	gotVal := fnVal.Kind()

	assert.Equal(t, wantVal, gotVal)
}

// Test has kind of

func TestIntValueHasKindOf(t *testing.T) {
	intVal := &aladino.IntValue{Val: 0}

	assert.True(t, intVal.HasKindOf(aladino.INT_VALUE))
}

func TestBoolValueHasKindOf(t *testing.T) {
	boolVal := &aladino.BoolValue{Val: true}

	assert.True(t, boolVal.HasKindOf(aladino.BOOL_VALUE))
}

func TestStringValueHasKindOf(t *testing.T) {
	strVal := &aladino.StringValue{Val: "Lorem Ipsum"}

	assert.True(t, strVal.HasKindOf(aladino.STRING_VALUE))
}

func TestTimeValueHasKindOf(t *testing.T) {
	timeVal := &aladino.TimeValue{Val: 1}

	assert.True(t, timeVal.HasKindOf(aladino.TIME_VALUE))
}

func TestArrayValueHasKindOf(t *testing.T) {
	arrayVal := &aladino.ArrayValue{Vals: []aladino.Value{}}

	assert.True(t, arrayVal.HasKindOf(aladino.ARRAY_VALUE))
}

func TestFunctionValueHasKindOf(t *testing.T) {
	fnVal := &aladino.FunctionValue{
		func(args []aladino.Value) aladino.Value {
			return &aladino.IntValue{Val: 0}
		},
	}

	assert.True(t, fnVal.HasKindOf(aladino.FUNCTION_VALUE))
}

// Test equals

func TestIntValueEquals_WhenDiffKinds(t *testing.T) {
	intVal := &aladino.IntValue{Val: 0}
	otherVal := &aladino.BoolValue{Val: true}

	assert.False(t, intVal.Equals(otherVal))
}

func TestIntValueEquals_WhenTrue(t *testing.T) {
	intVal := &aladino.IntValue{Val: 0}
	otherVal := &aladino.IntValue{Val: 0}

	assert.True(t, intVal.Equals(otherVal))
}

func TestIntValueEquals_WhenFalse(t *testing.T) {
	intVal := &aladino.IntValue{Val: 0}
	otherVal := &aladino.IntValue{Val: 1}

	assert.False(t, intVal.Equals(otherVal))
}

func TestBoolValueEquals_WhenDiffKinds(t *testing.T) {
	boolVal := &aladino.BoolValue{Val: true}
	otherVal := &aladino.IntValue{Val: 0}

	assert.False(t, boolVal.Equals(otherVal))
}

func TestBoolValueEquals_WhenTrue(t *testing.T) {
	boolVal := &aladino.BoolValue{Val: true}
	otherVal := &aladino.BoolValue{Val: true}

	assert.True(t, boolVal.Equals(otherVal))
}

func TestBoolValueEquals_WhenFalse(t *testing.T) {
	boolVal := &aladino.BoolValue{Val: true}
	otherVal := &aladino.BoolValue{Val: false}

	assert.False(t, boolVal.Equals(otherVal))
}

func TestStringValueEquals_WhenDiffKinds(t *testing.T) {
	strVal := &aladino.StringValue{Val: "Lorem Ipsum"}
	otherVal := &aladino.IntValue{Val: 0}

	assert.False(t, strVal.Equals(otherVal))
}

func TestStringValueEquals_WhenTrue(t *testing.T) {
	strVal := &aladino.StringValue{Val: "Lorem Ipsum"}
	otherVal := &aladino.StringValue{Val: "Lorem Ipsum"}

	assert.True(t, strVal.Equals(otherVal))
}

func TestStringValueEquals_WhenFalse(t *testing.T) {
	strVal := &aladino.StringValue{Val: "Lorem Ipsum #1"}
	otherVal := &aladino.StringValue{Val: "Lorem Ipsum #2"}

	assert.False(t, strVal.Equals(otherVal))
}

func TestTimeValueEquals_WhenDiffKinds(t *testing.T) {
	timeVal := &aladino.TimeValue{Val: 1}
	otherVal := &aladino.IntValue{Val: 0}

	assert.False(t, timeVal.Equals(otherVal))
}

func TestTimeValueEquals_WhenTrue(t *testing.T) {
	timeVal := &aladino.TimeValue{Val: 1}
	otherVal := &aladino.TimeValue{Val: 1}

	assert.True(t, timeVal.Equals(otherVal))
}

func TestTimeValueEquals_WhenFalse(t *testing.T) {
	timeVal := &aladino.TimeValue{Val: 0}
	otherVal := &aladino.TimeValue{Val: 1}

	assert.False(t, timeVal.Equals(otherVal))
}

func TestArrayValueEquals_WhenDiffKinds(t *testing.T) {
	arrayVal := &aladino.ArrayValue{Vals: []aladino.Value{}}
	otherVal := &aladino.IntValue{Val: 0}

	assert.False(t, arrayVal.Equals(otherVal))
}

func TestArrayValueEquals_WhenDiffLength(t *testing.T) {
	arrayVal := &aladino.ArrayValue{Vals: []aladino.Value{}}
	otherVal := &aladino.ArrayValue{Vals: []aladino.Value{&aladino.BoolValue{Val: true}}}

	assert.False(t, arrayVal.Equals(otherVal))
}

func TestArrayValueEquals_WhenDiffElems(t *testing.T) {
	arrayVal := &aladino.ArrayValue{Vals: []aladino.Value{&aladino.BoolValue{Val: false}}}
	otherVal := &aladino.ArrayValue{Vals: []aladino.Value{&aladino.BoolValue{Val: true}}}

	assert.False(t, arrayVal.Equals(otherVal))
}

func TestArrayValueEquals_WhenTrue(t *testing.T) {
	arrayVal := &aladino.ArrayValue{Vals: []aladino.Value{&aladino.BoolValue{Val: true}}}
	otherVal := &aladino.ArrayValue{Vals: []aladino.Value{&aladino.BoolValue{Val: true}}}

	assert.True(t, arrayVal.Equals(otherVal))
}

func TestArrayValueEquals_WhenFalse(t *testing.T) {
	arrayVal := &aladino.ArrayValue{Vals: []aladino.Value{&aladino.BoolValue{Val: true}}}
	otherVal := &aladino.ArrayValue{Vals: []aladino.Value{&aladino.BoolValue{Val: false}}}

	assert.False(t, arrayVal.Equals(otherVal))
}

func TestFunctionValueEquals_WhenTrue(t *testing.T) {
	fnVal := &aladino.FunctionValue{
		func(args []aladino.Value) aladino.Value {
			return &aladino.IntValue{Val: 0}
		},
	}

	otherVal := &aladino.FunctionValue{
		func(args []aladino.Value) aladino.Value {
			return &aladino.IntValue{Val: 0}
		},
	}

	assert.True(t, fnVal.Equals(otherVal))
}

func TestFunctionValueEquals_WhenFalse(t *testing.T) {
	fnVal := &aladino.FunctionValue{
		func(args []aladino.Value) aladino.Value {
			return &aladino.IntValue{Val: 0}
		},
	}

	otherVal := &aladino.IntValue{Val: 0}

	assert.False(t, fnVal.Equals(otherVal))
}
