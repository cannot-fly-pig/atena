/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// makeCsvCmd represents the makeCsv command
var makeCsvCmd = &cobra.Command{
	Use:   "make-csv [string to file name]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Args: cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		if path[len(path)-4:] != ".csv" {
			path += ".csv"
		}

		if Exists(args[0]) {
			fmt.Println(args[0] + " is already exist")
		} else {
			file, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
			defer file.Close()

			writer := csv.NewWriter(file)
			writer.Write([]string{"住所", "郵便番号(任意)", "名前"})
			writer.Flush()
		}
	},
}

func init() {
	rootCmd.AddCommand(makeCsvCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// makeCsvCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// makeCsvCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
