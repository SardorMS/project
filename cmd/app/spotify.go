package app

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/SardorMS/project/pkg/types"
	"github.com/google/uuid"
)

// Переменная адрес перенаправления вашего приложения.
var (
	defaultURI = "http://localhost:8888/callback"
	clientID   = "0bbb4ddc8b544c1393d5595175974de9"
)

// handleVerifyApp - метод для верификации приложения.
func (s *Server) handleVerify(writer http.ResponseWriter, request *http.Request) {

	if responseType := RequestURL(request, "response_type"); responseType != "code" {
		log.Println(responseType)
		http.Error(writer, "Response_Type must be: code!", http.StatusBadRequest)
		return
	}

	if codeChallengeMethod := RequestURL(request, "code_challenge_method"); codeChallengeMethod != "S256" {
		http.Error(writer, "Code_Challenge_Method must be: S256!", http.StatusBadRequest)
		return
	}

	redirectURI := RequestURL(request, "redirect_uri")
	decodedURI, err := url.PathUnescape(redirectURI)
	if decodedURI != defaultURI || err != nil {
		http.Error(writer, "Redirect URI is Not Matched or Empty!", http.StatusBadRequest)
		return
	}

	id := RequestURL(request, "client_id")
	if id != clientID {
		http.Error(writer, "Client ID is Wrong Or Empty!", http.StatusBadRequest)
		return

	}

	codeChallenge := RequestURL(request, "code_challenge")
	codeVerifier := RequestURL(request, "code_verifier")
	state := RequestURL(request, "state")

	scope := RequestURL(request, "scope")
	if scope == "" {
		http.Error(writer, "Scope is Empty!", http.StatusBadRequest)
		return
	}
	decodedScope, _ := url.PathUnescape(scope)
	parsedScope := strings.Split(decodedScope, " ")
	if err != nil {
		http.Error(writer, "Can not parse Scope:", http.StatusBadRequest)
		return
	}

	if codeChallenge == "" || state == "" || codeVerifier == "" {
		http.Error(writer, "Query Params is Empty!", http.StatusBadRequest)
		return
	}

	generatedCode := uuid.New().String()
	item := &types.DataApp{
		ClientID:      id,
		RedirectURI:   decodedURI,
		CodeChallenge: codeChallenge,
		CodeVerifier:  codeVerifier,
		State:         state,
		Scope:         parsedScope,
		Code:          generatedCode,
	}

	result, err := s.spotSvc.RegisterApp(request.Context(), item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	respondJSON(writer, result)
}

// handleUserRegister - метод для проверки поля статуса для пользователя.
func (s *Server) handleUserCheck(writer http.ResponseWriter, request *http.Request) {

	if status := RequestURL(request, "error"); status == "access_denied" {
		http.Error(writer, "Authorization failed: access_denied!", http.StatusBadRequest)
		return
	}

	if id := RequestURL(request, "client_id"); id != clientID || id == "" {
		http.Error(writer, "Client_id is Wrong Or Empty!", http.StatusBadRequest)
		return
	}

	code := RequestURL(request, "code")
	state := RequestURL(request, "state")

	if code == "" || state == "" {
		http.Error(writer, "Query Params is Empty!", http.StatusBadRequest)
		return
	}

	checkedState, err := s.spotSvc.CheckState(request.Context(), code, state)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	item := map[string]interface{}{"code": code, "state": checkedState}
	respondJSON(writer, item)
}

// handleTokenExchange - метод для обмена токена доступа для конкретного пользователя.
func (s *Server) handleTokenExchange(writer http.ResponseWriter, request *http.Request) {

	headerType := request.Header.Get("Content-Type")
	if headerType != "application/x-www-form-urlencoded" {
		http.Error(writer, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	if grantType := request.PostFormValue("grant_type"); grantType != "authorization_code" {
		log.Println(grantType)
		http.Error(writer, "Grant Type: authorization_code is not set!", http.StatusBadRequest)
		return
	}

	redirectURI := request.PostFormValue("redirect_uri")
	if redirectURI != defaultURI || redirectURI == "" {
		http.Error(writer, "Redirect URI is Not Matched Or Empty!", http.StatusBadRequest)
		return
	}

	id := request.PostFormValue("client_id")
	if id != clientID || id == "" {
		http.Error(writer, "Client ID is Wrong Or Empty!", http.StatusBadRequest)
		return
	}
	code := request.PostFormValue("code")
	codeVerifier := request.PostFormValue("code_verifier")

	if code == "" || codeVerifier == "" {
		http.Error(writer, "Query Params is Empty!", http.StatusBadRequest)
		return
	}

	item := &types.Token{
		ClientID:     id,
		RedirectURI:  redirectURI,
		Code:         code,
		CodeVerifier: codeVerifier,
	}

	result, err := s.spotSvc.Token(request.Context(), item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(writer, result)
}

// handleTokenRefresh - метод для продления или обновления токена для конкретного пользователя.
func (s *Server) handleTokenRefresh(writer http.ResponseWriter, request *http.Request) {

	headerType := request.Header.Get("Content-Type")
	if headerType != "application/x-www-form-urlencoded" {
		http.Error(writer, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	if grantType := request.PostFormValue("grant_type"); grantType != "refresh_token" {
		log.Println(grantType)
		http.Error(writer, "Grant Type: refresh_token is Not Set!", http.StatusBadRequest)
		return
	}

	if id := request.PostFormValue("client_id"); id != clientID {
		http.Error(writer, "Clietn ID is Wrong Or Empty!", http.StatusBadRequest)
		return
	}


	refreshToken := request.PostFormValue("refresh_token")
	if refreshToken == "" {
		http.Error(writer, "Refresh Token is Empty!", http.StatusBadRequest)
		return
	}

	item := &types.TokenRefresh{
		RefreshToken: refreshToken,
	}

	result, err := s.spotSvc.TokenUpdate(request.Context(), item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(writer, result)
}
