package prompt

import (
	"fmt"
	"os"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/zclconf/go-cty/cty"

	"github.com/criticalstack/crit/hack/tools/tfconfig/pkg/hclutil"
)

type ValuePrompter interface {
	Message() string
	Default() string
	Type() cty.Type
	Value() cty.Value
	SetValue(cty.Value)
}

func Input(p ValuePrompter) (err error) {
	prompt := &survey.Input{
		Message: p.Message(),
		Default: p.Default(),
	}
	var answer string
	if err := survey.AskOne(prompt, &answer); err != nil {
		return err
	}
	val, err := hclutil.ToValue(answer, p.Type())
	if err != nil {
		return err
	}
	p.SetValue(val)
	return nil
}

func Select(p ValuePrompter, opts []string) (err error) {
	prompt := &survey.Select{
		Message: p.Message(),
		Options: opts,
		VimMode: true,
	}
	if p.Default() != "" {
		if contains(opts, p.Default()) {
			prompt.Default = p.Default()
		}
	}
	var answer string
	if err := survey.AskOne(prompt, &answer); err != nil {
		return err
	}
	val, err := hclutil.ToValue(answer, p.Type())
	if err != nil {
		return err
	}
	p.SetValue(val)
	return nil
}

func Confirm(msg string) bool {
	prompt := &survey.Confirm{
		Message: msg,
		Default: true,
	}
	confirm := false
	err := survey.AskOne(prompt, &confirm)
	if err == terminal.InterruptErr {
		fmt.Println(err)
		os.Exit(1)
	}
	return confirm
}

func Confirmf(msg string, args ...interface{}) bool {
	return Confirm(fmt.Sprintf(msg, args...))
}

type CustomValuePrompt func(ValuePrompter) error

func NewSelectNumberPrompt(opts ...int) CustomValuePrompt {
	return func(p ValuePrompter) (err error) {
		options := make([]string, 0)
		for _, i := range opts {
			options = append(options, strconv.Itoa(i))
		}
		return Select(p, options)
	}
}

func NewSelectPrompt(opts ...string) CustomValuePrompt {
	return func(p ValuePrompter) (err error) {
		return Select(p, opts)
	}
}

func NewDeferredSelectPrompt(fn func() []string) CustomValuePrompt {
	return func(p ValuePrompter) (err error) {
		return Select(p, fn())
	}
}

func contains(ss []string, match string) bool {
	for _, s := range ss {
		if s == match {
			return true
		}
	}
	return false
}
