package imds

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type AzureInstance struct {
	Compute struct {
		AzEnvironment    string `json:"azEnvironment,omitempty"`
		ExtendedLocation struct {
			Type string `json:"type,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"extendedLocation,omitempty"`
		EvictionPolicy             string `json:"evictionPolicy,omitempty"`
		IsHostCompatibilityLayerVM string `json:"isHostCompatibilityLayerVm,omitempty"`
		LicenseType                string `json:"licenseType,omitempty"`
		Location                   string `json:"location,omitempty"`
		Name                       string `json:"name,omitempty"`
		Offer                      string `json:"offer,omitempty"`
		OsProfile                  struct {
			AdminUsername                 string `json:"adminUsername,omitempty"`
			ComputerName                  string `json:"computerName,omitempty"`
			DisablePasswordAuthentication string `json:"disablePasswordAuthentication,omitempty"`
		} `json:"osProfile,omitempty"`
		OsType           string `json:"osType,omitempty"`
		PlacementGroupID string `json:"placementGroupId,omitempty"`
		Plan             struct {
			Name      string `json:"name,omitempty"`
			Product   string `json:"product,omitempty"`
			Publisher string `json:"publisher,omitempty"`
		} `json:"plan,omitempty"`
		PlatformFaultDomain    string `json:"platformFaultDomain,omitempty"`
		PlatformSubFaultDomain string `json:"platformSubFaultDomain,omitempty"`
		PlatformUpdateDomain   string `json:"platformUpdateDomain,omitempty"`
		Priority               string `json:"priority,omitempty"`
		PublicKeys             []struct {
			KeyData string `json:"keyData,omitempty"`
			Path    string `json:"path,omitempty"`
		} `json:"publicKeys,omitempty"`
		Publisher         string `json:"publisher,omitempty"`
		ResourceGroupName string `json:"resourceGroupName,omitempty"`
		ResourceID        string `json:"resourceId,omitempty"`
		SecurityProfile   struct {
			SecureBootEnabled string `json:"secureBootEnabled,omitempty"`
			VirtualTpmEnabled string `json:"virtualTpmEnabled,omitempty"`
		} `json:"securityProfile,omitempty"`
		Sku            string `json:"sku,omitempty"`
		StorageProfile struct {
			DataDisks []struct {
				BytesPerSecondThrottle string `json:"bytesPerSecondThrottle,omitempty"`
				Caching                string `json:"caching,omitempty"`
				CreateOption           string `json:"createOption,omitempty"`
				DiskCapacityBytes      string `json:"diskCapacityBytes,omitempty"`
				DiskSizeGB             string `json:"diskSizeGB,omitempty"`
				Image                  struct {
					URI string `json:"uri,omitempty"`
				} `json:"image,omitempty"`
				IsSharedDisk string `json:"isSharedDisk,omitempty"`
				IsUltraDisk  string `json:"isUltraDisk,omitempty"`
				Lun          string `json:"lun,omitempty"`
				ManagedDisk  struct {
					ID                 string `json:"id,omitempty"`
					StorageAccountType string `json:"storageAccountType,omitempty"`
				} `json:"managedDisk,omitempty"`
				Name                 string `json:"name,omitempty"`
				OpsPerSecondThrottle string `json:"opsPerSecondThrottle,omitempty"`
				Vhd                  struct {
					URI string `json:"uri,omitempty"`
				} `json:"vhd,omitempty"`
				WriteAcceleratorEnabled string `json:"writeAcceleratorEnabled,omitempty"`
			} `json:"dataDisks,omitempty"`
			ImageReference struct {
				ID        string `json:"id,omitempty"`
				Offer     string `json:"offer,omitempty"`
				Publisher string `json:"publisher,omitempty"`
				Sku       string `json:"sku,omitempty"`
				Version   string `json:"version,omitempty"`
			} `json:"imageReference,omitempty"`
			OsDisk struct {
				Caching          string `json:"caching,omitempty"`
				CreateOption     string `json:"createOption,omitempty"`
				DiskSizeGB       string `json:"diskSizeGB,omitempty"`
				DiffDiskSettings struct {
					Option string `json:"option,omitempty"`
				} `json:"diffDiskSettings,omitempty"`
				EncryptionSettings struct {
					Enabled string `json:"enabled,omitempty"`
				} `json:"encryptionSettings,omitempty"`
				Image struct {
					URI string `json:"uri,omitempty"`
				} `json:"image,omitempty"`
				ManagedDisk struct {
					ID                 string `json:"id,omitempty"`
					StorageAccountType string `json:"storageAccountType,omitempty"`
				} `json:"managedDisk,omitempty"`
				Name   string `json:"name,omitempty"`
				OsType string `json:"osType,omitempty"`
				Vhd    struct {
					URI string `json:"uri,omitempty"`
				} `json:"vhd,omitempty"`
				WriteAcceleratorEnabled string `json:"writeAcceleratorEnabled,omitempty"`
			} `json:"osDisk,omitempty"`
			ResourceDisk struct {
				Size string `json:"size,omitempty"`
			} `json:"resourceDisk,omitempty"`
		} `json:"storageProfile,omitempty"`
		SubscriptionID         string `json:"subscriptionId,omitempty"`
		Tags                   string `json:"tags,omitempty"`
		Version                string `json:"version,omitempty"`
		VirtualMachineScaleSet struct {
			ID string `json:"id,omitempty"`
		} `json:"virtualMachineScaleSet,omitempty"`
		VMID           string `json:"vmId,omitempty"`
		VMScaleSetName string `json:"vmScaleSetName,omitempty"`
		VMSize         string `json:"vmSize,omitempty"`
		Zone           string `json:"zone,omitempty"`
	} `json:"compute,omitempty"`
	Network struct {
		Interface []struct {
			Ipv4 struct {
				IPAddress []struct {
					PrivateIPAddress string `json:"privateIpAddress,omitempty"`
					PublicIPAddress  string `json:"publicIpAddress,omitempty"`
				} `json:"ipAddress,omitempty"`
				Subnet []struct {
					Address string `json:"address,omitempty"`
					Prefix  string `json:"prefix,omitempty"`
				} `json:"subnet,omitempty"`
			} `json:"ipv4,omitempty"`
			Ipv6 struct {
				IPAddress []interface{} `json:"ipAddress,omitempty"`
			} `json:"ipv6,omitempty"`
			MacAddress string `json:"macAddress,omitempty"`
		} `json:"interface,omitempty"`
	} `json:"network,omitempty"`
}

func GetAzureInstance(ctx context.Context) (*AzureInstance, error) {
	var pt = &http.Transport{Proxy: nil}
	client := http.Client{Transport: otelhttp.NewTransport(pt), Timeout: 1 * time.Second}
	req, _ := http.NewRequestWithContext(
		ctx,
		"GET",
		"http://169.254.169.254/metadata/instance",
		nil,
	)
	req.Header.Add("Metadata", "True")
	q := req.URL.Query()
	q.Add("format", "json")
	q.Add("api-version", "2021-02-01")
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	ai := AzureInstance{}
	err = json.Unmarshal(body, &ai)
	if err != nil {
		return nil, err
	}
	return &ai, nil
}
