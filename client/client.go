package client

import (
	"bytes"
	"context"
	"encoding/json"
	"eshell/kibana"
	"io"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
)

func New(config elasticsearch.Config, username, password, kibanaUrl string) (*Client, error) {
	es, err := elasticsearch.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Client{Client: es,
		username:  username,
		password:  password,
		kibanaUrl: kibanaUrl,
	}, nil
}

type Client struct {
	*elasticsearch.Client
	username  string
	password  string
	kibanaUrl string
}

func (c *Client) Query(ctx *gin.Context, index string) {
	body, _ := io.ReadAll(ctx.Request.Body)
	reader := bytes.NewReader(body)
	res, err := c.Search(
		c.Search.WithIndex(index),
		c.Search.WithBody(reader),
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	body, err = io.ReadAll(res.Body)
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

func (c *Client) ListIndexes() []Index {
	var indexes []Index
	res, err := esapi.CatIndicesRequest{Format: "json"}.Do(context.Background(), c)
	if err != nil {
		return indexes
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &indexes)
	return indexes
}

func (c *Client) DataViews() []kibana.DataView {
	client := kibana.NewClient(c.kibanaUrl, c.username, c.password)
	client.Login()
	return client.DataViews()
}
