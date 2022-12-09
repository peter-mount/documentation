package util

import (
	"fmt"
	"strconv"
)

// The types of a Value returned by Type()
const (
	// Ignore 0 so use _ then if someone manually creates Value its an unknown type
	_ = iota
	VAR_NULL
	VAR_BOOL
	VAR_INT
	VAR_FLOAT
	VAR_COMPLEX
	VAR_STRING
)

type Value struct {
	varType    int
	boolVal    bool
	intVal     int64
	floatVal   float64
	complexVal complex128
	stringVal  string
}

var nullValue Value = Value{varType: VAR_NULL}
var falseValue Value = Value{varType: VAR_BOOL, boolVal: false}
var trueValue Value = Value{varType: VAR_BOOL, boolVal: true}

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
	case VAR_BOOL:
		return v.Bool() == b.Bool()
	case VAR_INT:
		return v.Int() == b.Int()
	case VAR_FLOAT:
		return v.Float() == b.Float()
	case VAR_COMPLEX:
		return v.Complex() == b.Complex()
	case VAR_STRING:
		return v.String() == b.String()
	default:
		return false
	}
}

// NullValue returns the Value for Null/nil
func NullValue() *Value {
	return &nullValue
}

// Type of this value
func (v *Value) Type() int {
	if v == nil {
		return VAR_NULL
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
	return &Value{varType: VAR_INT, intVal: i}
}

// FloatValue returns a Value for an float64
func FloatValue(i float64) *Value {
	return &Value{varType: VAR_FLOAT, floatVal: i}
}

// ComplexValue returns a Value for a complex128
func ComplexValue(i complex128) *Value {
	return &Value{varType: VAR_COMPLEX, complexVal: i}
}

// StringValue returns a Value for an string
func StringValue(i string) *Value {
	return &Value{varType: VAR_STRING, stringVal: i}
}

// IsNull returns true if the Value is null
func (v *Value) IsNull() bool {
	return v == nil || v.varType == VAR_NULL
}

// IsZero returns true if the Value is null, false, 0, 0.0 or "" dependent on it's type
func (v *Value) IsZero() bool {
	if v == nil {
		return false
	}

	switch v.varType {
	case VAR_NULL:
		return true
	case VAR_BOOL:
		return !v.boolVal
	case VAR_INT:
		return v.intVal == 0
	case VAR_FLOAT:
		return v.floatVal == 0.0
	case VAR_COMPLEX:
		return v.complexVal == 0+0i
	case VAR_STRING:
		return v.stringVal == ""
	default:
		return false
	}
}

// IsNumeric returns true if the Value is a number, i.e. int64 or float64
func (v *Value) IsNumeric() bool {
	return v != nil && (v.varType == VAR_INT || v.varType == VAR_FLOAT)
}

func (v *Value) IsComplex() bool {
	return v != nil && (v.varType == VAR_COMPLEX)
}

func (v *Value) IsNegative() bool {
	if v == nil {
		return false
	}

	switch v.varType {
	case VAR_INT:
		return v.intVal < 0
	case VAR_FLOAT:
		return v.floatVal < 0.0
	default:
		return false
	}
}

// Bool returns the value as a bool
func (v *Value) Bool() bool {
	if v == nil {
		return false
	}

	switch v.varType {
	case VAR_BOOL:
		return v.boolVal
	case VAR_INT:
		if v.intVal == 0 {
			return false
		}
		return true
	case VAR_FLOAT:
		if v.floatVal == 0 {
			return false
		}
		return true
	case VAR_COMPLEX:
		if v.complexVal == 0+0i {
			return false
		}
		return true
	case VAR_STRING:
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
	case VAR_BOOL:
		if v.boolVal {
			return 1
		}
		return 0
	case VAR_INT:
		return v.intVal
	case VAR_FLOAT:
		return int64(v.floatVal)
	case VAR_COMPLEX:
		return int64(real(v.complexVal))
	case VAR_STRING:
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
	case VAR_BOOL:
		if v.boolVal {
			return 1.0
		}
		return 0.0
	case VAR_INT:
		return float64(v.intVal)
	case VAR_FLOAT:
		return v.floatVal
	case VAR_COMPLEX:
		return real(v.complexVal)
	case VAR_STRING:
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

	if v.varType == VAR_COMPLEX {
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
	if v.varType == VAR_COMPLEX {
		return real(v.complexVal)
	}
	return v.Float()
}

// Imaginary returns the imaginary value of this Value.
// For real numbers this returns 0.0 but for complex numbers this returns the
// imaginary component.
func (v *Value) Imaginary() float64 {
	if v != nil && v.varType == VAR_COMPLEX {
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
	case VAR_BOOL:
		if v.boolVal {
			return "true"
		}
		return "false"
	case VAR_INT:
		return strconv.FormatInt(v.intVal, 10)
	case VAR_FLOAT:
		return strconv.FormatFloat(v.floatVal, 'f', 10, 64)
	case VAR_COMPLEX:
		return fmt.Sprintf("%v", v.complexVal)
	case VAR_STRING:
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
	if v.Type() == VAR_STRING || b.Type() == VAR_STRING {
		t = VAR_STRING
	} else if v.Type() == VAR_COMPLEX || b.Type() == VAR_COMPLEX {
		t = VAR_COMPLEX
	} else if v.Type() == VAR_FLOAT || b.Type() == VAR_FLOAT {
		t = VAR_FLOAT
	} else if v.Type() == VAR_INT || b.Type() == VAR_INT {
		t = VAR_INT
	}
	return t
}
