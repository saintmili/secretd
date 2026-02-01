package main

import (
	"github.com/saintmili/secretd/internal/app"
	"github.com/spf13/cobra"
)

func (c *CLI) listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists all the secrets",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := app.List(c.App)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
