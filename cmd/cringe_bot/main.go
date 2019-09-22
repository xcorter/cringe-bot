package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xcorter/cringe-bot/src/repository"
	"github.com/xcorter/cringe-bot/src/task"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"strconv"
)

func main() {
	fmt.Println("run!")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("TELEGRAM_API_KEY")
	client := getClient()
	bot, err := tgbotapi.NewBotAPIWithClient(apiKey, &client)
	checkErr(err)
	dbPath, err := os.Getwd()
	checkErr(err)
	dbPath = dbPath + "/storage.db"
	fmt.Println(dbPath)
	storage := repository.NewStorage(dbPath)
	tasks := task.NewTasks(storage)
	gocron.Every(1).Minute().Do(tasks.GetUpdates, *bot)
	times := getTimes()
	fmt.Println(times)
	for _, time := range times {
		gocron.Every(1).Day().At(time).Do(tasks.SendJokes, *bot)
	}
	<-gocron.Start()
	fmt.Print("bye bye")
}

func getClient() (result http.Client) {
	//creating the proxyURL
	proxyStr := "http://51.158.108.135:8811"
	proxyURL, err := url2.Parse(proxyStr)
	if err != nil {
		log.Println(err)
	}

	//adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	//adding the Transport object to the http Client
	return http.Client{
		Transport: transport,
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getTimes() []string {
	var result []string
	for i := 0; i < 24; i++ {
		hour := strconv.Itoa(i)
		if i < 10 {
			hour = "0" + hour
		}
		hour = hour + ":00"
		result = append(result, hour)
	}
	return result
}
