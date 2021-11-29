package githubtrending

import (
	"encoding/json"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"io/ioutil"
	"os"
)

func ReadJson(filename string) []Repositories {
	jsonFile, err := os.Open(filename)

	if err != nil{
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	repolist := make([]Repositories,0)

	err = json.Unmarshal([]byte(byteValue), & repolist)
	if err == nil {
	} else {
		fmt.Println(err)
	}
	return repolist
}

func DrawBar(filename string) {

	repoList := ReadJson(filename + ".json")
	languageMap := make(map[string]int)

	for _, repo := range repoList {
		_, ext := languageMap[repo.Language]
		if ext {
			languageMap[repo.Language] += 1
		}else{
			languageMap[repo.Language] = 1
		}
	}
	var namelist []string
	var numberlist []opts.BarData
	for name, data := range languageMap{
		namelist = append(namelist, name)
		numberlist = append(numberlist, opts.BarData{Value: data})
	}

	bar := charts.NewBar()

	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Bar"}),
		charts.WithXAxisOpts(opts.XAxis{AxisLabel: &opts.AxisLabel{ShowMinLabel: true, ShowMaxLabel: true, Interval: "0", Rotate: -40}}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Amount"}),
		charts.WithInitializationOpts(opts.Initialization{Width: "800px", Height: "600px"}),

	)

	// Put data into instance
	bar.SetXAxis(namelist).AddSeries("Github", numberlist)
	// Where the magic happens
	f, _ := os.Create(filename + "_bar.html")
	bar.Render(f)
}

func DrawWordCloud(filename string){
	repoList := ReadJson(filename + ".json")
	languageMap := make(map[string]int)

	for _, repo := range repoList {
		_, ext := languageMap[repo.Language]
		if ext {
			languageMap[repo.Language] += 1
		}else{
			languageMap[repo.Language] = 1
		}
	}

	var wordCloudDataList []opts.WordCloudData
	for name, amount := range languageMap{
		wordCloudDataList = append(wordCloudDataList, opts.WordCloudData{Name: name, Value: amount})
	}

	WordCloud := charts.NewWordCloud()

	WordCloud.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "WordCloud"}),
		charts.WithInitializationOpts(opts.Initialization{Width: "500px", Height: "500px"}),
	)

	WordCloud.AddSeries("wordcloud", wordCloudDataList)

	// Put data into instance
	// Where the magic happens
	f, _ := os.Create(filename + "_WordCloud.html")
	WordCloud.Render(f)
}

func DrawOverlap(filename string){

	repoList := ReadJson(filename + ".json")
	var namelist []string
	var forkMap []opts.BarData
	var starMap []opts.LineData

	for _, repo := range repoList {
		namelist = append(namelist, repo.Name)
		forkMap = append(forkMap, opts.BarData{Value: repo.Fork})
		starMap = append(starMap, opts.LineData{Value: repo.TotalStar})
	}


	bar := charts.NewBar()

	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Overlap"}),
		charts.WithXAxisOpts(opts.XAxis{AxisLabel: &opts.AxisLabel{ShowMinLabel: true, ShowMaxLabel: true, Interval: "0", Rotate: -30}}),
		charts.WithInitializationOpts(opts.Initialization{Width: "800px", Height: "600px"}),
		charts.WithLegendOpts(opts.Legend{Show: true}),
	)
	bar.SetXAxis(namelist).AddSeries("fork", forkMap)
	bar.SetSeriesOptions(charts.WithLabelOpts(opts.Label{Show: true}))


	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: filename}),
		charts.WithXAxisOpts(opts.XAxis{AxisLabel: &opts.AxisLabel{ShowMinLabel: true, ShowMaxLabel: true, Interval: "0", Rotate: -30}}),
		charts.WithInitializationOpts(opts.Initialization{Width: "800px", Height: "600px"}),
	)
	line.SetXAxis(namelist).AddSeries("star", starMap)
	line.SetSeriesOptions(charts.WithLabelOpts(opts.Label{Show: true}))

	bar.Overlap(line)
	f, _ := os.Create(filename + "_Overlap.html")
	bar.Render(f)
}


