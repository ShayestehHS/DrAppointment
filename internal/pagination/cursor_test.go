package pagination

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/huandu/go-sqlbuilder"
	"github.com/stretchr/testify/assert"
)

// ============================================================================
// CursorParams Validation Tests
// ============================================================================

func TestCursorParams_Validate(t *testing.T) {
	validCursor := encodeCursor("123")

	tests := []struct {
		name     string
		params   CursorParams
		wantErr  bool
		errMsg   string
		expected CursorParams // expected state after validation
	}{
		{
			name: "valid params without cursor",
			params: CursorParams{
				Cursor:   "",
				Ordering: "asc",
				Limit:    10,
				BaseURL:  "http://example.com/api",
			},
			wantErr: false,
			expected: CursorParams{
				Cursor:   "",
				Ordering: "asc",
				Limit:    10,
				BaseURL:  "http://example.com/api",
			},
		},
		{
			name: "valid params with cursor",
			params: CursorParams{
				Cursor:   validCursor,
				Ordering: "desc",
				Limit:    20,
				BaseURL:  "http://example.com/api",
			},
			wantErr: false,
			expected: CursorParams{
				Cursor:   validCursor,
				Ordering: "desc",
				Limit:    20,
				BaseURL:  "http://example.com/api",
			},
		},
		{
			name: "invalid ordering",
			params: CursorParams{
				Cursor:   "",
				Ordering: "invalid",
				Limit:    10,
				BaseURL:  "http://example.com/api",
			},
			wantErr: true,
			errMsg:  "ordering must be either 'asc' or 'desc'",
		},
		{
			name: "invalid cursor",
			params: CursorParams{
				Cursor:   "invalid-cursor!!!",
				Ordering: "asc",
				Limit:    10,
				BaseURL:  "http://example.com/api",
			},
			wantErr: true,
			errMsg:  "invalid cursor",
		},
		{
			name: "case insensitive ordering - uppercase ASC",
			params: CursorParams{
				Cursor:   "",
				Ordering: "ASC",
				Limit:    10,
				BaseURL:  "http://example.com/api",
			},
			wantErr: false,
			expected: CursorParams{
				Cursor:   "",
				Ordering: "asc",
				Limit:    10,
				BaseURL:  "http://example.com/api",
			},
		},
		{
			name: "case insensitive ordering - uppercase DESC",
			params: CursorParams{
				Cursor:   "",
				Ordering: "DESC",
				Limit:    10,
				BaseURL:  "http://example.com/api",
			},
			wantErr: false,
			expected: CursorParams{
				Cursor:   "",
				Ordering: "desc",
				Limit:    10,
				BaseURL:  "http://example.com/api",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Ordering, tt.params.Ordering)
			}
		})
	}
}

// ============================================================================
// CursorPaginator Paginate Tests
// ============================================================================

