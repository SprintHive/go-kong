package kong

import (
	"fmt"
	"net/http"
)

// CertificatesService handles communication with Kong's '/certificates' resource.
type CertificatesService struct {
	*service
}

// Certificates represents the object returned from Kong when querying for
// multiple certificates objects.
//
// In cases where the number of objects returned exceeds the maximum,
// Next holds the URI for the next set of results.
// i.e. "http://localhost:8001/certificates/?size=2&offset=4d924084-1adb-40a5-c042-63b19db421d1"
type Certificates struct {
	Data   []*Certificate `json:"data,omitempty"`
	Total  int            `json:"total,omitempty"`
	Next   string         `json:"next,omitempty"`
	Offset string         `json:"offset,omitempty"`
}

// CertificateRequest represents a Kong certificate object required when posting a new certificate to /certificate/
type CertificateRequest struct {
	Cert string `json:"cert,omitempty"`
	Key  string `json:"key,omitempty"`
	Snis string `json:"snis,omitempty"`
}

// Certificate represents a Kong certificate object received when querying /certificate/{id}
type Certificate struct {
	ID        string   `json:"id,omitempty"`
	Cert      string   `json:"cert,omitempty"`
	Key       string   `json:"key,omitempty"`
	Snis      []string `json:"snis,omitempty"`
	CreatedAt int64    `json:"created_at,omitempty"`
}

// Get queries for a single Kong certificate object, by name or id.
//
// Equivalent to GET /certificates/{sni or id}
func (s *CertificatesService) Get(certificateID string) (*Certificate, *http.Response, error) {
	u := fmt.Sprintf("certificates/%v", certificateID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	uResp := new(Certificate)
	resp, err := s.client.Do(req, uResp)
	if err != nil {
		return nil, resp, err
	}

	return uResp, resp, err
}

// Patch updates an existing Kong certificate object.
//
// Equivalent to PATCH /certificates/{name or id}
func (s *CertificatesService) Patch(certificate *CertificateRequest, certificateID string) (*http.Response, error) {
	u := fmt.Sprintf("certificates/%v", certificateID)

	req, err := s.client.NewRequest("PATCH", u, certificate)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err

}

// Delete deletes a single Kong certificates object, by name or id.
//
// Equivalent to DELETE /certificates/{name or id}
func (s *CertificatesService) Delete(certificateID string) (*http.Response, error) {
	u := fmt.Sprintf("certificates/%v", certificateID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// Post creates a new Kong certificates object.
//
// Equivalent to POST /certificates
func (s *CertificatesService) Post(certificate *CertificateRequest) (*http.Response, error) {
	req, err := s.client.NewRequest("POST", "certificates", certificate)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}

// GetAll retreives all certificates
//
// Equivalent to GET /certificates
func (s *CertificatesService) GetAll() (*Certificates, *http.Response, error) {
	u := fmt.Sprintf("certificates")

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	certificates := new(Certificates)
	resp, err := s.client.Do(req, certificates)
	if err != nil {
		return nil, resp, err
	}

	return certificates, resp, err
}
