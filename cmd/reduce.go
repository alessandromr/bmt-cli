/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/h2non/bimg"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var Verbose bool
var Format string
var Path string

// reduceCmd represents the reduce command
var reduceCmd = &cobra.Command{
	Use:   "reduce",
	Short: "Remove exif metedata and reduce image size",
	Long:  `Remove exif metedata and reduce image size`,
	Run: func(cmd *cobra.Command, args []string) {
		path, err := os.Getwd()
		if err != nil {
			log.Fatal("Error current path")
		}
		path = filepath.Join(path, Path)
		fmt.Println("reduce called")
		fmt.Println("Working on file format: ", Format)
		fmt.Println("Working on path: ", path)

		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal("Cannot get files")
		}
		for _, v := range files {
			if v.IsDir() {
				continue
			}
			if !strings.Contains(v.Name(), Format) {
				continue
			}
			fmt.Printf("Reducing %s\n", v.Name())
			buffer, err := bimg.Read(path + "/" + v.Name())
			if err != nil {
				log.Fatal(err)
			}

			converted, err := bimg.NewImage(buffer).Convert(bimg.WEBP)
			if err != nil {
				log.Fatal(err)
			}

			processed, err := bimg.NewImage(converted).Process(bimg.Options{Quality: 70, StripMetadata: true})
			if err != nil {
				log.Fatal(err)
			}
			filename := strings.Replace(v.Name(), "Format", "webp", 1)
			writeError := bimg.Write(fmt.Sprintf("%s/processed-%s", path, filename), processed)
			if writeError != nil {
				log.Fatal(writeError)
			}
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&Format, "format", "f", "png", "media format")
	rootCmd.PersistentFlags().StringVarP(&Path, "path", "p", "", "media path")
	rootCmd.AddCommand(reduceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reduceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reduceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
