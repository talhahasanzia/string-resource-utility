package main

import (
	"flag"
	"fmt"
	. "localize/parser"
	. "localize/reader"
	. "localize/record"
	. "localize/writer"
)

func main() {

	csvFile := flag.String("file", "", "csv file name with absolute path")

	platform := flag.String("platform", "", "target platform: ios/android/web")

	base := flag.String("base", "en", "base locale")

	output := flag.String("output", "./", "output directory for generated files")

	debugFlag := flag.Bool("debug", false, "enable detailed logs for debugging: true/false")

	overwrite := flag.Bool("overwrite", true, "overwrite existing file contents: true/false")

	flag.Parse()

	if *csvFile == "" {
		fmt.Println("Invalid input for file name/path.")
		return
	}

	fmt.Println("Reading from file:", *csvFile)

	data := ReadFile(csvFile)

	recordList := ParseData(data)

	channelMap := make(map[string]chan Record)

	for i, locale := range data[0] {
		if i == 0 {
			continue
		}
		channelMap[locale] = make(chan Record)
	}

	for i, locale := range data[0] {
		if i == 0 {
			continue
		}

		go WriteFile(*platform, locale, *base, channelMap[locale], *output, *debugFlag, *overwrite)

	}

	fmt.Println("Writing strings to files... ")

	for i, record := range recordList {
		channelMap[record.Locale] <- record

		if i == len(recordList)-1 {
			for k := range channelMap {
				CloseFile(*platform, k, *base, *output, *debugFlag)
			}
		}
	}
}
