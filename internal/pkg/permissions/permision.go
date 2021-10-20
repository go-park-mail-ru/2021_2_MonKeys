package permissions

// import (
// 	"dripapp/configs"
// 	// "dripapp/internal/pkg/responses"
// 	"net/http"
// 	"time"

// 	"github.com/gorilla/sessions"
// 	uuid "github.com/nu7hatch/gouuid"
// )

// func CheckAuthenticated(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {

// 			session := r.Context().Value(configs.SessionStaffID).(*sessions.Session)

// 			staffID, found := session.Values["userID"]
// 			if !found || staffID == -1 {
// 				responses.SendForbidden(w)
// 				return
// 			}

// 			next.ServeHTTP(w, r)
// 		})
// }

// func generateCsrfLogic(w http.ResponseWriter) {
// 	csrf, err := uuid.NewV4()
// 	if err != nil {
// 		responses.SendForbidden(w)
// 		return
// 	}
// 	timeDelta := time.Now().Add(time.Hour * 24 * 30)
// 	cookie1 := &http.Cookie{Name: "csrf", Value: csrf.String(), Path: "/", HttpOnly: true, Expires: timeDelta}

// 	http.SetCookie(w, cookie1)
// 	w.Header().Set("csrf", csrf.String())

// }

// func SetCSRF(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {
// 			generateCsrfLogic(w)
// 			next.ServeHTTP(w, r)
// 		})
// }

// func CheckCSRF(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {
// 			csrf := r.Header.Get("X-Csrf-Token")
// 			csrfCookie, err := r.Cookie("csrf")

// 			if err != nil || csrf == "" || csrfCookie.Value == "" || csrfCookie.Value != csrf {
// 				responses.SendSingleError("csrf-protection", w)
// 				return
// 			}
// 			generateCsrfLogic(w)
// 			next.ServeHTTP(w, r)
// 		})

// }