func TestCursorPaginator_Paginate(t *testing.T) {
	tests := []struct {
		name          string
		params        CursorParams
		expectedSQL   []string
		expectedArgs  []interface{}
		expectedLimit int
	}{
		{
			name: "forward without cursor",
			params: CursorParams{
				Cursor:   "",
				Ordering: "asc",
				Limit:    10,
				BaseURL:  "http://example.com/api",
			},
			expectedSQL:   []string{"LIMIT $1", "ORDER BY id ASC"},
			expectedArgs:  []interface{}{11}, // limit + 1
			expectedLimit: 11,
		},
		{
			name: "forward with cursor",
			params: CursorParams{
				Cursor:   encodeCursor("123"),
				Ordering: "asc",
				Limit:    10,
				BaseURL:  "http://example.com/api",
			},
			expectedSQL:   []string{"WHERE id > $1", "LIMIT $2", "ORDER BY id ASC"},
			expectedArgs:  []interface{}{"123", 11},
			expectedLimit: 11,
		},
		{
			name: "backward with cursor",
			params: CursorParams{
				Cursor:   encodeCursor("123"),
				Ordering: "desc",
				Limit:    10,
				BaseURL:  "http://example.com/api",
			},
			expectedSQL:   []string{"WHERE id < $1", "LIMIT $2", "ORDER BY id DESC"},
			expectedArgs:  []interface{}{"123", 11},
			expectedLimit: 11,
		},
		{
			name: "single item limit",
			params: CursorParams{
				Cursor:   "",
				Ordering: "asc",
				Limit:    1,
				BaseURL:  "http://example.com/api",
			},
			expectedSQL:   []string{"LIMIT $1", "ORDER BY id ASC"},
			expectedArgs:  []interface{}{2}, // limit + 1
			expectedLimit: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate the params first
			err := tt.params.Validate()
			assert.NoError(t, err)

			paginator := NewCursorPaginator[mockEntity](tt.params)
			sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
			sb.Select("*").From("test")

			err = paginator.Paginate(sb)
			assert.NoError(t, err)

			// Check that expected SQL clauses are present
			sql, args := sb.Build()
			for _, expectedClause := range tt.expectedSQL {
				assert.Contains(t, sql, expectedClause)
			}

			// Check that the correct values are in the args
			assert.Equal(t, len(tt.expectedArgs), len(args))
			for i, expectedArg := range tt.expectedArgs {
				assert.Equal(t, expectedArg, args[i])
			}
		})
	}
}

// ============================================================================
// CursorPaginator CreatePaginationResult Tests
// ============================================================================

func TestCursorPaginator_CreatePaginationResult(t *testing.T) {
	baseURL := "http://example.com/api"

	tests := []struct {
		name               string
		params             CursorParams
		items              []mockEntity
		totalCount         int
		expectedItemsCount int
		expectPrevious     bool
		expectNext         bool
		expectedFirstID    string
		expectedLastID     string
	}{
		{
			name: "forward with more items",
			params: CursorParams{
				Cursor:   "",
				Ordering: "asc",
				Limit:    10,
				BaseURL:  baseURL,
			},
			items:              generateMockItems(11), // limit + 1 to simulate hasMore
			totalCount:         25,
			expectedItemsCount: 10,
			expectPrevious:     false, // no cursor, first page
			expectNext:         true,  // has more items
			expectedFirstID:    "1",
			expectedLastID:     "10",
		},
		{
			name: "backward with more items",
			params: CursorParams{
				Cursor:   encodeCursor("20"),
				Ordering: "desc",
				Limit:    10,
				BaseURL:  baseURL,
			},
			items:              generateMockItemsReverse(11, 19), // Items 19, 18, 17, ..., 9
			totalCount:         25,
			expectedItemsCount: 10,
			expectPrevious:     true, // backward pagination shows previous
			expectNext:         true, // backward pagination shows next
			expectedFirstID:    "10", // first item after reversal
			expectedLastID:     "19", // last item after reversal
		},
		{
			name: "no more items",
			params: CursorParams{
				Cursor:   "",
				Ordering: "asc",
				Limit:    10,
				BaseURL:  baseURL,
			},
			items:              generateMockItems(10), // exactly limit items
			totalCount:         10,
			expectedItemsCount: 10,
			expectPrevious:     false, // no cursor, first page
			expectNext:         false, // no more items
			expectedFirstID:    "1",
			expectedLastID:     "10",
		},
		{
			name: "empty items",
			params: CursorParams{
				Cursor:   "",
				Ordering: "asc",
				Limit:    10,
				BaseURL:  baseURL,
			},
			items:              []mockEntity{},
			totalCount:         0,
			expectedItemsCount: 0,
			expectPrevious:     false,
			expectNext:         false,
		},
		{
			name: "single item",
			params: CursorParams{
				Cursor:   "",
				Ordering: "asc",
				Limit:    10,
				BaseURL:  baseURL,
			},
			items:              generateMockItems(1),
			totalCount:         1,
			expectedItemsCount: 1,
			expectPrevious:     false,
			expectNext:         false,
			expectedFirstID:    "1",
			expectedLastID:     "1",
		},
		{
			name: "backward without cursor",
			params: CursorParams{
				Cursor:   "",
				Ordering: "desc",
				Limit:    10,
				BaseURL:  baseURL,
			},
			items:              generateMockItems(5),
			totalCount:         5,
			expectedItemsCount: 5,
			expectPrevious:     true, // backward pagination shows previous
			expectNext:         true, // backward pagination shows next
			expectedFirstID:    "5",
			expectedLastID:     "1",
		},
		{
			name: "forward with cursor and has more",
			params: CursorParams{
				Cursor:   encodeCursor("5"),
				Ordering: "asc",
				Limit:    10,
				BaseURL:  baseURL,
			},
			items:              generateMockItems(11), // limit + 1 to simulate hasMore
			totalCount:         25,
			expectedItemsCount: 10,
			expectPrevious:     true, // has cursor, so previous exists
			expectNext:         true, // has more items, so next exists
			expectedFirstID:    "1",
			expectedLastID:     "10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validationErr := tt.params.Validate()
			assert.NoError(t, validationErr, "Validation should succeed for valid params")

			paginator := NewCursorPaginator[mockEntity](tt.params)
			result, resultErr := paginator.CreatePaginationResult(tt.items, tt.totalCount)

			assert.NoError(t, resultErr, "CreatePaginationResult should not return an error for valid params")
			assert.Equal(t, tt.expectedItemsCount, len(result.Items))
			assert.Equal(t, tt.totalCount, result.TotalCount)

			if tt.expectPrevious {
				assert.NotNil(t, result.Previous, "Expected previous link to exist")
			} else {
				assert.Nil(t, result.Previous, "Expected previous link to be nil")
			}

			if tt.expectNext {
				assert.NotNil(t, result.Next, "Expected next link to exist")
			} else {
				assert.Nil(t, result.Next, "Expected next link to be nil")
			}

			// Check first and last item IDs if items exist
			if len(result.Items) > 0 && tt.expectedFirstID != "" {
				assert.Equal(t, tt.expectedFirstID, result.Items[0].GetID())
			}
			if len(result.Items) > 0 && tt.expectedLastID != "" {
				assert.Equal(t, tt.expectedLastID, result.Items[len(result.Items)-1].GetID())
			}

			// Verify cursor encoding in next URL if next exists
			if result.Next != nil {
				nextURL, err := url.Parse(*result.Next)
				assert.NoError(t, err)
				cursor := nextURL.Query().Get("cursor")
				if cursor != "" {
					decodedCursor, err := decodeCursor(cursor)
					assert.NoError(t, err)
					if len(result.Items) > 0 {
						assert.Equal(t, result.Items[len(result.Items)-1].GetID(), decodedCursor)
					}
				}
			}
		})
	}
}

