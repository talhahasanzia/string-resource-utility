package parser

import . "localize/record"

func ParseData(data [][]string) []Record {

	var dataList []Record
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
			rec := Record{}
			if j == 0 {
				key = field
			} else {
				rec.Key = key
				rec.Value = field
				rec.Locale = m[j]
				dataList = append(dataList, rec)
			}

		}

	}
	return dataList
}
