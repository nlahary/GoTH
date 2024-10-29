package cookies

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var (
	secretKey            = []byte("secret")
	ErrCookieTooLarge    = errors.New("cookie size exceeds limit")
	ErrSignatureTooShort = errors.New("signature too short")
	ErrSignatureMismatch = errors.New("signature mismatch")
	ErrInvalidCookie     = errors.New("invalid cookie")
)

func Write(w http.ResponseWriter, cookie *http.Cookie) error {
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value)) // Encode the cookie value
	if len(cookie.String()) > 4096 {
		return ErrCookieTooLarge
	}
	http.SetCookie(w, cookie) // Set the cookie in the response
	return nil
}

func Read(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	decodedValue, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		switch {
		case err == base64.CorruptInputError(0):
			return "", base64.CorruptInputError(0)
		default:
			return "", ErrInvalidCookie
		}
	}
	return string(decodedValue), nil
}

func WriteSigned(w http.ResponseWriter, cookie *http.Cookie, secretKey []byte) error {
	mac := hmac.New(sha256.New, secretKey) // Generate a new HMAC with the secret key
	mac.Write([]byte(cookie.Name))         // Write the cookie name to the HMAC
	mac.Write([]byte(cookie.Value))        // Write the cookie value to the HMAC
	signature := mac.Sum(nil)              // Generate the HMAC

	cookie.Value = string(signature) + cookie.Value // Prepend the HMAC to the cookie value
	return Write(w, cookie)                         // Write the cookie to the response
}

func ReadSigned(r *http.Request, name string, secretKey []byte) (string, error) {
	signedCookie, err := Read(r, name)
	if err != nil {
		return "", err
	}
	if len(signedCookie) < sha256.Size {
		return "", ErrSignatureTooShort // Ensure the cookie is long enough to contain the HMAC
	}
	signature := signedCookie[:sha256.Size] // Extract the HMAC from the cookie
	value := signedCookie[sha256.Size:]     // Extract the value from the cookie

	mac := hmac.New(sha256.New, secretKey) // Generate a new HMAC with the secret key
	mac.Write([]byte(name))                // Write the cookie name to the HMAC
	mac.Write([]byte(value))               // Write the cookie value to the HMAC
	expectedSignature := mac.Sum(nil)      // Generate the HMAC

	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", ErrSignatureMismatch // Ensure the HMAC matches the expected value
	}
	return value, nil
}

func SetCartCookie(w http.ResponseWriter) string {
	value := uuid.New().String()
	cookie := &http.Cookie{
		Name:     "cartID",
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(1 * time.Hour),
	}

	err := WriteSigned(w, cookie, secretKey)
	if err != nil {
		switch err {
		case ErrCookieTooLarge:
			http.Error(w, "Cookie size exceeds limit", http.StatusBadRequest)
		default:
			http.Error(w, "Server error setting cart cookie", http.StatusInternalServerError)
		}
		log.Println(err)
		return ""
	}

	log.Println("Cart cookie created:", cookie.Value)
	return cookie.Value
}

func GetCartCookie(w http.ResponseWriter, r *http.Request) string {
	value, err := ReadSigned(r, "cartID", secretKey)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return SetCartCookie(w)
		case ErrSignatureMismatch, ErrSignatureTooShort:
			http.Error(w, "Invalid signature cookie", http.StatusBadRequest)
		default:
			http.Error(w, "Server error reading cart cookie", http.StatusInternalServerError)
		}
		log.Println(err)
	}
	log.Println("Cart cookie found:", value)

	return value
}
