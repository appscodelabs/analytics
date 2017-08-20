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
			output := cmd.Flag("output").Value.String()
			if err := DockerPrinter(DockerHubOrgs, output); err != nil {
				log.Fatalln(err)
			}
		},
	}

	cmd.Flags().StringSliceVar(&DockerHubOrgs, "docker-hub-orgs", DockerHubOrgs, "Array of Docker Hub organizations")
	cmd.Flags().String("output", "", "Directory used to store docker hub stats report")

	return cmd
}
