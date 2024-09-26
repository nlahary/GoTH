package middlewares

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

// Créer un middleware qui log chaque requête avec des détails
func DetailedLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusCode := http.StatusOK

		// Lire le corps de la requête
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Erreur lors de la lecture du corps de la requête:", err)
			return
		}
		// Rétablir le corps de la requête pour les prochains handlers
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		log.Printf(`{"method": "%s", "status_code": %d, "url": "%s", "body": "%s"}`, r.Method, statusCode, r.URL.String(), string(bodyBytes))
		next.ServeHTTP(w, r)
	})
}
