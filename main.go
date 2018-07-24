package main

import (
	"net/http"
	"time"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

func main() {
	page := int64(1)
	url := "https://venmo.com/api/v5/public?limit=1000"
	for {
		time.Sleep(100 * time.Millisecond)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting URL: %s\n", err)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting body: %s\n", err)
			continue
		}

		data := make(map[string]interface{})
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding data: %s\n", err)
			continue
		}

		pagingMap, cast := data["paging"].(map[string]interface{})
		if !cast {
			fmt.Fprintf(os.Stderr, "Error: couldn't find paging map in response\n")
			continue
		}

		prevStr, cast := pagingMap["previous"].(string)
		if !cast {
			fmt.Fprintf(os.Stderr, "Error: couldn't find paging map in response\n")
			continue
		}

		deferWrapper := func() bool {
			file, err := os.Create(strconv.FormatInt(page, 10) + ".json")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating JSON file: %s\n", err)
				return false
			}
			defer file.Close()

			encoder := json.NewEncoder(file)

			err = encoder.Encode(data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error encoding JSON: %s\n", err)
				return false
			}

			fmt.Printf("Encoded page: %s\n", url)
			return true
		}
		if deferWrapper() {
			url = prevStr
		}
	}
}
