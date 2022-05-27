package imds

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type GcpInstance struct {
	Attributes struct {
	} `json:"attributes"`
	CPUPlatform string `json:"cpuPlatform"`
	Description string `json:"description"`
	Disks       []struct {
		DeviceName string `json:"deviceName"`
		Index      int    `json:"index"`
		Mode       string `json:"mode"`
		Type       string `json:"type"`
	} `json:"disks"`
	Hostname string `json:"hostname"`
	ID       int64  `json:"id"`
	Image    string `json:"image"`
	Licenses []struct {
		ID string `json:"id"`
	} `json:"licenses"`
	MachineType       string `json:"machineType"`
	MaintenanceEvent  string `json:"maintenanceEvent"`
	Name              string `json:"name"`
	NetworkInterfaces []struct {
		AccessConfigs []struct {
			ExternalIP string `json:"externalIp"`
			Type       string `json:"type"`
		} `json:"accessConfigs"`
		DNSServers        []string      `json:"dnsServers"`
		ForwardedIps      []interface{} `json:"forwardedIps"`
		Gateway           string        `json:"gateway"`
		IP                string        `json:"ip"`
		IPAliases         []interface{} `json:"ipAliases"`
		Mac               string        `json:"mac"`
		Network           string        `json:"network"`
		Subnetmask        string        `json:"subnetmask"`
		TargetInstanceIps []interface{} `json:"targetInstanceIps"`
	} `json:"networkInterfaces"`
	Preempted        string `json:"preempted"`
	RemainingCPUTime int    `json:"remainingCpuTime"`
	Scheduling       struct {
		AutomaticRestart  string `json:"automaticRestart"`
		OnHostMaintenance string `json:"onHostMaintenance"`
		Preemptible       string `json:"preemptible"`
	} `json:"scheduling"`
	ServiceAccounts struct {
	} `json:"serviceAccounts"`
	Tags         []interface{} `json:"tags"`
	VirtualClock struct {
		DriftToken string `json:"driftToken"`
	} `json:"virtualClock"`
	Zone string `json:"zone"`
}

func GetGcpInstance(ctx context.Context) (*GcpInstance, error) {
	var pt = &http.Transport{Proxy: nil}
	client := http.Client{Transport: otelhttp.NewTransport(pt), Timeout: 1 * time.Second}
	req, _ := http.NewRequestWithContext(
		ctx,
		"GET",
		"http://metadata/computeMetadata/v1/instance/?recursive=true",
		nil,
	)
	req.Header.Add("Metadata-Flavor", "Google")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	gi := GcpInstance{}
	err = json.Unmarshal(body, &gi)
	if err != nil {
		return nil, err
	}
	return &gi, nil
}
