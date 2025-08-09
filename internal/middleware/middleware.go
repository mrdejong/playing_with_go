// Package middleware
package middleware

import (
	"awesome-go/internal/models"
	"awesome-go/internal/service"
	"awesome-go/internal/types"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
// 	return hmacSampleSecret, nil
// }, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
// if err != nil {
// 	log.Fatal(err)
// }
//
// if claims, ok := token.Claims.(jwt.MapClaims); ok {
// 	fmt.Println(claims["foo"], claims["nbf"])
// } else {
// 	fmt.Println(err)
// }

type Middlware struct {
	service *service.Service
}

func New(s *service.Service) Middlware {
	return Middlware{
		service: s,
	}
}

func (m *Middlware) decodeJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte("4E4bSGLxqMYcppapKkRQrhkaNga5pcnDfSy3QDbb"), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (m *Middlware) validateJWT(jwtToken *jwt.Token) (models.Session, bool) {
	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok {

		session, err := m.service.GetSessionByToken(claims["token"].(string))
		if err != nil {
			// There is no session, it might has cleared already?
			return models.Session{}, false
		}

		if time.Now().After(session.ExpiresOn) {
			fmt.Println("Expired in db")
			// Again already exprired
			return session, false
		}

		return session, true
	}
	return models.Session{}, false
}

func ctx(user models.User) context.Context {
	return context.WithValue(context.Background(), types.UserKey, user)
}

func (m *Middlware) validateAuth(c *fiber.Ctx) (models.Session, bool) {
	jwtToken := c.Cookies("auth", "")
	if jwtToken == "" {
		fmt.Println("No token")
		return models.Session{}, false
	}

	jt, err := m.decodeJWT(jwtToken)
	if err != nil {
		fmt.Printf("Can't decode %v\n", err)
		return models.Session{}, false
	}

	session, ok := m.validateJWT(jt)
	if !ok {
		fmt.Println("Session doesn't exist")
		return models.Session{}, false
	}

	return session, true
}

func (m *Middlware) LoadAuth(c *fiber.Ctx) error {
	c.SetUserContext(ctx(models.User{}))
	session, ok := m.validateAuth(c)
	if ok {
		c.SetUserContext(ctx(session.User))
	}
	return c.Next()
}

func (m *Middlware) RequireAuth(c *fiber.Ctx) error {
	_, ok := m.validateAuth(c)
	if !ok {
		return c.Redirect("/")
	}
	return c.Next()
}
