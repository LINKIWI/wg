package supercharged

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// HTTPClient is a Supercharged client that transacts requests and responses over HTTP.
type HTTPClient struct {
	backend *http.Client
	baseURL *url.URL
}

// NewHTTPClient creates a new client instance for a Supercharged-compliant HTTP server, with an
// optional http.Client backend for initiating requests.
func NewHTTPClient(baseURL string, backend *http.Client) (*HTTPClient, error) {
	if backend == nil {
		backend = &http.Client{}
	}

	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &HTTPClient{
		backend: backend,
		baseURL: parsed,
	}, nil
}

// Do executes an HTTP request against the server.
func (c *HTTPClient) Do(method string, endpoint string, data interface{}, response interface{}) *Error {
	if data == nil {
		data = struct{}{}
	}

	serializedData, err := json.Marshal(data)
	if err != nil {
		return Wrap(err)
	}

	resource, err := c.baseURL.Parse(endpoint)
	if err != nil {
		return Wrap(err)
	}

	// By Supercharged conventions, GET and HEAD requests put the payload in a request header, while
	// other methods put the payload in the HTTP body
	var req *http.Request
	switch method {
	case http.MethodHead, http.MethodGet:
		req, err = http.NewRequest(method, resource.String(), nil)
		if err != nil {
			return Wrap(err)
		}
		req.Header.Add("X-Supercharged-Data", string(serializedData))

	case http.MethodPost, http.MethodPut, http.MethodDelete:
		req, err = http.NewRequest(method, resource.String(), bytes.NewReader(serializedData))
		if err != nil {
			return Wrap(err)
		}

	default:
		return &Error{
			Code:    CodeNotFound,
			Message: "unsupported Supercharged HTTP method in request",
		}
	}

	httpResp, err := c.backend.Do(req)
	if err != nil {
		return Wrap(err)
	}

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return Wrap(err)
	}

	var scResp Response
	if err := json.Unmarshal(body, &scResp); err != nil {
		return Wrap(err)
	}

	if !scResp.Success {
		return &Error{
			Status:  httpResp.StatusCode,
			Code:    scResp.Code,
			Message: scResp.Message,
			Data:    scResp.Data,
		}
	}

	if response != nil {
		if err := json.Unmarshal(scResp.Data, response); err != nil {
			return Wrap(err)
		}
	}

	return nil
}
