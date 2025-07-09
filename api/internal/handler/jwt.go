package handler

import (
	"api/voyago/internal/config"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
	"time"
)

type JWTService struct {
	AccessCookie  *http.Cookie
	RefreshCookie *http.Cookie
	Error         error
}

type JWTBody struct {
	Id string
}

func GetInfoFromCookie(r *http.Request, cookieName string, JWTKey string) (*JWTBody, error) {
	token, err := GetTokenFromCookie(r, cookieName)
	if err != nil {
		return nil, err
	}
	claims, err := DecodeJWT(token, JWTKey)
	if err != nil {
		return nil, err
	}

	return &JWTBody{
		Id: claims["id"].(string),
	}, nil
}

func GetTokenFromCookie(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", fmt.Errorf("cookie %s не найден", cookieName)
		}
		return "", fmt.Errorf("ошибка при получении cookie: %w", err)
	}

	return cookie.Value, nil
}

func DecodeJWT(tokenString string, key string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка при разборе токена: %w, token %d,token after Parse %s", err, tokenString, token)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("недействительный токен")
	}

	return claims, nil
}

func GetNewCookies(id uuid.UUID, config config.Config) (jwtService *JWTService) {
	AccessCookie, RefreshCookie, err := createTokens(id, config)
	return &JWTService{
		AccessCookie:  AccessCookie,
		RefreshCookie: RefreshCookie,
		Error:         err,
	}
}

func createTokens(id uuid.UUID, config config.Config) (*http.Cookie, *http.Cookie, error) {
	token, err := createToken(id, []byte(config.JWTKey))
	if err != nil {
		log.Println("Create Token Error: " + err.Error())
		return nil, nil, err
	}

	refreshToken, err := createToken(id, []byte(config.JWTRefreshKey))
	if err != nil {
		log.Println("Create Token Error: " + err.Error())
		return nil, nil, err
	}

	JWTAccessTime, err := strconv.Atoi(config.JWTAccessTime)
	if err != nil {
		return nil, nil, err
	}

	JWTRefreshTime, err := strconv.Atoi(config.JWTRefreshTime)
	if err != nil {
		return nil, nil, err
	}

	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // True в проде (Https)
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Duration(JWTAccessTime) * time.Minute),
	}

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/refresh",
		HttpOnly: true,
		Secure:   false, // True в проде (Https)
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Duration(JWTRefreshTime) * time.Hour * 24),
	}
	return accessCookie, refreshCookie, nil
}

func createToken(id uuid.UUID, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
	})
	return token.SignedString(key)
}
