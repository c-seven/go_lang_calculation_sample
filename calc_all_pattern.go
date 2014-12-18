package main

import (
	"encoding/csv"
  "fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"code.google.com/p/go.text/encoding/japanese"
	"code.google.com/p/go.text/transform"
)

func noErrorOrPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func validate(records [][]string) {
	if len(records) == 0 {
		panic("ファイルが空っすね")
	}
}

func main() {
  if len(os.Args) < 2  {
     panic("ファイルを指定してください")
  }
	file, err := os.Open(os.Args[1])

	noErrorOrPanic(err)
	defer file.Close()

  reader := csv.NewReader(transform.NewReader(file, japanese.ShiftJIS.NewDecoder()))

	var records [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else {
			noErrorOrPanic(err)
		}

		records = append(records, record)
	}

	validate(records)
	header, data := records[:1][0], records[1:]

	// 転置
	var vs [][]float64 = make([][]float64,len(header))
	for _, row := range data {
		for col := 0; col < len(header); col++ {
			if v, err := strconv.ParseFloat(strings.Trim(row[col], " "), 64); err == nil {
				vs[col] = append(vs[col], v)
			}
		}
	}

  output, err := os.Create(file.Name() + ".calc.csv")
  writer := csv.NewWriter(transform.NewWriter(output, japanese.ShiftJIS.NewEncoder()))

	writer.Write(append(header,"計算結果"))

  toStrs := func(nums []float64) []string{
     ss := []string{}
     for _, n := range nums{
       ss = append(ss, fmt.Sprintf("%.4f", n))
     }
     return ss
  }

	multipleAll := func(nums []float64) {
		r := 1.0

		for _, n := range nums {
			r = r * n
		}
		writer.Write(toStrs(append(nums,r)))
	}

	var d []float64
	walk(vs, d, multipleAll)

  writer.Flush()
}

func walk(datas [][]float64, currentDatas []float64, proc func([]float64)) {
	currentLevelDatas, nextLevelDatas := datas[:1][0], datas[1:]
	if len(nextLevelDatas) == 0 {
		for _, each := range currentLevelDatas {
			ds := append(currentDatas, each)
			proc(ds)
		}
		return
	}

	for _, each := range currentLevelDatas {
		ds := append(currentDatas, each)
		walk(nextLevelDatas, ds, proc)
	}

}
