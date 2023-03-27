package main

import (
	"github.com/joho/godotenv"
	"github.com/mokan-r/place-for-your-thoughts/internal/service"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatal(`error loading env variables`)
	}

	s := service.New()
	s.Start()
}
