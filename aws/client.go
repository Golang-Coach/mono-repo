package aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/golang-coach/mono-repo/constants"
)

// DBClient return the instance of database client
func DBClient(logEnabled bool) *gorm.DB {
	connectionString := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True`,
		os.Getenv(constants.DatabaseUsername),
		os.Getenv(constants.DatabasePassword),
		os.Getenv(constants.DatabaseHostName),
		os.Getenv(constants.DatabaseName))
	db, err := gorm.Open("mysql", connectionString)
	if err != nil {
		panic("failed to connect database")
	}
	db.LogMode(logEnabled)
	return db
}

func Session() client.ConfigProvider {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv(constants.AwsRegionName))},
	)
	if err != nil {
		panic(err)
	}
	return sess
}
