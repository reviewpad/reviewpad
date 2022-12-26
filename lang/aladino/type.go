// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

type Type interface {
	Kind() string
	equals(other Type) bool
}

const (
	BOOL_TYPE          string = "BoolType"
	INT_TYPE           string = "IntType"
	STRING_TYPE        string = "StringType"
	FUNCTION_TYPE      string = "FunctionType"
	ARRAY_TYPE         string = "ArrayType"
	ARRAY_OF_TYPE      string = "ArrayOfType"
	JSON_TYPE          string = "JSONType"
	DYNAMIC_ARRAY_TYPE string = "DynamicArrayType"
)

type StringType struct{}

type IntType struct{}

type BoolType struct{}

type FunctionType struct {
	paramTypes []Type
	returnType Type
}

type ArrayOfType struct {
	elemType Type
}

type ArrayType struct {
	// static arrays
	elemsType []Type
}

type DynamicArrayType struct{}

type JSONType struct{}

func BuildStringType() *StringType { return &StringType{} }
func BuildIntType() *IntType       { return &IntType{} }
func BuildBoolType() *BoolType     { return &BoolType{} }

func BuildFunctionType(paramsTypes []Type, returnType Type) *FunctionType {
	return &FunctionType{paramsTypes, returnType}
}

func BuildArrayOfType(elemType Type) *ArrayOfType {
	return &ArrayOfType{elemType}
}

func BuildArrayType(elemsTypes []Type) *ArrayType {
	return &ArrayType{elemsTypes}
}

func BuildJSONType() *JSONType { return &JSONType{} }

func BuildDynamicArrayType() *DynamicArrayType { return &DynamicArrayType{} }

func (bTy *BoolType) Kind() string {
	return BOOL_TYPE
}

func (iTy *IntType) Kind() string {
	return INT_TYPE
}

func (sTy *StringType) Kind() string {
	return STRING_TYPE
}

func (fTy *FunctionType) Kind() string {
	return FUNCTION_TYPE
}

func (aTy *ArrayType) Kind() string {
	return ARRAY_TYPE
}

func (aTy *ArrayOfType) Kind() string {
	return ARRAY_OF_TYPE
}

func (jTy *JSONType) Kind() string {
	return JSON_TYPE
}

func (dATy *DynamicArrayType) Kind() string {
	return DYNAMIC_ARRAY_TYPE
}

// Equals
// equals on arrays
func equals(leftTys []Type, rightTys []Type) bool {
	if len(leftTys) != len(rightTys) {
		return false
	}

	for i, leftTy := range leftTys {
		rightTy := rightTys[i]

		if !leftTy.equals(rightTy) {
			return false
		}
	}

	return true
}

func (thisTy *BoolType) equals(thatTy Type) bool {
	return thatTy.Kind() == thisTy.Kind()
}

func (thisTy *StringType) equals(thatTy Type) bool {
	return thatTy.Kind() == thisTy.Kind()
}

func (thisTy *IntType) equals(thatTy Type) bool {
	return thatTy.Kind() == thisTy.Kind()
}

func (thisTy *FunctionType) equals(thatTy Type) bool {
	if thisTy.Kind() != thatTy.Kind() {
		return false
	}

	thatTyFunction := thatTy.(*FunctionType)
	argsCheck := equals(thisTy.paramTypes, thatTyFunction.paramTypes)
	retCheck := thisTy.returnType.equals(thatTyFunction.returnType)

	return argsCheck && retCheck
}

func (thisTy *ArrayType) equals(thatTy Type) bool {
	switch thatTy.Kind() {
	case ARRAY_TYPE:
		thatTyArray := thatTy.(*ArrayType)
		return equals(thatTyArray.elemsType, thisTy.elemsType)
	case ARRAY_OF_TYPE:
		thatTyArrayOf := thatTy.(*ArrayOfType)
		elems := make([]Type, len(thisTy.elemsType))
		for i := range thisTy.elemsType {
			elems[i] = thatTyArrayOf.elemType
		}
		return equals(elems, thisTy.elemsType)
	case DYNAMIC_ARRAY_TYPE:
		return true
	}
	return false
}

func (thisTy *ArrayOfType) equals(thatTy Type) bool {
	switch thatTy.Kind() {
	case ARRAY_TYPE:
		return thatTy.equals(thisTy)
	case ARRAY_OF_TYPE:
		thatTyArrayOf := thatTy.(*ArrayOfType)
		return thisTy.elemType.equals(thatTyArrayOf.elemType)
	case DYNAMIC_ARRAY_TYPE:
		return true
	}
	return false
}

func (thisTy *JSONType) equals(thatTy Type) bool {
	return thisTy.Kind() == thatTy.Kind()
}

func (thisTy *DynamicArrayType) equals(thatTy Type) bool {
	return thisTy.Kind() == thatTy.Kind() || thatTy.Kind() == ARRAY_OF_TYPE || thatTy.Kind() == ARRAY_TYPE
}
