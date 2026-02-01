package main

import (
	"errors"

	"github.com/saintmili/secretd/internal/app"
	"github.com/spf13/cobra"
)

func (c *CLI) addCmd() *cobra.Command {
	var title, username, password string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new secret",
		RunE: func(cmd *cobra.Command, args []string) error {
			if title == "" {
				return errors.New("title is required")
			}

			return app.Add(c.App)
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Entry title")
	cmd.Flags().StringVar(&username, "username", "", "Username")
	cmd.Flags().StringVar(&password, "password", "", "Password")

	_ = cmd.MarkFlagRequired("title")

	return cmd
}

