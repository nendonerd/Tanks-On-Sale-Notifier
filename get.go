package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	runewidth "github.com/mattn/go-runewidth"
)

type mapper = map[string]interface{}

var response mapper

type detail struct {
	price string
	end   string
}

func main() {
	// load the previous result from github artifacts, see https://github.com/actions/download-artifact
	list := crawl()
	res := extract(list)
	// sort and compare keys between prev result and curr result, if match then abort, else
	// save the result to github artifacts, see https://github.com/actions/upload-artifact
	printMap(res)
	// format the result to a twitter post
	// call twitter api to post, the api key is stored in github secrets
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

			if m[name] == (detail{}) || price < m[name].price {
				m[name] = detail{price, end}
			}

		}
	}
	return m
}

func printMap(m map[string]detail) {
	// Get the max width of names
	cellWidth := 0
	for name, _ := range m {
		width := runewidth.StringWidth(name)
		if width > cellWidth {
			cellWidth = width
		}
	}
	// Print results
	var b strings.Builder
	for name, d := range m {
		price := d.price
		end := d.end
		width := runewidth.StringWidth(name)
		space := strings.Repeat(" ", cellWidth+2-width)
		fmt.Fprintf(&b, "%s%s￥%s  至%s\n", name, space, price, end)

	}
	output := b.String()
	fmt.Println(output)
	fmt.Println(len(output))
}
