package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

func main() {
	// map with the desired scales and dimensions
	resolutions := map[string]uint{"100": 24, "125": 30, "150": 36, "175": 42, "200": 48, "400": 96}

	log.SetPrefix("multires: ")
	log.SetFlags(0)
	flag.Parse()
	args := flag.Args()
	var sourceFolder string
	if len(args) == 0 {
		fmt.Println("Please enter the folder name containig the svg files")
		fmt.Scanln(&sourceFolder)
	} else {
		sourceFolder = args[0]
	}

	svgFileNames, err := getSvgFileNames(sourceFolder)
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	for folderName, size := range resolutions {
		os.Mkdir(filepath.Join(sourceFolder, folderName), 0755)
		for _, svg := range svgFileNames {
			var svgFullName string = filepath.Join(sourceFolder, svg)
			var pngFileName string = strings.Replace(svg, ".svg", ".png", 1)
			var pngFullName string = filepath.Join(sourceFolder, folderName, pngFileName)
			err = toPng(svgFullName, pngFullName, int(size), int(size))
			if err != nil {
				log.Print(err)
			} else {
				count++
			}
		}
	}
	fmt.Printf("Created %d png\n", count)
}

func getSvgFileNames(sourceFolder string) ([]string, error) {
	files, err := os.ReadDir(sourceFolder)
	if err != nil {
		return nil, err
	}

	var svgFileNames []string
	for _, file := range files {
		if !file.IsDir() && (filepath.Ext(file.Name()) == ".svg") {
			svgFileNames = append(svgFileNames, file.Name())
		}
	}

	if len(svgFileNames) == 0 {
		return nil, fmt.Errorf("no svg files found in %s", sourceFolder)
	}
	return svgFileNames, nil
}

func toPng(svgFileName string, pngFileName string, w int, h int) error {
	in, err := os.Open(svgFileName)
	if err != nil {
		return fmt.Errorf("%s: %s", err, svgFileName)
	}
	defer in.Close()
	icon, err := oksvg.ReadIconStream(in, oksvg.StrictErrorMode)
	if err != nil {
		return fmt.Errorf("%s: %s", err, svgFileName)
	}

	icon.SetTarget(0, 0, float64(w), float64(h))
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, img, img.Bounds())), 1.0)
	if img != nil {
		file, err := os.Create(pngFileName)
		if err != nil {
			return fmt.Errorf("%s: %s", err, pngFileName)
		}
		defer file.Close()

		err = png.Encode(file, img)
		if err != nil {
			return fmt.Errorf("%s: %s", err, file.Name())
		}
	}
	return nil
}
