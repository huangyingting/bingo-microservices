package system

import (
	"context"
	"os"
	"runtime"
	"strings"

	bsv1 "bingo/api/bs/v1"
	"bingo/pkg/imds"

	"bingo/app/bs/internal/utils"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// FileExists checks if a file or directory exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info != nil
}

// GetStats get system statistic
func GetStats(ctx context.Context) (*bsv1.StatsResponse, error) {
	host, err := host.Info()
	if err != nil {
		return nil, err
	}

	cpu, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	mem, _ := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	/*
		partitions, _ := disk.Partitions(false)
		for _, partition := range partitions {
			Debug("%s %s %s %v\n", partition.Device, partition.Mountpoint, partition.Fstype, partition.Opts)
			usage, _ := disk.Usage(partition.Mountpoint)
			Debug("%d %d\n", usage.Used, usage.Total)
		}
	*/

	local_ip_arr, _ := utils.GetLocalIP()
	local_ip := strings.Join(local_ip_arr, ",")
	external_ip, _ := utils.GetExternalIP(utils.GoogleDns)
	if external_ip == "" {
		external_ip, _ = utils.GetExternalIP(utils.CloudflareDns)
	}

	// potential fix, read /proc/self/cgroup to tell if running from containerd
	// by looking at 0::/system.slice/containerd.service
	response := &bsv1.StatsResponse{
		Hostname:        host.Hostname,
		Os:              host.OS,
		Platform:        host.Platform,
		PlatformVersion: host.PlatformVersion,
		CpuModelName:    cpu[0].ModelName,
		CpuCores:        int32(runtime.NumCPU()),
		CpuCacheSize:    cpu[0].CacheSize,
		CpuMhz:          cpu[0].Mhz,
		GoArch:          runtime.GOARCH,
		GoVersion:       runtime.Version(),
		MemTotal:        mem.Total,
		ExternalIp:      external_ip,
		LocalIp:         local_ip,
		IsDocker:        FileExists("/.dockerenv"),
		IsKubernetes:    FileExists("/var/run/secrets/kubernetes.io"),
	}

	// get cloud instance metadata
	azure, _ := imds.GetAzureInstance(ctx)
	if azure != nil {
		response.Environment = "Azure"
		response.Location = azure.Compute.Location
		response.Zone = azure.Compute.Zone
		response.Name = azure.Compute.Name
		response.Size = azure.Compute.VMSize
		return response, nil
	}

	aws, _ := imds.GetAwsInstance(ctx)
	if aws != nil {
		response.Environment = "Aws"
		response.Location = aws.Region
		response.Zone = aws.AvailabilityZone
		response.Name = aws.InstanceID
		response.Size = aws.InstanceType
		return response, nil
	}

	gcp, _ := imds.GetGcpInstance(ctx)
	if aws != nil {
		response.Environment = "Google"
		response.Location = ""
		response.Zone = gcp.Zone
		response.Name = gcp.Name
		response.Size = gcp.MachineType
		return response, nil
	}

	return response, nil
}
