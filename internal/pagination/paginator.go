package pagination

import (
	"errors"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/huandu/go-sqlbuilder"
	"github.com/shayesteh1hs/DrAppointment/internal/api"
)

type Result[T api.PageEntityDTO] struct {
	Items      []T     `json:"items"`
	TotalCount int     `json:"total_count"`
	Previous   *string `json:"previous"`
	Next       *string `json:"next"`
}

type Params interface {
	Validate() error
}

type Paginator[T api.PageEntityDTO] interface {
	Paginate(sb *sqlbuilder.SelectBuilder) error
	CreatePaginationResult(items []T, totalCount int) (*Result[T], error)
	BindQueryParam(c *gin.Context) error
}

func validateBaseURL(baseURL string) error {
	if baseURL == "" {
		return errors.New("base url is required")
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return errors.New("invalid base url format")
	}

	if parsedURL.Host == "" {
		return errors.New("base url must contain a valid host")
	}

	return nil
}
