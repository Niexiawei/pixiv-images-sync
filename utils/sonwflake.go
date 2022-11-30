package utils

import "github.com/bwmarrin/snowflake"

var node *snowflake.Node

func InitSnowflake() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err.Error())
	}
}

func GetNextId() (nextId int64) {
	id := node.Generate()
	return id.Int64()
}
