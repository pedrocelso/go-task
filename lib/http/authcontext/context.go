package authcontext

import (
	"context"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
)

// AuthUser defines user attributes
type AuthUser struct {
	Name  string
	Email string
}

// Context holds all the necessary data
type Context struct {
	AppEngineCtx context.Context
	GinCtx       *gin.Context
	AuthUser     AuthUser
}

// NewAuthContext creates a new context with required data
func NewAuthContext(c *gin.Context) Context {
	authUser := c.Request.Context().Value("user")

	return Context{
		AppEngineCtx: appengine.NewContext(c.Request),
		GinCtx:       c,
		AuthUser: AuthUser{
			Name:  authUser.(*jwt.Token).Claims.(jwt.MapClaims)["name"].(string),
			Email: authUser.(*jwt.Token).Claims.(jwt.MapClaims)["email"].(string),
		},
	}
}
