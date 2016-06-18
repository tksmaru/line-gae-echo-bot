package main

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
	"google.golang.org/appengine/urlfetch"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const (
	CALLBACK_URI = "/callback"
	TEXT_URI     = "/text"
	TEXT_MESSAGE_QUEUE = "textqueue"
)

func init() {
	http.HandleFunc(CALLBACK_URI, callbackHandler)
	http.HandleFunc(TEXT_URI, textMessageHandler)
}

func textMessageHandler(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)

	to := r.FormValue("to")
	text := r.FormValue("text")

	log.Debugf(ctx, "to: %s, text: %s.", to, text)

	chi, _ := strconv.ParseInt(os.Getenv("CHANNEL_ID"), 10, 64)
	chs := os.Getenv("CHANNEL_SECRET")
	mid := os.Getenv("MID")

	client := urlfetch.Client(ctx)
	bot, err := linebot.NewClient(chi, chs, mid, linebot.WithHTTPClient(client))
	if err != nil {
		log.Errorf(ctx, "Client initialization failure. to: %s, text: %s, context: %s", to, text, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := bot.SendText([]string{to}, text)
	if err != nil {
		log.Errorf(ctx, "Message Send Failed. to: %s, text: %s, context: %s", to, text, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Infof(ctx, "Message send succeed. MessageID: %s, to: %s, text: %s", res.MessageID, to, text)
	w.WriteHeader(http.StatusOK)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)

	chi, err := strconv.ParseInt(os.Getenv("CHANNEL_ID"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	chs := os.Getenv("CHANNEL_SECRET")
	mid := os.Getenv("MID")

	bot, err := linebot.NewClient(chi, chs, mid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	received, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	for _, result := range received.Results {
		content := result.Content()
		if content != nil && content.IsMessage && content.ContentType == linebot.ContentTypeText {
			text, _ := content.TextContent()
			log.Debugf(ctx, "id: %s, text: %s, from: %s, to: %v", content.ID, text.Text, text.From, text.To)

			values := url.Values{}
			values.Set("to", text.From)
			values.Set("text", text.Text)
			t := taskqueue.NewPOSTTask(TEXT_URI, values)
			_, err = taskqueue.Add(ctx, t, TEXT_MESSAGE_QUEUE)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
	w.WriteHeader(http.StatusOK)
}
