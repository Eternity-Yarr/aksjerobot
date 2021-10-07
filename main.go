package main

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
)

func main() {

	token := "-"

	tr := &http.Transport{
		MaxIdleConnsPerHost: 5,
	}

	client := &http.Client{
		Transport: tr,
	}

	bot, err := tgbotapi.NewBotAPIWithClient(token, "https://api.telegram.org/bot%s/%s", client)

	if err != nil {
		panic(err)
		return
	}
	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.InlineQuery == nil {
			continue
		}
		if update.InlineQuery.From.ID != 89886125 {
			answerError(update.InlineQuery.ID, bot, "You ain't", "my master bitch!")
			continue
		}
		donwloadPng(bot, update.InlineQuery.ID, update.InlineQuery.Query)
	}
}

func answerError(queryId string, bot *tgbotapi.BotAPI, title string, message string) {
	article := tgbotapi.NewInlineQueryResultArticle(queryId, title, message)
	ic := tgbotapi.InlineConfig{
		InlineQueryID: queryId,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       []interface{}{article},
	}

	res, err := bot.AnswerInlineQuery(ic)
	if err != nil {
		fmt.Println(res)
		fmt.Println(err)
	}
}

func pngToJpeg(pngByte []byte) []byte {
	i, _ := png.Decode(bytes.NewReader(pngByte))
	var opt jpeg.Options
	opt.Quality = 80

	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, i, &opt)
	if err != nil {
		fmt.Println(err)
	}
	return buf.Bytes()
}

func donwloadPng(bot *tgbotapi.BotAPI, queryId string, ticker string) {
	url := "https://finviz.com/chart.ashx?s=l&ta=1&t=" + ticker
	pngB, err := DownloadFile(url)

	if err != nil {
		answerError(queryId, bot, "Error", "I fked up")
		return
	}

	fb := tgbotapi.FileBytes{
		"ticker",
		pngB,
	}
	c := tgbotapi.NewPhotoUpload(int64(89886125), fb)
	m, err := bot.Send(c)

	if err != nil {
		fmt.Println(m)
		fmt.Println(err)
		answerError(queryId, bot, "Error", "photo upload")
		return
	}

	sizes := *m.Photo
	biggest := sizes[0]
	for _, s := range sizes {
		if biggest.FileSize < s.FileSize {
			biggest = s
		}
	}
	image := tgbotapi.NewInlineQueryResultCachedPhoto("1", biggest.FileID)

	ans := tgbotapi.InlineConfig{
		InlineQueryID: queryId,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       []interface{}{image},
	}
	res, err := bot.AnswerInlineQuery(ans)

	if err != nil {
		fmt.Println(res)
		fmt.Println(err)
	}
}

func DownloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, err
}
