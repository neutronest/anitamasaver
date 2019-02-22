package main

import (
    "fmt"
    "os"
	"github.com/anaskhan96/soup"
    "github.com/neutronest/anitamasaver/anitama"
)

func main(){
    fmt.Println(anitama.ANITAMA_ROOT_URL)
    fmt.Println("Hello World")

    response, err := soup.Get(anitama.ANITAMA_CHANNEL_URL)
	if err != nil {
		os.Exit(1)
    }
    doc := soup.HTMLParse(response)
    areaArticleChannelDiv := doc.Find("div", "id", "area-article-channel")
    innerDiv := areaArticleChannelDiv.Find("div", "class", "inner")
    links := innerDiv.FindAll("a")
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
        articleType := infoParagraph.Find("span", "class", "channel")

        if title.Error != nil {
            continue
        }
        if subtitle.Error != nil {
            continue
        }
        if author.Error != nil {
            continue
        }
        fmt.Println(title.Text())
        fmt.Println(subtitle.Text())
        fmt.Println(author.Text())
        fmt.Println(articleType.Text())
        fmt.Println(anitama.ANITAMA_ROOT_URL + articleLink)
        fmt.Println()

    }


    
 }