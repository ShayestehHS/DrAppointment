package utils

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/shayesteh1hs/DrAppointment/internal/entity"
)

func GetImageBaseURL() *url.URL {
	baseURL := GetEnv("IMAGE_BASE_URL", "")

	// Ensure baseURL ends with a trailing slash
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	baseImageURL, err := url.Parse(baseURL)
	if err != nil {
		return &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("localhost:%d", GetEnvInt("SERVER_PORT", 8080)),
			Path:   "/",
		}
	}
	return baseImageURL
}

func GetFullImageURL[T entity.Image](image T) string {
	baseURL := GetImageBaseURL()
	return baseURL.ResolveReference(image.GetPath()).String()
}
