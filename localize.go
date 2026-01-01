package main

import (
	"flag"
	"fmt"
	. "localize/parser"
	. "localize/reader"
	. "localize/record"
	. "localize/writer"
)

const version = "v1.1"

func main() {

	csvFile := flag.String("file", "", "csv file name with absolute path")

	platform := flag.String("platform", "", "target platform: ios/android/web/flutter")

	output := flag.String("output", "./", "output directory for generated files")

	debugFlag := flag.Bool("debug", false, "enable detailed logs for debugging: true/false")

	overwrite := flag.Bool("overwrite", true, "overwrite existing file contents: true/false")

	versionFlag := flag.Bool("v", false, "print version information")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("localize version %s\n", version)
		return
	}

	if *csvFile == "" {
		fmt.Println("Invalid input for file name/path.")
		return
	}

	fmt.Println("Reading from file:", *csvFile)

	data := ReadFile(csvFile)

	recordList := ParseData(data)

	// Group records by locale
	localeRecords := make(map[string][]Record)
	for _, record := range recordList {
		localeRecords[record.Locale] = append(localeRecords[record.Locale], record)
	}

	// Calculate unique keys
	uniqueKeys := make(map[string]bool)
	for _, record := range recordList {
		uniqueKeys[record.Key] = true
	}

	fmt.Printf("CSV parsed successfully. Found %d unique keys across %d locales (%d total records).\n", len(uniqueKeys), len(localeRecords), len(recordList))

	fmt.Printf("Writing strings to files for %d locales...\n", len(localeRecords))

	// Write files sequentially for each locale
	for locale, records := range localeRecords {
		WriteFileSequential(*platform, locale, records, *output, *debugFlag, *overwrite)
		CloseFile(*platform, locale, *output, *debugFlag)
	}
}
