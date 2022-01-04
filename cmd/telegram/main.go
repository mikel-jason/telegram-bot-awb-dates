package main

import (
	"flag"
	"fmt"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sarcaustech/go-telegram-awb/pkg/awb"
	"log"
	"strings"
	"time"
)

func main() {
	botToken := flag.String("botToken", "", "Telegram Bot Token (see t.me/botfather)")

	buildingNo := flag.Int("buildingNo", -1, "Building number")
	streetCode := flag.Int("streetCode", -1, "Code for target street")

	flag.Parse()

	if len(*botToken) == 0 {
		log.Fatalln("No bot token provided")
	}

	for _, value := range []int{*buildingNo, *streetCode} {
		if value < 0 {
			log.Fatalln("Insufficient streetCode and/or buildingNo provided")
		}
	}

	bot, err := telegram.NewBotAPI(*botToken)
	if err != nil {
		log.Fatalln("Cannot connect to bot", err)
	}
	log.Println("Connected to bot", bot.Self.UserName)
	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	now := time.Now()
	nextWeek := time.Now()
	nextWeek.Add(7 * 24 * time.Hour) // during last week of a month, start fetching next month's dates

	fetcher := awb.Fetcher{
		BuildingNo: *buildingNo,
		StreetCode: *streetCode,
		StartMonth: int(now.Month()),
		StartYear:  now.Year(),
		EndMonth:   int(nextWeek.Month()),
		EndYear:    nextWeek.Year(),
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		dates, err := fetcher.Fetch()
		if err != nil {
			log.Println("Unable to fetch AWB data", err)
			_, _ = bot.Send(telegram.NewMessage(update.Message.Chat.ID, "Could not fetch data from AWB"))
			return
		}
		cnt := 0
		var messageParts []string
		messageParts = append(messageParts, "Next AWB dates:")
		messageParts = append(messageParts, "") // to add empty row

		for _, date := range dates {
			if cnt > 4 {
				break
			}
			var awbType string
			switch date.Type {
			case "grey":
				awbType = "Restmüll (schwarz)"
			case "blue":
				awbType = "Altpapier (blau)"
			case "wertstoff":
				awbType = "Wertstoff (gelb)"
			default:
				awbType = "Unbekannt"
			}

			if err != nil {
				log.Println("Could not parse date", err)
				_, _ = bot.Send(telegram.NewMessage(update.Message.Chat.ID, "Could not process request"))
				return
			}
			messageParts = append(messageParts, fmt.Sprintf("%s\t%s", date.Date.Format("Mon 02. Jan 2006"), awbType))
			cnt++
		}

		_, err = bot.Send(telegram.NewMessage(update.Message.Chat.ID, strings.Join(messageParts, "\n")))
		if err != nil {
			log.Println("Could not send message", err)
			_, _ = bot.Send(telegram.NewMessage(update.Message.Chat.ID, "Could not process request"))
			return
		}
	}
}
