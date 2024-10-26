package cookies

import (
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func SetCartCookie(w http.ResponseWriter) string {
	cartId := uuid.New().String()
	encodedCartId := base64.StdEncoding.EncodeToString([]byte(cartId))
	cookie := &http.Cookie{
		Name:     "cartID",
		Value:    encodedCartId,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Domain:   "localhost",
		Expires:  time.Now().Add(1 * time.Hour),
	}
	http.SetCookie(w, cookie)
	log.Println("Cart cookie created:", cookie.Value)
	return cookie.Value
}

func GetCartCookie(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("cartID")
	if err != nil {
		log.Println("Cart cookie not found:", err)
		return SetCartCookie(w)
	}
	decodedCartId, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		log.Println("Error decoding cart ID:", err)
		return SetCartCookie(w)
	}
	log.Println("Cart cookie found:", cookie.Value)
	return string(decodedCartId)
}
