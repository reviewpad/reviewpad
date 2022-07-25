// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildStringType_ExpectStringType(t *testing.T) {
	wantVal := &StringType{}
	gotVal := BuildStringType()

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildIntType_ExpectIntType(t *testing.T) {
	wantVal := &IntType{}
	gotVal := BuildIntType()

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildBoolType_ExpectBoolType(t *testing.T) {
	wantVal := &BoolType{}
	gotVal := BuildBoolType()

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildFunctionType_ExpectFunctionType(t *testing.T) {
	wantVal := &FunctionType{[]Type{&StringType{}}, &StringType{}}
	gotVal := BuildFunctionType([]Type{&StringType{}}, &StringType{})

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildArrayOfType_ExpectArrayOfType(t *testing.T) {
	wantVal := &ArrayOfType{&StringType{}}
	gotVal := BuildArrayOfType(&StringType{})

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildArrayType_ExpectArrayType(t *testing.T) {
	wantVal := &ArrayType{[]Type{&StringType{}, &IntType{}}}
	gotVal := BuildArrayType([]Type{&StringType{}, &IntType{}})

	assert.Equal(t, wantVal, gotVal)
}

func TestKind_WhenBoolType_ExpectBoolKind(t *testing.T) {
	wantVal := BOOL_TYPE
	gotVal := BuildBoolType().Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestKind_WhenIntType_ExpectIntKind(t *testing.T) {
	wantVal := INT_TYPE
	gotVal := BuildIntType().Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestKind_WhenStringType_ExpectStringKind(t *testing.T) {
	wantVal := STRING_TYPE
	gotVal := BuildStringType().Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestKind_WhenFunctionType_ExpectFunctionKind(t *testing.T) {
	wantVal := FUNCTION_TYPE
	gotVal := BuildFunctionType([]Type{&StringType{}}, &StringType{}).Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestKind_WhenArrayType_ExpectArrayKind(t *testing.T) {
	wantVal := ARRAY_TYPE
	gotVal := BuildArrayType([]Type{&StringType{}, &IntType{}}).Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestKind_WhenArrayOfType_ExpectArrayOfKind(t *testing.T) {
	wantVal := ARRAY_OF_TYPE
	gotVal := BuildArrayOfType(&StringType{}).Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestEquals_WhenArraysOfDiffSizes_ExpectFalse(t *testing.T) {
	leftTys := []Type{BuildBoolType(), BuildIntType()}
	rightTys := []Type{BuildBoolType()}

	assert.False(t, equals(leftTys, rightTys))
}

func TestEquals_WhenArraysContainDiffValues_ExpectFalse(t *testing.T) {
	leftTys := []Type{BuildIntType()}
	rightTys := []Type{BuildBoolType()}

	assert.False(t, equals(leftTys, rightTys))
}

func TestEquals_ExpectTrue(t *testing.T) {
	leftTys := []Type{BuildBoolType()}
	rightTys := []Type{BuildBoolType()}

	assert.True(t, equals(leftTys, rightTys))
}

func TestEquals_WhenBoolType_ExpectFalse(t *testing.T) {
	boolType := BuildBoolType()
	otherType := BuildStringType()

	assert.False(t, boolType.equals(otherType))
}

func TestEquals_WhenBoolType_ExpectTrue(t *testing.T) {
	boolType := BuildBoolType()
	otherType := BuildBoolType()

	assert.True(t, boolType.equals(otherType))
}

func TestEquals_WhenStringType_ExpectFalse(t *testing.T) {
	stringType := BuildStringType()
	otherType := BuildBoolType()

	assert.False(t, stringType.equals(otherType))
}

func TestEquals_WhenStringType_ExpectTrue(t *testing.T) {
	stringType := BuildStringType()
	otherType := BuildStringType()

	assert.True(t, stringType.equals(otherType))
}

func TestEquals_WhenIntType_ExpectFalse(t *testing.T) {
	intType := BuildIntType()
	otherType := BuildStringType()

	assert.False(t, intType.equals(otherType))
}

func TestEquals_WhenIntType_ExpectTrue(t *testing.T) {
	intType := BuildIntType()
	otherType := BuildIntType()

	assert.True(t, intType.equals(otherType))
}

func TestEquals_WhenFunctionType_ExpectFalse(t *testing.T) {
	functionType := BuildFunctionType([]Type{BuildStringType()}, BuildIntType())
	otherType := BuildIntType()

	assert.False(t, functionType.equals(otherType))
}

func TestEquals_WhenFunctionType_ExpectTrue(t *testing.T) {
	functionType := BuildFunctionType([]Type{BuildStringType()}, BuildIntType())
	otherType := BuildFunctionType([]Type{BuildStringType()}, BuildIntType())

	assert.True(t, functionType.equals(otherType))
}

func TestEquals_WhenArraysTypesHaveDiffSizes_ExpectFalse(t *testing.T) {
	arrayType := BuildArrayType([]Type{BuildStringType()})
	otherArrayType := BuildArrayType([]Type{BuildStringType(), BuildIntType()})

	assert.False(t, arrayType.equals(otherArrayType))
}

func TestEquals_WhenArraysTypesAreEqual_ExpectTrue(t *testing.T) {
	arrayType := BuildArrayType([]Type{BuildStringType()})
	otherArrayType := BuildArrayType([]Type{BuildStringType()})

	assert.True(t, arrayType.equals(otherArrayType))
}

func TestEquals_WhenArrayTypeAndArrayOfTypeHaveDiffKind_ExpectFalse(t *testing.T) {
	arrayType := BuildArrayType([]Type{BuildStringType()})
	arrayOfType := BuildArrayOfType(BuildIntType())

	assert.False(t, arrayType.equals(arrayOfType))
}

func TestEquals_WhenArrayTypeAndArrayOfTypeHaveSameKind_ExpectTrue(t *testing.T) {
	arrayType := BuildArrayType([]Type{BuildStringType()})
	arrayOfType := BuildArrayOfType(BuildStringType())

	assert.True(t, arrayType.equals(arrayOfType))
}

func TestEquals_WhenArrayType_ExpectFalse(t *testing.T) {
	arrayType := BuildArrayType([]Type{BuildStringType()})
	otherType := BuildStringType()

	assert.False(t, arrayType.equals(otherType))
}

func TestEquals_WhenArrayOfTypeAndArrayTypeHaveDiffKinds_ExpectFalse(t *testing.T) {
	arrayOfType := BuildArrayOfType(BuildIntType())
    arrayType := BuildArrayType([]Type{BuildStringType()})

    assert.False(t, arrayOfType.equals(arrayType))
}

func TestEquals_WhenArrayOfTypeAndArrayTypeHaveSameKind_ExpectTrue(t *testing.T) {
	arrayOfType := BuildArrayOfType(BuildIntType())
    arrayType := BuildArrayType([]Type{BuildIntType()})

    assert.True(t, arrayOfType.equals(arrayType))
}

func TestEquals_WhenArraysOfTypesHaveDiffKind_ExpectFalse(t *testing.T) {
	arrayOfType := BuildArrayOfType(BuildIntType())
	otherArrayOfType := BuildArrayOfType(BuildStringType())

	assert.False(t, arrayOfType.equals(otherArrayOfType))
}

func TestEquals_WhenArraysOfTypesAreEqual_ExpectTrue(t *testing.T) {
	arrayOfType := BuildArrayOfType(BuildIntType())
	otherArrayOfType := BuildArrayOfType(BuildIntType())

	assert.True(t, arrayOfType.equals(otherArrayOfType))
}

func TestEquals_WhenArrayOfType_ExpectFalse(t *testing.T) {
	arrayOfType := BuildArrayOfType(BuildIntType())
	otherType := BuildStringType()

	assert.False(t, arrayOfType.equals(otherType))
}
