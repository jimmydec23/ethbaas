package cmd

import (
	"ethbaas/internal/db"
	"ethbaas/pkg/chainclient"
	"ethbaas/pkg/projclient"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

type ChainCmd struct {
	argsName string
	projCli  *projclient.Client
	chainCli *chainclient.Client
}

func NewChainCmd(db *db.Client) *ChainCmd {
	e := &ChainCmd{
		projCli:  projclient.NewClient(db),
		chainCli: chainclient.NewClient(db),
	}
	return e
}

func (e *ChainCmd) rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain",
		Short: "Chain Operations.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(e.infoCmd())
	cmd.AddCommand(e.podsCmd())
	cmd.AddCommand(e.clusterCmd())
	return cmd
}

func (e *ChainCmd) infoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "chain info",
		Run: func(cmd *cobra.Command, args []string) {
			infoes, err := e.chainCli.Info(e.argsName)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("NetId\tCurrent\tHighest\tSyncing\tPeers\tDiff\tEnode")
			for _, info := range infoes {
				fmt.Printf(
					"%s\t%d\t%d\t%v\t%d\t%s\t%s\n",
					info.NetworkID, info.Current, info.Highest, info.Syncing,
					info.PeerCount, info.Diff, info.ENode,
				)
			}
		},
	}
	cmd.Flags().StringVarP(&e.argsName, "name", "n", "", "set project name")
	cmd.MarkFlagRequired("name")
	return cmd
}

func (e *ChainCmd) podsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pods",
		Short: "List chain pods",
		Run: func(cmd *cobra.Command, args []string) {
			proj, err := e.projCli.GetInModel(e.argsName)
			if err != nil {
				log.Fatal(err)
			}

			list, err := e.chainCli.Pods(proj)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Name\t\t\tStatus\tPorts\t\t\tService\tNodePorts")
			for _, p := range list {
				fmt.Printf("%s\t%s\t%v\t%s\t%v\n", p.Name, p.Status, p.Ports,
					p.Service, p.NodePorts,
				)
			}
		},
	}
	cmd.Flags().StringVarP(&e.argsName, "name", "n", "", "set project name")
	cmd.MarkFlagRequired("name")
	return cmd
}

func (e *ChainCmd) clusterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Connet nodes into a cluster.",
		Run: func(cmd *cobra.Command, args []string) {
			proj, err := e.projCli.GetInModel(e.argsName)
			if err != nil {
				log.Fatal(err)
			}

			if err := e.chainCli.Cluster(proj); err != nil {
				log.Fatal(err)
			}
			fmt.Println("Nodes has connected to a cluster.")
		},
	}
	cmd.Flags().StringVarP(&e.argsName, "name", "n", "", "set project name")
	cmd.MarkFlagRequired("name")
	return cmd
}
