package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/milindtheengineer/workout-tracker-server/config"
	"google.golang.org/api/idtoken"
)

func (app *App) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var loginInfo LoginInfo
	if err := json.NewDecoder(r.Body).Decode(&loginInfo); err != nil {
		app.logger.Error().Msgf("Decode: %v", err)
		return
	}
	payload, err := idtoken.Validate(context.Background(), loginInfo.Credential, config.AppConfig.GoogleToken)
	if err != nil {
		app.logger.Error().Msgf("Validate: %v", err)
		return
	}
	user, err := app.db.GetUserByEmail(payload.Claims["email"].(string))
	if err != nil {
		app.logger.Error().Msgf("GetUserByEmail: %v", err)
		return
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: strconv.Itoa(user.Id),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{expirationTime},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.SigningKey))
	if err != nil {
		app.logger.Error().Msgf("%v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		SameSite: http.SameSiteNoneMode,
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if err := DecodeJWT(token.Value); err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func DecodeJWT(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AppConfig.SigningKey), nil
	})
	if err != nil {
		return fmt.Errorf("DecodeJWT: %v", err)
	}

	if !token.Valid {
		return fmt.Errorf("DecodeJWT:: invalid token")
	}

	return nil
}
