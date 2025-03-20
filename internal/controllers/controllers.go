package controllers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"os"
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	internalLogger "github.com/iamcaleberic/randovec/internal/logger"
	internalModels "github.com/iamcaleberic/randovec/internal/models"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/grpc"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
	"go.uber.org/zap"
)

var logger = internalLogger.InitLogger()

// ImportData generates objects and vectors and writes them in batches into weaviate
func ImportData(ctx context.Context, client *weaviate.Client) error {
	var insertObjects []*models.Object

	numOfGeneratedObjects, _ := strconv.Atoi(os.Getenv("NUM_OBJECTS"))

	genratedObjs := GenerateData(numOfGeneratedObjects)

	for _, gobj := range genratedObjs {
		uuid := uuid.NewString()

		insertObj := &models.Object{
			Class: "RandClass",
			ID:    strfmt.UUID(uuid),
			Properties: map[string]interface{}{
				"content": gobj.Content,
			},
			Vector: gobj.Vector,
		}

		insertObjects = append(insertObjects, insertObj)
	}

	batchSize, _ := strconv.Atoi(os.Getenv("BATCH_SIZE"))

	objectsChunks := GetObjectsChunks(insertObjects, batchSize)

	for i, objectChunk := range objectsChunks {
		logger.Info("adding chunk of to db", zap.Int("chunk-size", batchSize))

		_, err := client.Batch().ObjectsBatcher().WithObjects(objectChunk...).Do(ctx)
		if err != nil {
			logger.Error("error imporing random vector data to weaviate", zap.Error(err))
			continue
			// return err
		}

		logger.Info("added objects to db", zap.Int("number-of-objects", i*batchSize))
	}
	return nil
}

// RandString generates random string
func RandString(size int) string {
	b := make([]byte, size)
	rand.Read(b)
	s := hex.EncodeToString(b)

	return s[:size]
}

// GenerateVectorData generates vectors
func GenerateVectorData(vectorSize int) []float32 {
	vector := make([]float32, vectorSize)
	for i := 0; i < len(vector); i++ {
		vector[i] = 0.12345
	}

	return vector
}

// GetObjectsChunks builds chunk/batches of obects passed
func GetObjectsChunks(objects []*models.Object, chunkSize int) [][]*models.Object {
	objectChunks := make([][]*models.Object, 0, (len(objects)+chunkSize-1)/chunkSize)

	for chunkSize < len(objects) {
		objects, objectChunks = objects[chunkSize:], append(objectChunks, objects[0:chunkSize:chunkSize])
	}

	objectChunks = append(objectChunks, objects)

	return objectChunks
}

// CreateWeaviateClient creates weaviate client
func CreateWeaviateClient() (*weaviate.Client, error) {
	cfg := weaviate.Config{
		Host:       os.Getenv("WEAVIATE_HTTP_ENDPONT"),
		Scheme:     "https",
		AuthConfig: auth.ApiKey{Value: os.Getenv("WEAVIATE_API_KEY")},
		GrpcConfig: &grpc.Config{
			Host:    os.Getenv("WEAVIATE_GRPC_ENDPONT"),
			Secured: true,
		},
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		logger.Info("failed to create client")
		return nil, err
	}

	return client, nil
}

// GetSchema queries weaviate schema
func GetSchema(ctx context.Context, client *weaviate.Client) error {

	schema, err := client.Schema().Getter().Do(ctx)
	if err != nil {
		logger.Info("failed to get schema")
		return err
	}

	logger.Info("schema", zap.Any("schema", schema))

	return nil

}

// CreateSchema creates schema
func CreateSchema(ctx context.Context, client *weaviate.Client) error {
	classes := []*models.Class{
		{
			Class: "RandClass",
			Properties: []*models.Property{
				{
					DataType:    []string{schema.DataTypeString.String()},
					Description: "Random text",
					Name:        "content",
				},
			},
			ReplicationConfig: &models.ReplicationConfig{
				AsyncEnabled: true,
				Factor:       3,
			},
		},
	}

	for _, class := range classes {
		err := client.Schema().ClassCreator().WithClass(class).Do(context.Background())
		if err != nil {
			logger.Error("error creating class", zap.Error(err))
			// return err
			continue
		}
	}

	return nil

}

// GenerateData generates data samples
func GenerateData(numberOfSamples int) []internalModels.DataObject {
	vectorSize, _ := strconv.Atoi(os.Getenv("VECTOR_SIZE"))
	var dataObjects []internalModels.DataObject
	for range numberOfSamples {
		obj := internalModels.DataObject{
			Content: RandString(10),
			Vector:  GenerateVectorData(vectorSize),
		}
		dataObjects = append(dataObjects, obj)
	}

	logger.Info("partial data:", zap.Any("partial", dataObjects[:5]))

	return dataObjects
}
