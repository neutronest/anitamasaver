package main

import (
    "fmt"
    "os"
    "strconv"
    "sync"
    // "time"
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
        publishTime := infoParagraph.Find("span", "class", "time")

        if title.Error != nil {
            continue
        }
        if subtitle.Error != nil {
            continue
        }
        if author.Error != nil {
            continue
        }
        if publishTime.Error != nil {
            continue
        }
        // fmt.Println(title.Text())
        // fmt.Println(subtitle.Text())
        // fmt.Println(author.Text())
        // fmt.Println(category.Text())
        // fmt.Println(publishTime.Text())
        // fmt.Println(anitama.ANITAMA_ROOT_URL + articleLink)
        // fmt.Println()

        articleMetadata := anitama.ArticleMetaData{
            Title: title.Text(),
            SubTitle: subtitle.Text(),
            Author: author.Text(),
            Category: category.Text(),
            Date: publishTime.Text(),
            Url: anitama.ANITAMA_ROOT_URL + articleLink}
        articleMetadatas = append(articleMetadatas, articleMetadata)
        
    }
    return articleMetadatas
}

func main(){

    articleMetadataChan := make(chan anitama.ArticleMetaData, 100)
    // anitama.CHANNEL_MAX_PAGINATION
    paginationNum := 10
    var wg sync.WaitGroup
    wg.Add(paginationNum)
    for idx := 1; idx <= paginationNum; idx++ {

        go func(channelPageId int) {
            defer wg.Done()
            channelUrl := anitama.ANITAMA_CHANNEL_URL + "/all/" + strconv.Itoa(channelPageId)
            fmt.Println("\n\n=====ChannelPageId", channelPageId)
            articleMetadatas := getArticleMetadatasFromChannel(channelUrl)
            for _, articleMetadata := range articleMetadatas {
                articleMetadataChan <- articleMetadata
            }
            fmt.Println("ChannelPageId ",channelPageId, "ended. ")
        }(idx)
        // if idx % 10 == 0 {
        //     time.Sleep(1000 * time.Millisecond)
        // }
        
    }
    
    for  {
        select {
        case articleMetadata := <-articleMetadataChan:
            metadataJson, err := json.Marshal(articleMetadata)
            if err != nil {
                fmt.Println("Axiba...")    
            }
            fmt.Println(string(metadataJson))
        }
    }
    wg.Wait()
    
    

    
 }