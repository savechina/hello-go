package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	// Create a new Cobra command.
	rootCmd := &cobra.Command{
		// Set the command's name.
		Use: "foo",

		// Set the command's description.
		Long: "A simple command that does nothing.",

		// Add a help command.
		// helpFunc: func(*Command, []string) {
		// 	fmt.Println("This is a help message.")
		// },
	}

	// Add a new command to the root command.
	addCmd := &cobra.Command{
		// Set the command's name.
		Use: "add",

		// Set the command's description.
		Long: "Adds two numbers together.",

		// Define the command's flags.
		// flags: []flag.Flag{
		// 	&flag.Flag{
		// 		// Set the flag's name.
		// 		Name: "first",

		// 		// Set the flag's description.
		// 		Usage: "The first number to add.",
		// 	},
		// 	&flag.Flag{
		// 		// Set the flag's name.
		// 		Name: "second",

		// 		// Set the flag's description.
		// 		Usage: "The second number to add.",
		// 	},
		// },

		// Define the command's action.
		Run: func(cmd *cobra.Command, args []string) {
			// Get the values of the flags.
			first, err := cmd.Flags().GetInt("first")
			if err != nil {
				fmt.Println("Error getting first number:", err)
				return
			}

			second, err := cmd.Flags().GetInt("second")
			if err != nil {
				fmt.Println("Error getting second number:", err)
				return
			}

			// Add the two numbers together.
			sum := first + second

			// Print the sum.
			fmt.Println("The sum is:", sum)
		},
	}

	// Add the add command to the root command.
	rootCmd.AddCommand(addCmd)

	// Start the command.
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		return
	}
}
