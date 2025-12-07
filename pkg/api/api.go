package api

import (
	"crypto/sha256"
	"database/sql"
	"net/http"
	"os"

	"github.com/len4ernova/go_final_project/pkg/services"
	"go.uber.org/zap"
)

type SrvHand struct {
	Logger *zap.Logger
	DB     *sql.DB
}

const pattern = "20060102"

// Init - конечные точки и вызов обработчиков.
func (h *SrvHand) Init(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/signin", h.authHandler)
	mux.HandleFunc("/api/nextdate", h.nextDayHandler)              // рассчитать следующую дату
	mux.HandleFunc("GET /api/task", h.getTaskHandler)              // вернуть задачу по ID
	mux.HandleFunc("POST /api/task", auth(h.addTaskHandler))       // добавить задачу в формате JSON
	mux.HandleFunc("PUT /api/task", h.putTaskHandler)              // изменить задачу по ID
	mux.HandleFunc("DELETE /api/task", h.deleteTaskHandler)        // удалить задачу по ID
	mux.HandleFunc("GET /api/tasks", auth(h.tasksHandler))         // вернуть список ближайших задач в формате JSON
	mux.HandleFunc("POST /api/task/done", auth(h.doneTaskHandler)) // задача выполнена

}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		passEnv := os.Getenv("TODO_PASSWORD")
		if len(passEnv) > 0 {
			var jwtCookie string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtCookie = cookie.Value
			}
			claims, err := services.ParseToken(jwtCookie)
			if err != nil {
				http.Error(w, "Неверный токен", http.StatusUnauthorized)
				return
			}

			hesh := sha256.Sum256([]byte(passEnv))

			var valid bool
			valid = hesh == claims.Hesh
			// валидация и проверки JWT-токена
			if valid {
				next.ServeHTTP(w, r)
				return
			}

			if !valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
