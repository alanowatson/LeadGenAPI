package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
    "github.com/joho/godotenv"
    "os"

	"github.com/alanowatson/LeadGenAPI/pkg/util"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey []byte

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
    if len(jwtKey) == 0 {
        log.Fatal("JWT_SECRET_KEY not set in .env file")
    }
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("JWT Key used for verification: %s", string(jwtKey))

        authHeader := r.Header.Get("Authorization")
        log.Printf("Auth header: %s", authHeader)

        tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return jwtKey, nil
        })

        if err != nil {
            log.Printf("Token parsing error: %v", err)
            util.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
            return
        }

        if !token.Valid {
            log.Printf("Token is invalid")
            util.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
            return
        }

        log.Printf("Token is valid")
        next.ServeHTTP(w, r)
    }
}

func GenerateToken(username string) (string, error) {
    log.Printf("JWT Key used for signing: %s", string(jwtKey))

    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["username"] = username
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return "", err
    }

    log.Printf("Generated token: %s", tokenString)
    return tokenString, nil
}
