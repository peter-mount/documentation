package util

import (
	"fmt"
	"strconv"
)

// The types of a Value returned by Type()
const (
	// Ignore 0 so use _ then if someone manually creates Value it's an unknown type
	_            = iota
	VarNull      // nil
	VarBool      // bool
	VarInt       // int
	VarFloat     // float64
	VarComplex   // complex
	VarString    // string
	VarInterface // interface
)

type Value struct {
	varType      int
	boolVal      bool
	intVal       int64
	floatVal     float64
	complexVal   complex128
	stringVal    string
	interfaceVal interface{}
}

var nullValue = Value{varType: VarNull}
var falseValue = Value{varType: VarBool, boolVal: false}
var trueValue = Value{varType: VarBool, boolVal: true}

func (v *Value) Same(b *Value) bool {
	if b == nil {
		return false
	}

	if v == b {
		return true
	}

	return v.varType == b.varType &&
		v.boolVal == b.boolVal &&
		v.intVal == b.intVal &&
		v.floatVal == b.floatVal &&
		v.complexVal == b.complexVal &&
		v.stringVal == b.stringVal
}

func (v *Value) Equal(b *Value) bool {
	if v == nil || v.IsNull() {
		return b == nil || b.IsNull()
	}
	switch v.Type() {
	case VarBool:
		return v.Bool() == b.Bool()
	case VarInt:
		return v.Int() == b.Int()
	case VarFloat:
		return v.Float() == b.Float()
	case VarComplex:
		return v.Complex() == b.Complex()
	case VarString:
		return v.String() == b.String()
	default:
		return false
	}
}

// NullValue returns the Value for Null/nil
func NullValue() *Value { return &nullValue }

func False() *Value { return &falseValue }

func True() *Value { return &trueValue }

// Type of this value
func (v *Value) Type() int {
	if v == nil {
		return VarNull
	}

	return v.varType
}

// BoolValue returns a Value for a bool
func BoolValue(i bool) *Value {
	if i {
		return &trueValue
	}
	return &falseValue
}

// IntValue returns a Value for an int64
func IntValue(i int64) *Value {
	return &Value{varType: VarInt, intVal: i}
}

// FloatValue returns a Value for an float64
func FloatValue(i float64) *Value {
	return &Value{varType: VarFloat, floatVal: i}
}

// ComplexValue returns a Value for a complex128
func ComplexValue(i complex128) *Value {
	return &Value{varType: VarComplex, complexVal: i}
}

// StringValue returns a Value for an string
func StringValue(i string) *Value {
	return &Value{varType: VarString, stringVal: i}
}

// IsNull returns true if the Value is null
func (v *Value) IsNull() bool {
	return v == nil || v.varType == VarNull
}

// IsZero returns true if the Value is null, false, 0, 0.0 or "" dependent on it's type
func (v *Value) IsZero() bool {
	if v == nil {
		return false
	}

	switch v.varType {
	case VarNull:
		return true
	case VarBool:
		return !v.boolVal
	case VarInt:
		return v.intVal == 0
	case VarFloat:
		return v.floatVal == 0.0
	case VarComplex:
		return v.complexVal == 0+0i
	case VarString:
		return v.stringVal == ""
	default:
		return false
	}
}

// IsNumeric returns true if the Value is a number, i.e. int64 or float64
func (v *Value) IsNumeric() bool {
	return v != nil && (v.varType == VarInt || v.varType == VarFloat)
}

func (v *Value) IsComplex() bool {
	return v != nil && (v.varType == VarComplex)
}

func (v *Value) IsNegative() bool {
	if v == nil {
		return false
	}

	switch v.varType {
	case VarInt:
		return v.intVal < 0
	case VarFloat:
		return v.floatVal < 0.0
	default:
		return false
	}
}

// AsBool converts Value into a Bool
func (v *Value) AsBool() *Value {
	return BoolValue(v.Bool())
}

