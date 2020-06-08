/*
 * Copyright © 2020. Ernesto Alejandro Santana Hidalgo <ernesto.alejandrosantana@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Ernesto Alejandro Santana Hidalgo <ernesto.alejandrosantana@gmail.com>
 * @copyright 	Ernesto Alejandro Santana Hidalgo <ernesto.alejandrosantana@gmail.com>
 * @license 	Apache-2.0
 *
 */

package client

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"time"
)

// DefaultRESTTimeout - default RPC timeout is one minute.
const DefaultRESTTimeout = 1 * time.Minute

// NetworkError - error type in case of errors related to http/transport
// for ex. connection refused, connection reset, dns resolution failure etc.
// All errors returned by storage-rest-server (ex errFileNotFound, errDiskNotFound) are not considered to be network errors.
type NetworkError struct {
	Err error
}

func (n *NetworkError) Error() string {
	return n.Err.Error()
}

// Client - http based RPC client.
type Client struct {
	httpClient          *http.Client
	httpIdleConnsCloser func()

	random *rand.Rand
}

// URL query separator constants
const (
	querySep = "?"
)

func (c *Client) CallRetryable(req *http.Request) (reply io.ReadCloser, err error) {
	var reqRetry = MaxRetry // Indicates how many times we can retry the request

	// Create a done channel to control 'ListObjects' go routine.
	doneCh := make(chan struct{}, 1)

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	for range newRetryTimer(reqRetry, DefaultRetryUnit, DefaultRetryCap, MaxJitter, doneCh) {
		// Instantiate a new request.
		var resp *http.Response

		// Initiate the request.
		resp, err = c.httpClient.Do(req)
		if err != nil {
			// For supported rest requests errors verify.
			if isHTTPReqErrorRetryable(err) {
				continue // Retry.
			}
			// For other errors, return here no need to retry.
			return nil, err
		}

		// For any known successful rest status, return quickly.
		for _, httpStatus := range successStatus {
			if httpStatus == resp.StatusCode {
				return resp.Body, err
			}
		}

		// Verify if rest status code is retryable.
		if isHTTPStatusRetryable(resp.StatusCode) {
			continue // Retry.
		}

		break
	}
	return nil, &NetworkError{errors.New("failed to fetch the resource: " + req.URL.String())}
}

// Close closes all idle connections of the underlying http client
func (c *Client) Close() {
	if c.httpIdleConnsCloser != nil {
		c.httpIdleConnsCloser()
	}
}

// List of success status.
var successStatus = []int{
	http.StatusOK,
	http.StatusNoContent,
	http.StatusPartialContent,
}

// DrainBody close non nil response with any response Body.
// convenient wrapper to drain any remaining data on response body.
//
// Subsequently this allows golang http RoundTripper
// to re-use the same connection for future requests.
func DrainBody(respBody io.ReadCloser) {
	// Callers should close resp.Body when done reading from it.
	// If resp.Body is not closed, the Client's underlying RoundTripper
	// (typically Transport) may not be able to re-use a persistent TCP
	// connection to the server for a subsequent "keep-alive" request.
	if respBody != nil {
		// Drain any remaining Body and then close the connection.
		// Without this closing connection would disallow re-using
		// the same connection for future uses.
		//  - http://stackoverflow.com/a/17961593/4465767
		defer respBody.Close()
		_, _ = io.Copy(ioutil.Discard, respBody)
	}
}

func newCustomDialContext(timeout time.Duration) func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		dialer := &net.Dialer{
			Timeout:   timeout,
			KeepAlive: timeout,
		}
		return dialer.DialContext(ctx, network, addr)
	}
}

// NewClient - returns new REST client.
func NewClient() *Client {
	// Transport is exactly same as Go default in https://golang.org/pkg/net/http/#RoundTripper
	// except custom DialContext and TLSClientConfig.
	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           newCustomDialContext(30 * time.Second),
		MaxIdleConns:          256,
		MaxIdleConnsPerHost:   256,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
		DisableCompression:    true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return &Client{
		httpClient:          &http.Client{Transport: tr},
		httpIdleConnsCloser: tr.CloseIdleConnections,
		// Introduce a new locked random seed.
		random: rand.New(&lockedRandSource{src: rand.NewSource(time.Now().UTC().UnixNano())}),
	}
}
