package cmd

import (
	"ethbaas/internal/server"

	"github.com/spf13/cobra"
)

type ServerCmd struct{}

func NewServerCmd() *ServerCmd {
	c := &ServerCmd{}
	return c
}

func (s *ServerCmd) rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Server Operations.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(s.startCmd())
	return cmd
}

func (s *ServerCmd) startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start server.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			s := server.NewServer()
			s.Start()
		},
	}
	return cmd
}
