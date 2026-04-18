package middleware_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raystack/salt/middleware"
	"github.com/raystack/salt/middleware/cors"
	"github.com/raystack/salt/middleware/requestid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func nopLogger() *slog.Logger { return slog.New(slog.DiscardHandler) }

func TestDefault(t *testing.T) {
	interceptors := middleware.Default(nopLogger())
	assert.Len(t, interceptors, 4)
}

func TestDefaultHTTP(t *testing.T) {
	chain := middleware.DefaultHTTP(nopLogger())
	assert.NotNil(t, chain)

	handler := chain(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	// Should have request ID in response
	assert.NotEmpty(t, rec.Header().Get(requestid.Header))
}

func TestChainHTTP(t *testing.T) {
	t.Run("chains in order", func(t *testing.T) {
		var order []string
		mw1 := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, "first")
				next.ServeHTTP(w, r)
			})
		}
		mw2 := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, "second")
				next.ServeHTTP(w, r)
			})
		}

		chain := middleware.ChainHTTP(mw1, mw2)
		handler := chain(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			order = append(order, "handler")
		}))

		req := httptest.NewRequest("GET", "/", nil)
		handler.ServeHTTP(httptest.NewRecorder(), req)
		assert.Equal(t, []string{"first", "second", "handler"}, order)
	})
}

func TestRecoveryHTTP(t *testing.T) {
	chain := middleware.DefaultHTTP(nopLogger())
	handler := chain(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestRequestIDHTTP(t *testing.T) {
	chain := middleware.DefaultHTTP(nopLogger())
	handler := chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := requestid.FromContext(r.Context())
		w.Write([]byte(id))
	}))

	t.Run("generates ID when missing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		body, _ := io.ReadAll(rec.Body)
		assert.NotEmpty(t, string(body))
		assert.Equal(t, string(body), rec.Header().Get(requestid.Header))
	})

	t.Run("propagates existing ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set(requestid.Header, "my-custom-id")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		body, _ := io.ReadAll(rec.Body)
		assert.Equal(t, "my-custom-id", string(body))
		assert.Equal(t, "my-custom-id", rec.Header().Get(requestid.Header))
	})
}

func TestCORSMiddleware(t *testing.T) {
	handler := cors.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("preflight returns 204", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/", nil)
		req.Header.Set("Origin", "http://example.com")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Equal(t, "http://example.com", rec.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("regular request gets CORS headers", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Origin", "http://example.com")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "http://example.com", rec.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("no Origin header skips CORS", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("includes Connect-specific headers", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/", nil)
		req.Header.Set("Origin", "http://example.com")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		headers := rec.Header().Get("Access-Control-Allow-Headers")
		require.Contains(t, headers, "Connect-Protocol-Version")
		require.Contains(t, headers, "Connect-Timeout-Ms")
	})
}
