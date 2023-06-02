// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package lang

import (
	"reflect"
)

type Value interface {
	HasKindOf(string) bool
	Kind() string
	Equals(Value) bool
	Type() Type
}

const (
	INT_VALUE        string = "IntValue"
	BOOL_VALUE       string = "BoolValue"
	STRING_VALUE     string = "StringValue"
	TIME_VALUE       string = "TimeValue"
	ARRAY_VALUE      string = "ArrayValue"
	FUNCTION_VALUE   string = "FunctionValue"
	JSON_VALUE       string = "JSONValue"
	DICTIONARY_VALUE string = "DictionaryValue"
)

// IntValue represents an integer value
type IntValue struct {
	Val int
}

func BuildIntValue(iVal int) *IntValue {
	return &IntValue{Val: iVal}
}

func (iVal *IntValue) Kind() string {
	return INT_VALUE
}

func (iVal *IntValue) HasKindOf(ty string) bool {
	return iVal.Kind() == ty
}

func (thisVal *IntValue) Equals(other Value) bool {
	if thisVal.Kind() != other.Kind() {
		return false
	}

	return thisVal.Val == other.(*IntValue).Val
}

func (iVal *IntValue) Type() Type {
	return BuildIntType()
}

// BoolValue represents a bool value
type BoolValue struct {
	// defaultValue
	Val bool
}

func BuildTrueValue() *BoolValue  { return BuildBoolValue(true) }
func BuildFalseValue() *BoolValue { return BuildBoolValue(false) }

func BuildBoolValue(bval bool) *BoolValue {
	return &BoolValue{Val: bval}
}

func (bVal *BoolValue) Kind() string {
	return BOOL_VALUE
}

func (thisVal *BoolValue) Equals(other Value) bool {
	if thisVal.Kind() != other.Kind() {
		return false
	}

	return thisVal.Val == other.(*BoolValue).Val
}

func (bVal *BoolValue) HasKindOf(ty string) bool {
	return bVal.Kind() == ty
}

func (bVal *BoolValue) Type() Type {
	return BuildBoolType()
}

// StringValue represents a string value
type StringValue struct {
	// defaultValue
	Val string
}

func BuildStringValue(sval string) *StringValue {
	return &StringValue{Val: sval}
}

func (sVal *StringValue) Kind() string {
	return STRING_VALUE
}

func (thisVal *StringValue) Equals(other Value) bool {
	if thisVal.Kind() != other.Kind() {
		return false
	}

	return thisVal.Val == other.(*StringValue).Val
}

func (sVal *StringValue) HasKindOf(ty string) bool {
	return sVal.Kind() == ty
}

func (sVal *StringValue) Type() Type {
	return BuildStringType()
}

type TimeValue struct {
	Val int
}

func BuildTimeValue(tVal int) *TimeValue {
	return &TimeValue{
		Val: tVal,
	}
}

func (tVal *TimeValue) Kind() string {
	return TIME_VALUE
}

func (tVal *TimeValue) HasKindOf(kind string) bool {
	return tVal.Kind() == kind
}

func (thisVal *TimeValue) Equals(other Value) bool {
	if thisVal.Kind() != other.Kind() {
		return false
	}

	return thisVal.Val == other.(*TimeValue).Val
}

func (tVal *TimeValue) Type() Type {
	return BuildIntType()
}

// ArrayValue represents an array value
type ArrayValue struct {
	// defaultValue
	Vals []Value
}

func BuildArrayValue(elVals []Value) *ArrayValue {
	return &ArrayValue{Vals: elVals}
}

func (aVal *ArrayValue) Kind() string {
	return ARRAY_VALUE
}

func (thisVal *ArrayValue) Equals(other Value) bool {
	if thisVal.Kind() != other.Kind() {
		return false
	}

	otherArray := other.(*ArrayValue)

	if len(thisVal.Vals) != len(otherArray.Vals) {
		return false
	}

	for i, val := range thisVal.Vals {
		otherVal := otherArray.Vals[i]
		if !val.Equals(otherVal) {
			return false
		}
	}

	return true
}

func (aVal *ArrayValue) HasKindOf(ty string) bool {
	return aVal.Kind() == ty
}

func (aVal *ArrayValue) Type() Type {
	types := make([]Type, len(aVal.Vals))
	for i, val := range aVal.Vals {
		types[i] = val.Type()
	}

	return BuildArrayType(types)
}

// FunctionValue represents a function value
type FunctionValue struct {
	// defaultValue
	Fn func(args []Value) Value
}

func BuildFunctionValue(fn func(args []Value) Value) *FunctionValue {
	return &FunctionValue{fn}
}

func (fVal *FunctionValue) Kind() string {
	return FUNCTION_VALUE
}

func (thisVal *FunctionValue) Equals(other Value) bool {
	// TODO: Implement a more reliable function equals
	return thisVal.Kind() == other.Kind()
}

func (fVal *FunctionValue) HasKindOf(ty string) bool {
	return fVal.Kind() == ty
}

// TODO: Implement a more reliable function type
// that returns the actual parameters and return type
func (fVal *FunctionValue) Type() Type {
	return BuildFunctionType([]Type{}, BuildBoolType())
}

// JSONValue represents a json value
type JSONValue struct {
	Val interface{}
}

func BuildJSONValue(value interface{}) *JSONValue {
	return &JSONValue{value}
}

func (jVal *JSONValue) Kind() string {
	return JSON_VALUE
}

func (jVal *JSONValue) HasKindOf(ty string) bool {
	return jVal.Kind() == ty
}

func (jVal *JSONValue) Equals(other Value) bool {
	if jVal.Kind() != other.Kind() {
		return false
	}

	return reflect.DeepEqual(jVal.Val, other.(*JSONValue).Val)
}

func (jVal *JSONValue) Type() Type {
	return BuildJSONType()
}

type DictionaryValue struct {
	Vals map[string]Value
}

func BuildDictionaryValue(vals map[string]Value) *DictionaryValue {
	return &DictionaryValue{Vals: vals}
}

func (dVal *DictionaryValue) Kind() string {
	return DICTIONARY_VALUE
}

func (dVal *DictionaryValue) HasKindOf(ty string) bool {
	return dVal.Kind() == ty
}

func (dVal *DictionaryValue) Equals(other Value) bool {
	if dVal.Kind() != other.Kind() {
		return false
	}

	otherDict := other.(*DictionaryValue)

	if len(dVal.Vals) != len(otherDict.Vals) {
		return false
	}

	for key, val := range dVal.Vals {
		otherVal := otherDict.Vals[key]
		if !val.Equals(otherVal) {
			return false
		}
	}

	return true
}

func (dVal *DictionaryValue) Type() Type {
	return BuildDictionaryType()
}
