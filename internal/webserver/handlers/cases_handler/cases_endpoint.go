package cases_handler

// Cases_endpoint handles the incoming requests to the /corona/v1/cases endpoint.

import (
	"assignment-2/internal/webserver/api_requests/cases_api"
	"assignment-2/internal/webserver/cache/country_cache"
	"assignment-2/internal/webserver/constants"
	"assignment-2/internal/webserver/utility/encode_struct"
	"assignment-2/internal/webserver/utility/logging"
	"assignment-2/internal/webserver/webhooks"
	"net/http"
	"path"
	"regexp"
	"strings"
)

// HandlerCases handler entrypoint.
func HandlerCases(w http.ResponseWriter, r *http.Request) {
	logging.LogRequest(r)
	// Check if the method is GET.
	if r.Method != http.MethodGet {
		http.Error(w, "Method is not supported. Currently only GET are supported.", http.StatusMethodNotAllowed)
		return
	}

	// Retrieves the search-terms from the path
	cleanPath := path.Clean(r.URL.Path)
	pathList := strings.Split(cleanPath, "/")

	// Check if the given path has a valid length.
	if len(pathList) != 5 {
		http.Error(w, constants.GetBadCasesRequestError().Error(), http.StatusBadRequest)
		return
	}

	// Defines the different queries in the url
	countryQuery := path.Base(cleanPath)

	// Check if country contains any numbers
	match, err := regexp.MatchString(constants.NoNumbersRegex, countryQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !match {
		http.Error(w, "country cannot contain numbers", http.StatusBadRequest)
		return
	}

	// Checks if the country is an alpha-3 code.
	match, err = regexp.MatchString(constants.AlphaCodeRegex, countryQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if match {
		// If it is an alpha-3 code, it gets the country name from the cache.
		countryQuery = strings.ToUpper(countryQuery)
		countryQuery, err = country_cache.GetCountry(countryQuery)
		if err != nil {
			switch err.Error() {
			case constants.CountryNotFoundError:
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		// Sanitizes the country name.
		if countryQuery != "US" && countryQuery != "us" {
			countryQuery = strings.Title(strings.ToLower(countryQuery))
		} else {
			countryQuery = strings.ToUpper(countryQuery)
		}
	}

	// Invokes the webhook package to see if any webhooks needs to be counted up or needs to be invoked.
	go webhooks.Invoke(countryQuery)

	// Retrieves the cases with the given country data from the api.
	casesResponseStruct, err := cases_api.GetResponseStruct(countryQuery)
	if err != nil {
		if constants.IsBadRequestError(err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if constants.IsNotFoundError(err) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Encodes the response struct to json and writes it to the response.
	err = encode_struct.EncodeStruct(w, casesResponseStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
