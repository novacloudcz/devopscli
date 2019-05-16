package main

import (
	"os"

	"github.com/novacloudcz/devopscli/cmd"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "DevOps cli"
	app.Usage = "..."
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		cmd.GitlabCmd(),
		cmd.DockerCmd(),
		cmd.AWSCmd(),
		cmd.SQLCmd(),
	}

	app.Run(os.Args)
}
