package main

import (
	"io/ioutil"
	"os"
	"fmt"
	"strings"
	"encoding/json"
	"strconv"
	"github.com/neutronest/anitamasaver/anitama"
	"github.com/bradfitz/slice"
)

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


}