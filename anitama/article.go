package anitama

type ArticleMetaData struct {
	
	Title string
	SubTitle string
	Author string
	Category string
	Date string
	Url string
}

type Article struct {

	MetaData ArticleMetaData
	Content	string
	RawContent string
}

