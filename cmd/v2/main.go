package main

import (
	"log"
	"os"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command/apply"
	"github.com/craftcms/nitro/command/initialize"
	"github.com/craftcms/nitro/command/start"
	"github.com/craftcms/nitro/command/stop"
	"github.com/craftcms/nitro/command/update"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/terminal"
)

var rootCommand = &cobra.Command{
	Use:          "nitro",
	Short:        "Local Craft CMS dev made easy",
	Long:         `Nitro is a command-line tool focused on making local Craft CMS development quick and easy.`,
	RunE:         rootMain,
	SilenceUsage: true,
}

func rootMain(command *cobra.Command, _ []string) error {
	return command.Help()
}

func init() {
	env, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// set any global flags
	flags := rootCommand.PersistentFlags()
	flags.StringP("environment", "e", env, "The environment")

	// create the docker client
	client, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	// create the "terminal" for capturing output
	term := terminal.New()

	// register all of the commands
	commands := []*cobra.Command{
		initialize.New(client, term),
		stop.New(client, term),
		start.New(client, term),
		//start.StartCommand,
		// destroy.DestroyCommand,
		// restart.RestartCommand,
		// ls.LSCommand,
		// composer.ComposerCommand,
		// npm.NPMCommand,
		//complete.CompleteCommand,
		apply.New(client, term),
		//context.ContextCommand,
		//exec.ExecCommand,
		//trust.TrustCommand,
		//db.DBCommand,
		update.New(),
	}

	// add the commands
	rootCommand.AddCommand(commands...)
}

func main() {
	// execute the root command
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
