package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jaytrairat/extract-account/constant"
	"github.com/spf13/cobra"
)

var (
	sourceFolder string
	configType   string
)

func cropImage(inputPath, outputPath string) error {
	configForCropping := constant.SLIP_TYPE[configType]
	if configForCropping != nil {
		inputFile, err := os.Open(inputPath)
		if err != nil {
			return err
		}
		defer inputFile.Close()

		img, _, err := image.Decode(inputFile)
		if err != nil {
			return err
		}

		outputFile, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer outputFile.Close()

		croppedImg := image.NewRGBA(image.Rect(0, 0, configForCropping["WIDTH"].(int), configForCropping["HEIGHT"].(int)))
		draw.Draw(croppedImg, croppedImg.Bounds(), img, image.Pt(configForCropping["X"].(int), configForCropping["Y"].(int)), draw.Src)

		ext := strings.ToLower(filepath.Ext(outputPath))
		if ext == ".png" {
			return png.Encode(outputFile, croppedImg)
		} else if ext == ".jpg" || ext == ".jpeg" {
			return jpeg.Encode(outputFile, croppedImg, nil)
		} else {
			return fmt.Errorf("Unsupported image format: %s", ext)
		}
	} else {
		return nil
	}
}

var rootCmd = &cobra.Command{
	Use:   "extract-account",
	Short: "Extract account from slip by cropping and OCR",
	Run: func(cmd *cobra.Command, args []string) {
		if sourceFolder == "" {
			fmt.Println("Folder not found")
			return
		}
		if configType == "" {
			fmt.Println("Type not found")
			return
		}
		fmt.Printf("%s :: start cropping\n", time.Now().Format("2006-01-02 03:04:05"))
		numberOfFile := 0
		filepath.Walk(sourceFolder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				ext := strings.ToLower(filepath.Ext(path))
				if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
					os.MkdirAll("./output", os.ModePerm)
					outputPath := filepath.Join("./output", filepath.Base(path))
					err := cropImage(path, outputPath)

					if err != nil {
						log.Printf("Error cropping image %s: %s\n", path, err)
					} else {
						numberOfFile++
					}
				}
			}
			return nil
		})
		fmt.Printf("%s :: %d image(s) cropped\n", time.Now().Format("2006-01-02 03:04:05"), numberOfFile)
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
