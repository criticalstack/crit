package process_test

import (
	"os/exec"
	"testing"

	"github.com/criticalstack/crit/pkg/util/process"
)

func TestIsCommandRunning(t *testing.T) {
	// Doesn't appear to work with CI
	t.Skip()
	running, err := process.IsCommandRunning("notarealprocessthatshouldberunning")
	if err != nil {
		t.Fatal(err)
	}
	if running {
		t.Fatalf("expected fake process to not report as running")
	}
	go func() {
		cmd := exec.Command("sleep", "infinity")
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := cmd.Process.Kill(); err != nil {
				t.Fatal(err)
			}
		}()
	}()

	running, err = process.IsCommandRunning("sleep")
	if err != nil {
		t.Fatal(err)
	}
	if !running {
		t.Fatalf("expected sleep process to be running")
	}
}