// Bool returns the value as a bool
func (v *Value) Bool() bool {
	if v == nil {
		return false
	}

	switch v.varType {
	case VarBool:
		return v.boolVal
	case VarInt:
		if v.intVal == 0 {
			return false
		}
		return true
	case VarFloat:
		if v.floatVal == 0 {
			return false
		}
		return true
	case VarComplex:
		if v.complexVal == 0+0i {
			return false
		}
		return true
	case VarString:
		if v.stringVal == "" || v.stringVal[0] == 'f' || v.stringVal[0] == 'F' {
			return false
		}
		return true
	default:
		return false
	}
}

// Int returns the value as an int64
func (v *Value) Int() int64 {
	if v == nil {
		return 0
	}

	switch v.varType {
	case VarBool:
		if v.boolVal {
			return 1
		}
		return 0
	case VarInt:
		return v.intVal
	case VarFloat:
		return int64(v.floatVal)
	case VarComplex:
		return int64(real(v.complexVal))
	case VarString:
		r, err := strconv.ParseInt(v.stringVal, 10, 64)
		if err == nil {
			return r
		}
		// Panic instead?
		return 0
	default:
		return 0
	}
}

// Float returns the value as a float64
func (v *Value) Float() float64 {
	if v == nil {
		return 0.0
	}

	switch v.varType {
	case VarBool:
		if v.boolVal {
			return 1.0
		}
		return 0.0
	case VarInt:
		return float64(v.intVal)
	case VarFloat:
		return v.floatVal
	case VarComplex:
		return real(v.complexVal)
	case VarString:
		r, err := strconv.ParseFloat(v.stringVal, 64)
		if err == nil {
			return r
		}
		// Panic instead?
		return 0.0
	default:
		return 0.0
	}
}

func (v *Value) Complex() complex128 {
	if v == nil {
		return complex(0, 0)
	}

	if v.varType == VarComplex {
		return v.complexVal
	}
	return complex(v.Float(), 0)
}

// Real returns the real value of this Value.
// For real numbers this is the same as Float() but for complex numbers this
// returns the real component.
func (v *Value) Real() float64 {
	if v == nil {
		return 0.0
	}
	if v.varType == VarComplex {
		return real(v.complexVal)
	}
	return v.Float()
}

// Imaginary returns the imaginary value of this Value.
// For real numbers this returns 0.0 but for complex numbers this returns the
// imaginary component.
func (v *Value) Imaginary() float64 {
	if v != nil && v.varType == VarComplex {
		return imag(v.complexVal)
	}
	return 0.0
}

// String returns the value as a string
func (v *Value) String() string {
	if v == nil {
		return "nil"
	}

	switch v.varType {
	case VarBool:
		if v.boolVal {
			return "true"
		}
		return "false"
	case VarInt:
		return strconv.FormatInt(v.intVal, 10)
	case VarFloat:
		return strconv.FormatFloat(v.floatVal, 'f', 10, 64)
	case VarComplex:
		return fmt.Sprintf("%v", v.complexVal)
	case VarString:
		return v.stringVal
	default:
		return ""
	}
}

// OperationType returns the type of the suggested value when performing some
// operation like addition or multiplication to keep the precision of the result.
// For example, if a Value is an Integer but the passed value is Float then
// we should use float.
func (v *Value) OperationType(b *Value) int {
	t := v.Type()
	if v.Type() == VarString || b.Type() == VarString {
		t = VarString
	} else if v.Type() == VarComplex || b.Type() == VarComplex {
		t = VarComplex
	} else if v.Type() == VarFloat || b.Type() == VarFloat {
		t = VarFloat
	} else if v.Type() == VarInt || b.Type() == VarInt {
		t = VarInt
	}
	return t
}

func InterfaceVal(v interface{}) *Value {
	return &Value{varType: VarInterface, interfaceVal: v}
}

func (v *Value) Interface() interface{} {
	if v == nil {
		return nil
	}

	switch v.varType {
	case VarBool:
		return v.boolVal
	case VarInt:
		return v.intVal
	case VarFloat:
		return v.floatVal
	case VarComplex:
		return v.complexVal
	case VarString:
		return v.stringVal
	default:
		return ""
	}
}
