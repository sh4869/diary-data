package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

const (
	BASEURL = "https://diary.sh4869.sh/indexes/"
)

type index struct {
	Url   string
	Title string
	Body  string
}

var rootCmd = &cobra.Command{
	Use:   "get-latest-diary",
	Short: "Get Latest Diary of diary.sh4869.sh",
	RunE: func(cmd *cobra.Command, args []string) error {
		latest, err := getLatestDiaryDay()
		if err != nil {
			return err
		}
		fmt.Println("DIARY_LATEST_DATE=" + latest)
		return nil
	},
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check",
	RunE: func(cmd *cobra.Command, args []string) error {
		latest, err := getLatestDiaryDay()
		if err != nil {
			return err
		}
		t, err := time.Parse("2006/01/02", latest)
		if err != nil {
			return err
		}
		b := t.After(time.Now().AddDate(0, 0, -2))
		fmt.Println("DIARY_UPDATED=" + strconv.FormatBool(b))
		return nil
	},
}

func getLatestDiaryDay() (string, error) {
	var r []byte
	year := time.Now().Year()
	for {
		resp, err := http.Get(BASEURL + strconv.Itoa(year) + ".json")
		if err != nil {
			return "", err
		}
		if resp.StatusCode != 404 {
			r, err = io.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}
			break
		}
		year -= 1
	}
	var i map[string]index
	json.Unmarshal(r, &i)
	keys := make([]string, 0, len(i))
	for k := range i {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	return keys[0], nil
}

func main() {
	rootCmd.AddCommand(checkCmd)
	rootCmd.Execute()
}
