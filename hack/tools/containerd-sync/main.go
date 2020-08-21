package main

import (
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use: "containerd-sync [cache-dir]",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(1)
		}
		dir := args[0]
		is, err := newImageSyncer(dir, viper.GetString("namespace"))
		if err != nil {
			log.Fatal(err)
		}
		if err := is.Sync(); err != nil {
			log.Fatal(err)
		}
		if viper.GetBool("watch") {
			log.Printf("watching %q ...", dir)
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					if err := is.Sync(); err != nil {
						log.Printf("sync error: %v", err)
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.Flags().String("namespace", "k8s.io", "containerd images namespace")
	rootCmd.Flags().Bool("watch", false, "watch cache dir for changes")
	viper.BindPFlags(rootCmd.Flags())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
