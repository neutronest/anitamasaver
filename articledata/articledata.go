package main

import (
	"io/ioutil"
	"os"
	"fmt"
	"strings"
	"encoding/json"
	"strconv"
	"github.com/anaskhan96/soup"
	"github.com/neutronest/anitamasaver/anitama"
	"github.com/bradfitz/slice"
)

func parseArticleContent(articleURL string, articleMetadata anitama.ArticleMetaData) {
	fmt.Println(articleURL)
	response, err := soup.Get(articleURL)
	if err != nil {
		fmt.Println("response error")
		os.Exit(1)
    }
	doc := soup.HTMLParse(response)
	mainDiv := doc.Find("div", "id", "mainer")
	innerDiv := mainDiv.Find("div", "class", "inner")
	// topArticleDiv := innerDiv.Find("div", "id", "area-top-article")
	// style := topArticleDiv.Attrs()["style"]
	// fmt.Println(topArticleDiv.Text())
	// fmt.Println(style)
	if innerDiv.Error != nil {
		fmt.Println("innerDiv error")
	}
	fmt.Println(innerDiv.Pointer.Data, innerDiv.NodeValue)
	contentDiv := innerDiv.Find("div", "id", "area-content-article")
	data := contentDiv.Children()
	for _, component := range data {

		if (component.NodeValue == "p") {
			if component.Children()[0].NodeValue == "img" {
				fmt.Println(component.Children()[0].Attrs()["data-src"])
			} else {
				fmt.Println(component.Text() + "\n")
			}
		} else if (component.NodeValue == "blockquote") {
			fmt.Println("[==BLOCKQUOTE WARNIGN===]")
			fmt.Println(component.FullText())
		} 
		
	}

}

func main() {
	fmt.Println("Hello")
	var articleMetadatas []anitama.ArticleMetaData

	articleMetadataJsonFile, err := os.Open("../article_metadata.json")
	if err != nil {
		fmt.Println(err)
	}
	defer articleMetadataJsonFile.Close()

	byteValue, _ := ioutil.ReadAll(articleMetadataJsonFile)
	json.Unmarshal(byteValue, &articleMetadatas)

	slice.Sort(articleMetadatas[:], func(i, j int) bool {

		leftStringArr := strings.Split(articleMetadatas[i].Id, "-")
		rightStringArr := strings.Split(articleMetadatas[j].Id, "-")

		leftPageId, _ := strconv.Atoi(leftStringArr[0])
		rightPageId, _ := strconv.Atoi(rightStringArr[0])
		leftPageItemId, _ := strconv.Atoi(leftStringArr[1])
		rightPageItemId, _ := strconv.Atoi(rightStringArr[1])

		if leftPageId != rightPageId {
			return leftPageId < rightPageId
		} else {
			return leftPageItemId < rightPageItemId
		}
	})
	
	for _, articleMetadata := range articleMetadatas {
		fmt.Println(articleMetadata.Title, articleMetadata.Id)
	}

	parseArticleContent(articleMetadatas[0].Url, articleMetadatas[0])


}