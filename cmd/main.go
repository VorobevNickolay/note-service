package main

import (
	"go.uber.org/zap"
	"note-service/internal/app"
	"note-service/internal/app/note"
	"note-service/internal/app/user"
	notepkg "note-service/internal/pkg/note"
	userpkg "note-service/internal/pkg/user"
	"time"
)

func main() {
	logger, _ := zap.NewProduction()

	userStore := userpkg.NewInMemoryStore()
	userService := userpkg.NewService(userStore)
	userRouter := user.NewRouter(userService, logger.Named("user-router"))

	noteStore := notepkg.NewInMemoryStore(logger.Named("note-store"))
	noteService := notepkg.NewService(noteStore)
	noteRouter := note.NewRouter(noteService, logger.Named("note-router"))
	noteExpService := notepkg.NewExpService(noteStore, 10*time.Second, logger.Named("note-exp-service"))
	go noteExpService.Run()

	router := app.NewRouter(logger.Named("router"), userRouter, noteRouter)
	router.SetUpRouter()
	router.Run()
}
