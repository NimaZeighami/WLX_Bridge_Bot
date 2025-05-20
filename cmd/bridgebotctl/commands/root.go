package commands

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "briegebotctl",
	Short: "Bridgebot CLI for managing database setup, migration, and seeding",
}

func Execute() error {
	return rootCmd.Execute()
}


// ! Warning
// ! consider that because of relative path, 
// ! these commands will not work if you run it from another directory, so you need to run it from the root of the project