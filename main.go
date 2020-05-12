package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	TOKEN    = ""
	KURS_URL = "https://www.nbrb.by/api/exrates/rates?periodicity=0"
	BYN      = "BYN"
	USD      = "USD"
)

type RateData struct {
	Abbreviation string  `json:"Cur_Abbreviation"`
	Rate         float32 `json:"Cur_OfficialRate"`
	Name         string  `json:"Cur_Name"`
}

func main() {
	bot, err := tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Command() == BYN {
			resp, err := http.Get(KURS_URL)
			if err != nil {
				log.Panic(err)
			}
			defer resp.Body.Close()

			var rates []RateData
			if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
				log.Panic(err)
			}

			for _, rate := range rates {
				if rate.Abbreviation == USD {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n%v\n", rate.Name, rate.Rate))
					if _, err := bot.Send(msg); err != nil {
						log.Panic(err)
					}
					break
				}
			}
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите комманду /BYN чтобы узнать текущий курс доллара.")
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
