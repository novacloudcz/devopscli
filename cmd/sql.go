package cmd

import (
	"fmt"
	"log"

	"github.com/novacloudcz/devopscli/tools"
	"github.com/urfave/cli"
)

// SQLCmd ...
func SQLCmd() cli.Command {
	return cli.Command{
		Name: "sql",
		Subcommands: []cli.Command{
			SQLDatabase(),
		},
	}
}

// SQLDatabase ...
func SQLDatabase() cli.Command {
	return cli.Command{
		Name: "database",
		Subcommands: []cli.Command{
			SQLClearTables(),
		},
	}
}

// SQLClearTables ...
func SQLClearTables() cli.Command {
	return cli.Command{
		Name: "clear-tables",
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
