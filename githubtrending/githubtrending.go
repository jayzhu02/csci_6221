// Package githubtrending
//Author: Jie Zhu
//Date: 2021.11.27
//This code is used for grabbing data from GitHub Trending.
//
package githubtrending

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var FileBase, _ = filepath.Abs(filepath.Dir("./"))

type Repositories struct {
	Name        string `json:"name"`
	TotalStar   string `json:"total_star"`
	Fork        string `json:"fork"`
	Link  		string `json:"link"`
	Description string `json:"description"`
	Language    string `json:"language"`
}

type RepositoriesList []Repositories

var All RepositoriesList

func TrendingContent(language string, since string) RepositoriesList {
	url := fmt.Sprintf("https://github.com/trending/%s?since=%s", language, since)
	req := gorequest.New()
	response, _, _ := req.Get(url).End()
	defer response.Body.Close()

	responseBytes, _ := ioutil.ReadAll(response.Body)

	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(responseBytes))
	doc.Find(".Box .Box-row").Each(func(i int, selection *goquery.Selection) {

		name := strings.TrimSpace(selection.Find("h1 a").Contents().Last().Text())

		relativeLink, _ := selection.Find("h1 a").Attr("href")

		link := "https://github.com" + relativeLink

		description := strings.TrimSpace(selection.Find("p").Text())

		var langIdx int
		spanSel := selection.Find("div>span")
		if spanSel.Size() == 2 {
			// language not exist
			langIdx = -1
		}

		// language
		if langIdx >= 0 {
			language = strings.TrimSpace(spanSel.Eq(langIdx).Text())
		} else {
			language = "unknown"
		}

		aSel := selection.Find("div>a")
		starStr := strings.TrimSpace(aSel.Eq(-2).Text())
		totalStar := strings.Replace(starStr, ",", "", -1)

		forkStr := strings.TrimSpace(aSel.Eq(-1).Text())
		fork := strings.Replace(forkStr, ",", "", -1)

		//fmt.Println("name: ", name, "\ndescription: ", description, "\ntotalStar: ", totalStar, "\nfork: ", fork, "\nLink: ", link, "\nlanguage: ", language)

		data := Repositories{
			Name:        name,
			Description: description,
			TotalStar:   totalStar,
			Fork:        fork,
			Link: 		 link,
			Language:    language,
		}
		All = append(All, data)

	})
	return All

}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func SaveText(result RepositoriesList, filename string) {

	var file *os.File
	fileName := "./static/" + filename + ".txt"

	if CheckFileIsExist(fileName) {
		fmt.Println(fileName)
		os.Remove(fileName)
	}

	file, _ = os.Create(fileName)
	for _, one := range result {
		string := fmt.Sprintf("Name: %s, Link: %s, TotalStar: %s, Language: %s, Description: %s",
			one.Name, one.Link, one.TotalStar, one.Language, one.Description)
		file.WriteString(string)
		file.WriteString("\n")
	}

	defer file.Close()
}

func SaveJson(result RepositoriesList, filename string) {

	var file *os.File
	fileName := FileBase + "//static/" + filename + ".json"

	if CheckFileIsExist(fileName) {
		os.Remove(fileName)
	}

	file, _ = os.Create(fileName)
	byteResponse, _ := json.MarshalIndent(result, "", " ")
	file.WriteString(string(byteResponse))
	defer file.Close()

}

func SaveCsv(result RepositoriesList, filename string) {

	var f *os.File
	var err error
	fileName := FileBase + "//static/" + filename + ".csv"
	if CheckFileIsExist(fileName) { //如果文件存在
		os.Remove(fileName)
	}
	f, err = os.Create(fileName) //创建文件
	if err != nil {
		f, err = os.Create(FileBase + "/trending.csv")
		fmt.Println(err, f.Name())
	}

	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f) //创建一个新的写入文件流
	var header = []string{"name", "link", "totalStar", "fork", "language", "description"}
	w.Write(header)
	var data [][]string
	for _, one := range result {
		var singleString []string
		singleString = append(singleString, one.Name, one.Link, one.TotalStar, one.Fork, one.Language, one.Description)
		data = append(data, singleString)

	}

	w.WriteAll(data) //写入数据
	w.Flush()
}

func TrendingStart(language string, since string) {
		TrendingContent(language, since)
		if language == ""{
			language = "Any"
		}

		filename := "githubtrending_" + language + "_" + since
		SaveText(All, filename)
		SaveJson(All, filename)
		SaveCsv(All, filename)
}

