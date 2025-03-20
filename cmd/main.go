package main

import (
	"context"

	ctrl "github.com/iamcaleberic/randovec/internal/controllers"
	intLogger "github.com/iamcaleberic/randovec/internal/logger"
	"github.com/iamcaleberic/randovec/internal/utils"

	"go.uber.org/zap"
)

var logger = intLogger.InitLogger()

func main() {

	logger.Info("validating env")
	err := ValidateEnv()
	if err != nil {
		panic(err)
	}

	logger.Info("attempting connection to weaviate instance")
	client, err := ctrl.CreateWeaviateClient()
	if err != nil {
		logger.Error("failed to create weaviate client", zap.Error(err))
	}

	ctx := context.Background()

	logger.Info("creating schema")

	_ = ctrl.CreateSchema(ctx, client)
	_ = ctrl.GetSchema(ctx, client)

	logger.Info("starting import")

	err = ctrl.ImportData(ctx, client)
	if err != nil {
		logger.Error("error importing data", zap.Error(err))
		return
	}

	logger.Info("import complete!")

}

func ValidateEnv() error {
	envVars := []string{
		"WEAVIATE_HTTP_ENDPONT",
		"WEAVIATE_GRPC_ENDPONT",
		"WEAVIATE_API_KEY",
	}

	for _, envVar := range envVars {
		_, err := utils.CheckEnv(envVar)
		if err != nil {
			return err
		}

	}

	return nil

}
