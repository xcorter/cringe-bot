package task

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xcorter/cringe-bot/src/joke"
	"github.com/xcorter/cringe-bot/src/repository"
	"log"
	"strconv"
)

type Tasks struct {
	storage repository.Storage
}

func (t *Tasks) GetUpdates(bot tgbotapi.BotAPI) {
	log.Printf("get updates")
	lastUpdateId := t.storage.GetLastUpdateId()
	log.Printf("last updateId: " + strconv.Itoa(lastUpdateId))
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
	joke, err := joke.GetJoke()
	if err != nil {
		log.Fatal(err)
		return
	}
	ids := t.storage.GetChatIds()
	for _, id := range ids {
		msg := tgbotapi.NewMessage(id, joke.Joke)
		bot.Send(msg)
	}
}

func NewTasks(storage repository.Storage) Tasks {
	return Tasks{storage: storage}
}
