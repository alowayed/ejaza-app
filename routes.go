package main

import "github.com/kataras/iris"

func certSubmit(ctx iris.Context) {
	var user User
	ctx.ReadJSON(&user)
	ctx.Writef("%s %s is %d years old and comes from %s", user.Firstname, user.Lastname, user.Age, user.City)
}
