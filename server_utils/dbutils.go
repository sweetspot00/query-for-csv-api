package server_utils

import (
	"encoding/csv"
	"io"
	"log"
	"os"

	"golang.org/x/text/encoding/charmap"
)

func WriteCsv(data [][]string) {

	f, err := os.Create("../result.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// 写入UTF-8 BOM
	f.WriteString("\xEF\xBB\xBF")

	w := csv.NewWriter(f)

	w.WriteAll(data)
	w.Flush()
}

func ReadCsv() ([][]string, map[string]int, error) {
	// Open the file
	csvfile, err := os.Open("../data.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	defer csvfile.Close()

	// Parse the file
	r := csv.NewReader(charmap.ISO8859_15.NewDecoder().Reader(csvfile))
	var data [][]string
	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		} else {
			data = append(data, record)
		}

	}

	col_name := make(map[string]int)
	for i, name := range data[0] {
		col_name[name] = i
	}
	//第一行是col name
	return data, col_name, nil
}
