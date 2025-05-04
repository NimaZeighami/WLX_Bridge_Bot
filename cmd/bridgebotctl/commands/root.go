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
