package user

import "note-service/internal/pkg/user"

func userModelFromUser(user user.User) Model {
	return Model{user.ID, user.Username}
}
