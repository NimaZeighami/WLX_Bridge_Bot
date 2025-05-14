//  The base/root command. It registers subcommands.
//  🧩 Other files like setup.go, migrate.go, and seed.go init() themselves here

package commands

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "briegebotctl",
	Short: "Bridgebot CLI for managing database setup, migration, and seeding",
}

func Execute() error {
	return rootCmd.Execute()
}


// ! consider that because of relative path, 
// ! the command will not work if you run it from another directory, so you need to run it from the root of the project