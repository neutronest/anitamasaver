package main

import (
	"time"
	"sync/atomic"
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



func parseArticleContent(articleURL string, articleMetadata anitama.ArticleMetaData) anitama.Article {
	fmt.Println("saving...", articleURL)
	response, err := soup.Get(articleURL)
	if err != nil {
		fmt.Println("response error")
		os.Exit(1)
    }
	doc := soup.HTMLParse(response)
	mainDiv := doc.Find("div", "id", "mainer")
	innerDiv := mainDiv.Find("div", "class", "inner")
	
	if innerDiv.Error != nil {
		fmt.Println(innerDiv.Error)
	}

	articleHtmlContent := ""
	articleRawContent := ""
	contentDiv := innerDiv.Find("div", "id", "area-content-article")
	data := contentDiv.Children()
	for _, component := range data {

		if (component.NodeValue == "p") {
			if len(component.Children()) == 0 {
				continue
			}
			if component.Children()[0].NodeValue == "img" {
				imgUrl := component.Children()[0].Attrs()["data-src"]
				articleHtmlContent += "<p><img src='" + imgUrl + "' /></p>"
				articleRawContent += "<img src='" + imgUrl + "' />"
			} else {
				paragraphContent := component.Text()
				articleHtmlContent += "<p> " + paragraphContent + " </p>"
				articleRawContent += component.Text() + "\n\n"
			}
		} else if (component.NodeValue == "blockquote") {
			if len(component.Children()) > 1 {
				articleHtmlContent += "<blockquote>"
				articleRawContent += "“\n\n"
				
				paragraphs := component.FindAll("p")
				for _, p := range paragraphs {
					paragraphContent := p.Text()
				articleHtmlContent += "<p> " + paragraphContent + " </p>"
				articleRawContent += paragraphContent + "\n\n"
				}
				articleRawContent += "”\n\n"
				articleHtmlContent += "</blockquote>"
				continue
			}
			paragraphContent := component.FullText()
			articleHtmlContent += "<blockquote><p> " + paragraphContent + " </p></blockquote>"
			articleRawContent += component.Text() + "\n\n"
		} 
		
	}
	return anitama.Article {
		MetaData: articleMetadata,
		Content: articleHtmlContent,
		RawContent: articleRawContent}
}

// Workaround: global vaiable is ugly..
var articles []anitama.Article
var articleMetaDataChan chan anitama.ArticleMetaData
var articleChan chan anitama.Article
var workCounter uint64

func asyncFeedArticleMetaDataPipeline(articleMetadatas []anitama.ArticleMetaData) {

	go func() {
		for _, articleMetadata := range articleMetadatas {
			articleMetaDataChan <- articleMetadata
		}
	}()
}

func asyncParseArticles() {
	
	go func() {
		for {
			select {
			case articleMetadata :=<- articleMetaDataChan:
				fmt.Println("consume ", articleMetadata)
				fmt.Println("consume ", len(articleChan))
				article := parseArticleContent(
					articleMetadata.Url, 
					articleMetadata)
				fmt.Println("done ", len(articleChan))
				articleChan <- article
				
			default:
			}
			time.Sleep(5000 * time.Millisecond)
			fmt.Println("len of metadata chan in parse ", len(articleMetaDataChan))
		}
		
	}()
}

func asyncFeedArticlesPipeline() {

	go func() {
		
		for {	
			select {
			case article :=<- articleChan:
				articles = append(articles, article)
				atomic.AddUint64(&workCounter, 1);
				time.Sleep(500 * time.Millisecond)
				fmt.Println("workCounter: ", workCounter)
			default:
			}

			if workCounter >0 && workCounter % 100 == 0 {
				fmt.Println("workCounter: ", workCounter)
			}
		}
	}()
}

func waitForPipelineEnded(finishCount uint64) {
	ok := true
	for ok {
		if workCounter == finishCount {
			for article := range articleChan {
				articles = append(articles, article)
			}
			ok = false
			fmt.Println("workCounter finished..")
		}
	}
}

func main() {
	fmt.Println("Now begin to saving articles...")
	var articleMetadatas []anitama.ArticleMetaData
	articleChan = make(chan anitama.Article, 1000)
	articleMetaDataChan = make(chan anitama.ArticleMetaData, 1000)

	articleMetaDataFile := "../output/article_metadata.json"
	articleFile := "../output/article_data.json"
	articleMetadataJsonFile, err := os.Open(articleMetaDataFile)
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
	fmt.Println("len of all articles: ", len(articleMetadatas))
	
	// single
	for _, articleMetaData := range articleMetadatas {
		article := parseArticleContent(
			articleMetaData.Url, 
			articleMetaData)
		articles = append(articles, article)
	}

	// asyncFeedArticleMetaDataPipeline(articleMetadatas)
	// asyncParseArticles()
	// asyncFeedArticlesPipeline()
	// waitForPipelineEnded(uint64(len(articleMetadatas)))
	// batchSize := 100

	// for id, _ := range articleMetadatas {
	// 	fmt.Println(id)
	// 	if id > 200 {
	// 		break
	// 	}
	// 	if id % batchSize == 0 {
	// 		go func(articleIdx int) {
	// 			fmt.Println("parsing workload ", articleIdx, " start..")
	// 			articleIdxStart := articleIdx
	// 			articleIdxEnd := articleIdxStart + batchSize
	// 			for articleId := articleIdxStart; articleId < articleIdxEnd; articleId++ {
	// 				article := parseArticleContent(
	// 					articleMetadatas[articleId].Url, 
	// 					articleMetadatas[articleId])
	// 				articleChan <- article
	// 				time.Sleep(500 * time.Millisecond)
	// 			}
	// 			fmt.Println("parsing workload ", id, " ended..")

	// 		}(id)
	// 	}
	// }

	
	articlesJson, err := json.Marshal(articles)
	if err != nil {
    }
    // for _, jsonData := range articleMetadataJson {
    f, err := os.OpenFile(articleFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
    if err != nil {
        fmt.Println(err)
    }
    if _, err := f.Write(articlesJson); err != nil {
        fmt.Println(err)
    }
    if err := f.Close(); err != nil {
        fmt.Println(err) 
	}
	fmt.Println("Done..")
}