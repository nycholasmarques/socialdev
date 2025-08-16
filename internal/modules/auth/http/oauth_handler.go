package auth_handlers

import (
	"context"
	"encoding/json"
	"net/http"

	configs "github.com/nycholasmarques/socialdev/config"
	services "github.com/nycholasmarques/socialdev/internal/modules/auth/service"
	pkg "github.com/nycholasmarques/socialdev/internal/pkg/auth"
	"github.com/nycholasmarques/socialdev/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OOAuthHandler struct {
	OAuthHandler services.AuthService
}

func NewOOAuthHandler(OAuthHandler services.AuthService) *OOAuthHandler {
	return &OOAuthHandler{
		OAuthHandler: OAuthHandler,
	}
}

var cfg = configs.LoadConfig()

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  cfg.GoogleRedirectURL,
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	randomState = "random"
)

func (h *OOAuthHandler) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.NewErrorResponse(w, http.StatusMethodNotAllowed, "Método não permitido", "")
		return
	}

	url := googleOauthConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *OOAuthHandler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.NewErrorResponse(w, http.StatusMethodNotAllowed, "Método não permitido", "")
		return
	}

	if r.FormValue("state") != randomState {
		utils.NewErrorResponse(w, http.StatusBadRequest, "Erro na autenticação", "")
		return
	}

	context := context.Background()

	token, err := googleOauthConfig.Exchange(context, r.FormValue("code"))
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Erro ao obter token", err.Error())
		return
	}

	client := googleOauthConfig.Client(context, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Erro ao obter informações do usuário", err.Error())
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Erro ao decodificar informações do usuário", "")
		return
	}

	email, ok := userInfo["email"].(string)
	if !ok {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Erro na autenticação", "")
		return
	}

	name, ok := userInfo["name"].(string)
	if !ok {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Erro na autenticação", "")
		return
	}

	user, err := h.OAuthHandler.LoginUser(email)
	if err != nil {
		user, err = h.OAuthHandler.CreateUserWithGoogle(name, email, "", userInfo["id"].(string), userInfo["picture"].(string))
		if err != nil {
			utils.NewErrorResponse(w, http.StatusInternalServerError, "Erro ao criar usuário", err.Error())
			return
		}
	}

	jwtToken, err := pkg.CreateToken(user.UserID, true)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, "Erro ao criar token", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": jwtToken,
		"user":  user,
	})
}
