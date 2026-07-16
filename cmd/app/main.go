package main

import (
	"log"

	"github.com/azmiagr/garudahacks-hackathon/internal/handler/rest"
	"github.com/azmiagr/garudahacks-hackathon/internal/repository"
	"github.com/azmiagr/garudahacks-hackathon/internal/service"
	"github.com/azmiagr/garudahacks-hackathon/pkg/bcrypt"
	"github.com/azmiagr/garudahacks-hackathon/pkg/config"
	"github.com/azmiagr/garudahacks-hackathon/pkg/database/mariadb"
	"github.com/azmiagr/garudahacks-hackathon/pkg/jwt"
	"github.com/azmiagr/garudahacks-hackathon/pkg/middleware"
	"github.com/azmiagr/garudahacks-hackathon/pkg/supabase"
)

func main() {
	config.LoadEnvironment()

	db, err := mariadb.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	err = mariadb.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	err = mariadb.Seed(db)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(db)
	bcrypt := bcrypt.Init()
	jwt := jwt.Init()
	supabase := supabase.Init()
	midtransConfig := config.LoadMidtransConfig()
	svc := service.NewService(repo, bcrypt, jwt, supabase, midtransConfig)

	middleware := middleware.Init(svc, jwt)
	r := rest.NewRest(svc, middleware)
	r.MountEndpoint()

	r.Run()
}
