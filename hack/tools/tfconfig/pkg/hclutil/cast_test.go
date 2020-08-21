package hclutil

import (
	"errors"
	"math"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zclconf/go-cty/cty"
)

func TestToStringTypes(t *testing.T) {
	tests := []struct {
		Name  string
		Value cty.Value
		Want  string
	}{
		{
			"string",
			cty.StringVal("test value"),
			"test value",
		},
		{
			"int",
			cty.NumberIntVal(10),
			"10",
		},
		{
			"float",
			cty.NumberFloatVal(1.5),
			"1.5",
		},
		{
			"bool",
			cty.BoolVal(true),
			"true",
		},
		{
			"list",
			cty.ListVal([]cty.Value{
				cty.StringVal("test value 1"),
				cty.StringVal("test value 2"),
			}),
			"[\"test value 1\", \"test value 2\"]",
		},
		{
			"tuple",
			cty.TupleVal([]cty.Value{
				cty.StringVal("test value"),
				cty.NumberIntVal(10),
				cty.NumberFloatVal(1.5),
				cty.BoolVal(true),
			}),
			"[\"test value\", 10, 1.5, true]",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got, err := ToString(test.Value)
			if err != nil {
				t.Fatalf("ToString returned error: %s", err)
			}

			if diff := cmp.Diff(test.Want, got); diff != "" {
				t.Errorf("wrong result: (-want +got)\n%s", diff)
			}
		})
	}
}

func TestToStringLists(t *testing.T) {
	tests := []struct {
		Name  string
		Value cty.Value
		Want  string
	}{
		{
			"list(string)",
			cty.ListVal([]cty.Value{
				cty.StringVal("test value 1"),
				cty.StringVal("test value 2"),
			}),
			"[\"test value 1\", \"test value 2\"]",
		},
		{
			"list(int)",
			cty.ListVal([]cty.Value{
				cty.NumberIntVal(10),
				cty.NumberIntVal(-10),
			}),
			"[10, -10]",
		},
		{
			"list(float)",
			cty.ListVal([]cty.Value{
				cty.NumberFloatVal(1.5),
				cty.NumberFloatVal(-1.5),
			}),
			"[1.5, -1.5]",
		},
		{
			"list(bool)",
			cty.ListVal([]cty.Value{
				cty.BoolVal(true),
				cty.BoolVal(false),
			}),
			"[true, false]",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got, err := ToString(test.Value)
			if err != nil {
				t.Fatalf("ToString returned error: %s", err)
			}

			if diff := cmp.Diff(test.Want, got); diff != "" {
				t.Errorf("wrong result: (-want +got)\n%s", diff)
			}
		})
	}
}

func TestToStringNumberConversion(t *testing.T) {
	tests := []struct {
		Name  string
		Value cty.Value
		Want  string
	}{
		{
			"int(0)",
			cty.NumberIntVal(0),
			"0",
		},
		{
			"int(max)",
			cty.NumberIntVal(1<<63 - 1),
			strconv.FormatFloat(1<<63-1, 'f', -1, 64),
		},
		{
			"int(min)",
			cty.NumberIntVal(-1 << 63),
			strconv.FormatFloat(-1<<63, 'f', -1, 64),
		},
		{
			"float(0)",
			cty.NumberFloatVal(0),
			strconv.FormatFloat(0, 'f', -1, 64),
		},
		{
			"float(max)",
			cty.NumberFloatVal(math.MaxFloat64),
			strconv.FormatFloat(math.MaxFloat64, 'f', -1, 64),
		},
		{
			"float(min)",
			cty.NumberFloatVal(math.SmallestNonzeroFloat64),
			strconv.FormatFloat(math.SmallestNonzeroFloat64, 'f', -1, 64),
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got, err := ToString(test.Value)
			if err != nil {
				t.Fatalf("ToString returned error: %s", err)
			}

			if diff := cmp.Diff(test.Want, got); diff != "" {
				t.Errorf("wrong result: (-want +got)\n%s", diff)
			}
		})
	}
}

func TestToStringEmptyValues(t *testing.T) {
	tests := []struct {
		Name  string
		Value cty.Value
		Want  string
	}{
		{
			"empty string",
			cty.StringVal(""),
			"",
		},
		{
			"nil",
			cty.NilVal,
			"",
		},
		{
			"nullval(string)",
			cty.NullVal(cty.String),
			"",
		},
		{
			"nullval(number)",
			cty.NullVal(cty.Number),
			"",
		},
		{
			"nullval(bool)",
			cty.NullVal(cty.Bool),
			"",
		},
		{
			"empty list(string)",
			cty.ListValEmpty(cty.String),
			"[]",
		},
		{
			"empty list(number)",
			cty.ListValEmpty(cty.Number),
			"[]",
		},
		{
			"empty list(bool)",
			cty.ListValEmpty(cty.Bool),
			"[]",
		},
		{
			"empty tuple",
			cty.EmptyTupleVal,
			"[]",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got, err := ToString(test.Value)
			if err != nil {
				t.Fatalf("ToString returned error: %s", err)
			}

			if diff := cmp.Diff(test.Want, got); diff != "" {
				t.Errorf("wrong result: (-want +got)\n%s", diff)
			}
		})
	}
}

