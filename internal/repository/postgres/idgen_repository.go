package postgres

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"digit-oss/noc-services/internal/config"
	"digit-oss/noc-services/internal/domain"
	"digit-oss/noc-services/internal/domain/idgen"
)

// IdGenRepository calls the external ID-generation service over HTTP.
type IdGenRepository struct {
	Client *http.Client
	Cfg    *config.Config
}

// GetID calls the idgen service to generate IDs.
func (r *IdGenRepository) GetID(
	requestInfo *domain.RequestInfo,
	tenantID string,
	idName string,
	count int,
) ([]string, error) {
	idRequests := make([]idgen.IdRequest, 0, count)
	for i := 0; i < count; i++ {
		idRequests = append(idRequests, idgen.IdRequest{
			IdName:   idName,
			TenantID: tenantID,
		})
	}
	reqBody := idgen.IdGenerationRequest{
		RequestInfo: requestInfo,
		IdRequests:  idRequests,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal idgen request: %w", err)
	}

	log.Printf("[IdGenRepository] Request Body: %s", string(body))

	url := r.Cfg.IdGenHost + r.Cfg.IdGenPath
	log.Printf("[IdGenRepository] POST URL: %s", url)
	resp, err := r.Client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("idgen POST %s: %w", url, err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read idgen response: %w", err)
	}

	log.Printf("[IdGenRepository] Response Status: %d, Response Body: %s", resp.StatusCode, string(respBody))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("idgen returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var idResp idgen.IdGenerationResponse
	if err := json.Unmarshal(respBody, &idResp); err != nil {
		return nil, fmt.Errorf("unmarshal idgen response: %w", err)
	}

	ids := make([]string, 0, len(idResp.IdResponses))
	for _, r := range idResp.IdResponses {
		ids = append(ids, r.ID)
	}
	return ids, nil
}
