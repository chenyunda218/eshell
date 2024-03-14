package main

import (
	"bytes"
	"context"
	"encoding/json"
	"eshell/controllers"
	"eshell/kibana"
	"fmt"
	"io"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
)

var es *elasticsearch.Client

func main() {
	controllers.Init()
	// es, _ = elasticsearch.NewClient(elasticsearch.Config{
	// 	// Username:               "elastic",
	// 	// Addresses:              []string{"https://172.16.101.56:9200"},
	// 	// Password:               "bU3Jr7As+=V1COcHVvWL",
	// 	// CertificateFingerprint: "5756c6555084fd74fa34e69557a85bc52a85599d9bb477e6b6c624ca0cf04c8a",
	// 	CertificateFingerprint: certificateFingerprint,
	// 	Username:               username,
	// 	Password:               password,
	// 	Addresses:              addresses,
	// })
	// r := gin.Default()
	// r.Use(server.CorsMiddleware())
	// r.GET("/data_views", func(ctx *gin.Context) {
	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"data_views": DataViews(),
	// 	})
	// })
	// r.GET("/indexes", func(ctx *gin.Context) {
	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"indexes": ListIndexes(es),
	// 	})
	// })
	// r.POST("/sessions", func(ctx *gin.Context) {
	// 	var req CreateSessionRequest
	// 	ctx.ShouldBindJSON(&req)
	// 	fmt.Println(req)
	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"message": "success",
	// 	})
	// })
	// r.POST("/search", func(ctx *gin.Context) {
	// 	index := ctx.Query("index")
	// 	Query(ctx, index)
	// })
	// r.Run()
}

type Doc struct {
	Name      string `json:"name"`
	Timestamp int64  `json:"timestamp"`
	Person    Person `json:"person"`
}

type Person struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Book   Book   `json:"book"`
	Gender string `json:"gender"`
}

type Book struct {
	Author string `json:"author"`
	Title  string `json:"title"`
}

func Indexing() {
	document := Doc{
		Name:      "my_document",
		Timestamp: 1234567890,
		Person: Person{
			Name: "my_name",
			Age:  123,
			Book: Book{Author: "my_author", Title: "my_title"},
		},
	}
	data, _ := json.Marshal(document)
	res, _ := es.Index("my_index", bytes.NewReader(data))
	fmt.Println(res)
}

type Index struct {
	Health       string `json:"health"`
	Status       string `json:"status"`
	Index        string `json:"index"`
	UUID         string `json:"uuid"`
	Pri          string `json:"pri"`
	Rep          string `json:"rep"`
	DocsCount    string `json:"docs.count"`
	DocsDeleted  string `json:"docs.deleted"`
	StoreSize    string `json:"store.size"`
	PriStoreSize string `json:"pri.store.size"`
}

type CreateSessionRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Query(ctx *gin.Context, index string) {
	body, _ := io.ReadAll(ctx.Request.Body)
	reader := bytes.NewReader(body)
	fmt.Println("query:", string(body))
	res, err := es.Search(
		es.Search.WithIndex(index),
		es.Search.WithBody(reader),
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	body, _ = io.ReadAll(res.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	jmap := make(map[string]interface{})
	json.Unmarshal(body, &jmap)
	ctx.JSON(http.StatusOK, jmap)
}

func ListIndexes(es *elasticsearch.Client) []Index {
	var indexes []Index
	res, err := esapi.CatIndicesRequest{Format: "json"}.Do(context.Background(), es)
	if err != nil {
		return indexes
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &indexes)
	return indexes
}

func DataViews() []kibana.DataView {
	client := kibana.NewClient("username", "password", "kibanaUrl")
	client.Login()
	return client.DataViews()
}
