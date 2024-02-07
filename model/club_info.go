package model

import (
	"encoding/json"
	"time"
)

type EsClubInfo struct {
	ClubID      string        `json:"club_id"`
	Created_at  time.Time     `json:"created_at"`
	Updated_at  time.Time     `json:"updated_at"`
	Deleted_at  time.Time     `json:"deleted_at"`
	ClubName    string        `json:"club_name"`
	CreatedBy   string        `json:"created_by"`
	ClubType    string        `json:"club_type"` // 俱乐部类型
	MatchStatus MatchStatus   `json:"match_status"`
	History     []MatchStatus `json:"history"`
}

type MatchStatus struct {
	Rank   int             `json:"rank"`
	Point  int             `json:"point"`
	Year   time.Time       `json:"year"`
	Raw    json.RawMessage `json:"raw"`
	Remark string          `json:"remark"`
}

const (
	ClubIndexName = "idx-itsm-club-info"
	ClubMapping   = `
	{
		"mappings": {
		  "properties": {
			"club_id": {
			  "type": "keyword"
			},
			"created_at": {
			  "type": "date"
			},
			"updated_at": {
			  "type": "date"
			},
			"deleted_at": {
			  "type": "date"
			},
			"club_name": {
			  "type": "text",
			  "analyzer": "ik_max_word",
			  "search_analyzer": "ik_smart"
			},
			"created_by": {
			  "type": "keyword"
			},
			"club_type": {
			  "type": "keyword"
			},
			"match_status": {
			  "properties": {
				"rank": {
				  "type": "integer"
				},
				"point": {
				  "type": "integer"
				},
				"year": {
				  "type": "date"
				},
				"raw": {
				  "type": "object"
				},
				"remark": {
				  "type": "text",
				  "analyzer": "ik_max_word",
				  "search_analyzer": "ik_smart"
				}
			  }
			},
			"history": {
			  "type": "nested",
			  "properties": {
				"rank": {
				  "type": "integer"
				},
				"point": {
				  "type": "integer"
				},
				"year": {
				  "type": "date"
				},
				"raw": {
				  "type": "object"
				},
				"remark": {
				  "type": "text",
				  "analyzer": "ik_max_word",
				  "search_analyzer": "ik_smart"
				}
			  }
			}
		  }
		}
	}
	  
	`
)
