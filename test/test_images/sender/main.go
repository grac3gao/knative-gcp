/*
Copyright 2019 Google LLC

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

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2/google"
)

type envConfig struct {
	BrokerURLEnvVar string `envconfig:"BROKER_URL" required:"true"`
	RetryEnvVar     string `envconfig:"RETRY"`
}

//// defaultRetry represents that there will be 4 iterations.
//// The duration starts from 30s and is multiplied by factor 2.0 for each iteration.
//var defaultRetry = wait.Backoff{
//	Steps:    4,
//	Duration: 30 * time.Second,
//	Factor:   2.0,
//	// The sleep at each iteration is the duration plus an additional
//	// amount chosen uniformly at random from the interval between 0 and jitter*duration.
//	Jitter: 1.0,
//}

func main() {

	resource := "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/email"
	req, err := http.NewRequest("GET", resource, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Metadata-Flavor", "Google")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	s := string(data)

	fmt.Printf("Get the resp: %v ", s)

	scope := "https://www.googleapis.com/auth/cloud-platform"
	ctx := context.Background()
	cred, err := google.FindDefaultCredentials(ctx, scope)
	if err != nil {
		fmt.Printf("Unable to get token: %v ", err)
	} else {
		s, _ := cred.TokenSource.Token()
		fmt.Printf("Get the token: %v ", s.Valid())
	}

	//"https://accounts.google.com"
	//v := ClientOptions()
	//provider, err := oidc.NewProvider(ctx, "https://container.googleapis.com/v1/projects/gracegao-knative-testing/locations/us-west2/clusters/knative-5")
	//if err != nil {
	//	fmt.Printf("Unable to get provider: %v ", err)
	//}
	//fmt.Printf("%v ", provider)
	//userInfo, err := provider.UserInfo(ctx, cred.TokenSource)
	//if err != nil {
	//	fmt.Printf("Unable to get userInfo: %v ", err)
	//} else {
	//	fmt.Printf("Get the userInfo: %v", userInfo.EmailVerified)
	//}

	//req, err := http.NewRequest("POST", p.userInfoURL, nil)

	//ctx := context.Background()
	//cloudresourcemanagerService, err := cloudresourcemanager.NewService(ctx)
	//projectId, _ := utils.ProjectID("", metadataClient.NewDefaultMetadataClient())
	//
	//rb := &cloudresourcemanager.TestIamPermissionsRequest{
	//	Permissions: []string{"pubsub.topics.publish"},
	//}
	//resp, err := cloudresourcemanagerService.Projects.TestIamPermissions(projectId, rb).Context(ctx).Do()
	//if err != nil {
	//	fmt.Printf("failed to test iam permissions" + err.Error())
	//} else {
	//	fmt.Printf("%v", resp.Permissions)
	//}

	//var env envConfig
	//if err := envconfig.Process("", &env); err != nil {
	//	panic(fmt.Sprintf("Failed to process env var: %s", err))
	//}
	//
	//brokerURL := env.BrokerURLEnvVar
	//needRetry := (env.RetryEnvVar == "true")
	//
	//ceClient, err := kncloudevents.NewDefaultClient(brokerURL)
	//if err != nil {
	//	fmt.Printf("Unable to create ceClient: %s ", err)
	//}
	//
	//// If needRetry is true, repeat sending Event with exponential backoff when there are some specific errors.
	//// In e2e test, sync problems could cause 404 and 5XX error, retrying those would help reduce flakiness.
	//success := true
	//span, err := sendEvent(ceClient, needRetry)
	//if !cev2.IsACK(err) {
	//	success = false
	//	fmt.Printf("failed to send event: %v", err)
	//}
	//
	//if err := writeTerminationMessage(map[string]interface{}{
	//	"success": success,
	//	"traceid": span.SpanContext().TraceID.String(),
	//}); err != nil {
	//	fmt.Printf("failed to write termination message, %s.\n", err)
	//}
}

//
//func sendEvent(ceClient cev2.Client, needRetry bool) (span *trace.Span, err error) {
//	send := func() error {
//		ctx := cev2.WithEncodingBinary(context.Background())
//		ctx, span = trace.StartSpan(ctx, "sender", trace.WithSampler(trace.AlwaysSample()))
//		defer span.End()
//		result := ceClient.Send(ctx, dummyCloudEvent())
//		return result
//	}
//
//	if needRetry {
//		err = retry.OnError(defaultRetry, isRetryable, send)
//	} else {
//		err = send()
//	}
//	return
//}
//
//func dummyCloudEvent() cev2.Event {
//	event := cev2.NewEvent(cev2.VersionV1)
//	event.SetID(lib.E2EDummyEventID)
//	event.SetType(lib.E2EDummyEventType)
//	event.SetSource(lib.E2EDummyEventSource)
//	event.SetData(cev2.ApplicationJSON, `{"source": "sender!"}`)
//	return event
//}
//
//func writeTerminationMessage(result interface{}) error {
//	b, err := json.Marshal(result)
//	if err != nil {
//		return err
//	}
//	return ioutil.WriteFile("/dev/termination-log", b, 0644)
//}
//
//// isRetryable determines if the err is an error which is retryable
//func isRetryable(err error) bool {
//	var result *cehttp.Result
//	if protocol.ResultAs(err, &result) {
//		// Potentially retry when:
//		// - 404 Not Found
//		// - 500 Internal Server Error, it is currently for reducing flakiness caused by Workload Identity credential sync up.
//		// We should remove it after https://github.com/google/knative-gcp/issues/1058 lands, as 500 error may indicate bugs in our code.
//		// - 503 Service Unavailable (with or without Retry-After) (IGNORE Retry-After)
//		sc := result.StatusCode
//		if sc == 404 || sc == 500 || sc == 503 {
//			log.Printf("got error: %s, retry sending event. \n", result.Error())
//			return true
//		}
//	}
//	return false
//}

//func ClientOptions() []option.ClientOption {
//	return nil
//}
