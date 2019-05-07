package authcontext

import (
	"context"

	"cloud.google.com/go/datastore"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// PersistenceClient A wrapper aound the datastore client
type PersistenceClient interface {
	AllocateIDs(ctx context.Context, keys []*datastore.Key) ([]*datastore.Key, error)
	Delete(ctx context.Context, key *datastore.Key) error
	Get(ctx context.Context, key *datastore.Key, dst interface{}) (err error)
	GetAll(ctx context.Context, q *datastore.Query, dst interface{}) (keys []*datastore.Key, err error)
	NewTransaction(ctx context.Context, opts ...datastore.TransactionOption) (t *datastore.Transaction, err error)
	Put(ctx context.Context, key *datastore.Key, src interface{}) (*datastore.Key, error)
}

// AuthUser defines user attributes
type AuthUser struct {
	Name  string
	Email string
}

// Context holds all the necessary data
type Context struct {
	AppEngineCtx    context.Context
	DataStoreClient PersistenceClient
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

	context := Context{
		AppEngineCtx:    ctx,
		DataStoreClient: client,
		GinCtx:          c,
	}

	if authUser != nil {
		context.AuthUser = AuthUser{
			Name:  authUser.(*jwt.Token).Claims.(jwt.MapClaims)["name"].(string),
			Email: authUser.(*jwt.Token).Claims.(jwt.MapClaims)["email"].(string),
		}
	}

	return &context, nil
}
