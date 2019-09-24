package task

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xcorter/cringe-bot/src/joke"
	"github.com/xcorter/cringe-bot/src/repository"
	"log"
)

type Tasks struct {
	storage repository.Storage
}

func (t *Tasks) GetUpdates(bot tgbotapi.BotAPI) {
	lastUpdateId := t.storage.GetLastUpdateId()
	updateConfig := tgbotapi.NewUpdate(lastUpdateId)
	updateConfig.Timeout = 60
	updates, err := bot.GetUpdates(updateConfig)
	if err != nil {
		log.Fatal(err)
	}

	for _, update := range updates {
		if update.Message == nil {
			t.storage.SaveUpdateId(update.UpdateID)
			continue
		}
		t.storage.SaveChatId(update.Message.Chat.ID)
		t.storage.SaveUpdateId(update.UpdateID)
	}
}

func (t *Tasks) SendJokes(bot tgbotapi.BotAPI) {
	log.Println("send joke")
	jokeObject, err := joke.GetJoke()
	fmt.Println("%+v", jokeObject)
	if err != nil {
		log.Fatal(err)
		return
	}
	ids := t.storage.GetChatIds()
	fmt.Println(ids)
	for _, id := range ids {
		msg := tgbotapi.NewMessage(id, jokeObject.Joke)
		bot.Send(msg)
	}
}

func NewTasks(storage repository.Storage) Tasks {
	return Tasks{storage: storage}
}
