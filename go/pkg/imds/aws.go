package imds

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type AwsInstance struct {
	AccountID               string      `json:"accountId"`
	Architecture            string      `json:"architecture"`
	AvailabilityZone        string      `json:"availabilityZone"`
	BillingProducts         interface{} `json:"billingProducts"`
	DevpayProductCodes      interface{} `json:"devpayProductCodes"`
	MarketplaceProductCodes interface{} `json:"marketplaceProductCodes"`
	ImageID                 string      `json:"imageId"`
	InstanceID              string      `json:"instanceId"`
	InstanceType            string      `json:"instanceType"`
	KernelID                interface{} `json:"kernelId"`
	PendingTime             time.Time   `json:"pendingTime"`
	PrivateIP               string      `json:"privateIp"`
	RamdiskID               interface{} `json:"ramdiskId"`
	Region                  string      `json:"region"`
	Version                 string      `json:"version"`
}

func GetAwsInstance(ctx context.Context) (*AwsInstance, error) {
	var pt = &http.Transport{Proxy: nil}
	client := http.Client{Transport: otelhttp.NewTransport(pt), Timeout: 1 * time.Second}
	req, _ := http.NewRequestWithContext(
		ctx,
		"GET",
		"http://169.254.169.254/latest/dynamic/instance-identity/document",
		nil,
	)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	ai := AwsInstance{}
	err = json.Unmarshal(body, &ai)
	if err != nil {
		return nil, err
	}
	return &ai, nil
}
