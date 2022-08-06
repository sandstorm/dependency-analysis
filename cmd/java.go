/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/sandstorm/dependency-analysis/analyser"

	"github.com/spf13/cobra"
)

// javaCmd represents the java command
var javaCmd = &cobra.Command{
	Use:   "java",
	Short: "Analyses a local Java source directory",
	Long: `Analyses a local Java source directory.

We only consider packages within the given directory. It can be the project
root directory or a package within the project to "zoom in".
`,
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := "."
		if len(args) > 0 {
			sourcePath = args[0]
		}  
		analyser.Analyse(sourcePath)
	},
}

func init() {
	rootCmd.AddCommand(javaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// javaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// javaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
