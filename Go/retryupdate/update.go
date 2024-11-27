//go:build !solution

package retryupdate

import (
	"errors"
	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	getRequest := kvapi.GetRequest{Key: key}
	getResponse, err := c.Get(&getRequest)
	var newValue string
	var value string
	var version uuid.UUID
	var apiErr *kvapi.APIError
	var authErr *kvapi.AuthError
	var conflictErr *kvapi.ConflictError
	var oldVersion uuid.UUID

	for {
		if errors.As(err, &apiErr) {
			if errors.As(apiErr.Err, &authErr) {
				return err
			} else if errors.As(apiErr.Err, &conflictErr) {
				if conflictErr.ExpectedVersion == oldVersion {
					return nil
				}
				getResponse, err = c.Get(&getRequest)
			} else if errors.Is(err, kvapi.ErrKeyNotFound) {
				newValue, err = updateFn(nil)
				if err != nil {
					return err
				}
				version = uuid.Must(uuid.NewV4())
				setRequest := kvapi.SetRequest{Key: key, Value: newValue, OldVersion: uuid.UUID{}, NewVersion: version}
				_, err = c.Set(&setRequest)
				if err == nil {
					return err
				}
			} else if apiErr.Method == "get" {
				getResponse, err = c.Get(&getRequest)
			} else if apiErr.Method == "set" {
				value = getResponse.Value
				version = getResponse.Version
				newValue, err = updateFn(&value)
				if err != nil {
					return err
				}
				newVersion := uuid.Must(uuid.NewV4())
				setRequest := kvapi.SetRequest{Key: key, Value: newValue, OldVersion: version, NewVersion: newVersion}
				_, err = c.Set(&setRequest)
				if err == nil {
					return err
				}
			}
		} else {
			value = getResponse.Value
			version = getResponse.Version
			newValue, err = updateFn(&value)
			if err != nil {
				return err
			}
			newVersion := uuid.Must(uuid.NewV4())
			setRequest := kvapi.SetRequest{Key: key, Value: newValue, OldVersion: version, NewVersion: newVersion}
			version = newVersion
			oldVersion = version
			_, err = c.Set(&setRequest)
			if err == nil {
				return err
			}
		}
	}
}
