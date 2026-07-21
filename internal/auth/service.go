package auth

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

func Authenticate(username string, password string) (string, error) { //db erstellen mit user (schreiben,dass der benutzer nur in seine dir gehenn kann)
	if username != "test" || password != "test" {
		return "", ErrInvalidCredentials
	}

	t, err := generateJWT(username)
	if err != nil {
		log.Printf("Could not generate jwt: %s", err)
		return "", nil
	}
	return t, nil

}

func generateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
	})

	pubToken, err := os.ReadFile("../../token.pem")
	if err != nil {
		log.Println("Could not read public key")
		return "", err
	}
	pub, err := jwt.ParseRSAPrivateKeyFromPEM(pubToken)
	if err != nil {
		return "", err
	}
	return token.SignedString(pub)
}

func ProtectedRoutes(handler func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//check auth
		cookie, err := r.Cookie("auth")
		if err != nil {
			w.WriteHeader(401)
			return
		}

		if !verifyJWT(cookie.Value) {
			w.WriteHeader(401)
			return
		}
		handler(w, r)
	}
}

func verifyJWT(cookie string) bool {
	_, err := jwt.Parse(cookie, func(t *jwt.Token) (any, error) {
		privateKey, err := os.ReadFile("../../token.pub")
		if err != nil {
			return nil, err
		}
		return jwt.ParseRSAPublicKeyFromPEM(privateKey)
	})
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
