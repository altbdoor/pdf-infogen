package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/signintech/gopdf"
)

// https://transform.tools/json-to-go
type CoordsJson []struct {
	Field    string    `json:"field"`
	Value    string    `json:"value"`
	Position []float64 `json:"position"`
	Block    bool      `json:"block"`
	CellSize float64   `json:"cellSize"`
	FontSize int       `json:"fontSize"`
}

func getCoords(path string, coordsData *CoordsJson) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println("coords.json file not found")
		os.Exit(1)
	}

	coordsFile, _ := os.ReadFile(path)
	json.Unmarshal(coordsFile, &coordsData)
}

func checkPdfPath(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Printf("%s file not found\n", path)
		os.Exit(1)
	}
}

//go:embed fonts/DejaVuSans.ttf
var dejavuSansFontData []byte

func main() {
	args := os.Args[1:]
	path, _ := os.Getwd()

	fmt.Printf("Running from %s\n", path)

	var coordsData CoordsJson
	getCoords(filepath.Join(path, "coords.json"), &coordsData)
	fmt.Println("coords.json file found")

	if len(args) != 1 {
		fmt.Println("Please provide the path to the input PDF file")
		os.Exit(1)
	}

	pdfPath := args[0]
	checkPdfPath(pdfPath)

	pdf := gopdf.GoPdf{}
	pdfConfig := gopdf.Config{PageSize: *gopdf.PageSizeA4}
	pdf.Start(pdfConfig)

	pdf.AddTTFFontData("dejavusans", dejavuSansFontData)

	pdf.AddPage()
	tpl := pdf.ImportPage(pdfPath, 1, "/MediaBox")
	pdf.UseImportedTemplate(tpl, 0, 0, pdfConfig.PageSize.W, pdfConfig.PageSize.H)

	_, isDebug := os.LookupEnv("DEBUG")
	if isDebug {
		pdf.SetLineWidth(0.5)
		pdf.SetStrokeColor(255, 0, 0)
	}

	for _, data := range coordsData {
		if data.Value == "" {
			continue
		}

		fmt.Printf("Filling in %s...\n", data.Field)
		pdf.SetXY(data.Position[0], data.Position[1])
		pdf.SetFont("dejavusans", "", data.FontSize)

		if data.Block {
			cellOptionBlock := gopdf.CellOption{Align: gopdf.Center | gopdf.Middle, CoefLineHeight: data.CellSize}
			cellRectBlock := gopdf.Rect{W: data.CellSize, H: data.CellSize}

			if isDebug {
				cellOptionBlock.Border = gopdf.AllBorders
			}

			valueSplit := strings.Split(data.Value, "")

			for _, char := range valueSplit {
				pdf.CellWithOption(
					&cellRectBlock,
					char,
					cellOptionBlock,
				)
			}
		} else {
			cellOptionInline := gopdf.CellOption{Align: gopdf.Left | gopdf.Middle, CoefLineHeight: data.CellSize}
			cellRectInline := gopdf.Rect{W: 1, H: data.CellSize}

			if isDebug {
				cellOptionInline.Border = gopdf.AllBorders
			}

			pdf.CellWithOption(&cellRectInline, data.Value, cellOptionInline)
		}
	}

	pdf.WritePdf("output.pdf")
	fmt.Println("Done")
}
