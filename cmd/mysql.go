package cmd

import (
	"fmt"
	"log"

	"github.com/inloop/goclitools"

	"github.com/novacloudcz/devopscli/tools"
	"github.com/urfave/cli"
)

func SQLCmd() cli.Command {
	return cli.Command{
		Name: "mysql",
		Subcommands: []cli.Command{
			SQLDatabase(),
		},
	}
}

// SQLDB ...
func SQLDatabase() cli.Command {
	return cli.Command{
		Name: "database",
		Subcommands: []cli.Command{
			SQLClear(),
			SQLDump(),
			SQLImport(),
		},
	}
}

// SQLClear ...
func SQLClear() cli.Command {
	return cli.Command{
		Name: "clear",
		Action: func(c *cli.Context) error {
			dbName := c.Args().First()
			if dbName == "" {
				return cli.NewExitError("missing database name as first argument", 1)
			}

			db := tools.NewDBWithEnv()

			tables, err := tableList(db, dbName)
			if err != nil {
				return cli.NewExitError(err, 1)
			}

			for _, table := range tables {
				err := tableDrop(db, table)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
			}

			return nil
		},
	}
}

// SQLDump ...
func SQLDump() cli.Command {
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

// SQLImport ...
func SQLImport() cli.Command {
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

func tableList(db *tools.DB, dbName string) ([]string, error) {
	tables := []string{}

	if dbName == "" {
		dbName = db.Database
	}

	query := fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = '%s';", dbName)

	rows, err := db.Client().Query(query)
	if err != nil {
		return tables, err
	}

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		tables = append(tables, name)
	}
	defer rows.Close()
	return tables, nil
}

func tableDrop(db *tools.DB, table string) error {
	query := fmt.Sprintf("DROP TABLE `%s`;", table)

	_, err := db.Client().Exec(query)

	return err
}
