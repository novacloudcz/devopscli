package cmd

import (
	"errors"
	"fmt"

	"github.com/novacloudcz/goclitools"
	"github.com/urfave/cli"
)

// AWSCmd ...
func AWSCmd() cli.Command {
	return cli.Command{
		Name: "aws",
		Subcommands: []cli.Command{
			AWSLambdaCmd(),
		},
	}
}

// AWSLambdaCmd ...
func AWSLambdaCmd() cli.Command {
	return cli.Command{
		Name:  "lambda",
		Flags: []cli.Flag{},
		Subcommands: []cli.Command{
			AWSLambdaDeployCmd(),
		},
	}

}

// AWSLambdaDeployCmd ...
func AWSLambdaDeployCmd() cli.Command {
	return cli.Command{
		Name: "deploy",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name,n",
				Value: "",
				Usage: "Function name",
			},
		},
		Action: func(c *cli.Context) error {

			if err := deployLambda(c.String("name"), "."); err != nil {
				return cli.NewExitError(err, 1)
			}

			return nil
		},
	}

}

func deployLambda(name, dir string) error {
	if name == "" {
		return errors.New("--name must be specified")
	}

	if err := goclitools.RunInteractiveInDir("zip -r -q -9 lambda-function-archive.zip .", dir); err != nil {
		return err
	}

	defer goclitools.RunInteractiveInDir("rm lambda-function-archive.zip", dir)

	deploycmd := fmt.Sprintf("aws lambda update-function-code --function-name %s --zip-file fileb://`pwd`/lambda-function-archive.zip", name)
	if err := goclitools.RunInteractiveInDir(deploycmd, dir); err != nil {
		return err
	}

	return nil
}
