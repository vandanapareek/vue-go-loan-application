package handler

import (
	"encoding/json"
	"fmt"
	"loan-api/errors"
	model "loan-api/internal/models"
	"loan-api/internal/service/user"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type userHandler struct {
	service user.UserService
}

func NewUserHandler(s user.UserService) *userHandler {
	return &userHandler{
		service: s,
	}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var jwtKey = []byte("secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Resp struct {
	UserInfo *model.UserInfo `json:"user_info"`
	Token    string          `json:"token"`
}

func (h userHandler) Login(writer http.ResponseWriter, request *http.Request) {

	reqDump, _ := httputil.DumpRequest(request, true)
	log.Printf("REQUEST:\n%s", string(reqDump))

	decoder := json.NewDecoder(request.Body)

	var req loginRequest
	err := decoder.Decode(&req)

	if err != nil {
		log.Printf("Error decoding request: %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	user, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		log.Printf("Login issue: %s", err)
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte("Unauthorized"))
		return
	}

	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error generating JWT token %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	userInfo, err := h.service.GetUserInfo(claims.Username)
	if err != nil {
		fmt.Println("err getting userInfo")
		errors.JSONError(writer, err, http.StatusUnprocessableEntity)
		return
	}
	resp := Resp{
		UserInfo: userInfo,
		Token:    tokenString,
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(resp)

	log.Printf("User %s succesfully logged in", req.Username)
	log.Printf("Generated JWT Token: %s", tokenString)
}
