// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/yieldbot/ybsensues/Godeps/_workspace/src/github.com/olivere/elastic/uritemplates"
)

// ExistsService checks if a document exists.
//
// See http://www.elastic.co/guide/en/elasticsearch/reference/current/docs-get.html
// for details.
type ExistsService struct {
	client     *Client
	pretty     bool
	id         string
	index      string
	typ        string
	parent     string
	preference string
	realtime   *bool
	refresh    *bool
	routing    string
}

// NewExistsService creates a new ExistsService.
func NewExistsService(client *Client) *ExistsService {
	return &ExistsService{
		client: client,
	}
}

// Id is the document ID.
func (s *ExistsService) Id(id string) *ExistsService {
	s.id = id
	return s
}

// Index is the name of the index.
func (s *ExistsService) Index(index string) *ExistsService {
	s.index = index
	return s
}

// Type is the type of the document (use `_all` to fetch the first
// document matching the ID across all types).
func (s *ExistsService) Type(typ string) *ExistsService {
	s.typ = typ
	return s
}

// Parent is the ID of the parent document.
func (s *ExistsService) Parent(parent string) *ExistsService {
	s.parent = parent
	return s
}

// Preference specifies the node or shard the operation should be
// performed on (default: random).
func (s *ExistsService) Preference(preference string) *ExistsService {
	s.preference = preference
	return s
}

// Realtime specifies whether to perform the operation in realtime or search mode.
func (s *ExistsService) Realtime(realtime bool) *ExistsService {
	s.realtime = &realtime
	return s
}

// Refresh the shard containing the document before performing the operation.
func (s *ExistsService) Refresh(refresh bool) *ExistsService {
	s.refresh = &refresh
	return s
}

// Routing is the specific routing value.
func (s *ExistsService) Routing(routing string) *ExistsService {
	s.routing = routing
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *ExistsService) Pretty(pretty bool) *ExistsService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *ExistsService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/{index}/{type}/{id}", map[string]string{
		"id":    s.id,
		"index": s.index,
		"type":  s.typ,
	})
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "1")
	}
	if s.parent != "" {
		params.Set("parent", s.parent)
	}
	if s.preference != "" {
		params.Set("preference", s.preference)
	}
	if s.realtime != nil {
		params.Set("realtime", fmt.Sprintf("%v", *s.realtime))
	}
	if s.refresh != nil {
		params.Set("refresh", fmt.Sprintf("%v", *s.refresh))
	}
	if s.routing != "" {
		params.Set("routing", s.routing)
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *ExistsService) Validate() error {
	var invalid []string
	if s.id == "" {
		invalid = append(invalid, "Id")
	}
	if s.index == "" {
		invalid = append(invalid, "Index")
	}
	if s.typ == "" {
		invalid = append(invalid, "Type")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *ExistsService) Do() (bool, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return false, err
	}

	// Get URL for request
	path, params, err := s.buildURL()
	if err != nil {
		return false, err
	}

	// Get HTTP response
	res, err := s.client.PerformRequest("HEAD", path, params, nil)
	if err != nil {
		return false, err
	}

	// Evaluate operation response
	switch res.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("elastic: got HTTP code %d when it should have been either 200 or 404", res.StatusCode)
	}
}