package main

import (
	"fmt"
	"log"
	"os"

	"github.com/saintmili/secretd/internal/app"
	"github.com/saintmili/secretd/internal/config"
	"github.com/saintmili/secretd/internal/lock"
	"github.com/spf13/cobra"
)

type CLI struct {
	App *app.App
}

func execute() error {
	cfg, warnings, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	for _, w := range warnings {
		fmt.Printf("⚠️  config.%s: %s\n", w.Field, w.Message)
	}

	appInstance, err := app.New(cfg)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}
	defer appInstance.Logger.Close()

	appInstance.Logger.Info("secretd starting")

	// only one instance at a time
	lock, err := lock.Aquire()
	if err != nil {
		log.Fatalln(err)
	}
	defer lock.Release()

	cli := &CLI{
		App: appInstance,
	}

	rootCmd := cli.newRootCmd()
	return rootCmd.Execute()
}

func (c *CLI) newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "secretd",
		Short:         "Local-first encrypted secret manager",
		Long: `secretd is a local-first, encrypted secret manager.

It stores secrets securely on disk and never syncs data automatically.
`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)

	// Register commands here
	cmd.AddCommand(c.initCmd())
	cmd.AddCommand(c.listCmd())
	cmd.AddCommand(c.addCmd())

	return cmd
}
