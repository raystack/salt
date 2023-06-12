package version_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/raystack/salt/version"
	"github.com/stretchr/testify/assert"
)

func TestGithubInfo(t *testing.T) {
	muxRouter := mux.NewRouter()
	server := httptest.NewServer(muxRouter)

	t.Run("should check for latest version availability by extracting correct version tag for valid json response on release URL", func(t *testing.T) {
		muxRouter.HandleFunc("/latest", func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(struct {
				TagName string `json:"tag_name"`
			}{
				TagName: "v0.0.2",
			})
			rw.Write(response)
		})
		info, err := version.ReleaseInfo("http://" + server.Listener.Addr().String() + "/latest")
		assert.Nil(t, err)
		assert.Equal(t, "v0.0.2", info.Version)
		info, err = version.ReleaseInfo("http://" + server.Listener.Addr().String() + "/latest")
		assert.Nil(t, err)
		assert.NotEqual(t, "v0.0.1", info.Version)
	})
}
func TestIsCurrentLatest(t *testing.T) {
	muxRouter := mux.NewRouter()
	server := httptest.NewServer(muxRouter)

	t.Run("should return true for current version  as the latest version ", func(t *testing.T) {
		muxRouter.HandleFunc("/latest", func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(struct {
				TagName string `json:"tag_name"`
			}{
				TagName: "v0.0.2",
			})
			rw.Write(response)
		})
		info, err := version.ReleaseInfo("http://" + server.Listener.Addr().String() + "/latest")
		assert.Nil(t, err)
		assert.Equal(t, "v0.0.2", info.Version)
		res, err := version.IsCurrentLatest("v0.0.2", info.Version)
		assert.Nil(t, err)
		assert.True(t, res)
	})
	t.Run("should return false for current version not same as the latest version", func(t *testing.T) {
		muxRouter.HandleFunc("/latest", func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(struct {
				TagName string `json:"tag_name"`
			}{
				TagName: "v0.0.2",
			})
			rw.Write(response)
		})
		info, err := version.ReleaseInfo("http://" + server.Listener.Addr().String() + "/latest")
		assert.Nil(t, err)
		assert.Equal(t, "v0.0.2", info.Version)
		res, err := version.IsCurrentLatest("v0.0.1", info.Version)
		assert.Nil(t, err)
		assert.False(t, res)
		res, err = version.IsCurrentLatest("", info.Version)
		assert.NotNil(t, err)
		assert.False(t, res)
		res, err = version.IsCurrentLatest("v0.0.3", "")
		assert.NotNil(t, err)
		assert.False(t, res)
	})
}

func TestUpdateNotice(t *testing.T) {
	muxRouter := mux.NewRouter()
	server := httptest.NewServer(muxRouter)
	t.Run("basic check for notify latest version", func(t *testing.T) {
		muxRouter.HandleFunc("/latest", func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(struct {
				TagName string `json:"tag_name"`
			}{
				TagName: "v0.0.1",
			})
			rw.Write(response)
		})
		info, err := version.ReleaseInfo("http://" + server.Listener.Addr().String() + "/latest")
		assert.Nil(t, err)
		assert.Equal(t, "v0.0.1", info.Version)
		res, err := version.IsCurrentLatest("v0.0.1", info.Version)
		assert.Nil(t, err)
		assert.True(t, res)
		s := version.UpdateNotice("v0.0.1", "raystack/optimus")
		assert.NotEqual(t, "v0.0.1", s)
	})
}
