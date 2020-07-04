package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

const defaultPicDir, indexUrlPC, userAgentPC = "./pics", "https://www.xuexi.cn/lgdata/index.json",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36"

var dir = flag.String("d", defaultPicDir, "图片的保存路径")

func getPc(dir string) {
	_, err := os.Stat(dir)
	if err != nil {
		if !os.IsExist(err) {
			log.Println("文件夹不存在，创建!")
			err := os.MkdirAll(dir, os.ModePerm)

			if err != nil {
				log.Fatalln("文件夹创建失败，退出")
				return
			}
		}
	}
	// 获得url数据
	client := &http.Client{}
	req, _ := http.NewRequest("GET", indexUrlPC, nil)
	req.Header.Set("User-Agent", userAgentPC)
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	//pageData first-slider data[{url,text,link,itemId}]
	var rowJson = string(body)
	var rootJson map[string]interface{}
	json.Unmarshal([]byte(rowJson), &rootJson)
	// walk on json node.
	var pageData = rootJson["pageData"].(map[string]interface{})
	var first = pageData["first-slider"].(map[string]interface{})
	var data = first["data"]

	// for parse image field to handle
	items := reflect.ValueOf(data)
	for i := 0; i < items.Len(); i++ {
		item := items.Index(i).Interface().(map[string]interface{})
		// useful fields
		//itemId := item["itemId"]
		//text := item["text"]
		//link := item["link"]
		url := item["url"].(string)
		lastDot := strings.LastIndex(url, ".")
		suffix := url[lastDot:len(url)]
		imageFilePath, _ := filepath.Abs(path.Join(dir, strconv.Itoa(i)+suffix))

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", userAgentPC)
		// download
		resp, _ := client.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		out, _ := os.Create(imageFilePath)

		io.Copy(out, bytes.NewReader(body))

		out.Close()
		resp.Body.Close()
	}
}

func main() {
	flag.Parse()
	fmt.Println("图片保存路径为:", *dir)
	getPc(*dir)
}
