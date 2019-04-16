package user

import (
	"fmt"

	"strings"

	"cloud.google.com/go/datastore"
	"github.com/golang/glog"
	"github.com/pedrocelso/go-task/lib/http/authcontext"
)

const (
	index           = `User`
	invalidUserData = `error: invalid User data`
)

// User defines user attributes
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Create an user
func Create(c *authcontext.Context, usr *User) (*User, error) {
	var output *User
	if usr == nil || usr.Email == `` {
		return nil, fmt.Errorf(invalidUserData)
	}

	output, _ = GetByEmail(c, usr.Email)

	if output == nil {
		key := datastore.NameKey(index, usr.Email, nil)
		_, err := c.DataStoreClient.Put(c.AppEngineCtx, key, usr)

		if err != nil {
			glog.Errorf("ERROR INSERTING USER: %v", err.Error())
			return nil, err
		}

		output, err = GetByEmail(c, usr.Email)
		if err != nil {
			glog.Errorf("ERROR GETTING USER OUTPUT: %v", err.Error())
			return nil, err
		}
		return output, nil
	}
	glog.Infof("User was previously saved: %v", usr.Email)
	return output, nil
}

// GetByEmail an user based on its Email
func GetByEmail(c *authcontext.Context, email string) (*User, error) {
	if email == `` {
		return nil, fmt.Errorf(invalidUserData)
	}
	userKey := datastore.NameKey(index, email, nil)
	var usr User
	err := c.DataStoreClient.Get(c.AppEngineCtx, userKey, &usr)

	if err != nil {
		if strings.HasPrefix(err.Error(), `datastore: no such entity`) {
			err = fmt.Errorf(`user '%v' not found`, email)
		}
		return nil, err
	}
	return &usr, nil
}

// GetUsers Fetches all users
func GetUsers(c *authcontext.Context) ([]User, error) {
	glog.Infof("GETTING ALL USERS FOR %s/%s", c.AuthUser.Name, c.AuthUser.Email)
	var output []User
	q := datastore.NewQuery(index)
	_, err := c.DataStoreClient.GetAll(c.AppEngineCtx, q, &output)

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
func Update(c *authcontext.Context, usr *User) (*User, error) {
	if usr == nil || usr.Email == `` {
		return nil, fmt.Errorf(invalidUserData)
	}

	output, _ := GetByEmail(c, usr.Email)
	if output != nil {
		key := datastore.NameKey(index, usr.Email, nil)
		_, err := c.DataStoreClient.Put(c.AppEngineCtx, key, usr)

		if err != nil {
			glog.Errorf("ERROR UPDATING USER: %v", err.Error())
			return nil, err
		}

		output, err = GetByEmail(c, usr.Email)
		if err != nil {
			glog.Errorf("ERROR GETTING USER OUTPUT: %v", err.Error())
			return nil, err
		}
		return output, nil
	}
	return nil, fmt.Errorf(`user '%v' not found`, usr.Email)
}

// Delete an user based on its email.
func Delete(c *authcontext.Context, email string) error {
	var output *User
	output, _ = GetByEmail(c, email)

	if output != nil {
		glog.Infof("Deleting user: %v", email)
		key := datastore.NameKey(index, email, nil)
		err := c.DataStoreClient.Delete(c.AppEngineCtx, key)

		if err != nil {
			glog.Errorf("ERROR DELETING USER: %v", err.Error())
			return err
		}
		return nil
	}
	return fmt.Errorf("user '%v' don't exist on the database", email)
}
