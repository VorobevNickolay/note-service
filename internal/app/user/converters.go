package user

import "note-service/internal/pkg/user"

func userModelFromUser(user user.User) UserModel {
	return UserModel{user.ID, user.Username}
}
