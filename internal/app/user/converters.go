package user

import "note-service/internal/pkg/user"

func userToUserResponse(user user.User) UserResponse {
	return UserResponse{user.ID, user.Username}
}
