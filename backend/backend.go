package backend

import (
	"github.com/twitter-remake/user/clients"
	"github.com/twitter-remake/user/repository"
)

type Dependency struct {
	clients *clients.Clients
	repo    *repository.Dependency
}

// New creates a new Backend struct
func New(clients *clients.Clients, repo *repository.Dependency) *Dependency {
	return &Dependency{
		clients: clients,
		repo:    repo,
	}
}
