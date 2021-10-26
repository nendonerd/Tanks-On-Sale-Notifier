package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	runewidth "github.com/mattn/go-runewidth"
)

var artifactDir string = "./artifact/"
var artifactName string = "tanks-info"
var artifactPath string = filepath.Join(artifactDir, artifactName)

type mapper = map[string]interface{}

var response mapper

type detail struct {
	Price string
	End   string
}

func main() {
	prev := load()
	list := crawl()
	curr := extract(list)
	isDiff := diff(prev, curr)
	if isDiff {
		save(curr)
		tweet := format2Tweet(curr)
		post(tweet)
	}
	printMap(curr)
}

func crawl() []interface{} {
	url := "https://shop.wot.360.cn/api/product/list?game_id=1&type=vehicles"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Referer", "https://shop.wot.360.cn/vehicles")
	client := &http.Client{}
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal([]byte(body), &response)
	data := response["data"].(mapper)
	list := data["list"].([]interface{})
	return list
}

func extract(list []interface{}) map[string]detail {
	// Extract useful info from json to map
	m := make(map[string]detail)
	for _, v := range list {
		item := v.(mapper)
		if strings.Contains(item["categories"].(string), "featured") {
			content := item["package_content"].([]interface{})[0].(mapper)
			name, ok := content["vehicle_name"].(string)
			if !ok {
				name = item["name"].(string)
			}

			price := item["price"].(string)
			end := item["nonselling_time"].(string)[5:16]

			if m[name] == (detail{}) || price < m[name].Price {
				m[name] = detail{price, end}
			}

		}
	}
	return m
}

func printMap(m map[string]detail) {
	// Get the max width of names
	cellWidth := 0
	for name := range m {
		width := runewidth.StringWidth(name)
		if width > cellWidth {
			cellWidth = width
		}
	}
	// Print results
	var b strings.Builder
	for name, d := range m {
		price := d.Price
		end := d.End
		width := runewidth.StringWidth(name)
		space := strings.Repeat(" ", cellWidth+2-width)
		fmt.Fprintf(&b, "%s%s￥%s  至%s\n", name, space, price, end)

	}
	output := b.String()
	fmt.Println(output)
	fmt.Println(len(output))
}

// save the result to github artifacts, see https://github.com/actions/upload-artifact
func save(m map[string]detail) {
	// 1. check if path exist
	err := os.MkdirAll(artifactDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	// 2. encoding/gob to serialize the map
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)
	err = e.Encode(m)
	if err != nil {
		log.Fatal(err)
	}
	// 3. save to path
	os.WriteFile(artifactPath, b.Bytes(), 0644)
}

// load the previous result from github artifacts, see https://github.com/actions/download-artifact
func load() map[string]detail {
	// 1. check if file exist
	// 2. if exist, read the file as map and returns it
	// 3. else return an empty map
	if _, err := os.Stat(artifactPath); os.IsNotExist(err) {
		m := make(map[string]detail)
		return m
	} else {
		data, err := os.ReadFile(artifactPath)
		if err != nil {
			log.Fatal(err)
		}

		reader := bytes.NewReader(data)
		var decodedMap map[string]detail
		d := gob.NewDecoder(reader)
		err = d.Decode(&decodedMap)
		return decodedMap
	}
}

// sort and compare keys between prev result and curr result, if match then abort, else
func diff(a, b map[string]detail) bool {
	// 1. sort keys of both map
	// 2. concat keys
	// 3. compare, if different return true else false
	tokenA := sortAndJoin(a)
	tokenB := sortAndJoin(b)
	if tokenA != tokenB {
		return true
	} else {
		return false
	}
}

func sortAndJoin(m map[string]detail) string {
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, "|")
}

// format the result to a twitter post
func format2Tweet(m map[string]detail) string {
	// 1. concat keys
	// 2. check whether exceed text limit, 280chars for twitter
	var b strings.Builder
	for name, d := range m {
		fmt.Fprintf(&b, "|%s:￥%s|", name, d.Price)
	}
	tweet := b.String()
	if runewidth.StringWidth(tweet) > 280 {
		tweet = tweet[:275] + "..."
	}
	return tweet
}

// call twitter api to post, the api key is stored in github secrets
func post(tweet string) {
	// 1. read apikey from github secrets
	// 2. call twitter api

}
