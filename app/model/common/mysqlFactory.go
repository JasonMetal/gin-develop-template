package common

import (
	"context"
)

// ConnectionObject 获取 业务 库客户端连接对象
func ConnectionObject(ctx context.Context) *MysqlClient {
	return NewMysqlClient(ctx, "film")
}
