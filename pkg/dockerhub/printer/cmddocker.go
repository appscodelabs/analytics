package printer

import (
	"github.com/appscode/log"
	"github.com/spf13/cobra"
)

var DockerHubOrgs []string

func NewCmdDockerHub() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "dockerhub",
		Short:             "Shows Dockerhub logs",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if err := DockerPrinter(DockerHubOrgs); err != nil {
				log.Fatalln(err)
			}
		},
	}

	cmd.Flags().StringArrayVar(&DockerHubOrgs, "docker-hub-orgs", DockerHubOrgs, "Array of Docker Hub organizations")

	return cmd
}
