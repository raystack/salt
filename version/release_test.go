package version

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReleaseInfo_Success(t *testing.T) {
	// Mock a successful GitHub API response
	mockResponse := `{
		"tag_name": "v1.2.3",
		"tarball_url": "https://example.com/tarball/v1.2.3"
	}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	info, err := ReleaseInfo(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, "v1.2.3", info.Version)
	assert.Equal(t, "https://example.com/tarball/v1.2.3", info.TarURL)
}

func TestReleaseInfo_Failure(t *testing.T) {
	// Mock a failed GitHub API response with a non-OK status code
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	info, err := ReleaseInfo(server.URL)
	assert.Error(t, err)
	assert.Nil(t, info)
}

func TestReleaseInfo_InvalidJSON(t *testing.T) {
	// Mock a response with invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	info, err := ReleaseInfo(server.URL)
	assert.Error(t, err)
	assert.Nil(t, info)
}

func TestIsCurrentLatest(t *testing.T) {
	// Test cases for version comparison
	tests := []struct {
		currVersion   string
		latestVersion string
		expected      bool
		shouldError   bool
	}{
		{"1.2.3", "1.2.2", true, false},
		{"1.2.3", "1.2.3", true, false},
		{"1.2.2", "1.2.3", false, false},
		{"invalid", "1.2.3", false, true},
		{"1.2.3", "invalid", false, true},
	}

	for _, test := range tests {
		result, err := IsCurrentLatest(test.currVersion, test.latestVersion)
		if test.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		}
	}
}

func TestUpdateNotice(t *testing.T) {
	// Mock a successful GitHub API response with a newer version
	mockResponse := `{
		"tag_name": "v2.0.0",
		"tarball_url": "https://example.com/tarball/v2.0.0"
	}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	notice := UpdateNotice("1.0.0", server.URL)
	assert.Equal(t, "A new release (v2.0.0) is available, consider updating the client.", notice)

	// Test with the current version being the latest
	notice = UpdateNotice("2.0.0", server.URL)
	assert.Equal(t, "", notice)
}

func TestUpdateNotice_ErrorHandling(t *testing.T) {
	// Mock a server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	notice := UpdateNotice("1.0.0", server.URL)
	assert.Equal(t, "", notice)
}
