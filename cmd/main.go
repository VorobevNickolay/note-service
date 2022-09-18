package main

import (
	"note-service/internal/app"
	"note-service/internal/app/note"
	"note-service/internal/app/user"
	notepkg "note-service/internal/pkg/note"
	userpkg "note-service/internal/pkg/user"
)

func main() {
	userStore := userpkg.NewInMemoryStore()
	userService := userpkg.NewService(userStore)
	userRouter := user.NewRouter(userService)
	noteStore := notepkg.NewInMemoryStore()
	noteService := notepkg.NewService(noteStore)
	noteRouter := note.NewRouter(noteService)
	router := app.NewRouter(userRouter, noteRouter)
	router.SetUpRouter()
	router.Run()
}
