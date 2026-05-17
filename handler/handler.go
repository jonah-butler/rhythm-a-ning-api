package handler

import "rhythmapi/repository"

type RhythmHandler struct {
	repo repository.IRhythmRepository
}

type UserHandler struct {
	repo repository.IUserRepository
}

func NewRhythmHandler(repo repository.IRhythmRepository) *RhythmHandler {
	return &RhythmHandler{repo: repo}
}

func NewUserHandler(repo repository.IUserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}
