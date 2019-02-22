package main

import (
    "fmt"
    "os"
    "encoding/json"
	"github.com/anaskhan96/soup"
    "github.com/neutronest/anitamasaver/anitama"
)

func getArticleMetadatasFromChannel(channelUrl string) []anitama.ArticleMetaData {
    response, err := soup.Get(channelUrl)
	if err != nil {
		os.Exit(1)
    }
    doc := soup.HTMLParse(response)
    areaArticleChannelDiv := doc.Find("div", "id", "area-article-channel")
    innerDiv := areaArticleChannelDiv.Find("div", "class", "inner")
    links := innerDiv.FindAll("a")

    articleMetadatas := []anitama.ArticleMetaData{}
    for _, articleMetadataElement := range links {
        articleLink := articleMetadataElement.Attrs()["href"]
        title := articleMetadataElement.Find("h1")
        subtitle := articleMetadataElement.Find("h3")
        
        infoParagraph := articleMetadataElement.Find("p", "class", "info-article-channel")
        descParagraph := articleMetadataElement.Find("p", "class", "desc")
        if infoParagraph.Error != nil {
            continue
        }
        if descParagraph.Error != nil {
            continue
        }
        author := infoParagraph.Find("span", "class", "author")
        category := infoParagraph.Find("span", "class", "channel")

        if title.Error != nil {
            continue
        }
        if subtitle.Error != nil {
            continue
        }
        if author.Error != nil {
            continue
        }
        // fmt.Println(title.Text())
        // fmt.Println(subtitle.Text())
        // fmt.Println(author.Text())
        // fmt.Println(category.Text())
        // fmt.Println(anitama.ANITAMA_ROOT_URL + articleLink)
        // fmt.Println()

        articleMetadata := anitama.ArticleMetaData{
            Title: title.Text(),
            SubTitle: subtitle.Text(),
            Author: author.Text(),
            Category: category.Text(),
            Date: "",
            Url: anitama.ANITAMA_ROOT_URL + articleLink}
        articleMetadatas = append(articleMetadatas, articleMetadata)
        
    }
    return articleMetadatas
}

func main(){

    articleMetadataChan := make(chan anitama.ArticleMetaData, 100)
    for idx := 1; idx <= 10; idx++ {

        go func(channelPageId int) {
            channelUrl := anitama.ANITAMA_CHANNEL_URL + "/all/" + string(channelPageId)
            articleMetadatas := getArticleMetadatasFromChannel(channelUrl)
            for _, articleMetadata := range articleMetadatas {
                articleMetadataChan <- articleMetadata
            }
        }(idx)
    }

    for {
        select {
        case articleMetadata := <-articleMetadataChan:
            metadataJson, err := json.Marshal(articleMetadata)
            if err != nil {
                fmt.Println("Axiba...")    
            }
            fmt.Println(string(metadataJson))
        }
    }

    
 }