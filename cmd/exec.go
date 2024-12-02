/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/snaka/ecs-exec-sh/ecs"
	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ecs.ExecuteCommand(clusterName, serviceName, containerName, command); err != nil {
			fmt.Println("Failed to execute command:", err)
			os.Exit(1)
		}
	},
}

var clusterName, serviceName, containerName, command string

func init() {
	rootCmd.AddCommand(execCmd)

	execCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "cluster name")
	execCmd.Flags().StringVarP(&serviceName, "service", "s", "", "service name")
	execCmd.Flags().StringVarP(&containerName, "container", "C", "", "container name")
	execCmd.Flags().StringVarP(&command, "command", "x", "/bin/sh", "command to execute")

	execCmd.MarkFlagRequired("cluster")
	execCmd.MarkFlagRequired("service")
	execCmd.MarkFlagRequired("container")
}
