package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/kom0055/gclone/pkg/config"
	"github.com/kom0055/gclone/pkg/options"
)

var (
	option options.Option
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gclone",
	Short: "clone repo to local path or sync all available repos to local path",
	Long: `to clone repo: gclone https://github.com/kom0055/gclone
to sync all available repos: gclone --sync-from-remote --remote-type=gh https://github.com `,

	Run: func(cmd *cobra.Command, args []string) {
		if err := option.Complete(); err != nil {
			log.Fatalln(err)
		}
		if err := option.Run(cmd.Context(), args); err != nil {
			log.Fatalln(err)
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVar(&option.CfgFilePath, "config", filepath.Join(config.HomeDir, config.CfgFilePath), "config file path")
	rootCmd.PersistentFlags().StringVar(&option.RepoRootPath, "repo-root-path", "", "repo root path")
	rootCmd.PersistentFlags().StringVar(&option.EcdsaPemFile, "ecdsa", "", "ecdsa pem file")
	rootCmd.PersistentFlags().StringVar(&option.EcdsaPemFilePasswd, "ecdsa-passwd", "", "ecdsa pem file passwd")
	rootCmd.PersistentFlags().StringVar(&option.RemoteGitlabAddr, "remote-gitlab-addr", "", "remote gitlab addr")
	rootCmd.PersistentFlags().BoolVar(&option.SyncFromRemote, "sync-from-remote", false, "sync all repos from remote")

	rootCmd.PersistentFlags().StringVar(&option.User, "user", "", "user name")
	rootCmd.PersistentFlags().StringVar(&option.Token, "token", "", "token or private key")

}
