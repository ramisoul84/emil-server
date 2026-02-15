package location

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
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

func GetFullClientInfo(ip string) (string, string) {
	if ip == "" {
		return "Unknown", "Unknown"
	}

	resp, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		return "Unknown", "Unknown"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Unknown", "Unknown"
	}

	var result IPAPIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "Unknown", "Unknown"
	}

	if result.Status != "success" {
		return "Unknown", "Unknown"
	}

	return result.Country, result.City
}

func GetRealClientIP(c *fiber.Ctx) string {
	// Check X-Real-IP header (set by nginx)
	if realIP := c.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Check X-Forwarded-For header
	if xff := c.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check CF-Connecting-IP (Cloudflare)
	if cfIP := c.Get("CF-Connecting-IP"); cfIP != "" {
		return cfIP
	}

	// Fallback to Fiber's IP method
	return c.IP()
}
