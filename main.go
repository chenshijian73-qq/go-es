package main

import (
	"encoding/json"
	"esplay/model"
	"esplay/utils"
	"fmt"
	"log"
	"time"
)

func main() {
	fmt.Println("Elasticsearch client test begin")
	// 创建 Elasticsearch 客户端
	esClient, err := utils.NewElasticsearchClient()
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}
	fmt.Println("Elasticsearch client created successfully")

	info, err := esClient.Ping()
	if err != nil {
		log.Fatalf("Elasticsearch连接失败： %v", err)
	}
	fmt.Printf("Elasticsearch连接成功： %s\n", info)



	// 删除数据
	err = utils.DeleteDocument(esClient, model.ClubIndexName, "3")
	if err != nil {
		log.Fatalf("Error deleting index doc: %s", err)
	}
	log.Println("Index doc deleted successfully")

	// 批量插入数据
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
}
