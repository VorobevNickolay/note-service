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
	userRouter := user.NewRouter(userStore)
	noteStore := notepkg.NewInMemoryStore()
	noteRouter := note.NewRouter(noteStore)
	router := app.NewRouter(userRouter, noteRouter)
	router.SetUpRouter()
	router.Run()
}
