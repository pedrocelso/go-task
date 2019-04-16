package authcontext

import (
	"context"

	"cloud.google.com/go/datastore"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthUser defines user attributes
type AuthUser struct {
	Name  string
	Email string
}

// Context holds all the necessary data
type Context struct {
	AppEngineCtx    context.Context
	DataStoreClient *datastore.Client
	GinCtx          *gin.Context
	AuthUser        AuthUser
}

// NewAuthContext creates a new context with required data
func NewAuthContext(c *gin.Context) (*Context, error) {
	authUser := c.Request.Context().Value("user")
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "go-rest-service")

	if err != nil {
		return nil, err
	}

	return &Context{
		AppEngineCtx:    ctx,
		DataStoreClient: client,
		GinCtx:          c,
		AuthUser: AuthUser{
			Name:  authUser.(*jwt.Token).Claims.(jwt.MapClaims)["name"].(string),
			Email: authUser.(*jwt.Token).Claims.(jwt.MapClaims)["email"].(string),
		},
	}, nil
}
