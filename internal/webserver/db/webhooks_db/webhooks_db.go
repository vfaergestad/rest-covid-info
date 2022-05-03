package webhooks_db

// Webhooks_db handles all communication with the webhooks' collection in the database.

import (
	"assignment-2/internal/webserver/constants"
	"assignment-2/internal/webserver/db"
	"assignment-2/internal/webserver/structs"
	"assignment-2/internal/webserver/utility/hash_util"
	"errors"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

// collection is the name of the collection in the database.
var collection = "webhooks"

// GetAllWebhooks returns all webhooks in the database.
func GetAllWebhooks() ([]structs.Webhook, error) {
	var resultSlice []structs.Webhook
	dbEmpty := true

	iter := db.GetClient().Collection(collection).Documents(db.GetContext())

	for {
		// Gets the next item in the collection
		doc, err := iter.Next()

		// Stops the loop if there is no more elements
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		dbEmpty = false

		// Converts the document to a webhook struct.
		var webhook structs.Webhook
		err = doc.DataTo(&webhook)
		if err != nil {
			return nil, err
		}

		resultSlice = append(resultSlice, webhook)

	}

	// Return an empty map if the database is empty
	if dbEmpty {
		return []structs.Webhook{}, errors.New(constants.WebhookDBIsEmpty)
	}

	return resultSlice, nil
}

// GetWebhook returns the webhook with the given id.
func GetWebhook(webhookId string) (structs.Webhook, error) {
	res := db.GetClient().Collection(collection).Doc(webhookId)
	doc, err := res.Get(db.GetContext())
	if err != nil {
		return structs.Webhook{}, errors.New(constants.WebhookNotFoundError)
	}

	// Converts the document to a webhook struct.
	var webhook structs.Webhook
	err = doc.DataTo(&webhook)
	if err != nil {
		return structs.Webhook{}, err
	}
	return webhook, nil
}

// AddWebhook adds a webhook with the given url, country, and calls, to the database.
// The count is set to 0, since it is a new webhook.
func AddWebhook(url string, country string, calls int) (string, error) {
	country = strings.Title(country)

	// Hashes the webhook information to create a unique id.
	webhookId := hash_util.HashWebhook(url, country, calls)

	res := db.GetClient().Collection(collection).Doc(webhookId)
	_, err := res.Get(db.GetContext())

	if status.Code(err) != codes.NotFound {
		return "", errors.New(constants.WebhookAlreadyExistingError)
	}

	_, err = res.Set(db.GetContext(), map[string]interface{}{
		"webhookId": webhookId,
		"url":       url,
		"country":   country,
		"calls":     calls,
		"count":     0,
	})
	if err != nil {
		return "", err
	} else {
		return webhookId, nil
	}
}

// UpdateWebhook updates the webhook with the given url, country, calls, and count.
func UpdateWebhook(url string, country string, calls int, count int) (string, error) {
	country = strings.Title(country)

	// Hashes the webhook information to get the webhook id.
	webhookId := hash_util.HashWebhook(url, country, calls)

	res := db.GetClient().Collection(collection).Doc(webhookId)
	_, err := res.Get(db.GetContext())

	if status.Code(err) == codes.NotFound {
		return "", errors.New(constants.WebhookNotFoundError)
	}

	_, err = res.Set(db.GetContext(), map[string]interface{}{
		"webhookId": webhookId,
		"url":       url,
		"country":   country,
		"calls":     calls,
		"count":     count,
	})
	if err != nil {
		return "", err
	} else {
		return webhookId, nil
	}
}

// DeleteWebhook deletes the webhook with the given id.
func DeleteWebhook(webhookId string) error {
	res := db.GetClient().Collection(collection).Doc(webhookId)
	_, err := res.Get(db.GetContext())
	if err != nil {
		return errors.New(constants.WebhookNotFoundError)
	}

	_, err = res.Delete(db.GetContext())
	if err != nil {
		return err
	} else {
		return nil
	}
}

// GetDBSize returns the number of webhooks in the database.
func GetDBSize() (int, error) {
	res, err := db.GetClient().Collection(collection).Documents(db.GetContext()).GetAll()
	if err != nil {
		return -1, err
	}
	return len(res), nil
}

func SetUpTestDatabase() []string {
	collection = "testing"
	webhooks, _ := GetAllWebhooks()
	for _, webhook := range webhooks {
		_ = DeleteWebhook(webhook.WebhookId)
	}

	id1, _ := AddWebhook("https://example.com", "Norway", 1)
	id2, _ := AddWebhook("https://example2.com", "Sweden", 2)
	id3, _ := AddWebhook("https://example3.com", "Denmark", 3)

	return []string{id1, id2, id3}

}
