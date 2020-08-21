package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/shurcooL/httpfs/filter"
	"github.com/shurcooL/vfsgen"
	"github.com/spf13/cobra"

	fsutil "github.com/criticalstack/crit/pkg/util/fs"
)

func main() {
	cmd := &cobra.Command{
		Use:  "gentmpls",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateDir := filter.Skip(fsutil.StripModTime(http.Dir(args[0])), func(path string, fi os.FileInfo) bool {
				return !fi.IsDir() && filepath.Ext(path) == ".go"
			})
			return vfsgen.Generate(templateDir, vfsgen.Options{
				BuildTags:    "!dev",
				VariableName: "Files",
				PackageName:  "cluster",
				Filename:     filepath.Join(args[1], "zz_generated.templates.go"),
			})
		},
	}
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
