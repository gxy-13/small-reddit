package snowflake

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	// 设置起始时间
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		fmt.Printf("snowflake time.parse failed, err:%v\n", err)
		return
	}
	snowflake.Epoch = st.UnixNano() / 1e6
	node, err = snowflake.NewNode(machineID)
	if err != nil {
		fmt.Printf("snowflake New Node failed,err:%v\n", err)
		return
	}
	return
}

// GenID 生成雪花ID
func GenID() int64 {
	return node.Generate().Int64()
}
