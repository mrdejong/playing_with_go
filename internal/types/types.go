// Package types
package types

import (
	"awesome-go/pkgs/srv"
)

type key int

var UserKey key

type UserForm struct {
	Name     srv.StringField `field:"name:required,min=5,max=100"`
	Email    srv.StringField `field:"email:required,email"`
	Password srv.StringField `field:"password:required"`
}

type AuthForm struct {
	Email    srv.StringField `field:"email:required,email"`
	Password srv.StringField `field:"password:required"`
}

type TodoForm struct {
	Title  srv.StringField `field:"title:required"`
	Status srv.StringField `field:"status:required"`
}
