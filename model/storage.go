package model

import "github.com/finishy1995/go-library/storage"

type SessionData struct {
	storage.Model
	Id     string `dynamo:",hash"` // session id
	UserId string
}

type NodeData struct {
	storage.Model
	SessionId   string `dynamo:",hash"`  //session id
	NodeId      string `dynamo:",range"` //node id
	Type        uint16
	Status      uint8
	Style       string
	OriginParam map[string]string
	Answer      string
	LastUpdate  int64 // 超过3mim未处理完成的节点认为ai处理超时 返回err
}

type SessionContentUpdateData struct {
	storage.Model
	Id               string   `dynamo:",hash" json:"id"` // session id
	FileID           string   `json:"fileID"`
	Author           string   `json:"author"`
	CreatedTime      int64    `json:"createdTime"`
	LastModifiedTime int64    `json:"lastModifiedTime"`
	ModifiedContent  string   `json:"modifiedContent"`
	Tags             []string `json:"tags"`
	Category         string   `json:"category"`
}
