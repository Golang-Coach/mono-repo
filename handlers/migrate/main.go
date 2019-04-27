package main

import (
	"context"

	"github.com/golang-coach/mono-repo/aws"
	"github.com/golang-coach/mono-repo/models"
)

func HandleRequest(_ context.Context) error {
	db := aws.DBClient(true)
	defer db.Close()

	err := db.AutoMigrate(&models.User{}, &models.Repository{}, &models.Categories{}, &models.Tags{}).Error

	if err != nil {
		panic(err)
	}
	return err
}

func main() {
	//lambda.Start(HandleRequest)
	HandleRequest(nil)
}
