// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package lang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildStringType(t *testing.T) {
	wantVal := &StringType{}
	gotVal := BuildStringType()

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildIntType(t *testing.T) {
	wantVal := &IntType{}
	gotVal := BuildIntType()

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildBoolType(t *testing.T) {
	wantVal := &BoolType{}
	gotVal := BuildBoolType()

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildFunctionType(t *testing.T) {
	wantVal := &FunctionType{[]Type{&StringType{}}, &StringType{}}
	gotVal := BuildFunctionType([]Type{&StringType{}}, &StringType{})

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildArrayOfType(t *testing.T) {
	wantVal := &ArrayOfType{&StringType{}}
	gotVal := BuildArrayOfType(&StringType{})

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildArrayType(t *testing.T) {
	wantVal := &ArrayType{[]Type{&StringType{}, &IntType{}}}
	gotVal := BuildArrayType([]Type{&StringType{}, &IntType{}})

	assert.Equal(t, wantVal, gotVal)
}

func TestKind_WhenBoolType(t *testing.T) {
	wantVal := BOOL_TYPE
	gotVal := BuildBoolType().Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestKind_WhenIntType(t *testing.T) {
	wantVal := INT_TYPE
	gotVal := BuildIntType().Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestKind_WhenStringType(t *testing.T) {
	wantVal := STRING_TYPE
	gotVal := BuildStringType().Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestKind_WhenFunctionType(t *testing.T) {
	wantVal := FUNCTION_TYPE
	gotVal := BuildFunctionType([]Type{&StringType{}}, &StringType{}).Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestKind_WhenArrayType(t *testing.T) {
	wantVal := ARRAY_TYPE
	gotVal := BuildArrayType([]Type{&StringType{}, &IntType{}}).Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestKind_WhenArrayOfType(t *testing.T) {
	wantVal := ARRAY_OF_TYPE
	gotVal := BuildArrayOfType(&StringType{}).Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestEquals_WhenArraysOfDiffSizes(t *testing.T) {
	leftTys := []Type{BuildBoolType(), BuildIntType()}
	rightTys := []Type{BuildBoolType()}

	assert.False(t, Equals(leftTys, rightTys))
}

func TestEquals_WhenArraysContainDiffValues(t *testing.T) {
	leftTys := []Type{BuildIntType()}
	rightTys := []Type{BuildBoolType()}

	assert.False(t, Equals(leftTys, rightTys))
}

func TestEquals_WhenEqualTypes(t *testing.T) {
	leftTys := []Type{BuildBoolType()}
	rightTys := []Type{BuildBoolType()}

	assert.True(t, Equals(leftTys, rightTys))
}

func TestEquals_WhenBoolTypeComparedToStringType(t *testing.T) {
	boolType := BuildBoolType()
	otherType := BuildStringType()

	assert.False(t, boolType.Equals(otherType))
}

func TestEquals_WhenBoolTypeComparedToSameType(t *testing.T) {
	boolType := BuildBoolType()
	otherType := BuildBoolType()

	assert.True(t, boolType.Equals(otherType))
}

func TestEquals_WhenStringTypeComparedToBoolType(t *testing.T) {
	stringType := BuildStringType()
	otherType := BuildBoolType()

	assert.False(t, stringType.Equals(otherType))
}

func TestEquals_WhenStringTypeComparedToSameType(t *testing.T) {
	stringType := BuildStringType()
	otherType := BuildStringType()

	assert.True(t, stringType.Equals(otherType))
}

func TestEquals_WhenIntTypeComparedToStringType(t *testing.T) {
	intType := BuildIntType()
	otherType := BuildStringType()

	assert.False(t, intType.Equals(otherType))
}

func TestEquals_WhenIntTypeComparedToSameType(t *testing.T) {
	intType := BuildIntType()
	otherType := BuildIntType()

	assert.True(t, intType.Equals(otherType))
}

func TestEquals_WhenFunctionTypeComparedToIntType(t *testing.T) {
	functionType := BuildFunctionType([]Type{BuildStringType()}, BuildIntType())
	otherType := BuildIntType()

	assert.False(t, functionType.Equals(otherType))
}

func TestEquals_WhenFunctionTypeComparedToSameType(t *testing.T) {
	functionType := BuildFunctionType([]Type{BuildStringType()}, BuildIntType())
	otherType := BuildFunctionType([]Type{BuildStringType()}, BuildIntType())

	assert.True(t, functionType.Equals(otherType))
}

func TestEquals_WhenArraysTypesHaveDiffSizes(t *testing.T) {
	arrayType := BuildArrayType([]Type{BuildStringType()})
	otherArrayType := BuildArrayType([]Type{BuildStringType(), BuildIntType()})

	assert.False(t, arrayType.Equals(otherArrayType))
}

func TestEquals_WhenArraysTypesAreEqual(t *testing.T) {
	arrayType := BuildArrayType([]Type{BuildStringType()})
	otherArrayType := BuildArrayType([]Type{BuildStringType()})

	assert.True(t, arrayType.Equals(otherArrayType))
}

func TestEquals_WhenArrayTypeAndArrayOfTypeHaveDiffKind(t *testing.T) {
	arrayType := BuildArrayType([]Type{BuildStringType()})
	arrayOfType := BuildArrayOfType(BuildIntType())

	assert.False(t, arrayType.Equals(arrayOfType))
}

func TestEquals_WhenArrayTypeAndArrayOfTypeHaveSameKind(t *testing.T) {
	arrayType := BuildArrayType([]Type{BuildStringType()})
	arrayOfType := BuildArrayOfType(BuildStringType())

	assert.True(t, arrayType.Equals(arrayOfType))
}

func TestEquals_WhenArrayTypeComparedToStringType(t *testing.T) {
	arrayType := BuildArrayType([]Type{BuildStringType()})
	otherType := BuildStringType()

	assert.False(t, arrayType.Equals(otherType))
}

func TestEquals_WhenArrayOfTypeAndArrayTypeHaveDiffKinds(t *testing.T) {
	arrayOfType := BuildArrayOfType(BuildIntType())
	arrayType := BuildArrayType([]Type{BuildStringType()})

	assert.False(t, arrayOfType.Equals(arrayType))
}

func TestEquals_WhenArrayOfTypeAndArrayTypeHaveSameKind(t *testing.T) {
	arrayOfType := BuildArrayOfType(BuildIntType())
	arrayType := BuildArrayType([]Type{BuildIntType()})

	assert.True(t, arrayOfType.Equals(arrayType))
}

func TestEquals_WhenArraysOfTypesHaveDiffKind(t *testing.T) {
	arrayOfType := BuildArrayOfType(BuildIntType())
	otherArrayOfType := BuildArrayOfType(BuildStringType())

	assert.False(t, arrayOfType.Equals(otherArrayOfType))
}

func TestEquals_WhenArraysOfTypesAreEqual(t *testing.T) {
	arrayOfType := BuildArrayOfType(BuildIntType())
	otherArrayOfType := BuildArrayOfType(BuildIntType())

	assert.True(t, arrayOfType.Equals(otherArrayOfType))
}

func TestEquals_WhenArrayOfTypeComparedToStringType(t *testing.T) {
	arrayOfType := BuildArrayOfType(BuildIntType())
	otherType := BuildStringType()

	assert.False(t, arrayOfType.Equals(otherType))
}

func TestEquals_WhenArrayTypeComparedToDynamicArrayType(t *testing.T) {
	arrayType := BuildArrayType([]Type{BuildStringType(), BuildIntType()})
	dynamicArrayType := BuildDynamicArrayType()

	assert.True(t, arrayType.Equals(dynamicArrayType))
}

func TestEquals_WhenArrayOfTypeComparedToDynamicArrayType(t *testing.T) {
	arrayOfType := BuildArrayOfType(BuildStringType())
	dynamicArrayType := BuildDynamicArrayType()

	assert.True(t, arrayOfType.Equals(dynamicArrayType))
}
