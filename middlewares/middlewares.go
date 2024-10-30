package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	kafka "github.com/Nathanael-FR/website/kafka"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

type logEntry struct {
	Method     string `json:"method"`
	StatusCode int    `json:"status_code"`
	URL        string `json:"url"`
	Body       string `json:"body"`
}

func (rec *statusRecorder) WriteHeader(status int) {
	rec.status = status
	rec.ResponseWriter.WriteHeader(status)
}

func DetailedLoggingMiddleware(next http.Handler, p *kafka.Producer) http.Handler {
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

		logEntry := logEntry{
			Method:     r.Method,
			StatusCode: rec.status,
			URL:        r.URL.String(),
			Body:       string(bodyBytes),
		}

		logBytes, err := json.Marshal(logEntry)
		if err != nil {
			log.Println("Error while serializing log entry:", err)
			return
		}
		err = p.SendMessage(string(logBytes))
		if err != nil {
			log.Println("Error while sending log entry to Kafka:", err)
			return
		}

		log.Printf(`{"method": "%s", "status_code": %d, "url": "%s", "body": "%s"}`, r.Method, rec.status, r.URL.String(), string(bodyBytes))
	})
}
