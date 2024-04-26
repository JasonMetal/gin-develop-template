package common

import (
	"context"
	myError "develop-template/app/error"
	"github.com/JasonMetal/submodule-support-go.git/bootstrap"
	"gorm.io/gorm"
)

// MysqlClient mysql连接对象
type MysqlClient struct {
	DB     *gorm.DB
	Master bool
}

func NewMysqlClient(ctx context.Context, name string) *MysqlClient {
	mc := new(MysqlClient)
	gdb, err := mc.WithDBContext(ctx, name) //gdb.GetMysqlInstance(name)

	CheckMysqlError(err)
	mc.DB = gdb

	return mc
}

func CheckMysqlError(err error) myError.Error {
	if err == nil || err == gorm.ErrRecordNotFound {
		return nil
	}
	e := myError.NewMysqlError()
	e.AppendCause().SetCause(err)
	e.SetMessage(err.Error())

	bootstrap.CheckError(e, "mysql")

	return e
}

func (mc MysqlClient) WithDBContext(ctx context.Context, name string) (*gorm.DB, error) {
	instance, err := bootstrap.GetMysqlInstance(name)
	if err != nil {
		return nil, err
	}

	return instance.DB, err
}
