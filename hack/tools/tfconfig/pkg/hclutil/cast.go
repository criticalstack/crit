package hclutil

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

func ToString(v cty.Value) (s string, err error) {
	if isZero(v) {
		return "", nil
	}

	switch t := v.Type(); {
	case t == cty.String:
		return v.AsString(), nil
	case t == cty.Number:
		var val float64
		if err := gocty.FromCtyValue(v, &val); err != nil {
			return "", err
		}
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	case t == cty.Bool:
		var val bool
		if err := gocty.FromCtyValue(v, &val); err != nil {
			return "", err
		}
		return strconv.FormatBool(val), nil
	case t.IsListType(), t.IsTupleType():
		var ss []string
		v.ForEachElement(func(key cty.Value, val cty.Value) (stop bool) {
			s, err = ToString(val)
			if err != nil {
				return true
			}
			if val.Type() == cty.String {
				s = strconv.Quote(s)
			}
			ss = append(ss, s)
			return
		})
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("[%s]", strings.Join(ss, ", ")), nil
	default:
		return "", errors.Errorf("received unknown type %v: %q", v.Type().FriendlyName(), v.GoString())
	}
}

func isZero(v cty.Value) bool {
	return v == cty.NilVal || v == cty.NullVal(cty.String) || v == cty.NullVal(cty.Number) || v == cty.NullVal(cty.Bool)
}

func ToValue(s string, t cty.Type) (cty.Value, error) {
	if s == "" {
		return cty.NullVal(t), nil
	}

	switch {
	case t == cty.String:
		return gocty.ToCtyValue(s, cty.String)
	case t == cty.Number:
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return cty.NilVal, err
		}
		return gocty.ToCtyValue(n, cty.Number)
	case t == cty.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return cty.NilVal, err
		}
		return gocty.ToCtyValue(b, cty.Bool)
	case t == cty.DynamicPseudoType:
		return gocty.ToCtyValue(s, cty.String)
	case t.IsListType():
		ss, err := toStringSlice(s)
		if err != nil {
			return cty.NilVal, err
		}
		if len(ss) == 0 {
			return cty.ListValEmpty(t.ElementType()), nil
		}
		var vs []cty.Value
		for _, s := range ss {
			v, err := ToValue(s, t.ElementType())
			if err != nil {
				return cty.NilVal, err
			}
			vs = append(vs, v)
		}
		return cty.ListVal(vs), nil
	case t.IsTupleType():
		ss, err := toStringSlice(s)
		if err != nil {
			return cty.NilVal, err
		}
		if len(ss) != t.Length() {
			return cty.NilVal, errors.Errorf("expected %d tuple elements, received %d", t.Length(), len(ss))
		}
		var vs []cty.Value
		for i, s := range ss {
			v, err := ToValue(s, t.TupleElementType(i))
			if err != nil {
				return cty.NilVal, err
			}
			vs = append(vs, v)
		}
		return cty.TupleVal(vs), nil
	default:
		return cty.NilVal, errors.Errorf("received unknown type %v: %q", t.FriendlyName(), s)
	}
}

func toStringSlice(s string) (ss []string, err error) {
	var list []interface{}
	if err := json.Unmarshal([]byte(s), &list); err != nil {
		return nil, err
	}
	for _, item := range list {
		ss = append(ss, fmt.Sprintf("%v", item))
	}
	return
}
