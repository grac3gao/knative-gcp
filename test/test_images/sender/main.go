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
	//"database/sql"
	"fmt"
	"time"

	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	//"time"
	//
	//"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	"github.com/kelseyhightower/envconfig"
)

type envConfig struct {
	DBName string `envconfig:"DB_NAME" required:"true"`
	DBUser string `envconfig:"DB_USER" required:"true"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		panic(fmt.Sprintf("Failed to process env var: %s", err))
	}

	//panic(env.dbName)

	cfg := mysql.Cfg("gracegao-knative-testing:us-central1:sqlinstance-mysql", env.DBUser, "")
	cfg.DBName = env.DBName
	cfg.ParseTime = true

	const timeout = 10 * time.Second
	cfg.Timeout = timeout
	cfg.ReadTimeout = timeout
	cfg.WriteTimeout = timeout

	db, err := mysql.DialCfg(cfg)
	if err != nil {
		panic("couldn't dial: " + err.Error())
	}
	defer db.Close()
	var now time.Time
	fmt.Println(db.QueryRow("SELECT NOW()").Scan(&now))
	fmt.Println(db.Stats())
	fmt.Println(now)
	//_, err := sql.Open("mysql", env.dbUser + ":@/" + env.dbName)
	//if err != nil {
	//	panic(err)
	//}
	//defer db.Close()
	//fmt.Sprintf("%v", db.Stats())
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
//
//func main() {
//
//	var env envConfig
//	if err := envconfig.Process("", &env); err != nil {
//		panic(fmt.Sprintf("Failed to process env var: %s", err))
//	}
//
//	brokerURL := env.BrokerURLEnvVar
//	needRetry := (env.RetryEnvVar == "true")
//
//	ceClient, err := kncloudevents.NewDefaultClient(brokerURL)
//	if err != nil {
//		fmt.Printf("Unable to create ceClient: %s ", err)
//	}
//
//	// If needRetry is true, repeat sending Event with exponential backoff when there are some specific errors.
//	// In e2e test, sync problems could cause 404 and 5XX error, retrying those would help reduce flakiness.
//	success := true
//	span, err := sendEvent(ceClient, needRetry)
//	if !cev2.IsACK(err) {
//		success = false
//		fmt.Printf("failed to send event: %v", err)
//	}
//
//	if err := writeTerminationMessage(map[string]interface{}{
//		"success": success,
//		"traceid": span.SpanContext().TraceID.String(),
//	}); err != nil {
//		fmt.Printf("failed to write termination message, %s.\n", err)
//	}
//}
//
//func sendEvent(ceClient cev2.Client, needRetry bool) (span *trace.Span, err error) {
//	send := func() error {
//		ctx := cev2.WithEncodingBinary(context.Background())
//		ctx, span = trace.StartSpan(ctx, "sender", trace.WithSampler(trace.AlwaysSample()))
//		defer span.End()
//		result := ceClient.Send(ctx, sampleCloudEvent())
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
//func sampleCloudEvent() cev2.Event {
//	event := cev2.NewEvent(cev2.VersionV1)
//	event.SetID(lib.E2ESampleEventID)
//	event.SetType(lib.E2ESampleEventType)
//	event.SetSource(lib.E2ESampleEventSource)
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
