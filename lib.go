package lib

import (
	"encoding/csv"
	"fmt"
	"github.com/signintech/gopdf"
)

//func make_fromCSV(csv, path string) error {
//
//}

func make_fromCSV(name, address, path string) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 283.5, H: 419.5}})
	pdf.AddPage()

	err := pdf.AddTTFFont("test", "/usr/share/fonts/truetype/fonts-japanese-gothic.ttf")
	if err != nil {
		return error
	}

	pdf.setFont("test", 16)
	pdf.SetY(34.0157)
	x := []float64{92.3, 99.3, 106.3, 113.9, 120.7, 127.5, 134.3}

	for i := 0; i < 7; i++ {
		pdf.SetX(x[i])
		pdf.Cell(nil, address[i])
	}

	pdf.WritePdf(path)

	return nil
}
