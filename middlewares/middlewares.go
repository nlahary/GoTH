package middlewares

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(status int) {
	rec.status = status
	rec.ResponseWriter.WriteHeader(status)
}

func DetailedLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Lire le corps de la requête
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Erreur lors de la lecture du corps de la requête:", err)
			return
		}
		// Rétablir le corps de la requête pour les prochains handlers
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Utiliser le ResponseWriter personnalisé
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		// Exécuter le prochain handler
		next.ServeHTTP(rec, r)

		log.Printf(`{"method": "%s", "status_code": %d, "url": "%s", "body": "%s"}`, r.Method, rec.status, r.URL.String(), string(bodyBytes))
	})
}
