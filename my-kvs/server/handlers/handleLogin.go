package handlers

import (
	"github.com/golang-jwt/jwt"
	"my-kvs/server/logger"
	"net/http"
	"strings"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var jwtKey = []byte("my_secret_key")

var loginDetails = map[string]string{
	"user_a": "passwordA",
	"user_b": "passwordB",
	"user_c": "passwordC",
	"admin":  "Password1",
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	if r.Method == http.MethodGet {
		username, password, ok := r.BasicAuth()
		if !ok {
			logger.ErrorLogger.Println("Error parsing basic auth", ok)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		expectedPassword, ok := loginDetails[username]
		if !ok {
			logger.ErrorLogger.Printf("Unknown user: %s", username)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if password != expectedPassword {
			logger.ErrorLogger.Println("Provided password is incorrect")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Create the JWT claims, which includes the username and no expiry time
		claims := &Claims{
			Username: username,
			StandardClaims: jwt.StandardClaims{
				Issuer: "UserJWTService",
			},
		}

		// Declare the token with the algorithm used for signing,
		// That is we create a token from the claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Create the JWT token - do this by signing the token
		// using a secure private key (jwtKey) it will create
		// {header}.{payload}.{signature}
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			logger.ErrorLogger.Printf("Error creating the token: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Set the client cookie for "token" as the JWT we just generated
		// we don't set an expiry time
		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: tokenString,
		})

		_, err = w.Write([]byte("Bearer " + tokenString))
		if err != nil {
			logger.ErrorLogger.Printf("error writing response %v: %s", tokenString, err.Error())
		}
	} else {
		logger.ErrorLogger.Printf("Wrong HTTP method (should be GET), not %s\n", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func Authorise(w http.ResponseWriter, r *http.Request) (bool, string) {
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		logger.ErrorLogger.Printf("missing bearer token")
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte("missing bearer token"))
		if err != nil {
			logger.ErrorLogger.Printf("error writing 'missing bearer token' response: %s", err.Error())
		}
		return false, ""
	}

	tokenSlice := strings.Split(reqToken, "Bearer ")
	if len(tokenSlice) != 2 || tokenSlice[0] != "" {
		logger.ErrorLogger.Printf("invalid bearer token")
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("invalid bearer token"))
		if err != nil {
			logger.ErrorLogger.Printf("error writing 'invalid bearer token' response: %s", err.Error())
		}
		return false, ""
	}

	tokenString := tokenSlice[1]

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT token string
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			logger.ErrorLogger.Printf("Error signature invalid %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return false, ""
		}
		logger.ErrorLogger.Printf("Error processing JWT token %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return false, ""
	}
	if !token.Valid {
		logger.ErrorLogger.Printf("Error invalid token %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return false, ""
	}

	return true, claims.Username
}
