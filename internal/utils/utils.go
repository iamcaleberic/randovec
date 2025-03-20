package utils

import (
	"fmt"
	"os"

	internalLogger "github.com/iamcaleberic/randovec/internal/logger"

	"go.uber.org/zap"
)

var logger = internalLogger.InitLogger()

// uniform error for env releted errors
func CheckEnv(envName string) (string, error) {
	value := os.Getenv(envName)
	if value == "" {
		logger.Error("failed because required an environment variable is not set: ", zap.String("env_var", envName))
		return "", fmt.Errorf("failed because required an environment variable is not set: env_var: %v ", envName)
	}
	return value, nil
}