func TestToValueTypes(t *testing.T) {
	tests := []struct {
		Name  string
		Value string
		Type  cty.Type
		Want  cty.Value
	}{
		{
			"string",
			"test value",
			cty.String,
			cty.StringVal("test value"),
		},
		{
			"int",
			"10",
			cty.Number,
			cty.NumberIntVal(10),
		},
		{
			"float",
			"1.5",
			cty.Number,
			cty.NumberFloatVal(1.5),
		},
		{
			"bool",
			"true",
			cty.Bool,
			cty.BoolVal(true),
		},
		{
			"list",
			"[\"test value 1\", \"test value 2\"]",
			cty.List(cty.String),
			cty.ListVal([]cty.Value{
				cty.StringVal("test value 1"),
				cty.StringVal("test value 2"),
			}),
		},
		{
			"tuple",
			"[\"test value\", 10, 1.5, true]",
			cty.Tuple([]cty.Type{cty.String, cty.Number, cty.Number, cty.Bool}),
			cty.TupleVal([]cty.Value{
				cty.StringVal("test value"),
				cty.NumberIntVal(10),
				cty.NumberFloatVal(1.5),
				cty.BoolVal(true),
			}),
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got, err := ToValue(test.Value, test.Type)
			if err != nil {
				t.Fatalf("ToValue returned error: %s", err)
			}

			diff := cmp.Diff(test.Want, got, cmp.Comparer(func(x, y cty.Value) bool {
				return x.RawEquals(y) || y.RawEquals(x)
			}))
			if diff != "" {
				t.Errorf("wrong result: (-want +got)\n%s", diff)
			}
		})
	}
}

func TestToValueLists(t *testing.T) {
	tests := []struct {
		Name  string
		Value string
		Type  cty.Type
		Want  cty.Value
	}{
		{
			"list(string)",
			"[\"test value 1\", \"test value 2\"]",
			cty.List(cty.String),
			cty.ListVal([]cty.Value{
				cty.StringVal("test value 1"),
				cty.StringVal("test value 2"),
			}),
		},
		{
			"list(int)",
			"[10, -10]",
			cty.List(cty.Number),
			cty.ListVal([]cty.Value{
				cty.NumberIntVal(10),
				cty.NumberIntVal(-10),
			}),
		},
		{
			"list(float)",
			"[1.5, -1.5]",
			cty.List(cty.Number),
			cty.ListVal([]cty.Value{
				cty.NumberFloatVal(1.5),
				cty.NumberFloatVal(-1.5),
			}),
		},
		{
			"list(bool)",
			"[true, false]",
			cty.List(cty.Bool),
			cty.ListVal([]cty.Value{
				cty.BoolVal(true),
				cty.BoolVal(false),
			}),
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got, err := ToValue(test.Value, test.Type)
			if err != nil {
				t.Fatalf("ToValue returned error: %s", err)
			}

			diff := cmp.Diff(test.Want, got, cmp.Comparer(func(x, y cty.Value) bool {
				return x.RawEquals(y) || y.RawEquals(x)
			}))
			if diff != "" {
				t.Errorf("wrong result: (-want +got)\n%s", diff)
			}
		})
	}
}

func TestToValueEmptyString(t *testing.T) {
	tests := []struct {
		Name  string
		Value string
		Type  cty.Type
		Want  cty.Value
	}{
		{
			"string",
			"",
			cty.String,
			cty.NullVal(cty.String),
		},
		{
			"number",
			"",
			cty.Number,
			cty.NullVal(cty.Number),
		},
		{
			"bool",
			"",
			cty.Bool,
			cty.NullVal(cty.Bool),
		},
		{
			"empty list(string)",
			"[]",
			cty.List(cty.String),
			cty.ListValEmpty(cty.String),
		},
		{
			"empty list(number)",
			"[]",
			cty.List(cty.Number),
			cty.ListValEmpty(cty.Number),
		},
		{
			"empty list(bool)",
			"[]",
			cty.List(cty.Bool),
			cty.ListValEmpty(cty.Bool),
		},
		{
			"empty tuple",
			"[]",
			cty.EmptyTuple,
			cty.EmptyTupleVal,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got, err := ToValue(test.Value, test.Type)
			if err != nil {
				t.Fatalf("ToValue returned error: %s", err)
			}

			diff := cmp.Diff(test.Want, got, cmp.Comparer(func(x, y cty.Value) bool {
				return x.Equals(y) == cty.True || y.Equals(x) == cty.True
			}))
			if diff != "" {
				t.Errorf("wrong result: (-want +got)\n%s", diff)
			}
		})
	}
}

func TestToValueInputTypeMismatch(t *testing.T) {
	tests := []struct {
		Name  string
		Value string
		Type  cty.Type
		Want  cty.Value
		Error error
	}{
		{
			"string into number",
			"test value",
			cty.Number,
			cty.NilVal,
			errors.New(""),
		},
		{
			"string into bool",
			"test value",
			cty.Bool,
			cty.NilVal,
			errors.New(""),
		},
		{
			"number into string",
			"10",
			cty.String,
			cty.StringVal("10"),
			nil,
		},
		{
			"number into bool",
			"10",
			cty.Bool,
			cty.NilVal,
			errors.New(""),
		},
		{
			"bool into string",
			"true",
			cty.String,
			cty.StringVal("true"),
			nil,
		},
		{
			"bool into number",
			"true",
			cty.Number,
			cty.NilVal,
			errors.New(""),
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got, err := ToValue(test.Value, test.Type)
			if err != nil && test.Error == nil {
				t.Fatalf("ToValue returned error: %s", err)
			}

			diff := cmp.Diff(test.Want, got, cmp.Comparer(func(x, y cty.Value) bool {
				return x.Equals(y) == cty.True || y.Equals(x) == cty.True
			}))
			if diff != "" {
				t.Errorf("wrong result: (-want +got)\n%s", diff)
			}
		})
	}
}
