package idgenerator

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
)

type Snowflake struct {
	node *snowflake.Node
}

func NewSnowflake(nodeID int64) *Snowflake {
	node, err := snowflake.NewNode(nodeID)

	if err != nil {
		panic(fmt.Sprintf("Failed to create snowflake node: %v", err))
	}

	return &Snowflake{node}
}

func (s *Snowflake) Generate() int64 {
	return s.node.Generate().Int64()
}
