// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by unknown module path version unknown version DO NOT EDIT.
package api

import (
	"encoding/json"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for AccessTokenTokenType.
const (
	Bearer AccessTokenTokenType = "Bearer"
)

// Defines values for HealthCheckStatusStatus.
const (
	Available   HealthCheckStatusStatus = "available"
	Starting    HealthCheckStatusStatus = "starting"
	Stopping    HealthCheckStatusStatus = "stopping"
	Unavailable HealthCheckStatusStatus = "unavailable"
)

// AS defines model for AS.
type AS = string

// AccessCredentials defines model for AccessCredentials.
type AccessCredentials struct {
	// ClientId ID of the control service requesting authentication.
	ClientId string `json:"client_id"`

	// ClientSecret Secret that authenticates the control service.
	ClientSecret string `json:"client_secret"`
}

// AccessToken defines model for AccessToken.
type AccessToken struct {
	// AccessToken The encoded JWT token
	AccessToken string `json:"access_token"`

	// ExpiresIn Validity duration of this token in seconds.
	ExpiresIn int `json:"expires_in"`

	// TokenType Type of returned access token. Currently always Bearer.
	TokenType AccessTokenTokenType `json:"token_type"`
}

// AccessTokenTokenType Type of returned access token. Currently always Bearer.
type AccessTokenTokenType string

// CertificateChain defines model for CertificateChain.
type CertificateChain struct {
	// AsCertificate Base64 encoded AS certificate.
	AsCertificate []byte `json:"as_certificate"`

	// CaCertificate Base64 encoded CA certificate.
	CaCertificate []byte `json:"ca_certificate"`
}

// CertificateChainPKCS7 Certificate chain containing the the new AS certificate and the issuing
// CA certificate encoded in a degenerate PKCS#7 data structure.
type CertificateChainPKCS7 = []byte

// HealthCheckStatus defines model for HealthCheckStatus.
type HealthCheckStatus struct {
	Status HealthCheckStatusStatus `json:"status"`
}

// HealthCheckStatusStatus defines model for HealthCheckStatus.Status.
type HealthCheckStatusStatus string

// Problem Error message encoded as specified in
// [RFC7807](https://tools.ietf.org/html/rfc7807)
type Problem struct {
	// CorrelationId Identifier to correlate multiple error messages to the same case.
	CorrelationId *openapi_types.UUID `json:"correlation_id,omitempty"`

	// Detail A human readable explanation specific to this occurrence of the problem that is helpful to locate the problem and give advice on how to proceed. Written in English and readable for engineers, usually not suited for non technical stakeholders and not localized.
	Detail *string `json:"detail,omitempty"`

	// Instance A URI reference that identifies the specific occurrence of the problem, e.g. by adding a fragment identifier or sub-path to the problem type.
	Instance *string `json:"instance,omitempty"`

	// Status The HTTP status code generated by the server for this occurrence of the problem.
	Status int `json:"status"`

	// Title A short summary of the problem type. Written in English and readable for engineers, usually not suited for non technical stakeholders and not localized.
	Title string `json:"title"`

	// Type A URI reference that uniquely identifies the problem type in the context of the provided API.
	Type string `json:"type"`
}

// RenewalRequest defines model for RenewalRequest.
type RenewalRequest struct {
	// Csr Base64 encoded renewal request as described below.
	//
	// The renewal requests consists of a CMS SignedData structure that
	// contains a PKCS#10 defining the parameters of the requested
	// certificate.
	//
	// The following must hold for the CMS structure:
	//
	// - The `certificates` field in `SignedData` MUST contain an existing
	//   and verifiable certificate chain that authenticates the private
	//   key that was used to sign the CMS structure. It MUST NOT contain
	//   any other certificates.
	//
	// - The `eContentType` is set to `id-data`. The contents of `eContent`
	//   is the ASN.1 DER encoded PKCS#10. This ensures backwards
	//   compatibility with PKCS#7, as described in
	//   [RFC5652](https://tools.ietf.org/html/rfc5652#section-5.2.1)
	//
	// - The `SignerIdentifier` MUST be the choice `IssuerAndSerialNumber`,
	//   thus, `version` in `SignerInfo` must be 1, as required by
	//   [RFC5652](https://tools.ietf.org/html/rfc5652#section-5.3)
	Csr []byte `json:"csr"`
}

// RenewalResponse defines model for RenewalResponse.
type RenewalResponse struct {
	CertificateChain RenewalResponse_CertificateChain `json:"certificate_chain"`
}

// RenewalResponse_CertificateChain defines model for RenewalResponse.CertificateChain.
type RenewalResponse_CertificateChain struct {
	union json.RawMessage
}

// PostAuthTokenJSONRequestBody defines body for PostAuthToken for application/json ContentType.
type PostAuthTokenJSONRequestBody = AccessCredentials

// PostCertificateRenewalJSONRequestBody defines body for PostCertificateRenewal for application/json ContentType.
type PostCertificateRenewalJSONRequestBody = RenewalRequest

// AsCertificateChain returns the union data inside the RenewalResponse_CertificateChain as a CertificateChain
func (t RenewalResponse_CertificateChain) AsCertificateChain() (CertificateChain, error) {
	var body CertificateChain
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromCertificateChain overwrites any union data inside the RenewalResponse_CertificateChain as the provided CertificateChain
func (t *RenewalResponse_CertificateChain) FromCertificateChain(v CertificateChain) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeCertificateChain performs a merge with any union data inside the RenewalResponse_CertificateChain, using the provided CertificateChain
func (t *RenewalResponse_CertificateChain) MergeCertificateChain(v CertificateChain) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsCertificateChainPKCS7 returns the union data inside the RenewalResponse_CertificateChain as a CertificateChainPKCS7
func (t RenewalResponse_CertificateChain) AsCertificateChainPKCS7() (CertificateChainPKCS7, error) {
	var body CertificateChainPKCS7
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromCertificateChainPKCS7 overwrites any union data inside the RenewalResponse_CertificateChain as the provided CertificateChainPKCS7
func (t *RenewalResponse_CertificateChain) FromCertificateChainPKCS7(v CertificateChainPKCS7) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeCertificateChainPKCS7 performs a merge with any union data inside the RenewalResponse_CertificateChain, using the provided CertificateChainPKCS7
func (t *RenewalResponse_CertificateChain) MergeCertificateChainPKCS7(v CertificateChainPKCS7) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

func (t RenewalResponse_CertificateChain) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *RenewalResponse_CertificateChain) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}