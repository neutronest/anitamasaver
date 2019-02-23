package main


import (
    "fmt"
    "os"
    "strconv"
    "time"
    "sync"
    "sync/atomic"
    "encoding/json"
	"github.com/anaskhan96/soup"
    "github.com/neutronest/anitamasaver/anitama"
)

var muList sync.Mutex = sync.Mutex{}
var metadatas []anitama.ArticleMetaData 
var workCounter uint64


func getArticleMetadatasFromChannel(channelUrl string, pageId int) []anitama.ArticleMetaData {
    response, err := soup.Get(channelUrl)
	if err != nil {
		os.Exit(1)
    }
    doc := soup.HTMLParse(response)
    areaArticleChannelDiv := doc.Find("div", "id", "area-article-channel")
    innerDiv := areaArticleChannelDiv.Find("div", "class", "inner")
    links := innerDiv.FindAll("a")

    articleMetadatas := []anitama.ArticleMetaData{}
    for articleIdInPage, articleMetadataElement := range links {
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

        articleMetadata := anitama.ArticleMetaData{
            Id: strconv.Itoa(pageId) + "-" + strconv.Itoa(articleIdInPage),
            Title: title.Text(),
            SubTitle: subtitle.Text(),
            Author: author.Text(),
            Category: category.Text(),
            Date: publishTime.Text(),
            Url: anitama.ANITAMA_ROOT_URL + articleLink}
        articleMetadatas = append(articleMetadatas,articleMetadata)
        
        
    }
    return articleMetadatas
}

func aritlceMetaDataToJsonPipeline(
    articleMetadataChan chan anitama.ArticleMetaData) {
    // articleMetadatas []anitama.ArticleMetaData) {
    
    go func() {
        for  {
            select {
            case articleMetadata :=<- articleMetadataChan:
                metadatas = append(metadatas, articleMetadata)
            }
        }
    }()
}

func main(){

    articleMetadataChan := make(chan anitama.ArticleMetaData, 1000)
    // anitama.CHANNEL_MAX_PAGINATION
    paginationNum := anitama.CHANNEL_MAX_PAGINATION
    batchSize := 50
    for idx := 0; idx < paginationNum; idx++ {
        if (idx % batchSize == 0) {
            go func(channelPageId int) {
                for pageId := channelPageId; pageId < channelPageId + batchSize; pageId ++ {
                    channelUrl := anitama.ANITAMA_CHANNEL_URL + "/all/" + strconv.Itoa(pageId)
                    fmt.Println("\n\n=====ChannelPageId", pageId)
                    fmt.Println(channelUrl)
                    articleMetadatas := getArticleMetadatasFromChannel(channelUrl, pageId)
                    for _, articleMetadata := range articleMetadatas {
                        articleMetadataChan <- articleMetadata
                    }
                    fmt.Println("ChannelPageId ",pageId, "ended. ")
                }
                atomic.AddUint64(&workCounter, uint64(batchSize));
            }(idx+1)
        }    
    }

    aritlceMetaDataToJsonPipeline(articleMetadataChan)

    ok := true
    for ok {
        if (int(workCounter) == paginationNum && len(articleMetadataChan) == 0) {
            fmt.Println("Finish", workCounter);
            time.Sleep(5000 * time.Millisecond)
            ok = false
        }
        fmt.Println("workCounter", workCounter)
        time.Sleep(3000 * time.Millisecond)
    }

    fmt.Println("len of chan", len(articleMetadataChan))
    metadataJson, err := json.Marshal(metadatas)
    fmt.Println("json", metadataJson)
    fmt.Println("data", metadatas)
    if err != nil {
        fmt.Println("Axiba...")    
    }
    // for _, jsonData := range articleMetadataJson {
    f, err := os.OpenFile("./article_metadata.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
    if err != nil {
        fmt.Println(err)
    }
    if _, err := f.Write(metadataJson); err != nil {
        fmt.Println(err)
    }
    if err := f.Close(); err != nil {
        fmt.Println(err) 
    }
    // }
 }