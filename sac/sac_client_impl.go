package sac

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"sac-cli/sac/dto"

	"golang.org/x/oauth2/clientcredentials"
	"gopkg.in/resty.v1"
)

var ErrorPermissionDenied = fmt.Errorf("permission denied")
var ErrorNotFound = fmt.Errorf("not found")
var ErrConflict = fmt.Errorf("already exist")

type SecureAccessCloudClientImpl struct {
	Setting *SecureAccessCloudSettings
	Client  *resty.Client
}

func NewSecureAccessCloudClientImpl(setting *SecureAccessCloudSettings) SecureAccessCloudClient {
	return &SecureAccessCloudClientImpl{Client: nil, Setting: setting}
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Application API
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *SecureAccessCloudClientImpl) FindApplicationByID(id string) (*dto.ApplicationDTO, error) {
	endpoint := s.Setting.BuildAPIPrefixURL() + "/v2/applications/" + id

	var application dto.ApplicationDTO
	err := s.performGetRequest(endpoint, &application)

	if err != nil {
		return &dto.ApplicationDTO{}, err
	}

	return &application, nil
}

func (s *SecureAccessCloudClientImpl) FindApplicationByName(name string) (*dto.ApplicationDTO, error) {
	endpoint := s.Setting.BuildAPIPrefixURL() + "/v2/applications" + "?filter=" + url.QueryEscape(name)

	var applications dto.ApplicationPageDTO
	err := s.performGetRequest(endpoint, &applications)

	if err != nil {
		return &dto.ApplicationDTO{}, err
	}

	if applications.NumberOfElements == 0 {
		return &dto.ApplicationDTO{}, ErrorNotFound
	}

	// Return the first policy if more than one found
	return &applications.Content[0], nil
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Policy API
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *SecureAccessCloudClientImpl) FindPolicyByName(name string) (dto.PolicyDTO, error) {
	endpoint := s.Setting.BuildAPIPrefixURL() + "/v2/policies" + "?filter=" + url.QueryEscape(name)

	var policies dto.PoliciesPageDTO
	err := s.performGetRequest(endpoint, &policies)

	if err != nil {
		return dto.PolicyDTO{}, err
	}

	if policies.NumberOfElements == 0 {
		return dto.PolicyDTO{}, ErrorNotFound
	}

	// Return the first policy if more than one found
	return policies.Content[0], nil
}

func (s *SecureAccessCloudClientImpl) FindPoliciesByNames(names []string) ([]dto.PolicyDTO, error) {
	var results []dto.PolicyDTO

	for _, name := range names {
		policyDTO, err := s.FindPolicyByName(name)
		if err != nil {
			return results, err
		}

		results = append(results, policyDTO)
	}

	return results, nil
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Private Functions
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *SecureAccessCloudClientImpl) getClient() *resty.Client {
	if s.Client != nil {
		return s.Client
	}

	cfg := clientcredentials.Config{
		ClientID:     s.Setting.ClientID,
		ClientSecret: s.Setting.ClientSecret,
		TokenURL:     s.Setting.BuildOAuthTokenURL(),
		Scopes:       []string{},
	}

	oauthClient := cfg.Client(context.Background())
	client := resty.New().SetRetryCount(0).SetTimeout(1 * time.Minute)
	// http://godoc.org/golang.org/x/oauth2 implements `httpRoundTripper` interface
	// Set the oauthClient transport
	client.SetTransport(oauthClient.Transport)

	s.Client = client
	return s.Client
}

func (s *SecureAccessCloudClientImpl) performGetRequest(endpoint string, obj interface{}) error {
	// 1. Get Authorization Token
	client := s.getClient()

	// 2. Perform the GET request
	response, err := client.NewRequest().Get(endpoint)
	if err != nil {
		return err
	}

	if response.StatusCode() == http.StatusNotFound {
		return ErrorNotFound
	}

	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed with status-code: %d and body: %s", response.StatusCode(), response.String())
	}

	// 3. Convert to Commit model
	err = json.Unmarshal(response.Body(), &obj)
	if err != nil {
		return err
	}

	return nil
}

func (s *SecureAccessCloudClientImpl) performModifyRequest(method string, endpoint string, requestObj interface{}, responseObj interface{}) error {
	// 1. Get Authorization Token
	client := s.getClient()

	// 2. Perform the request
	request := client.NewRequest()
	var response *resty.Response
	var err error

	// 2.1. Marshal the request
	if requestObj != nil {
		body, err := json.Marshal(requestObj)
		if err != nil {
			return err
		}

		request = request.SetBody(body)

	}
	request.SetHeader("Content-Type", "application/json")

	switch method {
	case http.MethodPost:
		response, err = request.Post(endpoint)
	case http.MethodPut:
		response, err = request.Put(endpoint)
	default:
		return errors.New("unsupported http method: " + method)
	}

	if response.StatusCode() == http.StatusNotFound {
		return ErrorNotFound
	}

	if !isSuccess(response.StatusCode()) {
		return fmt.Errorf("failed with status-code: %d and body: %s", response.StatusCode(), response.String())
	}

	// 4. Unmarshal response body
	if responseObj != nil {
		err = json.Unmarshal(response.Body(), &responseObj)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SecureAccessCloudClientImpl) performPostRequest(endpoint string, body, obj interface{}) error {
	// 1. Get Authorization Token
	client := s.getClient()

	// 2. Perform the POST request
	response, err := client.NewRequest().SetBody(body).Post(endpoint)
	if err != nil {
		return err
	}

	if response.StatusCode() == http.StatusConflict {
		return ErrConflict
	}

	if response.StatusCode() != http.StatusCreated {
		return fmt.Errorf("failed with status-code: %d and body: %s", response.StatusCode(), response.String())
	}

	// 3. Convert to Commit model
	err = json.Unmarshal(response.Body(), obj)
	if err != nil {
		return err
	}

	return nil
}

func (s *SecureAccessCloudClientImpl) performDeleteRequest(endpoint string) error {
	// 1. Get Authorization Token
	client := s.getClient()

	// 2. Perform the DELETE request
	response, err := client.NewRequest().Delete(endpoint)
	if err != nil {
		return err
	}

	if response.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("failed with status-code: %d and body: %s", response.StatusCode(), response.String())
	}

	return nil
}

func isSuccess(status int) bool {
	return status >= http.StatusOK && status <= http.StatusOK+99
}
