/*
Copyright 2020 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package authcheck provides utilities to check authentication configuration for data plane resources.
// File authcheck contains functions to run customized checks inside of a Pod.
package authcheck

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	nethttp "net/http"
	"regexp"

	"golang.org/x/oauth2/google"
)

const (
	// Resource is used as the path to get email from metadata server.
	resource = "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/email"
	// Scope is used as the scope to get token from default credential.
	scope = "https://www.googleapis.com/auth/cloud-platform"
)

// Regex for a valid google service account email.
var email_regexp = regexp.MustCompile(`^[a-z][a-z0-9-]{5,29}@[a-z][a-z0-9-]{5,29}.iam.gserviceaccount.com$`)

func AuthenticationCheck(ctx context.Context, authType AuthTypes, response nethttp.ResponseWriter) {
	var err error
	if authType == Secret {
		err = AuthenticationCheckForSecret(ctx, authType)
	} else if authType == WorkloadIdentityGSA {
		err = AuthenticationCheckForWorkloadIdentityGSA(ctx, authType)
	}

	if err != nil {
		b, err := json.Marshal(map[string]interface{}{
			"error": err.Error(),
		})
		if err != nil {
			response.WriteHeader(nethttp.StatusOK)
			return
		}
		ioutil.WriteFile("/dev/termination-log", b, 0644)
		response.WriteHeader(nethttp.StatusUnauthorized)
		return
	}
	response.WriteHeader(nethttp.StatusOK)
}

// AuthCheckForSecret performs the authentication check for Pod in secret mode.
func AuthenticationCheckForSecret(ctx context.Context, authType AuthTypes) error {
	cred, err := google.FindDefaultCredentials(ctx, scope)
	if err != nil {
		return nil
	}

	s, err := cred.TokenSource.Token()
	if err != nil {
		return nil
	}
	if !s.Valid() {
		return fmt.Errorf("using %s mode, token is not valid, "+
			"probably due to the key stored in the Kubernetes Secret is expired or revoked", authType)
	}
	return nil
}

// AuthenticationCheckForWorkloadIdentityGSA performs the authentication check for Pod in workload-identity-gsa mode.
func AuthenticationCheckForWorkloadIdentityGSA(ctx context.Context, authType AuthTypes) error {
	req, err := http.NewRequest("GET", resource, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("Metadata-Flavor", "Google")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	email, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	// Check if email is a valid google service account email.
	if match := email_regexp.FindStringSubmatch(string(email)); len(match) == 0 {
		// The format of google service account email is service-account-name@project-id.iam.gserviceaccount.com

		// Service account name must be between 6 and 30 characters (inclusive),
		// must begin with a lowercase letter, and consist of lowercase alphanumeric characters that can be separated by hyphens.

		// Project IDs must start with a lowercase letter and can have lowercase ASCII letters, digits or hyphens,
		// must be between 6 and 30 characters.
		return fmt.Errorf("using %s mode, Pod is not authenticated with a valid Google Service Account email, "+
			"probably due to mis-configuration between the Kubernetes Service Account and the Google Service Account", authType)
	}
	return nil
}
