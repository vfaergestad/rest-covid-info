package webhooks

import (
	"assignment-2/internal/webserver/cache/country_cache"
	"assignment-2/internal/webserver/constants"
	"assignment-2/internal/webserver/db/webhooks_db"
	"assignment-2/internal/webserver/structs"
	"assignment-2/internal/webserver/utility/request"
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func Invoke(country string) {
	// Checks if given alpha-code is a three letter string.
	match, err := regexp.MatchString(constants.AlphaCodeRegex, country)
	if err != nil {
		log.Println(err)
	} else if match {
		country, err = country_cache.GetCountry(country)
		if err != nil {
			log.Println(err)
		}
	}

	err = checkAndInvokeWebhooks(country)
	if err != nil {
		log.Println(err)
	}
}

func checkAndInvokeWebhooks(country string) error {
	webhooks, err := webhooks_db.GetAllWebhooks()
	if err != nil {
		return err
	}
	for _, w := range webhooks {

		if w.Country == country {
			w.Count = w.Count + 1
			if w.Count >= w.Calls {
				go callWebhook(w)
				w.Count = 0
			}
			_, err = webhooks_db.UpdateWebhook(w.Url, w.Country, w.Calls, w.Count)
		}
	}
	return nil
}

func callWebhook(webhook structs.Webhook) {
	body := webhook
	body.Url = ""
	result, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
	}
	res, err := request.PostRequest(webhook.Url, strings.NewReader(string(result)))
	if err != nil {
		log.Println(err)
		return
	}

	// Read the response
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Something is wrong with invocation response. Error:", err)
		return
	}

	log.Println("Webhook invoked: " + webhook.Url + ". Received status code " + strconv.Itoa(res.StatusCode) +
		" and response: " + string(response))
}
