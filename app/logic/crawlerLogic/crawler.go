/**
* @Author: qsx
* @Description:
* @File: controller
* @Module: crawler
* @Date: 2022-10-12 16:13:33
 */
package crawlerLogic

import (
	"bufio"
	_ "develop-template/app/entity/resp/typeRespEntity"
	configLib "github.com/olebedev/config"
	"github.com/robfig/cron/v3"
	"log"
	"math/rand"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	//"bytes"
	"context"
	"github.com/JasonMetal/submodule-support-go.git/helper/time"
	"os"
	//"log"
	"develop-template/app/constant"
	myError "develop-template/app/error"
	_ "develop-template/app/error/common"
	//errMsg "develop-template/app/error/crawlerError"
	"fmt"

	"github.com/gin-gonic/gin"
	time2 "time"
)

var crawlerErr myError.Error

type logic struct {
	Ctx        context.Context
	GCtx       *gin.Context
	userId     uint32
	CurlConfig *configLib.Config
}

func NewLogic(ctx context.Context) *logic {
	//func NewLogic(ctx *gin.Context) *logic {
	return &logic{Ctx: ctx}
}

// SetInit 设置初始化数据
func (l *logic) SetInit(userId uint32) {
	l.userId = userId
}

// @Title       ：CreateDataCron
// @Description ：
// @Return      ：res
// @Time        ：2023-01-09 10:11:21
// @Author      ：user
func (l *logic) CreateVodDetailDataCron() (res interface{}) {
	c := cron.New(cron.WithSeconds())
	fmt.Println("==========================Run...", time2.Now().Format(constant.DefaultDateFormat))
	//每小时的每一分每一秒 一次
	rule1 := "*/30 * * * * *"
	c.AddFunc(rule1, func() {
		r := l.GetAllTagsData()
		fmt.Printf("\n Run models.GetAllTagsData..., %#v\n%s\n", r, time2.Now().Format(constant.DefaultDateFormat))
	})

	rule := "*/3 * * * * *"
	c.AddFunc(rule, func() {
		r := l.GetVodDetailV2()
		fmt.Printf("\n Run models.GetVodDetail..., %#v\n%s\n", r, time2.Now().Format(constant.DefaultDateFormat))
	})
	c.Start()
	// 每10秒执行一次
	t1 := time2.NewTimer(time2.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time2.Second * 10)
		}
	}
	return
}

// @Title       ：GetAllTagsData
// @Description ：//https://github.com/gocolly/colly.git
// @Param       ：params
// @Return      ：interface{}
// @Return      ：myError.Error
// @Time        ：2022-10-27 13:56:48
// @Author      ：user
func (l *logic) GetAllTagsData() (res interface{}) {
	r, _ := l.GetTagsData()
	//写入文件
	res, _ = l.writeFile("complete the task GetAllTagsData at "+time.Date(0, constant.DefaultDateFormat), "GetAllTagsData")
	return map[string]interface{}{
		"return":       r,
		"ResWriteFile": res,
	}
}

func (l *logic) GetTagsData() (resPon interface{}, crErr myError.Error) {
	return
}

func (l *logic) GetRandomNum(number int) int {
	rand.Seed(time2.Now().UnixNano())
	return rand.Intn(number)
}

func (l *logic) GetRandomFloatNum() float64 {
	rand.Seed(time2.Now().UnixNano())
	return rand.Float64()
}

// writeFile
// @Description:
// @receiver l
// @param str

// @Title       ：GetVodDetailV2
// @Description ： 随机请求头
// @Return      ：res
// @Time        ：2023-01-09 11:28:29
// @Author      ：user
func (l *logic) GetVodDetailV2() (res interface{}) {
	var insertGetId = 1
	defer func() {
		r := recover()
		fmt.Println("============panic GetVodDetail============", r)
	}()
	//number := 0
	//......
	//执行shell文件
	command := `/home/www/go-develop-template/shell/crawler.sh`
	fmt.Printf("Execute Shell:%s ", command)
	cmd := exec.Command("/bin/bash", "-c", command)
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Execute Shell:%s failed with error:%s", command, err.Error())
		log.Fatal(err)
	}
	fmt.Printf("Execute Shell:%s finished with output:\n%s", command, string(out))

	return map[string]interface{}{
		"return":       insertGetId,
		"ResWriteFile": res,
	}
}

// 最终方案-全兼容
func GetCurrentAbPath() string {
	dir := GetCurrentAbPathByExecutable()
	if strings.Contains(dir, getTmpDir()) {
		return GetCurrentAbPathByCaller()
	}
	return dir
}

// 获取系统临时目录，兼容go run
func getTmpDir() string {
	dir := os.Getenv("TEMP")
	if dir == "" {
		dir = os.Getenv("TMP")
	}
	res, _ := filepath.EvalSymlinks(dir)
	return res
}

// 获取当前执行文件绝对路径
func GetCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func GetCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

//
//  writeFile
//  @Description:
//  @receiver l
//  @param str
//

func (l *logic) writeFile(str string, fileName string) (int, error) {
	outputFile, outputError := os.OpenFile("output"+fileName+".dat", os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		fmt.Printf("An error occurred with file opening or creation\n")
		return 0, nil
	}
	defer outputFile.Close()
	outputWriter := bufio.NewWriter(outputFile)
	outputString := str + "\n"
	i, err := outputWriter.WriteString(outputString)
	if err != nil {
		fmt.Printf("An error \n", err)
	}
	outputWriter.Flush()
	return i, err
}
