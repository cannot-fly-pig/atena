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
	"github.com/cheggaaa/pb/v3"
	m "github.com/ktnyt/go-moji"
	"github.com/signintech/gopdf"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var name, address0, address1, csv_path, path, code, font string

// makeCmd represents the make command
var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "はがきの宛名を作成",
	Long: `オプションで与えられた住所情報または
	csvファイルからはがきの宛名を作成します`,
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case name == "" && address0 == "" && csv_path == "":
			cmd.Help()

		case (name == "" || address0 == "") && csv_path == "":
			cmd.Help()

		case (name != "" && address0 != "") && csv_path == "":
			make_fromName(name, address0, address1, code, path)

		case (name == "" && address0 == "") && csv_path != "":
			make_fromcsv(csv_path, path)

		}
	},
}

func init() {
	rootCmd.AddCommand(makeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// makeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// makeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	makeCmd.PersistentFlags().StringVar(&name, "name", "", "宛先の名前")
	makeCmd.PersistentFlags().StringVar(&address0, "address1", "", "宛先の住所の前半")
	makeCmd.PersistentFlags().StringVar(&address1, "address2", "", "宛先の住所の後半")
	makeCmd.PersistentFlags().StringVar(&code, "code", "", "宛先の郵便番号")
	makeCmd.PersistentFlags().StringVar(&csv_path, "csv", "", "宛先のリストのcsvファイルのパス")
	makeCmd.PersistentFlags().StringVar(&font, "font", "", "宛先のフォントのpath")
	makeCmd.PersistentFlags().StringVar(&path, "output", "", "出力するpdfファイルのパス")
}

func mm2pt(n float64) float64 {
	return n / 1000 * 2835
}

func moji(s string) []string {
	s = m.Convert(s, m.ZE, m.HE)
	s = m.Convert(s, m.ZS, m.HS)
	list := strings.Split(s, "")
	for i := 0; i < len(list); i++ {
		switch {
		case list[i] == "0":
			list[i] = "〇"

		case list[i] == "1":
			list[i] = "一"

		case list[i] == "2":
			list[i] = "二"

		case list[i] == "3":
			list[i] = "三"

		case list[i] == "4":
			list[i] = "四"

		case list[i] == "5":
			list[i] = "五"

		case list[i] == "6":
			list[i] = "六"

		case list[i] == "7":
			list[i] = "七"

		case list[i] == "8":
			list[i] = "八"

		case list[i] == "9":
			list[i] = "九"

		case list[i] == "-" || list[i] == "−" || list[i] == "ー":
			list[i] = "丨"

		}
	}
	return list
}

func make_fromName(name, address0, address1, code, path string) {

	name_size := 32
	address_size := 14

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: mm2pt(100), H: mm2pt(148)}})
	pdf.AddPage()

	_ = pdf.AddTTFFont("test", font)
	//	if err != nil {
	//		return err
	//	}

	// 郵便番号api
	if len(strings.Split(code, "")) != 7 {
		request, _ := http.NewRequest("GET", "https://zipcoda.net/api/", nil)
		values := url.Values{}
		values.Add("address", address0+address1)
		request.URL.RawQuery = values.Encode()

		client := new(http.Client)
		resp, _ := client.Do(request)
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		txt := string(body)
		index := strings.Index(txt, "zipcode")
		code = txt[index+10 : index+17]
	}

	//郵便番号印刷

	pdf.SetFont("test", "", 12)
	pdf.SetY(mm2pt(13))
	x := []float64{44.6, 52.3, 60.0, 68.3, 75.8, 83.3, 90.8}

	for i := 0; i < 7; i++ {
		pdf.SetX(mm2pt(x[i]))
		pdf.Cell(nil, strings.Split(code, "")[i])
	}

	//名前印刷

	pdf.SetFont("test", "", name_size)
	y := mm2pt(35)
	name = m.Convert(name, m.ZS, m.HS)
	name_list := strings.Split(name, " ")
	for i := 0; i < len(strings.Split(name_list[0], "")); i++ {
		pdf.SetX(141.75 - float64(name_size/2))
		pdf.SetY(y)
		pdf.Cell(nil, strings.Split(name_list[0], "")[i])
		y += 36
	}

	y += 18

	orig_y := y

	max_l := 0
	for i := 1; i < len(name_list); i++ {
		if max_l < len(strings.Split(name_list[i], "")) {
			max_l = len(strings.Split(name_list[i], ""))
		}
	}

	for i := 1; i < len(name_list); i++ {

		y = orig_y

		for j := 0; j < len(strings.Split(name_list[i], "")); j++ {
			pdf.SetX(141.75 - float64(name_size/2) - float64((i-1)*(name_size+2)))
			pdf.SetY(y)
			pdf.Cell(nil, strings.Split(name_list[i], "")[j])
			y += 36
		}

		pdf.SetX(141.75 - float64(name_size/2) - float64((i-1)*(name_size+2)))
		pdf.SetY(orig_y + float64(max_l*36) + 18)
		pdf.Cell(nil, "様")

	}

	//住所印刷

	pdf.SetFont("test", "", address_size)
	if address1 != "" {
		y = mm2pt(30)

		address0_list := moji(address0)
		for i := 0; i < len(address0_list); i++ {
			pdf.SetX(mm2pt(90))
			pdf.SetY(y)
			pdf.Cell(nil, address0_list[i])
			y += float64(address_size + 2)
		}

		y = mm2pt(30) + float64(address_size+2)*float64(16-len(strings.Split(address1, "")))

		address1_list := moji(address1)
		for i := 0; i < len(address1_list); i++ {
			pdf.SetX(mm2pt(90) - float64(address_size+4))
			pdf.SetY(y)
			pdf.Cell(nil, address1_list[i])
			y += float64(address_size + 2)
		}

	} else {
		y = mm2pt(30)

		address0_list := moji(address0)
		for i := 0; i < len(address0_list); i++ {
			pdf.SetX(mm2pt(90) - float64(address_size/2))
			pdf.SetY(y)
			pdf.Cell(nil, address0_list[i])
			y += float64(address_size + 2)
		}
	}

	pdf.WritePdf(path)
	//	return nil
}

func make_fromcsv(csv_path, path string) {
	file, err := os.Open(csv_path)
	defer file.Close()

	reader := csv.NewReader(file)

	if path[len(path)-4:] == ".pdf" {
		path = path[:len(path)-4]
	}

	lines, err := reader.ReadAll()

	bar := pb.StartNew(len(lines) - 1)

	n := 0
	for _, line := range lines {
		if err != nil {
			break
		}
		if n != 0 {
			make_fromName(line[2], line[0], line[1], line[3], string(path)+strconv.Itoa(n)+".pdf")
			bar.Increment()
		}
		n++
	}
	bar.Finish()
}
