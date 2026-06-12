package internal

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// reflectField pairs a field's name with its value for JSON encoding.
type reflectField struct {
	name  string
	value reflect.Value
}

// collectExportedFields returns all exported leaf fields, flattening anonymous embedded structs.
func collectExportedFields(val reflect.Value) []reflectField {
	t := val.Type()
	var fields []reflectField

	for i := range t.NumField() {
		sf := t.Field(i)
		fv := val.Field(i)

		if sf.Anonymous && fv.Kind() == reflect.Struct {
			fields = append(fields, collectExportedFields(fv)...)
			continue
		}

		if !sf.IsExported() {
			continue
		}

		fields = append(fields, reflectField{sf.Name, fv})
	}

	return fields
}

// MarshalToJSON serialises exported fields of a struct to a JSON byte slice.
// Supported field types: string, bool, int*, uint*, float*.
// Pointer-to-struct is accepted; nil pointer returns an error.
// Anonymous embedded structs are flattened into the output object.
func MarshalToJSON(v any) ([]byte, error) {
	val := reflect.ValueOf(v)
	if !val.IsValid() {
		return nil, fmt.Errorf("marshalToJSON: cannot marshal nil value")
	}

	for val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return nil, fmt.Errorf("marshalToJSON: cannot marshal nil pointer")
		}

		val = val.Elem()
	}

	encoded, err := encodeJSONValue(val)
	if err != nil {
		return nil, fmt.Errorf("marshalToJSON: %w", err)
	}

	return []byte(encoded), nil
}

func encodeJSONValue(v reflect.Value) (string, error) {
	switch v.Kind() {
	case reflect.Struct:
		return encodeStructValue(v)

	case reflect.Slice, reflect.Array:
		return encodeArrayValue(v)

	default:
		return encodeFieldValue(v)
	}
}

func encodeStructValue(v reflect.Value) (string, error) {
	fields := collectExportedFields(v)

	var b strings.Builder
	b.WriteByte('{')
	written := 0

	for _, f := range fields {
		encoded, err := encodeJSONValue(f.value)
		if err != nil {
			return "", fmt.Errorf("field %s: %w", f.name, err)
		}

		if encoded == "" {
			continue
		}

		if written > 0 {
			b.WriteString(", ")
		}

		b.WriteByte('"')
		b.WriteString(f.name)
		b.WriteString(`": `)
		b.WriteString(encoded)
		written++
	}

	b.WriteByte('}')
	return b.String(), nil
}

func encodeArrayValue(v reflect.Value) (string, error) {
	var b strings.Builder
	b.WriteByte('[')
	written := 0

	for i := range v.Len() {
		encoded, err := encodeJSONValue(v.Index(i))
		if err != nil {
			return "", fmt.Errorf("index %d: %w", i, err)
		}

		if encoded == "" {
			continue
		}

		if written > 0 {
			b.WriteString(", ")
		}
		b.WriteString(encoded)
		written++
	}

	b.WriteByte(']')
	return b.String(), nil
}

func encodeFieldValue(v reflect.Value) (string, error) {
	switch v.Kind() {
	case reflect.String:
		return encodeJSONString(v.String()), nil
	case reflect.Bool:
		if v.Bool() {
			return "true", nil
		}
		return "false", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64), nil
	default:
		return "", nil
	}
}

func encodeJSONString(s string) string {
	var b strings.Builder
	b.WriteByte('"')

	for _, c := range s {
		switch c {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		default:
			b.WriteRune(c)
		}
	}

	b.WriteByte('"')

	return b.String()
}

// UnmarshalFromJSON deserialises a JSON byte slice into a dynamically instantiated value.
// It handles structs and slices of structs. Supported field types: string, bool, int*, uint*, float*.
func UnmarshalFromJSON(data []byte, t reflect.Type) (reflect.Value, error) {
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return reflect.Value{}, fmt.Errorf("unmarshalFromJSON: %w", err)
	}

	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	val := reflect.New(t).Elem()
	if err := populateFromAny(val, raw); err != nil {
		return reflect.Value{}, fmt.Errorf("unmarshalFromJSON: %w", err)
	}

	return val, nil
}

func populateFromAny(v reflect.Value, raw any) error {
	if raw == nil {
		return nil
	}

	t := v.Type()
	switch t.Kind() {
	case reflect.Slice:
		return populateSlice(v, raw)
	case reflect.Struct:
		return populateStruct(v, raw)
	default:
		return setFieldValue(v, t.Kind(), raw)
	}
}

func populateSlice(v reflect.Value, raw any) error {
	rawSlice, ok := raw.([]any)
	if !ok {
		return fmt.Errorf("expected array, got %T", raw)
	}

	t := v.Type()
	slice := reflect.MakeSlice(t, len(rawSlice), len(rawSlice))
	for i, rawElem := range rawSlice {
		if err := populateFromAny(slice.Index(i), rawElem); err != nil {
			return fmt.Errorf("index %d: %w", i, err)
		}
	}
	v.Set(slice)
	return nil
}

func populateStruct(v reflect.Value, raw any) error {
	rawMap, ok := raw.(map[string]any)
	if !ok {
		return fmt.Errorf("expected object, got %T", raw)
	}

	t := v.Type()
	for key, rawVal := range rawMap {
		field, ok := t.FieldByName(key)
		if !ok || !field.IsExported() {
			continue
		}

		if err := populateFromAny(v.FieldByName(key), rawVal); err != nil {
			return fmt.Errorf("field %s: %w", key, err)
		}
	}
	return nil
}

func setFieldValue(field reflect.Value, kind reflect.Kind, raw any) error {
	switch kind {
	case reflect.String:
		s, ok := raw.(string)
		if !ok {
			return fmt.Errorf("cannot assign %T to string", raw)
		}
		field.SetString(s)

	case reflect.Bool:
		b, ok := raw.(bool)
		if !ok {
			return fmt.Errorf("cannot assign %T to bool", raw)
		}
		field.SetBool(b)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f, ok := raw.(float64)
		if !ok {
			return fmt.Errorf("cannot assign %T to int", raw)
		}
		field.SetInt(int64(math.Round(f)))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		f, ok := raw.(float64)
		if !ok {
			return fmt.Errorf("cannot assign %T to uint", raw)
		}
		field.SetUint(uint64(math.Round(f)))

	case reflect.Float32, reflect.Float64:
		f, ok := raw.(float64)
		if !ok {
			return fmt.Errorf("cannot assign %T to float", raw)
		}
		field.SetFloat(f)

	default:
		return nil
	}

	return nil
}
