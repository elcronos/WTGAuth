package main

//  Generate RSA signing files via shell (adjust as needed):
//
//  $ openssl genrsa -out app.rsa 1024
//  $ openssl rsa -in app.rsa -pubout > app.rsa.pub
//
// Code borrowed and modified from the following sources:

//

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	db "./db/database"
	. "./db/dataobjects"
	"crypto/rsa"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	_ "github.com/adam-hanna/goLang-jwt-auth-example/db"
	"golang.org/x/crypto/bcrypt"
)

const (
	// For simplicity these files are in the same folder as the app binary.
	// You shouldn't do this in production.
	privKeyPath = "./keys/app.rsa"
	pubKeyPath  = "./keys/app.rsa.pub"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func initKeys() {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

func StartServer() {
	// Non-Protected Endpoint(s)
	http.HandleFunc("/auth/login", LoginHandler)
	http.HandleFunc("/auth/signup", SignupHandler)

	// Protected Endpoints
	http.Handle("/auth/validate", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(ProtectedHandler)),
	))
	//Options
	log.Println("Now listening...")
	http.ListenAndServe(":8888", nil)
}

func main() {
	initKeys()
	StartServer()
}

func CORS(w http.ResponseWriter) (http.ResponseWriter){
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	return w
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{"Gained access to protected resource"}
	JsonResponse(response, w)
}

func OptionsAccess(w http.ResponseWriter, r *http.Request) (http.ResponseWriter){
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers","Content-Type, Access-Control-Allow-Headers")
		w.Header().Set("Allow", "GET,POST,OPTIONS")
	}
	return w
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	w = CORS(w)
	w = OptionsAccess(w, r)
	var user User
	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Bad request parameter")
			return
		}
		//Insert new user
		p := user.Password
		bs, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)

		user.Password = string(bs)
		user.CreatedAt, user.UpdatedAt = time.Now(),time.Now()
		db.DB.Create(&user)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var userCredentials UserCredentials
	w = CORS(w)
	w = OptionsAccess(w,r)

	if r.Method == "POST"{
		err := json.NewDecoder(r.Body).Decode(&userCredentials)

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "Error in request")
			return
		}
		var user User
		db.DB.Where("username = ?", strings.ToLower(userCredentials.Username)).First(&user)

		if len(strings.ToLower(user.Username)) > 0 {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userCredentials.Password))

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}else{
				w.WriteHeader(http.StatusOK)
			}

			token := jwt.New(jwt.SigningMethodRS256)
			claims := make(jwt.MapClaims)
			claims["user"] = userCredentials.Username
			claims["role"] = user.Role
			claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
			claims["iat"] = time.Now().Unix()
			token.Claims = claims

			tokenString, err := token.SignedString(signKey)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "Error while signing the token")
				fatal(err)
			}

			response := Token{tokenString}
			JsonResponse(response, w)
		} else{
			w.WriteHeader(http.StatusForbidden)
			fmt.Println("Error logging in")
			fmt.Fprint(w, "Invalid credentials")
			return
		}
	}

}

func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w = CORS(w)
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})

	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(w, "The token is expired")
				return
			case jwt.ValidationErrorSignatureInvalid:
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(w, "The sign does not match")
				return
			default:
				fmt.Fprintln(w, "The token is invalid")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		default:
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "The token is invalid")
			return
		}
	}

	if err == nil {
		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

}

func JsonResponse(response interface{}, w http.ResponseWriter) {
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
