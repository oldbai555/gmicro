/**
 * @Author: zjj
 * @Date: 2024/12/12
 * @Desc:
**/

package main

import "os"

var (
	pbDir      = os.Getenv("PB_DIR")
	outputDir  = os.Getenv("OUTPUT_DIR")
	goOut      = os.Getenv("GO_OUT")
	gitPath    = os.Getenv("GIT_PATH")
	pbNameList []string
)

func init() {
	CmdByGenClient.Flags().StringSliceVarP(&pbNameList, "pb_name_list", "p", []string{}, "输入proto名称,多个用,隔开")
	CmdByGenServer.Flags().StringSliceVarP(&pbNameList, "pb_name_list", "p", []string{}, "输入proto名称,多个用,隔开")
}
