package trace

import (
	"sync"
)

func init() {
	Ctx = &Context{}
}

var Ctx *Context

type Context struct {
	gidToTraceIdMap sync.Map
}

func (c *Context) RemoveGTrace(gid int64) {
	c.gidToTraceIdMap.Delete(gid)
}

func (c *Context) GetCurGTrace(gid int64) (string, bool) {
	traceId, ok := c.gidToTraceIdMap.Load(gid)
	if !ok {
		return "", false
	}
	return traceId.(string), ok
}

func (c *Context) SetCurGTrace(gid int64, traceId string) {
	c.gidToTraceIdMap.Store(gid, traceId)
}

func GenTraceId() string {
	return NewTraceId()
}
