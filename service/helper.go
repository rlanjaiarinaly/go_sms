package service

import (
	"net/http"

	"github.com/gorilla/schema"
)

func Helper(dest interface{}, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	decoder := schema.NewDecoder()
	if err := decoder.Decode(dest, r.PostForm); err != nil {
		return err
	}
	return nil

}
