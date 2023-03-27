package service

import (
	"github.com/jmoiron/sqlx"
	"github.com/mokan-r/place-for-your-thoughts/internal/controller"
	"github.com/mokan-r/place-for-your-thoughts/pkg/postgres"
	"github.com/sirupsen/logrus"
)

type Service struct {
	Client *sqlx.DB
}

func New() *Service {
	client, err := postgres.New(&postgres.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "mdaniell",
		Password: "mdaniell",
		DBName:   "superheroes",
		SSLMode:  "disable",
	})

	if err != nil {
		logrus.Fatal(err)
	}

	return &Service{Client: client}
}

func (s *Service) Start() {
	h := controller.New(s.Client)
	h.Run()
}
