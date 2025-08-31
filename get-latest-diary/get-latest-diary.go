package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

const (
	BASEURL = "https://diary.sh4869.sh/indexes/"
)

type DiaryInfo struct {
	Url   string
	Title string
	Body  string
}

var rootCmd = &cobra.Command{
	Use:   "get-latest-diary",
	Short: "Get Latest Diary of diary.sh4869.sh",
	RunE: func(cmd *cobra.Command, args []string) error {
		indexes, err := getLatestDiaryDay()
		if err != nil {
			return err
		}
		info := getTwoWeekInfo(indexes)
		f := formatTwoWeekInfo(info)
		fmt.Println("RESULT<<EOF\n" + f + "\nEOF")
		return nil
	},
}

func getLatestDiaryDay() (map[string]DiaryInfo, error) {
	var r []byte
	year := time.Now().Year()
	for {
		resp, err := http.Get(BASEURL + strconv.Itoa(year) + ".json")
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 404 {
			r, err = io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			break
		}
		year -= 1
	}

	var indexes map[string]DiaryInfo
	err := json.Unmarshal(r, &indexes)
	if err != nil {
		return nil, err
	}
	return indexes, nil
}

func getTwoWeekInfo(indexes map[string]DiaryInfo) map[string]*DiaryInfo {
	start := time.Now()
	result := map[string]*DiaryInfo{}
	for i := 0; i < 14; i++ {
		t := start.AddDate(0, 0, -i)
		key := t.Format("2006/01/02")
		if v, ok := indexes[key]; ok {
			result[key] = &v
		}
	}
	return result
}

var weekdays = []string{"S", "M", "T", "W", "T", "F", "S"}

func formatTwoWeekInfo(indexes map[string]*DiaryInfo) string {
	// cli calendar のように投稿されている日は ☑、そうじゃない日は空を表示する
	result := ""
	// 曜日ヘッダー表示
	for _, day := range weekdays {
		result += fmt.Sprintf("%3s", day)
	}
	result += "\\n"

	// 1日の曜日を取得
	start := time.Now().AddDate(0, 0, -14)
	weekday := int(start.Weekday())

	for i := 1; i < weekday; i++ {
		result += "  "
	}
	for i := 0; i < 14; i++ {
		t := start.AddDate(0, 0, i)
		key := t.Format("2006/01/02")

		if _, ok := indexes[key]; ok {
			result += "  <" + indexes[key].Url + "|o>"
		} else {
			result += "  x"
		}
		if t.Weekday() == time.Saturday {
			result += t.Format("|  01/02\\n")
		}
	}
	return result
}

func main() {
	rootCmd.Execute()
}
