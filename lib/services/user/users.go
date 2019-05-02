package user

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"strings"

	"cloud.google.com/go/datastore"
	"github.com/golang/glog"
	"github.com/pedrocelso/go-task/lib/http/authcontext"
)

const (
	index           = `User`
	invalidUserData = `error: invalid User data`
	// Pepper for user authentication
	Pepper          = `AAUIAbhABJb*&!^^@$%^6756nVBZZVBvnGHVAjhM<PE($#^$&())NnwpiwnOW?"|">?UwoUWBK`
)

// Basic defines basic user attributes
type Basic struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Full defines full user attributes
type Full struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Password     string `json:"password"`
	CreationTime int64  `json:"creationTime"`
	UpdateTime   int64  `json:"updateTime"`
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// Create an user
func Create(ctx *authcontext.Context, usr *Full) (*Basic, error) {
	if usr == nil || usr.Email == `` {
		return nil, fmt.Errorf(invalidUserData)
	}

	output, _ := GetByEmail(ctx, usr.Email)

	if output == nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%s%s", usr.Password, Pepper)), bcrypt.DefaultCost)

		if err != nil {
			return nil, err
		}

		usr.Password = string(hashedPassword)
		usr.CreationTime = makeTimestamp()

		key := datastore.NameKey(index, usr.Email, nil)
		_, err = ctx.DataStoreClient.Put(ctx.AppEngineCtx, key, usr)

		if err != nil {
			glog.Errorf("ERROR INSERTING USER: %v", err.Error())
			return nil, err
		}

		output, err = GetByEmail(ctx, usr.Email)
		if err != nil {
			glog.Errorf("ERROR GETTING USER OUTPUT: %v", err.Error())
			return nil, err
		}
		return output, nil
	}
	return output, fmt.Errorf(`User '%v' already exists`, usr.Email)
}

// GetByEmail an user based on its Email
func GetByEmail(ctx *authcontext.Context, email string) (*Basic, error) {
	var usr Basic

	usrFull, err := GetFullByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	usr.Email = usrFull.Email
	usr.Name = usrFull.Name

	return &usr, nil
}

// GetFullByEmail full User info based on its email
func GetFullByEmail(ctx *authcontext.Context, email string) (*Full, error) {
	if email == `` {
		return nil, fmt.Errorf(invalidUserData)
	}
	userKey := datastore.NameKey(index, email, nil)
	var usr Full
	err := ctx.DataStoreClient.Get(ctx.AppEngineCtx, userKey, &usr)

	if err != nil {
		if strings.HasPrefix(err.Error(), `datastore: no such entity`) {
			err = fmt.Errorf(`user '%v' not found`, email)
		}
		return nil, err
	}
	return &usr, nil
}

// GetUsers Fetches all users
func GetUsers(ctx *authcontext.Context) ([]Basic, error) {
	var output []Basic
	q := datastore.NewQuery(index)
	_, err := ctx.DataStoreClient.GetAll(ctx.AppEngineCtx, q, &output)

	if err != nil {
		glog.Errorf("error fetching all users")
		return nil, err
	}

	if len(output) <= 0 {
		return nil, fmt.Errorf("no users found")
	}
	return output, nil
}

// Update user data
func Update(ctx *authcontext.Context, usr *Full) (*Basic, error) {
	if usr == nil || usr.Email == `` {
		return nil, fmt.Errorf(invalidUserData)
	}

	output, _ := GetByEmail(ctx, usr.Email)
	if output != nil {
		key := datastore.NameKey(index, usr.Email, nil)
		_, err := ctx.DataStoreClient.Put(ctx.AppEngineCtx, key, usr)

		if err != nil {
			glog.Errorf("ERROR UPDATING USER: %v", err.Error())
			return nil, err
		}

		output, err = GetByEmail(ctx, usr.Email)
		if err != nil {
			glog.Errorf("ERROR GETTING USER OUTPUT: %v", err.Error())
			return nil, err
		}
		return output, nil
	}
	return nil, fmt.Errorf(`user '%v' not found`, usr.Email)
}

// Delete an user based on its email.
func Delete(ctx *authcontext.Context, email string) error {
	output, _ := GetByEmail(ctx, email)

	if output != nil {
		glog.Infof("Deleting user: %v", email)
		key := datastore.NameKey(index, email, nil)
		err := ctx.DataStoreClient.Delete(ctx.AppEngineCtx, key)

		if err != nil {
			glog.Errorf("ERROR DELETING USER: %v", err.Error())
			return err
		}
		return nil
	}
	return fmt.Errorf("user '%v' don't exist on the database", email)
}
