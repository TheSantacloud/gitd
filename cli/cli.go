package cli

import (
	"fmt"
	"github.com/dormunis/gitd/adapters"
	"github.com/dormunis/gitd/taskmanagers/taskmanager"
	"os"

	"github.com/spf13/cobra"
	// TODO: use viper as well
)

var settings adapters.Settings

var rootCmd = &cobra.Command{
	Use:   "github.com/dormunis/gitd",
	Short: "github.com/dormunis/gitd is a CLI for managing tasks",
	Long: `github.com/dormunis/gitd is a CLI for managing tasks.
    It is designed to work with task managers like Todoist, etc.
    It is also designed to work with archive managers like Obsidian, Notion, etc.`,
}

// TODO: create registration process for adapters

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review phase",
	Long:  `Review all tasks and purge old and irrelevant tasks from task manager`,
}

var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge tasks",
	Long:  `Purge old and irrelevant tasks from task manager`,
	Run: func(cmd *cobra.Command, args []string) {
		taskManager, err := taskmanager.Initialize(taskmanager.Todoist, settings)
		if err != nil {
			// TODO: handle errors more gracefully
			fmt.Println(err)
			os.Exit(1)
		}

		timespanString, err := cmd.Flags().GetString("timespan")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		timespan, err := adapters.NewTimeSpan(timespanString)
		if err != nil {
			fmt.Printf("Invalid timespan %s, error: %s", timespanString, err)
			os.Exit(1)
		}

		Purge(taskManager, *timespan)
	},
}

func init() {
	settings = adapters.GetSettings()
	rootCmd.AddCommand(reviewCmd)
	reviewCmd.AddCommand(purgeCmd)
	purgeCmd.PersistentFlags().String("timespan", "1 month", "timespan to review")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
