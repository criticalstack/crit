package hclutil

import (
	"sort"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform/configs"
	"github.com/zclconf/go-cty/cty"
)

// Variable represents an abstract of the coalesced state of input variables
// and defined variables for a given set of terraform scripts. This is used to
// prompt a command-line user to supply values for variables, therefore it
// implements the prompt.ValuePrompter interface.
type Variable struct {
	*configs.Variable

	value cty.Value
}

func (v *Variable) IsNull() bool {
	return v.value.IsNull() && v.Variable.Default.IsNull()
}

func (v *Variable) Value() cty.Value {
	if !v.value.IsNull() {
		return v.value
	}
	return v.Variable.Default
}

func (v *Variable) SetValue(value cty.Value) {
	v.value = value
}

func (v *Variable) Default() string {
	if !v.IsNull() {
		val, err := ToString(v.Value())
		if err != nil {
			panic(err)
		}
		return val
	}
	return ""
}

func (v *Variable) Message() string {
	if v.DescriptionSet {
		return v.Description
	}
	return v.Name
}

func (v *Variable) Type() cty.Type {
	return v.Variable.Type
}

type Variables map[string]*Variable

func (vars Variables) String() string {
	return string(vars.Bytes())
}

func (vars Variables) Bytes() []byte {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	file := hclwrite.NewFile()
	for _, k := range keys {
		v := vars[k]
		if v.IsNull() {
			continue
		}
		file.Body().SetAttributeValue(v.Name, v.Value())
	}
	return file.Bytes()
}
