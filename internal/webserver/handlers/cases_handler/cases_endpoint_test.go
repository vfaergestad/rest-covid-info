package cases_handler

import (
	"assignment-2/internal/webserver/cache/country_cache"
	"assignment-2/internal/webserver/constants"
	"assignment-2/internal/webserver/db"
	"assignment-2/internal/webserver/mock_apis"
	"assignment-2/internal/webserver/structs"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var casesEndpoint *httptest.Server

func TestMain(m *testing.M) {
	casesMock := httptest.NewServer(http.HandlerFunc(mock_apis.HandlerCases))
	defer casesMock.Close()
	countriesMock := httptest.NewServer(http.HandlerFunc(mock_apis.HandlerCountries))
	defer countriesMock.Close()

	constants.SetTestCasesApiUrl(casesMock.URL)
	constants.SetTestCountryAPIUrl(countriesMock.URL)
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

	err = country_cache.InitCache()
	if err != nil {
		panic(err)
	}

	casesEndpoint = httptest.NewServer(http.HandlerFunc(HandlerCases))
	defer casesEndpoint.Close()

	m.Run()
}

func TestRequestsFromCases(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name               string
		args               args
		method             string
		expectedStatusCode int
		expectedResponse   structs.CasesResponse
	}{
		{
			name: "Valid Get Request With Common Name",
			args: args{
				url: casesEndpoint.URL + "/corona/v1/cases/Norway",
			},
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
			expectedResponse: structs.CasesResponse{
				Country:    "Norway",
				Date:       "2020-01-01",
				Confirmed:  1,
				Recovered:  2,
				Deaths:     3,
				GrowthRate: 4,
			},
		},
		{
			name: "Valid Get Request With Country Code",
			args: args{
				url: casesEndpoint.URL + "/corona/v1/cases/NOR",
			},
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
			expectedResponse: structs.CasesResponse{
				Country:    "Norway",
				Date:       "2020-01-01",
				Confirmed:  1,
				Recovered:  2,
				Deaths:     3,
				GrowthRate: 4,
			},
		},
		{
			name: "Invalid Get Request With To Long Path",
			args: args{
				url: casesEndpoint.URL + "/corona/v1/cases/NOR/Test",
			},
			method:             http.MethodGet,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: structs.CasesResponse{
				Country:    "",
				Date:       "",
				Confirmed:  0,
				Recovered:  0,
				Deaths:     0,
				GrowthRate: 0,
			},
		},
		{
			name: "Invalid Get Request With To Numbers in Country",
			args: args{
				url: casesEndpoint.URL + "/corona/v1/cases/NOR123",
			},
			method:             http.MethodGet,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: structs.CasesResponse{
				Country:    "",
				Date:       "",
				Confirmed:  0,
				Recovered:  0,
				Deaths:     0,
				GrowthRate: 0,
			},
		},
		{
			name: "Invalid Get Request With Invalid Country Code",
			args: args{
				url: casesEndpoint.URL + "/corona/v1/cases/Veg",
			},
			method:             http.MethodGet,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: structs.CasesResponse{
				Country:    "",
				Date:       "",
				Confirmed:  0,
				Recovered:  0,
				Deaths:     0,
				GrowthRate: 0,
			},
		},
		{
			name: "Invalid Get Request With Valid Country In But Invalid In Case API",
			args: args{
				url: casesEndpoint.URL + "/corona/v1/cases/Taiwan",
			},
			method:             http.MethodGet,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: structs.CasesResponse{
				Country:    "",
				Date:       "",
				Confirmed:  0,
				Recovered:  0,
				Deaths:     0,
				GrowthRate: 0,
			},
		},
		{
			name: "Invalid Post Request",
			args: args{
				url: casesEndpoint.URL + "/corona/v1/cases/Norway",
			},
			method:             http.MethodPost,
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{}
			req, err := http.NewRequest(tt.method, tt.args.url, nil)
			if err != nil {
				t.Errorf("Error creating request: %v", err)
			}
			res, err := client.Do(req)
			if err != nil {
				t.Errorf("Error making %s request: %s", tt.method, err)
			}
			if res.StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected %d response, got %d", tt.expectedStatusCode, res.StatusCode)
			}

			if tt.expectedResponse != (structs.CasesResponse{}) {
				var actual structs.CasesResponse
				_ = json.NewDecoder(res.Body).Decode(&actual)

				if actual != tt.expectedResponse {
					t.Errorf("Expected %v, got %v", tt.expectedResponse, actual)
				}
			}
		})
	}
}

func TestInvalidPostRequestFromCases(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name               string
		args               args
		expectedStatusCode int
	}{
		{
			name: "Invalid Post Request",
			args: args{
				url: casesEndpoint.URL + "/corona/v1/cases/Norway",
			},
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{}
			res, err := client.Post(tt.args.url, "application/json", nil)
			if err != nil {
				t.Errorf("Error making POST request: %s", err)
			}
			if res.StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected %d response, got %d", tt.expectedStatusCode, res.StatusCode)
			}
		})
	}

}
