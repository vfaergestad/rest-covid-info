package webhooks_db

import (
	"assignment-2/internal/webserver/constants"
	"assignment-2/internal/webserver/db"
	"assignment-2/internal/webserver/structs"
	"assignment-2/internal/webserver/utility/hash_util"
	"errors"
	"testing"
)

var webhooksIds []string

func TestMain(m *testing.M) {
	constants.SetTestServiceAccountLocation()

	err := db.InitializeFirestore()
	if err != nil {
		panic(err)
	}

	defer func() {
		err = db.CloseFirestore()
		if err != nil {
			panic(err)
		}
	}()

	webhooksIds = SetUpTestDatabase()
	m.Run()
}

func TestGetWebhook(t *testing.T) {
	type args struct {
		webhookId string
	}
	tests := []struct {
		name             string
		args             args
		expectedResponse structs.Webhook
		expectedError    error
	}{
		{
			name: "Test Get Valid Webhook",
			args: args{
				webhookId: webhooksIds[2],
			},
			expectedResponse: structs.Webhook{
				WebhookId: webhooksIds[2],
				Url:       "https://example3.com",
				Country:   "Denmark",
				Calls:     3,
			},
		},
		{
			name: "Test Get Invalid Webhook",
			args: args{
				webhookId: "invalid_id",
			},
			expectedResponse: structs.Webhook{},
			expectedError:    errors.New(constants.WebhookNotFoundError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualResponse, actualError := GetWebhook(tt.args.webhookId)
			if actualError != nil && actualError.Error() != tt.expectedError.Error() {
				t.Errorf("GetWebhook() actualError = %v, want %v", actualError, tt.expectedError)
				return
			}
			if actualResponse != tt.expectedResponse {
				t.Errorf("GetWebhook() actualResponse = %v, want %v", actualResponse, tt.expectedResponse)
			}
		})
	}
}

func TestDeleteWebhook(t *testing.T) {
	type args struct {
		webhookId string
	}
	tests := []struct {
		name          string
		args          args
		expectedError error
	}{
		{
			name: "Test Delete Valid Webhook",
			args: args{
				webhookId: webhooksIds[2],
			},
			expectedError: nil,
		},
		{
			name: "Test Delete Invalid Webhook",
			args: args{
				webhookId: "invalid_id",
			},
			expectedError: errors.New(constants.WebhookNotFoundError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualError := DeleteWebhook(tt.args.webhookId)
			if actualError != nil && actualError.Error() != tt.expectedError.Error() {
				t.Errorf("DeleteWebhook() actualError = %v, want %v", actualError, tt.expectedError)
				return
			}
		})
	}

	t.Cleanup(func() {
		SetUpTestDatabase()
	})
}

func TestAddWebhook(t *testing.T) {
	type args struct {
		url     string
		country string
		calls   int
	}
	tests := []struct {
		name          string
		args          args
		expectedId    string
		expectedError error
	}{
		{
			name: "Test Add Valid Webhook",
			args: args{
				url:     "https://example4.com",
				country: "France",
				calls:   1,
			},
			expectedId:    hash_util.HashWebhook("https://example4.com", "France", 1),
			expectedError: nil,
		},
		{
			name: "Test Add Existing Webhook",
			args: args{
				url:     "https://example3.com",
				country: "Denmark",
				calls:   3,
			},
			expectedId:    "",
			expectedError: errors.New(constants.WebhookAlreadyExistingError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualId, actualError := AddWebhook(tt.args.url, tt.args.country, tt.args.calls)
			if actualError != nil && actualError.Error() != tt.expectedError.Error() {
				t.Errorf("AddWebhook() actualError = %v, want %v", actualError, tt.expectedError)
				return
			}
			if actualId != tt.expectedId {
				t.Errorf("AddWebhook() actualId = %v, want %v", actualId, tt.expectedId)
			}
		})
	}
	t.Cleanup(func() {
		SetUpTestDatabase()
	})
}

func TestUpdateWebhook(t *testing.T) {
	type args struct {
		webhookId string
		url       string
		country   string
		calls     int
		count     int
	}
	tests := []struct {
		name          string
		args          args
		expectedId    string
		expectedError error
	}{
		{
			name: "Test Update Valid Webhook",
			args: args{
				url:     "https://example3.com",
				country: "Denmark",
				calls:   3,
				count:   1,
			},
			expectedId:    hash_util.HashWebhook("https://example3.com", "Denmark", 3),
			expectedError: nil,
		},
		{
			name: "Test Update Invalid Webhook",
			args: args{
				url:     "https://example8.com",
				country: "Norway545",
				calls:   12,
				count:   12,
			},
			expectedId:    "",
			expectedError: errors.New(constants.WebhookNotFoundError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualId, actualError := UpdateWebhook(tt.args.url, tt.args.country, tt.args.calls, tt.args.count)
			if actualError != nil && actualError.Error() != tt.expectedError.Error() {
				t.Errorf("UpdateWebhook() actualError = %v, want %v", actualError, tt.expectedError)
				return
			}
			if actualId != tt.expectedId {
				t.Errorf("UpdateWebhook() actualId = %v, want %v", actualId, tt.expectedId)
			}
		})
	}
	t.Cleanup(func() {
		SetUpTestDatabase()
	})
}

func TestGetDBSize(t *testing.T) {
	tests := []struct {
		name          string
		expectedCount int
		expectedError error
	}{
		{
			name:          "Test Get DB Size",
			expectedCount: 3,
			expectedError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualCount, actualError := GetDBSize()
			if actualError != nil && actualError.Error() != tt.expectedError.Error() {
				t.Errorf("GetDBSize() actualError = %v, want %v", actualError, tt.expectedError)
				return
			}
			if actualCount != tt.expectedCount {
				t.Errorf("GetDBSize() actualCount = %v, want %v", actualCount, tt.expectedCount)
			}
		})
	}
}
