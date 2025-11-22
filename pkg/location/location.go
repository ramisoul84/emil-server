package location

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type IPAPIResponse struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	AS          string  `json:"as"`
	Query       string  `json:"query"`
	Message     string  `json:"message"`
}

func GetFullClientInfo(c echo.Context) (*IPAPIResponse, error) {
	ip := GetRealClientIP(c)

	fmt.Println("IP", ip)

	if ip == "" {
		return nil, fmt.Errorf("could not determine client IP")
	}

	resp, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		return nil, fmt.Errorf("failed to call IP-API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result IPAPIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("IP-API failed: %s", result.Message)
	}

	return &result, nil
}

func GetRealClientIP(c echo.Context) string {
	// First try X-Real-IP (set by nginx)
	if realIP := c.Request().Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Then try X-Forwarded-For
	if xff := c.Request().Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0]) // First IP is the original client
		}
	}

	// Fallback
	return c.RealIP()
}
