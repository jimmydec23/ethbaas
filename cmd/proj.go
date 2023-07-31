package cmd

import (
	"ethbaas/internal/db"
	"ethbaas/internal/model"
	"ethbaas/pkg/projclient"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

type ProjCmd struct {
	argsName      string
	argsNodeCount int
	argsPort      int32
	pcli          *projclient.Client
}

func newProjCmd(db *db.Client) *ProjCmd {
	return &ProjCmd{
		pcli: projclient.NewClient(db),
	}
}

func (p *ProjCmd) rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proj",
		Short: "Project operations.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(p.initCmd())
	cmd.AddCommand(p.listCmd())
	cmd.AddCommand(p.startCmd())
	cmd.AddCommand(p.stopCmd())
	cmd.AddCommand(p.deleteCmd())
	return cmd
}

func (p *ProjCmd) initCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "init",
		Short: "Init a project",
		Run: func(cmd *cobra.Command, args []string) {
			proj := &model.Project{
				Name:          p.argsName,
				NodeCount:     p.argsNodeCount,
				FirstNodePort: p.argsPort,
			}
			if err := p.pcli.Init(proj); err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Project %s initialized.\n", proj.Name)
		},
	}
	cmd.Flags().StringVarP(&p.argsName, "name", "n", "", "set project name")
	cmd.MarkFlagRequired("name")

	cmd.Flags().IntVarP(&p.argsNodeCount, "nodeCount", "c", 1, "set node count")

	cmd.Flags().Int32VarP(&p.argsPort, "port", "p", 30545, "first node's nodePort")
	return cmd
}

func (p *ProjCmd) listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List projects",
		Run: func(cmd *cobra.Command, args []string) {
			total, list, err := p.pcli.List()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Project Amount:", total)
			if total != 0 {
				fmt.Println("Name\tNodes\tPorts\t\tRunning\tCreated")
				for _, item := range list {
					t := time.Unix(item.Created, 0)
					tf := t.Format(time.RFC3339)
					fmt.Printf(
						"%s\t%d\t%s\t%v\t%s\n",
						item.Name, item.NodeCount, item.NodePort, item.Running, tf,
					)
				}
			}
		},
	}
	return cmd
}

func (p *ProjCmd) startCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "start",
		Short: "Start a project",
		Run: func(cmd *cobra.Command, args []string) {
			if err := p.pcli.Start(p.argsName); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Project %s started.\n", p.argsName)
		},
	}
	cmd.Flags().StringVarP(&p.argsName, "name", "n", "", "set project name")
	cmd.MarkFlagRequired("name")
	return cmd
}

func (p *ProjCmd) stopCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop a project",
		Run: func(cmd *cobra.Command, args []string) {
			if err := p.pcli.Stop(p.argsName); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Project %s stoped.\n", p.argsName)
		},
	}
	cmd.Flags().StringVarP(&p.argsName, "name", "n", "", "set project name")
	cmd.MarkFlagRequired("name")
	return cmd
}

func (p *ProjCmd) deleteCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a project",
		Run: func(cmd *cobra.Command, args []string) {
			if err := p.pcli.Delete(p.argsName); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Project %s deleted.\n", p.argsName)
		},
	}
	cmd.Flags().StringVarP(&p.argsName, "name", "n", "", "set project name")
	cmd.MarkFlagRequired("name")
	return cmd
}
