package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"

	"github.com/kom0055/git-mirror/pkg/options"
	"github.com/kom0055/git-mirror/pkg/utils"
)

var (
	option = options.Option{}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gclone",
	Short: "clone repo to local path or sync all available repos to local path",
	Long: `to clone repo: gclone https://github.com/kom0055/git-mirror
to sync all available repos: gclone --sync-from-remote --remote-type=gh https://github.com `,

	Run: func(cmd *cobra.Command, args []string) {
		if err := os.MkdirAll(utils.DefaultTmpPath, 0755); err != nil {
			log.Fatalf("failed to create tmp dir: %v\n", err)
		}
		defer func() {
			_ = os.RemoveAll(utils.DefaultTmpPath)
		}()

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
	option = options.Option{
		Source: options.BasicOpt{
			EcdsaPemFile:       os.Getenv("SOURCE_ECDSA_PEM_FILE"),
			EcdsaPemFilePasswd: os.Getenv("SOURCE_ECDSA_PEM_FILE_PASSWD"),
			RemoteGitlabAddr:   os.Getenv("SOURCE_REMOTE_GITLAB_ADDR"),
			User:               os.Getenv("SOURCE_USER"),
			Token:              os.Getenv("SOURCE_TOKEN"),
			Proto:              os.Getenv("SOURCE_PROTO"),
			GroupName:          os.Getenv("SOURCE_GROUP_NAME"),
		},
		Dest: options.BasicOpt{
			EcdsaPemFile:       os.Getenv("DEST_ECDSA_PEM_FILE"),
			EcdsaPemFilePasswd: os.Getenv("DEST_ECDSA_PEM_FILE_PASSWD"),
			RemoteGitlabAddr:   os.Getenv("DEST_REMOTE_GITLAB_ADDR"),
			User:               os.Getenv("DEST_USER"),
			Token:              os.Getenv("DEST_TOKEN"),
			Proto:              os.Getenv("DEST_PROTO"),
			GroupName:          os.Getenv("DEST_GROUP_NAME"),
		},
		Worker: 8,
	}

	rootCmd.PersistentFlags().IntVar(&option.Worker, "worker", option.Worker, "worker num")

	rootCmd.PersistentFlags().StringVar(&option.Source.EcdsaPemFile, "source-ecdsa", option.Source.EcdsaPemFile, "ecdsa pem file")
	rootCmd.PersistentFlags().StringVar(&option.Source.EcdsaPemFilePasswd, "source-ecdsa-passwd", option.Source.EcdsaPemFilePasswd, "ecdsa pem file passwd")
	rootCmd.PersistentFlags().StringVar(&option.Source.RemoteGitlabAddr, "source-remote-gitlab-addr", option.Source.RemoteGitlabAddr, "remote gitlab addr")
	rootCmd.PersistentFlags().StringVar(&option.Source.User, "source-user", option.Source.User, "user name")
	rootCmd.PersistentFlags().StringVar(&option.Source.Token, "source-token", option.Source.Token, "token or private key")
	rootCmd.PersistentFlags().StringVar(&option.Source.Proto, "source-proto", option.Source.Proto, "proto: git, ssh, http or https")

	rootCmd.PersistentFlags().StringVar(&option.Dest.EcdsaPemFile, "dest-ecdsa", option.Dest.EcdsaPemFile, "ecdsa pem file")
	rootCmd.PersistentFlags().StringVar(&option.Dest.EcdsaPemFilePasswd, "dest-ecdsa-passwd", option.Dest.EcdsaPemFilePasswd, "ecdsa pem file passwd")
	rootCmd.PersistentFlags().StringVar(&option.Dest.RemoteGitlabAddr, "dest-remote-gitlab-addr", option.Dest.RemoteGitlabAddr, "remote gitlab addr")
	rootCmd.PersistentFlags().StringVar(&option.Dest.User, "dest-user", option.Dest.User, "user name")
	rootCmd.PersistentFlags().StringVar(&option.Dest.Token, "dest-token", option.Dest.Token, "token or private key")
	rootCmd.PersistentFlags().StringVar(&option.Dest.Proto, "dest-proto", option.Dest.Proto, "proto: git, ssh, http or https")

}
