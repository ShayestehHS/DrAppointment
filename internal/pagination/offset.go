package pagination

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/huandu/go-sqlbuilder"
	"github.com/shayesteh1hs/DrAppointment/internal/api"
	"github.com/shayesteh1hs/DrAppointment/internal/utils"
)

var _ Params = (*LimitOffsetParams)(nil)

type LimitOffsetParams struct {
	Page         int        `form:"page,default=1" binding:"min=1"`
	Limit        int        `form:"limit,default=10" binding:"min=1,max=100"`
	BaseURL      string     `form:"-"`
	ClientParams url.Values `form:"-"`
}

func (p *LimitOffsetParams) Validate() error {
	if err := validateBaseURL(p.BaseURL); err != nil {
		return err
	}

	return nil
}

func (p *LimitOffsetParams) GetOffset() int {
	return p.Limit * (p.Page - 1)
}

func (p *LimitOffsetParams) BindQueryParam(c *gin.Context) error {
	if err := c.ShouldBindQuery(p); err != nil {
		return fmt.Errorf("invalid cursor parameters: %w", err)
	}

	p.ClientParams = c.Request.URL.Query()
	p.ClientParams.Del("page")
	p.ClientParams.Del("limit")

	p.BaseURL = utils.BuildBaseURL(c)
	return p.Validate()
}

type LimitOffsetPaginator[T api.PageEntityDTO] struct {
	params LimitOffsetParams
}

func NewLimitOffsetPaginator[T api.PageEntityDTO](params LimitOffsetParams) *LimitOffsetPaginator[T] {
	return &LimitOffsetPaginator[T]{params: params}
}

func NewOffsetPaginator[T api.PageEntityDTO]() *LimitOffsetPaginator[T] {
	return &LimitOffsetPaginator[T]{params: LimitOffsetParams{}}
}

func (p *LimitOffsetPaginator[T]) GetParams() LimitOffsetParams {
	return p.params
}

func (p *LimitOffsetPaginator[T]) Paginate(sb *sqlbuilder.SelectBuilder) error {
	sb.Limit(p.params.Limit)
	sb.Offset(p.params.GetOffset())
	return nil
}

func (p *LimitOffsetPaginator[T]) BindQueryParam(c *gin.Context) error {
	return p.params.BindQueryParam(c)
}

func (p *LimitOffsetPaginator[T]) CreatePaginationResult(items []T, totalCount int) (*Result[T], error) {
	result := &Result[T]{
		Items:      items,
		TotalCount: totalCount,
	}

	totalPages := (totalCount + p.params.Limit - 1) / p.params.Limit

	if p.params.Page > 1 {
		prevPage := p.params.Page - 1
		prevURL, err := p.buildURL(prevPage)
		if err != nil {
			return nil, err
		}
		result.Previous = &prevURL
	}

	if p.params.Page < totalPages {
		nextPage := p.params.Page + 1
		nextURL, err := p.buildURL(nextPage)
		if err != nil {
			return nil, err
		}
		result.Next = &nextURL
	}

	return result, nil
}

func (p *LimitOffsetPaginator[T]) buildURL(page int) (string, error) {
	if p.params.BaseURL == "" {
		return "", errors.New("base url is required")
	}

	// Parse existing URL to preserve query parameters
	u, err := url.Parse(p.params.BaseURL)
	if err != nil {
		log.Printf("failed to parse base URL: %v", err)
		return "", errors.New("failed to parse base url")
	}

	params := p.params.ClientParams
	if params == nil {
		// If ClientParams is nil, extract query parameters from the parsed URL
		params = u.Query()
	}
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("limit", fmt.Sprintf("%d", p.params.Limit))

	u.RawQuery = params.Encode()
	return u.String(), nil
}
