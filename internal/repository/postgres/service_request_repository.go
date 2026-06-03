package postgres

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// ServiceRequestRepository is a generic HTTP POST helper for calling
// external DIGIT micro-services. Mirrors ServiceRequestRepository.java.
type ServiceRequestRepository struct {
	Client *http.Client
}

// FetchResult makes an HTTP POST to the given URI, serialising request as JSON,
// and deserialising the response body into the supplied response pointer.
//
// The response parameter must be a pointer to the target struct
// (e.g., *domain.BPAResponse).
func (r *ServiceRequestRepository) FetchResult(uri string, request interface{}, response interface{}) error {
	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	log.Printf("[ServiceRequestRepository] POST %s body=%s", uri, string(body))

	resp, err := r.Client.Post(uri, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("POST %s: %w", uri, err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response from %s: %w", uri, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("service %s returned status %d: %s", uri, resp.StatusCode, string(respBody))
	}

	if err := json.Unmarshal(respBody, response); err != nil {
		return fmt.Errorf("unmarshal response from %s: %w", uri, err)
	}

	return nil
}