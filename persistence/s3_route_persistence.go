package persistence

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tonytw1/gauges/model"
	"io"
	"log"
	"strconv"
	"strings"
)

type RoutePersistence interface {
	PersistRoutes(asJson []byte) (*s3.PutObjectOutput, error)
	LoadPersistedRoutes() []model.Route
}

type S3RoutePersistence struct {
	s3Client *s3.Client
	bucket   string
	key      string
}

func NewS3RoutePersistence(bucket string, key string) RoutePersistence {
	// Setup S3 client
	region := "eu-west-2"
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}
	return S3RoutePersistence{s3Client: s3.NewFromConfig(cfg), bucket: bucket, key: key}
}

func (svc S3RoutePersistence) PersistRoutes(asJson []byte) (*s3.PutObjectOutput, error) {
	log.Print("Persisting routes to bucket: '" + svc.bucket + "'" + " key: '" + svc.key + "'")
	body := string(asJson)
	log.Print("Persisting: " + body)
	output, err := svc.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(svc.bucket),
		Key:    aws.String(svc.key),
		Body:   strings.NewReader(body),
	})
	if err != nil {
		log.Print("Failed to persist routes: " + err.Error())
		return output, err
	}
	log.Print("Persisted routes")
	return output, nil
}

func (svc S3RoutePersistence) LoadPersistedRoutes() []model.Route {
	persistedRoutes := make([]model.Route, 0)
	log.Print("Loading persisted routes from bucket: '" + svc.bucket + "'" + " key: '" + svc.key + "'")
	persistedRoutesObject, err := svc.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(svc.bucket),
		Key:    aws.String(svc.key),
	})
	if err != nil {
		log.Print("Could not load persisted routes: " + err.Error())
	} else {
		if persistedRoutesObject.Body != nil {
			buf := new(strings.Builder)
			_, err := io.Copy(buf, persistedRoutesObject.Body)
			if err == nil {
				persistedJson := buf.String()
				log.Print("Persisted JSON: " + persistedJson)

				err = json.NewDecoder(strings.NewReader(persistedJson)).Decode(&persistedRoutes)
				if err != nil {
					log.Print("Could not decode persisted routes: " + err.Error())
				} else {
					log.Print("Loaded " + strconv.Itoa(len(persistedRoutes)) + " persisted routes")
				}
			} else {
				log.Print("Could not read persisted routes")
			}
		} else {
			log.Print("Persisted routes object body is nil")
		}
	}
	return persistedRoutes
}
