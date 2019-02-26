# AnitamaSaver

## Update

Today we are glad to hear that the Anitama will keep serving. Thanks everyoue.

As the situation, the project will be deprecated. Any developed features later can only be regard as Golang practice. 

---

Anitama is one of the professional ACG criticism websites in China. Recently, we were regreted to hear that the website will stop operations for some reasons, and it will be closed therefore we cannot visit the articles any more. In order to store the articles, the knowledge and thoughts for both authors and prudent readers, we build AnitamaSaver, a temporal project to save all the historical articles automatically.

## Roadmap

* [x] Save the article metadatas from channel page as json (Title, subtitle, author, categories)
* [x] Save each article raw-html-data as json (sync)
* [ ] Save article data asynchronously
* [ ] save the series article data
* [ ] Save the media data (mainly images) for each article and maintain the image-article mapping
* [ ] Parse the raw html article data for better represent rendering

## Install & Run

First make sure that you have prepare the golang environment.

Then just:
```
sh run.sh
```

This scipt will install dependencies and begin the saving task automatically.Each article metadata and content will be downloaded as json file in output/ folder. Notice that now we support sync download style only and it will cost you several hours. Asynchronous task is in the plan.

The Image data and series article saving task is still to be done.

## Other

If you cannot install libraries by the "golang/x/net" failed mistake, you can download this lib manually by "git clone https://github.com/golang/net.git" in the correct path. It depends on your $GOPATH.