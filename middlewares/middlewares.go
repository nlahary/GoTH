package middlewares

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/nlahary/website/models"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(status int) {
	rec.status = status
	rec.ResponseWriter.WriteHeader(status)
}

func DetailedLoggingMiddleware(next http.Handler, l *models.HttpLogger) http.HandlerFunc {
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
		startTimeRequest := time.Now().UnixNano()
		next.ServeHTTP(rec, r)
		duration := time.Duration(time.Now().UnixNano() - startTimeRequest).Microseconds()

		logHttp := models.LogHTTP{
			Method:       r.Method,
			StatusCode:   rec.status,
			URL:          r.URL.String(),
			Body:         string(bodyBytes),
			ResponseTime: duration,
		}

		l.Producer.SendMessage(logHttp)

		log.Printf(`{"method": "%s", "status_code": %d, "url": "%s", "body": "%s", "response_time": "%d"}`, logHttp.Method, logHttp.StatusCode, logHttp.URL, logHttp.Body, logHttp.ResponseTime)
	})
}
