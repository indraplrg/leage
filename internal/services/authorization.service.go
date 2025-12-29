package services

import "share-notes-app/internal/repositories"

type AuthorizationService interface {}

type authorizationService struct {
	repo repositories.AuthorizationRepositories 
}

func NewAuthorizationsService(repo repositories.AuthorizationRepositories) AuthorizationService {
	return &authorizationService{repo:repo}
}