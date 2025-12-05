package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/len4ernova/go_final_project/pkg/services"
)

func (h *SrvHand) authHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Sugar().Info("START  /api/signin", r.Method)
	passUser := r.FormValue("password")

	passEnv := os.Getenv("TODO_PASSWORD")

	if passEnv == passUser {
		encr, err := services.EncryptPass(services.Key, passUser)
		if err != nil {
			writeJson(w, reterror{Error: fmt.Sprintf("didn't get encrypt pass: %v", err)})
			return
		}
		jwttoken, err := services.GenerateJWT(encr)
		if err != nil {
			writeJson(w, reterror{Error: fmt.Sprintf("didn't generate jwt token: %v", err)})
			return
		}
		data := struct {
			Token string `json:"token"`
		}{
			Token: jwttoken,
		}

		fmt.Println("jwt token", jwttoken)
		writeJson(w, data)
	}
	writeJson(w, reterror{Error: "invalid password"})

}
