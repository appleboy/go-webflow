package webflow

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strconv"

	"github.com/astaxie/beego/httplib"
)

// API for webflow structure
type API struct {
	Token   string
	Version string
	Debug   bool
}

// Param for http get parameter
type Param struct {
	Page    int
	PerPage int
}

var (
	// ErrorMissingTokenOrVersion for missing config
	ErrorMissingTokenOrVersion = errors.New("missing webflow token or version")
)

// New for web flow API object
func New(token, version string, debug bool) (*API, error) {
	if token == "" || version == "" {
		return nil, ErrorMissingTokenOrVersion
	}

	api := &API{
		Token:   token,
		Version: version,
		Debug:   debug,
	}

	return api, nil
}

// GetAllItemsFromCollection Get All Items For a Collection
func (api *API) GetAllItemsFromCollection(id string, resp interface{}, params ...Param) error {
	url := fmt.Sprintf(allItemsFromCollectionURL, id)
	err := api.fetchData(url, resp, params...)

	return err
}

func (api *API) fetchData(url string, res interface{}, params ...Param) error {
	req := httplib.Get(url).Debug(api.Debug)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	// Add authorization header
	req.Header("Authorization", "Bearer "+api.Token)
	req.Header("Accept-Version", api.Version)

	offset := 0
	limit := 10
	if len(params) > 0 {
		param := params[0]
		if param.PerPage > 0 {
			limit = param.PerPage
		}
		if param.Page > 1 {
			offset = (param.Page - 1) * limit
		}
	}

	req.Param("offset", strconv.Itoa(offset))
	req.Param("limit", strconv.Itoa(limit))
	err := req.ToJSON(res)

	return err
}
