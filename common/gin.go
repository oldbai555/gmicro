/**
 * @Author: zjj
 * @Date: 2024/6/3
 * @Desc:
**/

package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func NewLogFormatter(svr string) func(param gin.LogFormatterParams) string {
	return func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}
		hint := param.Keys["hint"]
		if hint == nil {
			hint = "none"
		}
		v := fmt.Sprintf("[%s] [GIN] <%v> %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
			svr,
			hint,
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			param.ErrorMessage,
		)
		return v
	}
}
