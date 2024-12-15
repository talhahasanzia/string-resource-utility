package writer

import (
	"bufio"
	"fmt"
	. "localize/record"
	"log"
	"os"
)

func WriteFile(platform string, locale string, base string, channel chan Record, output string, debugFlag bool, overwrite bool) {

	fileName := getFilename(platform, locale, base, output)
	dirName := getDirname(platform, locale, base, output)

	os.MkdirAll(dirName, os.ModePerm)

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	defer file.Close()

	if overwrite {
		file.Truncate(0)
	}

	datawriter := bufio.NewWriter(file)

	preText := getPretext(platform)

	datawriter.WriteString(preText)
	datawriter.Flush()

	if debugFlag {
		fmt.Printf("Written %s to %s\n", preText, fileName)
	}

	for data := range channel {

		if data.Key != "" {

			formatted := getFormattedEntry(platform, data.Key, data.Value)
			datawriter.WriteString(formatted)
			datawriter.Flush()

			if debugFlag {
				fmt.Printf("\nWritten %s to %s", formatted, fileName)
			}

		}
	}
}

func CloseFile(platform string, locale string, base string, output string, debugFlag bool) {

	fileName := getFilename(platform, locale, base, output)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed open file: %s", err)
	}

	defer file.Close()

	datawriter := bufio.NewWriter(file)

	postText := getPosttext(platform)
	datawriter.WriteString(postText)
	datawriter.Flush()
	if debugFlag && postText != "" {
		fmt.Printf("\nWritten %s to %s", postText, fileName)
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

func getFilename(platform string, locale string, base string, output string) string {
	dir := getDirname(platform, locale, base, output)
	if platform == "ios" {
		return dir + "/Localizable.strings"
	} else if platform == "android" {
		return dir + "/strings.xml"
	} else {
		if locale == base {
			return dir + "/strings.ts"
		} else {
			return dir + "/strings_" + locale + ".ts"
		}
	}
}

func getDirname(platform string, locale string, base string, output string) string {
	if platform == "ios" {
		if locale == base {
			return output + "/base.lproj"
		} else {
			return output + "/" + locale + ".lproj"
		}
	} else if platform == "android" {
		if locale == base {
			return output + "/values"
		} else {
			return output + "/values-" + locale
		}
	} else {
		return output
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
