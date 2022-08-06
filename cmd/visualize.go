/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/sandstorm/dependency-analysis/analyser"

	"github.com/spf13/cobra"
)

// visualizeCmd represents the visualize command
var visualizeCmd = &cobra.Command{
	Use:   "visualize",
	Short: "Renders the dependencies into an image.",
	Long: `TODO … feature … output … limitations … tipps`,
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := "."
		if len(args) > 0 {
			sourcePath = args[0]
		}  
		analyser.Analyse(sourcePath)
	},
}

func init() {
	javaCmd.AddCommand(visualizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// visualizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// visualizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
