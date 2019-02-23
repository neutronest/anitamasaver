package anitama

type ArticleMetaData struct {
	
	Id	string
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

