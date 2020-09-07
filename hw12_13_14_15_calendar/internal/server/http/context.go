package http_server

import (
	"calendar/internal/repository"
	"context"
	"log"
)

func getUserId(ctx context.Context) repository.ID {
	userId, ok := ctx.Value(userIdKey).(repository.ID)

	if !ok {
		log.Println("userId is missing: ", userId)
	}

	return userId
}

func getRepository(ctx context.Context) repository.BaseRepo {
	repo, ok := ctx.Value(repositoryKey).(repository.BaseRepo)

	if !ok {
		log.Println("repository is missing: ", repo)
	}

	return repo
}