package main

import (
	"github.com/urfave/cli/v2"
	"github.com/ytake/kfchc/command"
	"github.com/ytake/kfchc/config"
	"github.com/ytake/kfchc/log"
	"os"
)

func main() {
	l := log.NewLogger()
	defer l.Provider.Sync()
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "connectors:health_check",
				Aliases: []string{"c:conns"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     config.FlagJsonConfigPath,
						Usage:    "specify the Kafka Connect config file",
						Required: true,
					},
				},
				Usage: "Kafka Connect connectors and clients using the REST interface.",
				Action: func(context *cli.Context) error {
					ch := &command.HealthCheckHandle{
						Logger: l}
					return ch.Run(context)
				},
			},
			{
				Name:    "connectors:gen_server_config",
				Aliases: []string{"c:gsc"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     config.FlagOutputPath,
						Usage:    "specify the output path of config file",
						Required: true,
					},
				},
				Usage: "generating Kafka Connect connectors config.",
				Action: func(context *cli.Context) error {
					ch := &command.GenerateServerConfig{}
					return ch.Run(context)
				},
			},
		}}
	app.Name = `kfchc`
	err := app.Run(os.Args)
	if err != nil {
		l.RuntimeFatalError("kfchc command error", err)
	}
}
