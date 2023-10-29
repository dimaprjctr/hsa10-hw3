package handlers

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Handlers struct {
	mongoClient *mongo.Client
	esClient    *elasticsearch.Client
}

func NewHandlers(mongoClient *mongo.Client, esClient *elasticsearch.Client) Handlers {
	return Handlers{
		mongoClient,
		esClient,
	}
}

func (h *Handlers) InsertTestData(w http.ResponseWriter, r *http.Request) {
	if err := h.insertToMongo(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status": false}`))
		return
	}

	if err := h.insertToEs(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status": false}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": true}`))
}

func (h *Handlers) insertToMongo(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	collection := h.mongoClient.Database("testing").Collection("numbers")
	_, err := collection.InsertOne(ctx, bson.D{{"value", rand.Intn(10000)}})

	return err
}

func (h *Handlers) insertToEs(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	docJSON, err := json.Marshal(map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
	})
	if err != nil {
		log.Infof("Error serializing document to JSON: %s\n", err)
		return err
	}
	req := esapi.IndexRequest{
		Index:      "test-index",
		DocumentID: "test-document-id",
		Body:       strings.NewReader(string(docJSON)),
	}
	esRes, err := req.Do(ctx, h.esClient)
	if err != nil {
		log.Infof("Error sending request: %s", err)
		return err
	}
	defer esRes.Body.Close()

	if esRes.IsError() {
		log.Infof("Getting err on response: %s", err)
		return err
	}
	return nil
}
