package process

import (
	"path/filepath"

	"github.com/prometheus/procfs"
)

func IsCommandRunning(name string) (bool, error) {
	procs, err := procfs.AllProcs()
	if err != nil {
		return false, err
	}
	for _, p := range procs {
		cmdline, err := p.CmdLine()
		if err != nil {
			return false, err
		}
		if len(cmdline) == 0 {
			continue
		}
		if filepath.Base(cmdline[0]) == name {
			return true, nil
		}
	}
	return false, nil
}
