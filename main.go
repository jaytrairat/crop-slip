package main

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"image/jpeg"
	"image/png"

	"github.com/fatih/color"
	"github.com/jaytrairat/crop-slip/constant"
	"github.com/spf13/cobra"
)

var (
	sourceFolder string
	configType   string
)

func cropImage(inputPath, outputPath string) bool {
	configForCropping := constant.SLIP_TYPE[configType]
	if configForCropping != nil {
		inputFile, err := os.Open(inputPath)
		if err != nil {
			return false
		}
		defer inputFile.Close()

		img, _, err := image.Decode(inputFile)
		if err != nil {
			return false
		}

		outputFile, err := os.Create(outputPath)
		if err != nil {
			return false
		}
		defer outputFile.Close()

		croppedImg := image.NewRGBA(image.Rect(0, 0, configForCropping["WIDTH"].(int), configForCropping["HEIGHT"].(int)))
		draw.Draw(croppedImg, croppedImg.Bounds(), img, image.Pt(configForCropping["X"].(int), configForCropping["Y"].(int)), draw.Src)

		ext := strings.ToLower(filepath.Ext(outputPath))
		if ext == ".png" {
			png.Encode(outputFile, croppedImg)
		} else if ext == ".jpg" || ext == ".jpeg" {
			jpeg.Encode(outputFile, croppedImg, nil)
		} else {
			return false
		}
		strings.ToLower(filepath.Ext(outputPath))
		return true
	} else {
		return false
	}
}

var rootCmd = &cobra.Command{
	Use:   "crop-slip",
	Short: "Crop account from slip",
	Run: func(cmd *cobra.Command, args []string) {
		red := color.New(color.FgRed)
		green := color.New(color.FgGreen)
		blue := color.New(color.FgBlue)
		green.Printf("%s :: start cropping\n", time.Now().Format("2006-01-02 03:04:05"))
		numberOfFile := 0
		filepath.Walk(sourceFolder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				ext := strings.ToLower(filepath.Ext(path))

				if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
					outputFolderName := fmt.Sprintf("./%s_extracted", regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAllString(sourceFolder, ""))
					os.MkdirAll(outputFolderName, os.ModePerm)
					outputPath := filepath.Join(outputFolderName, filepath.Base(path))
					completed := cropImage(path, outputPath)

					if !completed {
						red.Printf("%s :: Error cropping image at %s\n", time.Now().Format("2006-01-02 03:04:05"), path)
					} else {
						blue.Printf("%s :: Image cropped image at %s\n", time.Now().Format("2006-01-02 03:04:05"), path)
						numberOfFile++
					}
				}
			}
			return nil
		})
		green.Printf("%s :: %d image(s) cropped\n", time.Now().Format("2006-01-02 03:04:05"), numberOfFile)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	var bankList []string
	for key := range constant.SLIP_TYPE {
		bankList = append(bankList, key)
	}
	rootCmd.Flags().StringVarP(&sourceFolder, "source", "s", "", "Source folder containing images")
	rootCmd.Flags().StringVarP(&configType, "type", "t", "", fmt.Sprintf("Type of account :: %s", strings.Join(bankList, ", ")))
	rootCmd.MarkFlagRequired("source")
	rootCmd.MarkFlagRequired("type")

	Execute()
}
