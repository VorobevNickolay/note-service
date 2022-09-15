package main

import (
	"note-service/internal/app"
	"note-service/internal/app/user"
	userpkg "note-service/internal/pkg/user"
)

func main() {
	userStore := userpkg.NewInMemoryStore()
	userRouter := user.NewRouter(userStore)
	router := app.NewRouter(userRouter)
	router.SetUpRouter()
	router.Run()
}
