package hclutil

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/terraform/configs"
	"github.com/zclconf/go-cty/cty"
)

func ReadVars(path string) (Variables, error) {
	parser := configs.NewParser(nil)
	modules, diags := parser.LoadConfigDir(path)
	if diags.HasErrors() {
		return nil, diags
	}
	variables := make(map[string]*Variable)
	for name, v := range modules.Variables {
		variables[name] = &Variable{Variable: v}
	}
	files, err := ReadDir(path)
	if err != nil {
		return nil, err
	}
	for name, f := range files {
		if strings.ToLower(filepath.Ext(name)) != ".tfvars" {
			continue
		}
		var varfile map[string]cty.Value
		diags := gohcl.DecodeBody(f.Body, nil, &varfile)
		if diags.HasErrors() {
			return nil, diags
		}
		for name, v := range varfile {
			if variable, ok := variables[name]; ok {
				variable.value = v
			}
		}
	}
	return variables, nil
}

func ReadDir(path string) (map[string]*hcl.File, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	p := hclparse.NewParser()
	for _, f := range files {
		switch strings.ToLower(filepath.Ext(f.Name())) {
		case ".tf", ".hcl", ".tfvars":
			_, diags := p.ParseHCLFile(filepath.Join(path, f.Name()))
			if diags.HasErrors() {
				return nil, diags
			}
		}
	}
	return p.Files(), nil
}

func WriteVars(path string, vars Variables) error {
	return ioutil.WriteFile(path, vars.Bytes(), 0644)
}
