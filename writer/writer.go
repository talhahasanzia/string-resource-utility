package writer

import (
	"bufio"
	"fmt"
	. "localize/record"
	"log"
	"os"
)

func WriteFile(platform string, locale string, channel chan Record, debugFlag bool) {

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

		if data.Key != "" {

			formatted := getFormattedEntry(platform, data.Key, data.Value)
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
