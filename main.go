// Reference: https://www.thepolyglotdeveloper.com/2017/03/authenticate-a-golang-api-with-json-web-tokens/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JwtToken struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string `json:"message"`
}

func createTokenEndpoint(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
	})
	tokenStr, err := token.SignedString([]byte("secret"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(JwtToken{Token: tokenStr})
}

func protectedEndpoint(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	tok, err := jwt.Parse(params["token"][0], func(tok *jwt.Token) (interface{}, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return []byte("secret"), nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if claims, ok := tok.Claims.(jwt.MapClaims); ok && tok.Valid {
		json.NewEncoder(w).Encode(User{
			Username: claims["username"].(string),
			Password: claims["password"].(string),
		})
	} else {
		json.NewEncoder(w).Encode(Exception{
			Message: "Invalid authorization token.",
		})
	}
}

func handleJWTValidation(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("authorization")
		if authHeader != "" {
			bearerTok := strings.Split(authHeader, " ")
			if len(bearerTok) == 2 {
				tok, err := jwt.Parse(bearerTok[1], func(tok *jwt.Token) (interface{}, error) {
					if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("there war an error")
					}
					return []byte("secret"), nil
				})
				if err != nil {
					json.NewEncoder(w).Encode(Exception{
						Message: err.Error(),
					})
					return
				}

				if tok.Valid {
					ctx := context.Background()
					val := context.WithValue(ctx, "decoded", tok.Claims)
					r.WithContext(val)
				} else {
					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
				}
			}
		} else {
			json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
		}
	})
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}
