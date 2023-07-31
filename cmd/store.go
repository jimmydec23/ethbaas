package cmd

import (
	"ethbaas/internal/db"
	"ethbaas/pkg/contractclient"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

type StoreCmd struct {
	storeCli  *contractclient.StoreClient
	argsName  string
	argsProj  string
	argsPath  string
	argsKey   string
	argsValue string
}

func NewStoreCmd(db *db.Client) *StoreCmd {
	e := &StoreCmd{
		storeCli: contractclient.NewStoreClient(db),
		argsName: "store",
	}
	return e
}

func (s *StoreCmd) rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store",
		Short: "Store contract operations.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(s.deployCmd())
	cmd.AddCommand(s.queryCmd())
	cmd.AddCommand(s.writeCmd())
	return cmd
}

func (s *StoreCmd) deployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy store contract",
		Run: func(cmd *cobra.Command, args []string) {
			if err := s.storeCli.Deploy(s.argsProj, s.argsName, s.argsPath); err != nil {
				log.Fatal(err)
			}
			log.Printf("Contract %s deployed.\n", s.argsName)
		},
	}
	cmd.Flags().StringVarP(&s.argsProj, "proj", "p", "", "project name")
	cmd.MarkFlagRequired("proj")

	cmd.Flags().StringVarP(&s.argsPath, "path", "", "", "contract path ")
	cmd.MarkFlagRequired("path")
	return cmd
}

func (s *StoreCmd) queryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Query store contract",
		Run: func(cmd *cobra.Command, args []string) {
			value, err := s.storeCli.Query(s.argsName, s.argsKey)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Query contract result:", value)
		},
	}
	cmd.Flags().StringVarP(&s.argsKey, "key", "k", "", "contract key")
	cmd.MarkFlagRequired("key")
	return cmd
}

func (s *StoreCmd) writeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "write",
		Short: "write to store contract",
		Run: func(cmd *cobra.Command, args []string) {
			tx, err := s.storeCli.Write(s.argsName, s.argsKey, s.argsValue)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Contract writed, tx:", tx)
		},
	}
	cmd.Flags().StringVarP(&s.argsKey, "key", "k", "", "contract key")
	cmd.MarkFlagRequired("key")
	cmd.Flags().StringVarP(&s.argsValue, "value", "v", "", "contract value")
	cmd.MarkFlagRequired("value")
	return cmd
}
