package util

import (
	"fmt"
	"sort"
)

func BuildArgumentListFromMap(baseArguments, overrideArguments map[string]string) (command []string) {
	args := make(map[string]string)
	for k, v := range baseArguments {
		args[k] = v
	}
	for k, v := range overrideArguments {
		args[k] = v
	}
	for k, v := range args {
		command = append(command, fmt.Sprintf("--%s=%s", k, v))
	}
	sort.Strings(command)
	return command
}
