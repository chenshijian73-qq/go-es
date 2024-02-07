package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"esplay/model"
	"log"
	"strings"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

func NewElasticsearchClient() (*elasticsearch.Client, error) {
	config := GetConfig()
	log.Printf("Elasticsearch config: %+v", config.Es)
	esClient, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses: []string{config.Es.Address},
			Username:  config.Es.Username,
			Password:  config.Es.Password,
		},
	)
	return esClient, err
}

func CreateIndex(es *elasticsearch.Client, indexName, mapping string) error {
	// 检查索引是否已存在
	exists, err := IndexExists(es, indexName)
	if err != nil {
		log.Printf("Error checking if index exists: %s", err)
		return err
	}
	if exists {
		log.Printf("Index '%s' already exists, skipping creation.\n", indexName)
		return errors.New("Index already exists")
	}
	// 创建索引
	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  bytes.NewReader([]byte(mapping)),
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Printf("Error creating index: %s", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error creating index: %s", res.Status())
		return errors.New("Index creation failed")
	}
	log.Printf("Index created: %s", indexName)
	return nil
}

func IndexExists(es *elasticsearch.Client, indexName string) (bool, error) {
	req := esapi.IndicesExistsRequest{Index: []string{indexName}}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	return res.StatusCode == 200, nil
}

// 插入文档
func InsertIndexDocument(es *elasticsearch.Client, indexName string, data model.EsClubInfo) error {
	doc, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling document: %s", err)
		return err
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: data.ClubID,
		Body:       bytes.NewReader(doc),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Printf("Error indexing document: %s", err)
	}
	defer res.Body.Close()
	return nil
}

func BulkInsertIndexDocument[T any](esClient *elasticsearch.Client, indexName string, docs []T) error {
	bulkReq, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index: indexName,
		Client: esClient,
	})
	if err != nil {
		log.Printf("creat bulkReq err: %s", bulkReq)
		return err
	}
	defer bulkReq.Close(context.Background())

	for _, doc := range docs {
		// 将结构体转换为 JSON 字节切片
		docBytes, err := json.Marshal(doc)
		if err != nil {
			log.Fatalf("Error encoding document as JSON: %s", err)
		}
		err = bulkReq.Add(context.Background(), esutil.BulkIndexerItem{
			Action: "index",
			Body:   bytes.NewReader(docBytes),
		})
		if err != nil {
			log.Printf("bulkReq.Add error: %s", err)
			return err
		}
	}
	return nil
}

// 通过 ID 查询
func GetDocumentByID(client *elasticsearch.Client, indexName string, id string) (map[string]interface{}, error) {
	resp, err := client.Get(indexName, id)
	if err != nil {
		log.Printf("get document by id failed, err:%v\n", err)
		return nil, err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Printf("decode response body failed, err:%v\n", err)
		return nil, err
	}

	sourceData, ok := result["_source"].(map[string]interface{})
	if !ok {
		log.Println("_source field not found in the response")
		return nil, errors.New("_source field not found in the response")
	}
	return sourceData, nil
}

// 查询所有
func SearchAllIndexDocument(client *elasticsearch.Client, indexName string) (docs []interface{}, err error) {
	// 搜索文档
	query := `{ "query": { "match_all": {} } }`
	res, err := client.Search(
		client.Search.WithIndex(indexName),
		client.Search.WithBody(strings.NewReader(query)),
	)
	defer res.Body.Close()
	if err != nil || res.IsError() {
		log.Printf("search document failed, err:%v\n", err)
		return
	}
	var r map[string]interface{}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %s\n", err)
		return nil, err
	}
	docs = r["hits"].(map[string]interface{})["hits"].([]interface{})
	return docs, nil
}

// 通过 DSL 查询
func SearchDocumentsByDSL(es *elasticsearch.Client, indexName string, queryDslBody []byte) (docs []interface{}, err error) {
	var buf bytes.Buffer
	buf.Write(queryDslBody)

	req := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  &buf,
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Printf("Error searching documents: %s", err)
		return nil, err
	}
	defer res.Body.Close()

	var r map[string]interface{}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %s\n", err)
		return nil, err
	}
	if r["hits"] == nil {
		return nil, errors.New("no 'hits' field in the response")
	}
	docs = r["hits"].(map[string]interface{})["hits"].([]interface{})
	return docs, nil
}

// 删除数据
func DeleteDocument(es *elasticsearch.Client, indexName string, docID string) error {
	req := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: docID,
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error deleting document: %s\n", err)
		return err
	}
	return nil
}
