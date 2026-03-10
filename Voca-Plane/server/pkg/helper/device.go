package helper

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/mileusna/useragent"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type ServerInfo struct {
	ServerName     string `json:"server_host_name"`
	ServerOS       string `json:"server_os"`
	ServerRAM      string `json:"server_ram_usage"`
	ServerCPU      string `json:"server_cpu_model"`
	InternetStatus string `json:"internet_status"`
	Suggestion     string `json:"suggestion"`
}

type ClientInfo struct {
	UserIP      string `json:"user_ip"`
	OS          string `json:"os"`
	Browser     string `json:"browser"`
	Version     string `json:"browser_version"`
	DeviceType  string `json:"device_type"`
	DeviceBrand string `json:"device_brand"`
	IsBot       bool   `json:"is_bot"`
	RawAgent    string `json:"user_agent"`
}

type DeviceDetails struct {
	Server ServerInfo `json:"server"`
	Client ClientInfo `json:"client"`
}

// GetServerInfo mengambil metrik resource server saat ini
func GetServerInfo() ServerInfo {
	v, _ := mem.VirtualMemory()
	ramUsage := fmt.Sprintf("%.2f GB Total (terpakai: %.1f%%)", float64(v.Total)/1e9, v.UsedPercent)

	cpuInfo, _ := cpu.Info()
	cpuModel := "Unknown"
	if len(cpuInfo) > 0 {
		cpuModel = cpuInfo[0].ModelName
	}

	start := time.Now()
	client := http.Client{Timeout: 2 * time.Second}

	var internetStatus string
	var latency time.Duration

	resp, err := client.Get("http://www.google.com")
	if err != nil {
		internetStatus = "Disconnected"
	} else {
		defer resp.Body.Close()
		latency = time.Since(start)
		internetStatus = fmt.Sprintf("Stable (%v)", latency.Truncate(time.Millisecond))
	}

	hostname, _ := os.Hostname()

	suggestion := "Server operating normally."
	if v.UsedPercent > 85 {
		suggestion = "High RAM usage detected."
	} else if latency > 500*time.Millisecond {
		suggestion = "High internet latency."
	}

	return ServerInfo{
		ServerName:     hostname,
		ServerOS:       runtime.GOOS,
		ServerRAM:      ramUsage,
		ServerCPU:      cpuModel,
		InternetStatus: internetStatus,
		Suggestion:     suggestion,
	}
}

// GetClientInfo memproses User Agent dan IP client
func GetClientInfo(clientIP string, uaString string) ClientInfo {
	// Menggunakan library mileusna/useragent untuk parsing yang lebih akurat
	ua := useragent.Parse(uaString)

	deviceType := "desktop"
	if ua.Mobile {
		deviceType = "mobile"
	} else if ua.Tablet {
		deviceType = "tablet"
	}

	// Deteksi Brand Fallback
	// Jika library gagal mendeteksi brand dari UA yang tidak standar,
	// kita tambahkan logika manual sederhana sebagai cadangan.
	brand := ua.Device
	if brand == "" {
		brand = manualBrandDetection(uaString)
	}

	return ClientInfo{
		UserIP:      clientIP,
		OS:          ua.OS,
		Browser:     ua.Name,
		Version:     ua.Version,
		DeviceType:  deviceType,
		DeviceBrand: brand,
		IsBot:       ua.Bot,
		RawAgent:    uaString,
	}
}

// manualBrandDetection adalah cadangan jika library utama tidak mengenali model tertentu
func manualBrandDetection(ua string) string {
	ua = strings.ToLower(ua)
	switch {
	case strings.Contains(ua, "infinix") || strings.Contains(ua, "x68") || strings.Contains(ua, "x67"):
		return "Infinix"
	case strings.Contains(ua, "tecno") || strings.Contains(ua, "kj5"):
		return "Tecno"
	case strings.Contains(ua, "itel"):
		return "Itel"
	case strings.Contains(ua, "samsung") || strings.Contains(ua, "sm-"):
		return "Samsung"
	case strings.Contains(ua, "oppo") || strings.Contains(ua, "cph"):
		return "Oppo"
	case strings.Contains(ua, "xiaomi") || strings.Contains(ua, "redmi") || strings.Contains(ua, "poco"):
		return "Xiaomi"
	default:
		return "Generic Device"
	}
}