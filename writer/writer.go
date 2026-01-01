package writer

import (
	"bufio"
	"fmt"
	. "localize/record"
	"log"
	"os"
)

func WriteFileSequential(platform string, locale string, records []Record, output string, debugFlag bool, overwrite bool) {

	fileName := getFilename(platform, locale, output)
	dirName := getDirname(platform, locale, output)

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
		fmt.Printf("\nWritten %s to %s\n", preText, fileName)
	}

	for _, data := range records {

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


func CloseFile(platform string, locale string, output string, debugFlag bool) {

	fileName := getFilename(platform, locale, output)
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
	switch platform {
	case "ios":
		return "\"" + key + "\" =" + " \"" + value + "\";\n"
	case "android":
		return "<string name=\"" + key + "\">" + value + "</string>\n"
	case "flutter":
		return "\"" + key + "\": \"" + value + "\",\n"
	default:
		return key + ": \"" + value + "\",\n"
	}
}

func getFilename(platform string, locale string, output string) string {
	dir := getDirname(platform, locale, output)
	switch platform {
	case "ios":
		return dir + "/Localizable.strings"
	case "android":
		return dir + "/strings.xml"
	case "flutter":
		return dir + "/app_" + locale + ".arb"
	default:
		return dir + "/strings_" + locale + ".ts"
	}
}

func getDirname(platform string, locale string, output string) string {
	switch platform {
	case "ios":
		return output + "/" + locale + ".lproj"
	case "android":
		return output + "/values-" + locale
	case "flutter":
		return output + "/l10n"
	default:
		return output
	}
}

func getPretext(platform string) string {
	switch platform {
	case "ios":
		return ""
	case "android":
		return "<resources>\n"
	case "flutter":
		return "{\n"
	default:
		return "const LOCALIZED_STRINGS = {\n"
	}
}

func getPosttext(platform string) string {
	switch platform {
	case "ios":
		return ""
	case "android":
		return "</resources>"
	case "flutter":
		return "\"eof\":\"eof\"\n}\n"
	default:
		return "}\n\nexport default LOCALIZED_STRINGS;"
	}
}
