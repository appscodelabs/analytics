package cmds

import (
	"github.com/appscode/analytics/pkg/analytics"
	"github.com/appscode/analytics/pkg/dockerhub"
	"github.com/appscode/analytics/pkg/server"
	"github.com/appscode/analytics/pkg/spreadsheet"
	"github.com/appscode/log"
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
)

func NewCmdServer(version string) *cobra.Command {
	srv := hostfacts.Server{
		WebAddress:      ":9844",
		OpsAddress:      ":56790",
		DockerHubOrgs:   map[string]string{},
		EnableAnalytics: true,
	}
	cmd := &cobra.Command{
		Use:               "run",
		Short:             "Run server",
		DisableAutoGenTag: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			if srv.EnableAnalytics {
				analytics.Enable()
			}
			analytics.SendEvent("analytics", "started", version)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			analytics.SendEvent("analytics", "stopped", version)
		},
		Run: func(cmd *cobra.Command, args []string) {
			spreadsheet.Client_secret_path = cmd.Flag("client-secret-dir").Value.String()
			if err := dockerhub.CollectAnalytics(srv.DockerHubOrgs); err != nil {
				log.Fatalln(err)
			}

			c := cron.New()
			c.AddFunc("@every 4h", func() {
				if err := dockerhub.CollectAnalytics(srv.DockerHubOrgs); err != nil {
					log.Errorln(err)
				}
			})
			c.Start()

			srv.ListenAndServe()
		},
	}

	cmd.Flags().StringVar(&srv.WebAddress, "web-address", srv.WebAddress, "Http server address")
	cmd.Flags().StringVar(&srv.CACertFile, "cacert-file", srv.CACertFile, "File containing CA certificate")
	cmd.Flags().StringVar(&srv.CertFile, "cert-file", srv.CertFile, "File container server TLS certificate")
	cmd.Flags().StringVar(&srv.KeyFile, "key-file", srv.KeyFile, "File containing server TLS private key")

	cmd.Flags().StringToStringVar(&srv.DockerHubOrgs, "docker-hub-orgs", srv.DockerHubOrgs, "Map of Docker Hub organizations to Google spreadsheets")
	cmd.Flags().String("client-secret-dir", "/tmp/secrets/credentials", "Directory used to store client secrets and access tokens")

	cmd.Flags().StringVar(&srv.OpsAddress, "ops-addr", srv.OpsAddress, "Address to listen on for web interface and telemetry.")
	cmd.Flags().BoolVar(&srv.EnableAnalytics, "analytics", srv.EnableAnalytics, "Send analytical events to Google Analytics")
	return cmd
}
