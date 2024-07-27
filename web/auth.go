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

type contextKey string

var contextKeyUserID = contextKey("userID")

const userIDKey = "userID"

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

func (app *App) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userID, err := app.decodeJWT(token.Value)
		if err != nil {
			app.logger.Error().Msgf("decode: %w", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *App) decodeJWT(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AppConfig.SigningKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("DecodeJWT: %v", err)
	}
	if !token.Valid {
		return "", fmt.Errorf("DecodeJWT:: invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("DecodeJWT:: Invalid claims")
	}
	userID, ok := claims["UserId"].(string)
	if !ok {
		return "", fmt.Errorf("DecodeJWT:: User ID not found in claims")
	}
	return userID, nil
}
