package anitama

type ArticleMetaData struct {
	
	Title string
	SubTitle string
	Author string
	Type string
	Date string
}

type Article struct {

	MetaData ArticleMetaData
	Content	string
	RawContent string
}

