/**
 * @Author: zjj
 * @Date: 2025/1/9
 * @Desc:
**/

package mysql

import (
	utils "gmicro/common"
	"gmicro/pkg/gerr"
	"gorm.io/gorm"
)

type trInfo struct {
	txDb      *gorm.DB
	createdAt uint32
	updatedAt uint32
}

func (g *GormEngine) Begin() (string, error) {
	g.trLock.Lock()
	defer g.trLock.Unlock()
	begin := g.db.Begin()
	err := begin.Error
	if err != nil {
		return "", err
	}
	txId := utils.GenUUID()
	g.trMap[txId] = &trInfo{
		txDb:      begin,
		createdAt: utils.TimeNow(),
		updatedAt: utils.TimeNow(),
	}
	return txId, nil
}

func (g *GormEngine) Rollback(trId string) error {
	g.trLock.Lock()
	defer g.trLock.Unlock()
	info := g.trMap[trId]
	if info == nil {
		return gerr.NewInvalidArg("%s not found", trId)
	}
	delete(g.trMap, trId)
	return info.txDb.Rollback().Error
}

func (g *GormEngine) Commit(trId string) error {
	g.trLock.Lock()
	defer g.trLock.Unlock()
	info := g.trMap[trId]
	if info == nil {
		return gerr.NewInvalidArg("%s not found", trId)
	}
	delete(g.trMap, trId)
	return info.txDb.Commit().Error
}

// 事务超时...
