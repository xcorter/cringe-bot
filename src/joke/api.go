package joke

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"html"
	"net/http"
	"strings"
)

func GetJoke() (result Response, err error) {
	url := "https://www.anekdot.ru/random/anekdot/"
	response, err := http.Get(url)
	if err != nil {
		return Response{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return Response{}, errors.New("status code error")
	}

	//var joke string
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return Response{}, errors.New("error when read document")
	}

	selection := doc.Find(".content .topicbox .text").First()
	htmlJoke, err := selection.Html()
	if err != nil {
		return Response{}, err
	}

	joke := processJoke(htmlJoke)

	return Response{Joke: joke}, nil
}

func processJoke(joke string) (result string) {
	joke = strings.Replace(joke, "<br/>", "\n", -1)
	joke = html.UnescapeString(joke)
	return joke
}

type Response struct {
	Joke string
}
