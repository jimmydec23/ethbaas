package cmd

import (
	"ethbaas/internal/db"
	"os"

	"github.com/spf13/cobra"
)

type Command struct {
	dbClient    *db.Client
	rootCmd     *cobra.Command
	projCmd     *ProjCmd
	chainCmd    *ChainCmd
	contractCmd *ContractCmd
	storeCmd    *StoreCmd
	serverCmd   *ServerCmd
}

func NewCommand(db *db.Client) *Command {
	c := &Command{
		dbClient: db,
	}
	c.setup()
	return c
}

func (c *Command) setup() {
	var rootCmd = &cobra.Command{
		Use:   "ethbaas",
		Short: "A generator of baas using ethereum.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	c.rootCmd = rootCmd

	projCmd := newProjCmd(c.dbClient)
	rootCmd.AddCommand(projCmd.rootCmd())
	c.projCmd = projCmd

	chainCmd := NewChainCmd(c.dbClient)
	rootCmd.AddCommand(chainCmd.rootCmd())
	c.chainCmd = chainCmd

	contractCmd := NewContractCmd(c.dbClient)
	rootCmd.AddCommand(contractCmd.rootCmd())
	c.contractCmd = contractCmd

	storeCmd := NewStoreCmd(c.dbClient)
	rootCmd.AddCommand(storeCmd.rootCmd())
	c.storeCmd = storeCmd

	serverCmd := NewServerCmd()
	rootCmd.AddCommand(serverCmd.rootCmd())
	c.serverCmd = serverCmd
}

func (c *Command) Execute() {
	if err := c.rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
