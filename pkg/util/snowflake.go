package util

import "github.com/bwmarrin/snowflake"

func GenSnowflakeIdInt64() (int64, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return 0, err
	}
	id := node.Generate()
	return id.Int64(), nil
}
