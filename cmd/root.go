package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"

	"github.com/kom0055/gclone/pkg/options"
	"github.com/kom0055/gclone/pkg/utils"
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
		if err := os.MkdirAll(utils.DefaultTmpPath, 0755); err != nil {
			log.Fatalf("failed to create tmp dir: %v\n", err)
		}
		defer func() {
			_ = os.RemoveAll(utils.DefaultTmpPath)
		}()
		//if err := option.Complete(); err != nil {
		//	log.Fatalf("failed to complete option: %v\n", err)
		//}
		if err := option.Mirror(cmd.Context()); err != nil {
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
	rootCmd.PersistentFlags().StringVar(&option.Source.EcdsaPemFile, "source-ecdsa", "", "ecdsa pem file")
	rootCmd.PersistentFlags().StringVar(&option.Source.EcdsaPemFilePasswd, "source-ecdsa-passwd", "", "ecdsa pem file passwd")
	rootCmd.PersistentFlags().StringVar(&option.Source.RemoteGitlabAddr, "source-remote-gitlab-addr", "", "remote gitlab addr")
	rootCmd.PersistentFlags().StringVar(&option.Source.User, "source-user", "", "user name")
	rootCmd.PersistentFlags().StringVar(&option.Source.Token, "source-token", "", "token or private key")
	rootCmd.PersistentFlags().StringVar(&option.Source.Proto, "source-proto", "", "proto: git, ssh, http or https")

	rootCmd.PersistentFlags().StringVar(&option.Dest.EcdsaPemFile, "dest-ecdsa", "", "ecdsa pem file")
	rootCmd.PersistentFlags().StringVar(&option.Dest.EcdsaPemFilePasswd, "dest-ecdsa-passwd", "", "ecdsa pem file passwd")
	rootCmd.PersistentFlags().StringVar(&option.Dest.RemoteGitlabAddr, "dest-remote-gitlab-addr", "", "remote gitlab addr")
	rootCmd.PersistentFlags().StringVar(&option.Dest.User, "dest-user", "", "user name")
	rootCmd.PersistentFlags().StringVar(&option.Dest.Token, "dest-token", "", "token or private key")
	rootCmd.PersistentFlags().StringVar(&option.Dest.Proto, "dest-proto", "", "proto: git, ssh, http or https")

}
