package cmd

import (
	"fmt"
	"os"

	"goAPI/cmd/book_keeper"

	"github.com/spf13/cobra"
)

var name string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		a := book_keeper.App{}
		a.Execute()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Define a new flag for the root command
	rootCmd.PersistentFlags().StringVar(&name, "name", "Charles", "A parameter that will be passed to the main function")

	// Other flag definitions can be added here as needed
}
