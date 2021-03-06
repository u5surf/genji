package shell

import (
	"fmt"

	"github.com/agnivade/levenshtein"
	"github.com/genjidb/genji"
	"github.com/genjidb/genji/document"
)

var commands = []struct {
	Name        string
	Options     string
	Description string
}{
	{".exit", ``, "Exit this program."},
	{".help", ``, "List all commands."},
	{".tables", ``, "List names of tables."},
	{".indexes", `[table_name]`, "Display all indexes or the indexes of the given table name."},
}

// runTablesCmd shows all tables.
func runTablesCmd(db *genji.DB, cmd []string) error {
	if len(cmd) > 1 {
		return fmt.Errorf("usage: .tables")
	}

	res, err := db.Query("SELECT table_name FROM __genji_tables")
	if err != nil {
		return err
	}
	defer res.Close()

	return res.Iterate(func(d document.Document) error {
		var tableName string
		err = document.Scan(d, &tableName)
		if err != nil {
			return err
		}
		fmt.Println(tableName)
		return nil
	})
}

// displayTableIndex prints all indexes that the given table contains.
func displayTableIndex(db *genji.DB, tableName string) error {
	err := db.View(func(tx *genji.Tx) error {
		t, err := tx.GetTable(tableName)
		if err != nil {
			return err
		}

		indexes, err := t.Indexes()
		if err != nil {
			return err
		}

		for _, idx := range indexes {
			fmt.Printf("%s on %s(%s)\n", idx.Opts.IndexName, idx.Opts.TableName, idx.Opts.Path)
		}

		return nil
	})

	return err
}

// displayAllIndexes shows all indexes that the database contains.
func displayAllIndexes(db *genji.DB) error {
	err := db.View(func(tx *genji.Tx) error {
		indexes, err := tx.ListIndexes()
		if err != nil {
			return err
		}

		for _, idx := range indexes {
			fmt.Printf("%s on %s(%s)\n", idx.IndexName, idx.TableName, idx.Path)
		}

		return nil
	})

	return err
}

// runIndexesCmd executes all indexes of the database or all indexes of the given table.
func runIndexesCmd(db *genji.DB, in []string) error {
	switch len(in) {
	case 1:
		// If the input is ".indexes"
		return displayAllIndexes(db)
	case 2:
		// If the input is ".indexes <tableName>"
		return displayTableIndex(db, in[1])
	}

	return fmt.Errorf("usage: .indexes [tablename]")
}

// runHelpCmd shows all available commands.
func runHelpCmd() error {
	for _, c := range commands {
		// spaces indentation for readability.
		spaces := 25
		indent := spaces - len(c.Name) - len(c.Options)
		fmt.Printf("%s %s %*s %s\n", c.Name, c.Options, indent, "", c.Description)
	}

	return nil
}

// displaySuggestions shows suggestions.
func displaySuggestions(in string) error {
	var suggestions []string
	for _, c := range commands {
		d := levenshtein.ComputeDistance(c.Name, in)
		// input should be at least half the command size to get a suggestion.
		if d < (len(c.Name) / 2) {
			suggestions = append(suggestions, c.Name)
		}
	}

	if len(suggestions) == 0 {
		return fmt.Errorf("Unknown command %q. Enter \".help\" for help.", in)
	}

	fmt.Printf("\"%s\" is not a command. Did you mean: ", in)
	for i := range suggestions {
		if i > 0 {
			fmt.Printf(", ")
		}

		fmt.Printf("%q", suggestions[i])
	}

	fmt.Println()
	return nil
}
