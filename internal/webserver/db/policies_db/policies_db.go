package policies_db

// policies_db handles all communications to the policy collection in the database.

import (
	"assignment-2/internal/webserver/constants"
	"assignment-2/internal/webserver/db"
	"assignment-2/internal/webserver/structs"
	"assignment-2/internal/webserver/utility/hash_util"
	"errors"
	"time"
)

// collection is the name of the collection in the database.
const collection = "policies"

var testMode = false

// GetPolicy returns the policy from the given country and scope as a
func GetPolicy(country string, scope string) (structs.PolicyResponse, error) {
	if testMode {
		return structs.PolicyResponse{}, errors.New(constants.TestModeActiveError)
	}
	// Hashes the country and scope to get the policy ID.
	id := hash_util.HashPolicy(country, scope)

	// Gets the policy from the database.
	res := db.GetClient().Collection(collection).Doc(id)
	doc, err := res.Get(db.GetContext())
	if err != nil {
		return structs.PolicyResponse{}, errors.New(constants.PolicyNotInCacheError)
	}

	// Converts the document to a policy struct.
	var policy structs.PolicyResponse
	err = doc.DataTo(&policy)
	if err != nil {
		return structs.PolicyResponse{}, err
	}

	// Checks if the policy is expired.
	if time.Since(policy.Time).Hours() > constants.CacheExpire {
		// Removes the policy from the database if it is expired.
		_ = DeletePolicy(country, scope)
		return structs.PolicyResponse{}, errors.New(constants.ExpiredCacheEntry)
	}

	return policy, nil

}

// AddPolicy adds the given policy to the database.
func AddPolicy(policy structs.PolicyResponse) error {
	if testMode {
		return errors.New(constants.TestModeActiveError)
	}

	// Hashes the country and scope to get the policy ID.
	id := hash_util.HashPolicy(policy.CountryCode, policy.Scope)
	_, err := db.GetClient().Collection(collection).Doc(id).Set(db.GetContext(), map[string]interface{}{
		"countryCode": policy.CountryCode,
		"scope":       policy.Scope,
		"stringency":  policy.Stringency,
		"policies":    policy.Policies,
		"time":        time.Now(),
	})
	if err != nil {
		return err
	} else {
		return nil
	}
}

// DeletePolicy deletes the policy corresponding with the given country and scope from the database.
func DeletePolicy(country string, scope string) error {
	if testMode {
		return errors.New(constants.TestModeActiveError)
	}

	// Hashes the country and scope to get the policy ID.
	id := hash_util.HashPolicy(country, scope)

	_, err := db.GetClient().Collection(collection).Doc(id).Delete(db.GetContext())
	if err != nil {
		return err
	} else {
		return nil
	}
}

// SetTestMode sets the database to not respond to any requests, for testing purposes.
func SetTestMode() {
	testMode = true
}
