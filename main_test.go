package main

import (
	"encoding/json"
	"esplay/model"
	"esplay/utils"
	"log"
	"testing"
	"time"

	"github.com/mottaquikarim/esquerydsl"
)

// 测试 ES 连接
func Test_EsConnection(t *testing.T) {
	esClient, err := utils.NewElasticsearchClient()
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}
	log.Printf("Elasticsearch client created successfully")

	info, err := esClient.Ping()
	if err != nil {
		log.Fatalf("Elasticsearch连接失败： %v", err)
	}
	log.Printf("Elasticsearch连接成功： %s\n", info)
}

// 插入单条数据
func Test_InsertIndexDocument(t *testing.T) {
	esClient, err := utils.NewElasticsearchClient()
	if  err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	data := model.EsClubInfo{
		ClubID: "3",
		ClubName: "测试俱乐部3",
		Created_at: time.Now(),
		CreatedBy: "xiaoming3",
		ClubType: "football3",
		MatchStatus: model.MatchStatus{
			Rank:   1,
			Point:  100,
			Year:   time.Now(),
			Raw:    json.RawMessage(`{"key1": "value1","key2": "value2"}`),
			Remark: "Great performance",
		},
		History: []model.MatchStatus{
			{
				Rank:   2,
				Point:  90,
				Year:   time.Now().AddDate(0, -1, 0),
				Raw:    json.RawMessage(`{"key1": "hello world","key3": "value3"}`),
				Remark: "Previous match",
			},
		},
	}
	err = utils.InsertIndexDocument(esClient, model.ClubIndexName, data)
	if err != nil {
		log.Fatalf("Error creating index doc: %s", err)
	}
	log.Printf("Index doc created successfully")
}

// 批量插入数据
func Test_BulkInsert(t *testing.T) {
	esClient, err := utils.NewElasticsearchClient()
	if  err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	var docs []model.EsClubInfo
	names := []string{"Alice", "John", "Mary"}
	for i, name := range names {
		docs = append(docs, model.EsClubInfo{
			ClubID: string(rune(i)),
			ClubName: name,
			Created_at: time.Now(),
			CreatedBy: "xiaoming3",
			ClubType: "football3",
			MatchStatus: model.MatchStatus{
				Rank:   1,
				Point:  100,
				Year:   time.Now(),
				Raw:    json.RawMessage(`{"key1": "value1","key2": "value2"}`),
				Remark: "Great performance",
			},
			History: []model.MatchStatus{
				{
					Rank:   2,
					Point:  90,
					Year:   time.Now().AddDate(0, -1, 0),
					Raw:    json.RawMessage(`{"key1": "hello world","key3": "value3"}`),
					Remark: "Previous match",
				},
			},
		})
	}
	err = utils.BulkInsertIndexDocument(esClient, model.ClubIndexName, docs)
	if err != nil {
		log.Fatalf("Error bulk inserting index doc: %s", err)
	}
	log.Printf("Bulk insert completed")
}

// 根据 ID 查询文档
func Test_GetDocumentByID(t *testing.T) {
	esClient, err := utils.NewElasticsearchClient()
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}
	log.Printf("Elasticsearch client created successfully")

	// 根据 id 查询
	doc, err := utils.GetDocumentByID(esClient, model.ClubIndexName, "3")
	log.Printf("Document retrieved successfully: %v", doc)
}

// 查询所有
func Test_SearchAllDocuments(t *testing.T) {
	esClient, err := utils.NewElasticsearchClient()
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}
	log.Printf("Elasticsearch client created successfully")

	docs, err := utils.SearchAllIndexDocument(esClient, model.ClubIndexName)
	for _, hit := range docs {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}
	log.Printf("Search completed")
}

// 根据 DSL 查询
func Test_SearchDocumentsWithDSL(t *testing.T) {
	esClient, err := utils.NewElasticsearchClient()
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}
	log.Printf("Elasticsearch client created successfully")
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"created_at": map[string]interface{}{
								"gte": time.Now().Add(-6 * time.Hour), 
								"lte": time.Now(),        
							},
						},
					},
					{
						"match": map[string]interface{}{
							"club_name": "Mary", 
						},
					},
				},
			},
		},
	}

	body, _ := json.Marshal(query)
	docs, err := utils.SearchDocumentsByDSL(esClient, model.ClubIndexName, body)
	if err != nil {
		log.Fatalf("Error searching documents: %s", err)
	}
	for _, hit := range docs {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}
	log.Printf("Search completed")
}

// 删除数据
func Test_DeleteDocument(t *testing.T) {
	esClient, err := utils.NewElasticsearchClient()
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}
	log.Printf("Elasticsearch client created successfully")
	err = utils.DeleteDocument(esClient, model.ClubIndexName, "3")
	if err != nil {
		log.Fatalf("Error deleting index doc: %s", err)
	}
	log.Println("Index doc deleted successfully")
}

// 测试 ESQueryDSL 包
func Test_ESQueryDSL(t *testing.T) {
	esClient, err := utils.NewElasticsearchClient()
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}
	log.Printf("Elasticsearch client created successfully")
	body, _ := json.Marshal(esquerydsl.QueryDoc{
		Index: model.ClubIndexName,
		Sort:  []map[string]string{
			{"club_id.keyword": "asc"},
		},
		And: []esquerydsl.QueryItem{
			{
				Field: "club_name",
				Value: "Mary",
				Type:  esquerydsl.Match,
			},
		},
	})
	docs, err := utils.SearchDocumentsByDSL(esClient, model.ClubIndexName, body)
	if err != nil {
		log.Fatalf("Error searching documents: %s", err)
	}
	for _, hit := range docs {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}
	log.Printf("Search completed")
}