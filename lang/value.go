// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package lang

import "reflect"

type Value interface {
	HasKindOf(string) bool
	Kind() string
	Equals(Value) bool
}

const (
	INT_VALUE      string = "IntValue"
	BOOL_VALUE     string = "BoolValue"
	STRING_VALUE   string = "StringValue"
	TIME_VALUE     string = "TimeValue"
	ARRAY_VALUE    string = "ArrayValue"
	FUNCTION_VALUE string = "FunctionValue"
	JSON_VALUE     string = "JSONValue"
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
