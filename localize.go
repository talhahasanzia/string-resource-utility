package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type record struct {
	key    string
	value  string
	locale string
}

func main() {

	csvFile := flag.String("file", "csv", "csv file name with path")

	platform := flag.String("platform", "", "target platform: ios/android/web")

	debugFlag := flag.Bool("debug", false, "enable detailed logs for debugging: true/false")

	flag.Parse()

	fmt.Println("Reading from file:", *csvFile)

	f, err := os.Open(*csvFile)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	recordList := parseData(data)

	channelMap := make(map[string]chan record)

	defer func() {
		// All data sent, setting post text to complete the OS specific file formatting

		for i, locale := range data[0] {
			if i == 0 {
				continue
			}
			// Send empty record to close file with end texts
			channelMap[locale] <- record{}
			close(channelMap[locale])
		}

	}()

	for i, locale := range data[0] {
		if i == 0 {
			continue
		}
		channelMap[locale] = make(chan record)
	}

	for i, locale := range data[0] {
		if i == 0 {
			continue
		}

		go writeFile(*platform, locale, channelMap[locale], *debugFlag)

	}

	fmt.Println("Writing strings to files... ")

	for _, record := range recordList {
		channelMap[record.locale] <- record
	}

}

func parseData(data [][]string) []record {

	var dataList []record
	m := make(map[int]string)
	for i, line := range data {
		if i == 0 {

			for j, field := range line {

				m[j] = field

			}

			continue

		}
		var key string
		for j, field := range line {
			rec := record{}
			if j == 0 {
				key = field
			} else {
				rec.key = key
				rec.value = field
				rec.locale = m[j]
				dataList = append(dataList, rec)
			}

		}

	}
	return dataList
}

func getTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func writeFile(platform string, locale string, channel chan record, debugFlag bool) {

	fileName := getFilename(platform, locale)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	defer file.Close()
	defer fmt.Printf("\nClosing %s", fileName)

	datawriter := bufio.NewWriter(file)

	preText := getPretext(platform)

	datawriter.WriteString(preText)
	datawriter.Flush()

	if debugFlag {
		fmt.Printf("\nWritten %s to %s", preText, fileName)
	}

	for data := range channel {

		if data.key != "" {

			formatted := getFormattedEntry(platform, data.key, data.value)
			datawriter.WriteString(formatted)
			datawriter.Flush()

			if debugFlag {
				fmt.Printf("\nWritten %s to %s", formatted, fileName)
			}

		} else {

			postText := getPosttext(platform)
			datawriter.WriteString(postText)
			datawriter.Flush()
			if debugFlag && postText != "" {
				fmt.Printf("\nWritten %s to %s", postText, fileName)
			}

		}
	}

}

func getFormattedEntry(platform string, key string, value string) string {
	if platform == "ios" {
		return "\"" + key + "\" =" + " \"" + value + "\";\n"
	} else if platform == "android" {
		return "<string name=\"" + key + "\">" + value + "</string>\n"
	} else {
		return key + ": \"" + value + "\",\n"
	}
}

func getFilename(platform string, locale string) string {
	if platform == "ios" {
		return "Localized_" + locale + ".strings"
	} else if platform == "android" {
		return "strings_" + locale + ".xml"
	} else {
		return "strings_" + locale + ".ts"
	}
}

func getPretext(platform string) string {
	if platform == "ios" {
		return ""
	} else if platform == "android" {
		return "<resources>\n"
	} else {
		return "const LOCALIZED_STRINGS = {\n"
	}
}

func getPosttext(platform string) string {
	if platform == "ios" {
		return ""
	} else if platform == "android" {
		return "</resources>"
	} else {
		return "}\n\nexport default LOCALIZED_STRINGS;"
	}
}