// ============================================================================
// Helper Functions Tests
// ============================================================================

func TestEncodeCursor(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "string input",
			input:    "123",
			expected: base64.RawURLEncoding.EncodeToString([]byte("123")),
		},
		{
			name:     "integer input",
			input:    456,
			expected: base64.RawURLEncoding.EncodeToString([]byte("456")),
		},
		{
			name:     "empty string",
			input:    "",
			expected: base64.RawURLEncoding.EncodeToString([]byte("")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := encodeCursor(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDecodeCursor(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		shouldError bool
	}{
		{
			name:        "valid cursor - 123",
			input:       base64.RawURLEncoding.EncodeToString([]byte("123")),
			expected:    "123",
			shouldError: false,
		},
		{
			name:        "valid cursor - 456",
			input:       base64.RawURLEncoding.EncodeToString([]byte("456")),
			expected:    "456",
			shouldError: false,
		},
		{
			name:        "valid cursor - empty",
			input:       base64.RawURLEncoding.EncodeToString([]byte("")),
			expected:    "",
			shouldError: false,
		},
		{
			name:        "invalid base64",
			input:       "invalid-base64!!!",
			expected:    "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := decodeCursor(tt.input)
			if tt.shouldError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestReverseSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []mockEntity
		expected []mockEntity
	}{
		{
			name:     "three items",
			input:    []mockEntity{{ID: "1"}, {ID: "2"}, {ID: "3"}},
			expected: []mockEntity{{ID: "3"}, {ID: "2"}, {ID: "1"}},
		},
		{
			name:     "two items",
			input:    []mockEntity{{ID: "1"}, {ID: "2"}},
			expected: []mockEntity{{ID: "2"}, {ID: "1"}},
		},
		{
			name:     "single item",
			input:    []mockEntity{{ID: "1"}},
			expected: []mockEntity{{ID: "1"}},
		},
		{
			name:     "empty slice",
			input:    []mockEntity{},
			expected: []mockEntity{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to avoid modifying the original
			input := make([]mockEntity, len(tt.input))
			copy(input, tt.input)

			reverseSlice(input)
			assert.Equal(t, tt.expected, input)
		})
	}
}

// ============================================================================
// BindQueryParam Tests
// ============================================================================

func TestCursorParams_BindQueryParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestURL     string
		host           string
		expectedParams CursorParams
		expectError    bool
		errorMsg       string
	}{
		{
			name:       "valid params with all fields",
			requestURL: "/api/test?cursor=" + encodeCursor("123") + "&ordering=asc&limit=20&filter=test&extra=value",
			host:       "example.com",
			expectedParams: CursorParams{
				Cursor:   encodeCursor("123"),
				Ordering: "asc",
				Limit:    20,
				BaseURL:  "http://example.com/api/test",
			},
			expectError: false,
		},
		{
			name:       "params with defaults",
			requestURL: "/api/test",
			host:       "example.com",
			expectedParams: CursorParams{
				Cursor:   "",
				Ordering: "asc", // default
				Limit:    10,    // default
				BaseURL:  "http://example.com/api/test",
			},
			expectError: false,
		},
		{
			name:        "invalid limit - zero",
			requestURL:  "/api/test?limit=0",
			host:        "example.com",
			expectError: true,
			errorMsg:    "invalid cursor parameters",
		},
		{
			name:        "invalid limit - too high",
			requestURL:  "/api/test?limit=101",
			host:        "example.com",
			expectError: true,
			errorMsg:    "invalid cursor parameters",
		},
		{
			name:        "invalid cursor",
			requestURL:  "/api/test?cursor=invalid!!!",
			host:        "example.com",
			expectError: true,
			errorMsg:    "invalid cursor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, _ := http.NewRequest("GET", tt.requestURL, nil)
			req.Host = tt.host
			c.Request = req

			params := CursorParams{}
			err := params.BindQueryParam(c)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedParams.Cursor, params.Cursor)
				assert.Equal(t, tt.expectedParams.Ordering, params.Ordering)
				assert.Equal(t, tt.expectedParams.Limit, params.Limit)
				assert.Equal(t, tt.expectedParams.BaseURL, params.BaseURL)
				assert.NotNil(t, params.ClientParams)

				// Check that pagination params are removed from ClientParams
				assert.Empty(t, params.ClientParams.Get("cursor"))
				assert.Empty(t, params.ClientParams.Get("ordering"))
				assert.Empty(t, params.ClientParams.Get("limit"))

				// Check that other params are preserved
				if tt.requestURL == "/api/test?cursor="+encodeCursor("123")+"&ordering=asc&limit=20&filter=test&extra=value" {
					assert.Equal(t, "test", params.ClientParams.Get("filter"))
					assert.Equal(t, "value", params.ClientParams.Get("extra"))
				}
			}
		})
	}
}

