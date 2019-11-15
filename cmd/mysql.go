package cmd

import (
	"fmt"

	"github.com/novacloudcz/goclitools"

	"github.com/novacloudcz/devopscli/tools"
	"github.com/urfave/cli"
)

// MySQLCmd ...
func MySQLCmd() cli.Command {
	return cli.Command{
		Name: "mysql",
		Subcommands: []cli.Command{
			MySQLDatabase(),
		},
	}
}

// MySQLDatabase ...
func MySQLDatabase() cli.Command {
	return cli.Command{
		Name: "database",
		Subcommands: []cli.Command{
			MySQLDump(),
			MySQLImport(),
		},
	}
}

// MySQLDump ...
func MySQLDump() cli.Command {
	return cli.Command{
		Name: "dump",
		Action: func(c *cli.Context) error {

			db := tools.NewDBWithEnv()

			username := db.Username
			password := db.Password
			hostname := db.Hostname

			dbName := c.Args().First()
			if dbName == "" {
				dbName = db.Database
			}

			cmd := fmt.Sprintf("mysqldump -u %s -p%s -h %s %s", username, password, hostname, dbName)
			if err := goclitools.RunInteractive(cmd); err != nil {
				return cli.NewExitError(err, 1)
			}

			return nil
		},
	}
}

// MySQLImport ...
func MySQLImport() cli.Command {
	return cli.Command{
		Name: "import",
		Action: func(c *cli.Context) error {

			db := tools.NewDBWithEnv()

			username := db.Username
			password := db.Password
			hostname := db.Hostname

			dbName := c.Args().First()
			if dbName == "" {
				dbName = db.Database
			}

			cmd := fmt.Sprintf("mysql -u %s -p%s -h %s %s", username, password, hostname, dbName)
			if err := goclitools.RunInteractive(cmd); err != nil {
				return cli.NewExitError(err, 1)
			}

			return nil
		},
	}
}
