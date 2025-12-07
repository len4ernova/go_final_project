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
	h.Logger.Sugar().Info("START  /api/signin", r.Method)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJson(w, reterror{Error: fmt.Sprintf("didn't get body: %v", err)})
		return
	}
	var ps password
	err = json.Unmarshal(body, &ps)
	if err != nil {
		writeJson(w, reterror{Error: fmt.Sprintf("didn't get body: %v", err)})
		return
	}
	passUser := ps.Password
	passEnv := os.Getenv("TODO_PASSWORD")
	if passEnv == passUser {
		// хеш пароля
		hesh := sha256.Sum256([]byte(passEnv))
		// генерация токена
		jwttoken, err := services.GenerateJWT(hesh)
		if err != nil {
			writeJson(w, reterror{Error: fmt.Sprintf("didn't generate jwt token: %v", err)})
			return
		}
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
		writeJson(w, data)
		return
	}
	writeJson(w, reterror{Error: "invalid password"})

}
