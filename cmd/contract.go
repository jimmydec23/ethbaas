package cmd

import (
	"ethbaas/internal/db"
	"ethbaas/internal/model"
	"ethbaas/pkg/contractclient"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/spf13/cobra"
)

type ContractCmd struct {
	argsName    string
	argsProj    string
	argsKey     string
	argsValue   string
	argsMethod  string
	argsInputs  string
	argsABI     string
	argsBIN     string
	contractCli *contractclient.Client
}

func NewContractCmd(dbClient *db.Client) *ContractCmd {
	c := &ContractCmd{
		contractCli: contractclient.NewClient(dbClient),
	}
	return c
}

func (c *ContractCmd) rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract",
		Short: "Contract Operation.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(c.listCmd())
	cmd.AddCommand(c.deployCmd())
	cmd.AddCommand(c.queryCmd())
	cmd.AddCommand(c.writeCmd())
	return cmd
}

func (c *ContractCmd) listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List contract",
		Run: func(cmd *cobra.Command, args []string) {
			_, list, err := c.contractCli.List()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Name\tProject\tAddress\t\t\t\t\t\tCreated")
			for _, item := range list {
				t := time.Unix(item.Created, 0)
				tf := t.Format(time.RFC3339)
				fmt.Printf("%s\t%s\t%s\t%s\n", item.Name, item.Proj, item.Address, tf)
			}
		},
	}
	return cmd
}

func (c *ContractCmd) deployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy contract",
		Run: func(cmd *cobra.Command, args []string) {
			abiB, err := ioutil.ReadFile(c.argsABI)
			if err != nil {
				log.Fatal(err)
			}
			binB, err := ioutil.ReadFile(c.argsBIN)
			if err != nil {
				log.Fatal(err)
			}

			cont := &model.Contract{
				Name: c.argsName,
				ABI:  string(abiB),
				BIN:  string(binB),
			}
			addr, err := c.contractCli.Deploy(c.argsProj, cont)
			if err != nil {
				log.Fatal(err)
			}
			_ = addr
		},
	}
	cmd.Flags().StringVarP(&c.argsName, "name", "n", "", "contract name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&c.argsProj, "proj", "p", "", "project name")
	cmd.MarkFlagRequired("proj")
	cmd.Flags().StringVarP(&c.argsABI, "abi", "a", "", "contract abi path ")
	cmd.MarkFlagRequired("abi")
	cmd.Flags().StringVarP(&c.argsBIN, "bin", "b", "", "contract bin path ")
	cmd.MarkFlagRequired("bin")
	return cmd
}

func (c *ContractCmd) queryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "query contract",
		Run: func(cmd *cobra.Command, args []string) {
			res, err := c.contractCli.Query(c.argsName, c.argsMethod, c.argsInputs)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Call method %s return: %s\n", c.argsMethod, res)
		},
	}
	cmd.Flags().StringVarP(&c.argsName, "name", "n", "", "contract name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&c.argsMethod, "method", "m", "", "contract method")
	cmd.MarkFlagRequired("method")
	cmd.Flags().StringVarP(&c.argsInputs, "input", "i", "", "contract method inputs: a,b")
	return cmd
}

func (c *ContractCmd) writeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "write",
		Short: "Write to contarct",
		Run: func(cmd *cobra.Command, args []string) {
			tx, err := c.contractCli.Write(c.argsName, c.argsMethod, c.argsInputs)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Tx sent: %s\n", tx)
		},
	}
	cmd.Flags().StringVarP(&c.argsName, "name", "n", "", "contract name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&c.argsMethod, "method", "m", "", "contract method")
	cmd.MarkFlagRequired("method")
	cmd.Flags().StringVarP(&c.argsInputs, "input", "i", "", "contract method inputs: a,b")
	return cmd
}
