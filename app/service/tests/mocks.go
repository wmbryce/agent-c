package tests

import (
	"net/http"

	"github.com/wmbryce/agent-c/app/types"
)

// MockStore implements store.SqlStore for testing
type MockStore struct {
	Creds     *types.ModelCredentials
	CredsErr  error
	Models    []types.Model
	CreateErr error
}

func (m *MockStore) CreateModel(model *types.Model) (*types.Model, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	return model, nil
}

func (m *MockStore) GetModels() ([]types.Model, error) {
	return m.Models, nil
}

func (m *MockStore) GetModelCredentials(modelKey string) (*types.ModelCredentials, error) {
	return m.Creds, m.CredsErr
}

func (m *MockStore) Close() {}

// MockHTTPClient implements service.HTTPClient for testing
type MockHTTPClient struct {
	Response *http.Response
	Err      error
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}
