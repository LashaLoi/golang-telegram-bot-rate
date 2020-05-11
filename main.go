package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	TOKEN    = "1085475952:AAEI0EEPj60dl5D6cHGy8fP1ct-KXB1pRRY"
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

			var kurses []RateData
			if err := json.NewDecoder(resp.Body).Decode(&kurses); err != nil {
				log.Panic(err)
			}

			for _, kurs := range kurses {
				if kurs.Abbreviation == USD {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n%v\n", kurs.Name, kurs.Rate))
					if _, err := bot.Send(msg); err != nil {
						log.Panic(err)
					}
					break
				}
			}
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter /BYN to check current kursExchange")
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}

	}
}
