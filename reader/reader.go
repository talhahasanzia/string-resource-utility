package reader

import (
	"encoding/csv"
	"log"
	"os"
)

func ReadFile(csvFile *string) [][]string {
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
	return data
}
