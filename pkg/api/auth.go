package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/len4ernova/go_final_project/pkg/services"
)

type password struct {
	Password string `json:"password"`
}

// аутентификация
func (h *SrvHand) authHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Sugar().Info("START  /api/signin ", r.Method)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Sugar().Error(fmt.Sprintf("didn't get body: %v", err))
		writeJson(w, reterror{Error: fmt.Sprint("didn't get body request")}, http.StatusBadRequest)
		return
	}
	var ps password
	err = json.Unmarshal(body, &ps)
	if err != nil {
		h.Logger.Sugar().Error(fmt.Sprintf("didn't get body: %v", err))
		writeJson(w, reterror{Error: fmt.Sprint("didn't get body request")}, http.StatusBadRequest)
		return
	}
	passUser := ps.Password
	passEnv := os.Getenv("TODO_PASSWORD")
	//fmt.Println("passUser", passUser, "passEnv", passEnv)
	if passEnv == passUser {
		// хеш пароля
		hesh := sha256.Sum256([]byte(passEnv))
		// генерация токена
		jwttoken, err := services.GenerateJWT(hesh)
		if err != nil {
			h.Logger.Sugar().Error(fmt.Sprintf("didn't generate jwt token: %v", err))
			writeJson(w, reterror{Error: fmt.Sprint("didn't generate jwt token")}, http.StatusInternalServerError)
			return
		}
		// fmt.Println("jwttoken=", jwttoken)
		data := struct {
			Token string `json:"token"`
		}{
			Token: jwttoken,
		}
		cookie := http.Cookie{
			Name:  "token",
			Value: jwttoken,
			//Path: "/",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, &cookie)
		h.Logger.Sugar().Info("token was set")
		writeJson(w, data, http.StatusOK)
		return
	}
	h.Logger.Sugar().Error("invalid password")
	writeJson(w, reterror{Error: "invalid password"}, http.StatusUnauthorized)

}
