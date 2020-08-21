package exec

import (
	"fmt"
	"os"
	"testing"
)

func TestPrefixWriter(t *testing.T) {
	out := NewPrefixWriter(os.Stdout, "\t")
	defer out.Close()

	for i := 0; i < 20; i++ {
		out.Write([]byte(fmt.Sprintf("line %d", i)))
		if i%3 == 0 {
			out.Write([]byte("\n"))
		}
	}
}
