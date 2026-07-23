package auth

import (
	"crypto/rsa"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var pubToken *rsa.PrivateKey

func Init() error {
	privateKey, err := os.ReadFile("../../token.pem")
	if err != nil {
		log.Println("Could not read public key")
		return err
	}

	pub, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return err
	}

	pubToken = pub
	return nil
}
func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
	})

	return token.SignedString(pubToken)
}

func ProtectedRoutes(handler func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//check auth
		cookie, err := r.Cookie("auth")
		if err != nil {
			w.WriteHeader(401)
			return
		}

		if !VerifyJWT(cookie.Value) {
			w.WriteHeader(401)
			return
		}
		handler(w, r)
	}
}

func VerifyJWT(cookie string) bool {
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

func CheckPasswordHash(hashedPw string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(password)) == nil
}
