package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tonytw1/gauges/model"
	"io"
	"log"
	"strconv"
	"strings"
)

func PersistRoutes(s3Client *s3.Client, bucket string, key string, asJson []byte) (*s3.PutObjectOutput, error) {
	log.Print("Persisting routes to bucket: '" + bucket + "'" + " key: '" + key + "'")
	body := string(asJson)
	log.Print("Persisting: " + body)
	output, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(body),
	})
	if err != nil {
		log.Print("Failed to persist routes: " + err.Error())
		return output, err
	}
	log.Print("Persisted routes")
	return output, nil
}

func LoadPersistedRoutes(bucket string, key string, s3Client *s3.Client) []model.Route {
	persistedRoutes := make([]model.Route, 0)
	log.Print("Loading persisted routes from bucket: '" + bucket + "'" + " key: '" + key + "'")
	persistedRoutesObject, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
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
				fmt.Sprint(err.Error())
			}
		} else {
			log.Print("Persisted routes object body is nil")
		}
	}
	return persistedRoutes
}
