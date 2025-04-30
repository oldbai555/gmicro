package trace

import (
	"gmicro/common"
)

func NewTraceId() string {
	return common.GenUUID()
}
