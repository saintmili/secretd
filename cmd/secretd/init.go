package main

import (
	"github.com/saintmili/secretd/internal/app"
	"github.com/spf13/cobra"
)

func (c *CLI) initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := app.Init(c.App)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
