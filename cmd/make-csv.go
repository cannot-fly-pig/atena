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

var make_csvpath string

// makeCsvCmd represents the makeCsv command
var makeCsvCmd = &cobra.Command{
	Use:   "make-csv [path to csv]",
	Short: "atena make のためのcsvファイルのテンプレートを作成",
	Long:  "atena make のためのcsvファイルのテンプレートを作成",

	Run: func(cmd *cobra.Command, args []string) {
		if len(make_csvpath) < 5 || make_csvpath[len(make_csvpath)-4:] != ".csv" {
			make_csvpath += ".csv"
		}

		if Exists(make_csvpath) {
			fmt.Println(make_csvpath + " is already exist")
		} else {
			file, _ := os.OpenFile(make_csvpath, os.O_WRONLY|os.O_CREATE, 0666)
			defer file.Close()

			writer := csv.NewWriter(file)
			writer.Write([]string{"住所1", "住所2", "名前", "郵便番号"})
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
	makeCsvCmd.PersistentFlags().StringVar(&make_csvpath, "output", "./address-list.csv", "出力するcsvファイルのパス")
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