// ============================================================================
// BuildURL Tests
// ============================================================================

func TestCursorPaginator_BuildURL(t *testing.T) {
	tests := []struct {
		name         string
		params       CursorParams
		id           string
		ordering     string
		expectError  bool
		errorMsg     string
		checkContent []string
	}{
		{
			name: "with existing query params",
			params: CursorParams{
				Cursor:   "",
				Ordering: "asc",
				Limit:    10,
				BaseURL:  "http://example.com/api",
				ClientParams: func() url.Values {
					cp := url.Values{}
					cp.Add("limit", "10")
					cp.Add("ordering", "asc")
					cp.Add("filter", "test")
					return cp
				}(),
			},
			id:          "123",
			ordering:    "asc",
			expectError: false,
			checkContent: []string{
				"http://example.com/api?",
				"ordering=asc",
				"limit=10",
				"filter=test",
			},
		},
		{
			name: "empty base URL",
			params: CursorParams{
				Cursor:   "",
				Ordering: "asc",
				Limit:    10,
				BaseURL:  "",
			},
			id:          "123",
			ordering:    "asc",
			expectError: true,
			errorMsg:    "base url is required",
		},
		{
			name: "with special characters in ID",
			params: CursorParams{
				Cursor:   "",
				Ordering: "asc",
				Limit:    10,
				BaseURL:  "http://example.com/api",
			},
			id:          "id with spaces",
			ordering:    "asc",
			expectError: false,
			checkContent: []string{
				"cursor=",
				"ordering=asc",
				"limit=10",
			},
		},
		{
			name: "with nil client params",
			params: CursorParams{
				Cursor:       "",
				Ordering:     "asc",
				Limit:        10,
				BaseURL:      "http://example.com/api",
				ClientParams: nil,
			},
			id:          "123",
			ordering:    "asc",
			expectError: false,
			checkContent: []string{
				"cursor=",
				"ordering=asc",
				"limit=10",
			},
		},
		{
			name: "with empty client params",
			params: CursorParams{
				Cursor:       "",
				Ordering:     "asc",
				Limit:        10,
				BaseURL:      "http://example.com/api",
				ClientParams: url.Values{},
			},
			id:          "456",
			ordering:    "desc",
			expectError: false,
			checkContent: []string{
				"ordering=desc",
				"limit=10",
			},
		},
		{
			name: "with complex client params",
			params: CursorParams{
				Cursor:   "",
				Ordering: "asc",
				Limit:    10,
				BaseURL:  "http://example.com/api",
				ClientParams: func() url.Values {
					cp := url.Values{}
					cp.Add("filter", "name")
					cp.Add("filter", "specialty") // Multiple values
					cp.Add("sort", "created_at")
					return cp
				}(),
			},
			id:          "789",
			ordering:    "asc",
			expectError: false,
			checkContent: []string{
				"ordering=asc",
				"limit=10",
				"sort=created_at",
				"filter=name",
				"filter=specialty",
			},
		},
	}

	specialCharTests := []string{
		"id with spaces",
		"id-with-special-chars!@#$%",
		"id_with_unicode_你好",
		"id/with/slashes",
		"id=with=equals",
	}

	// Add special character tests
	for _, testID := range specialCharTests {
		tests = append(tests, struct {
			name         string
			params       CursorParams
			id           string
			ordering     string
			expectError  bool
			errorMsg     string
			checkContent []string
		}{
			name: fmt.Sprintf("special chars ID: %s", testID),
			params: CursorParams{
				Cursor:   "",
				Ordering: "asc",
				Limit:    10,
				BaseURL:  "http://example.com/api",
			},
			id:          testID,
			ordering:    "asc",
			expectError: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.params.Validate()
			paginator := NewCursorPaginator[mockEntity](tt.params)

			result, err := paginator.buildURL(tt.id, tt.ordering)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)

				// Check expected content
				for _, content := range tt.checkContent {
					assert.Contains(t, result, content)
				}

				// Verify cursor encoding
				u, err := url.Parse(result)
				assert.NoError(t, err)
				cursor := u.Query().Get("cursor")
				decodedCursor, err := decodeCursor(cursor)
				assert.NoError(t, err)
				assert.Equal(t, tt.id, decodedCursor)

				// Verify ordering and limit
				assert.Equal(t, tt.ordering, u.Query().Get("ordering"))
				assert.Equal(t, fmt.Sprintf("%d", tt.params.Limit), u.Query().Get("limit"))

				// For complex client params test, verify multiple filter values
				if tt.name == "with complex client params" {
					filterValues := u.Query()["filter"]
					assert.Contains(t, filterValues, "name")
					assert.Contains(t, filterValues, "specialty")
				}
			}
		})
	}
}
