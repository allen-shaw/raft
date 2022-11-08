package main

import (
	"encoding/json"
	"github.com/AllenShaw19/raft"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StoreServer struct {
	e       *gin.Engine
	m       *raft.NodeManager
	address string
	t       *Table
}

func NewStoreServer(m *raft.NodeManager, address string, table *Table) *StoreServer {
	gin.SetMode(gin.ReleaseMode)

	s := &StoreServer{
		m:       m,
		address: address,
		t:       table,
	}
	s.e = gin.Default()
	s.init()

	return s
}
func (s *StoreServer) Run() {
	go func() {
		err := s.e.Run(s.address)
		if err != nil {
			panic(err)
		}
	}()
}

func (s *StoreServer) init() {
	s.e.Use(gin.Logger(), gin.Recovery())
	g := s.e.Group("/store")
	g.POST("/get", s.get)
	g.POST("/set", s.set)
	//g.DELETE("/delete", s.delete)
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
type GetRequest struct {
	GroupID string `json:"group_id"`
	Key     string `json:"key"`
}
type GetResponse struct {
	Value string `json:"value"`
}

type SetRequest struct {
	GroupID string `json:"group_id"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}

func (s *StoreServer) get(c *gin.Context) {
	req := &GetRequest{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	value := s.t.Get(req.Key)
	resp := &Response{
		Code: 0,
		Msg:  "ok",
		Data: GetResponse{
			Value: value,
		},
	}
	c.JSON(http.StatusOK, resp)
}

func (s *StoreServer) set(c *gin.Context) {
	req := &SetRequest{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	node := s.m.GetNode(req.GroupID)

	bytes, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	resp := &Response{Code: 0, Msg: "ok"}
	err = node.Apply(bytes)
	if err != nil {
		resp.Code = 10001
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}
