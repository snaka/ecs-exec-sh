/*
Copyright © 2024 Shinji NAKAMATSU

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/snaka/ecs-exec-sh/ecs"
	"github.com/spf13/cobra"
)

var (
	clusterName   string
	serviceName   string
	containerName string
	command       string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ecs-exec-sh",
	Short: "Interactively select cluster, service, and container to execute shell commands.",
	Long: `ecs-exec-sh lets you easily run commands in your ECS containers.
Choose the cluster, service, and container interactively, without needing to know their names`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if err := ecs.ExecuteCommand(clusterName, serviceName, containerName, command); err != nil {
			fmt.Println("Failed to execute command:", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "cluster name")
	rootCmd.Flags().StringVarP(&serviceName, "service", "s", "", "service name")
	rootCmd.Flags().StringVarP(&containerName, "container", "C", "", "container name")
	rootCmd.Flags().StringVarP(&command, "command", "x", "/bin/sh", "command to execute")
}
