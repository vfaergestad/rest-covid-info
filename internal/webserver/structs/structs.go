package structs

import "time"

type CasesResponse struct {
	Country    string  `json:"country"`
	Date       string  `json:"date"`
	Confirmed  int     `json:"confirmed"`
	Recovered  int     `json:"recovered"`
	Deaths     int     `json:"deaths"`
	GrowthRate float64 `json:"growth_rate"`
}

type CountryCacheEntry struct {
	AlphaCode   string
	CountryName string
	Time        time.Time
}

type PolicyResponse struct {
	CountryCode string    `json:"country_code"`
	Scope       string    `json:"scope"`
	Stringency  float64   `json:"stringency"`
	Policies    int       `json:"policies"`
	Time        time.Time `json:"-"`
}

type Webhook struct {
	WebhookId  string `json:"webhook_id"`
	Url        string `json:"url,omitempty"`
	Country    string `json:"country"`
	Calls      int    `json:"calls"`
	StartCount int    `json:"-"`
}

type CasesApiResponse struct {
	Data Data `json:"data"`
}

type Data struct {
	Country CountryStruct `json:"country"`
}

type CountryStruct struct {
	Name       string           `json:"name"`
	MostRecent MostRecentStruct `json:"mostRecent"`
}

type MostRecentStruct struct {
	Date       string  `json:"date"`
	Confirmed  int     `json:"confirmed"`
	Recovered  int     `json:"recovered"`
	Deaths     int     `json:"deaths"`
	GrowthRate float64 `json:"growthRate"`
}

type PolicyApiResponse struct {
	PolicyActions  []PolicyAction         `json:"policyActions"`
	StringencyData map[string]interface{} `json:"stringencyData"`
}

type PolicyAction struct {
	PolicyTypeCode string `json:"policy_type_code"`
}
