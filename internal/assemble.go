package internal

import (
	"net/http"
	"time"

	"github.com/airgap-solution/crypto-wallet-rest/internal/config"
	cryptowalletrest "github.com/airgap-solution/crypto-wallet-rest/openapi/servergen/go"
)

var readTimeout = time.Second * 10

func Assemble(cfg config.Config, servicer cryptowalletrest.DefaultAPIServicer) *http.Server {
	ctrl := cryptowalletrest.NewDefaultAPIController(servicer)
	router := cryptowalletrest.NewRouter(ctrl)
	srv := &http.Server{Addr: cfg.ListenAddr, Handler: corsMiddleware(router), ReadTimeout: readTimeout}
	return srv
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}
