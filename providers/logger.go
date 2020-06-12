package providers

import (
	"bytes"
	"encoding/json"
	"github.com/cihub/seelog"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

// 数据库日志记录
var Logger seelog.LoggerInterface
var OnceL sync.Once

// 初始化函数
func GetLogger() seelog.LoggerInterface{
	OnceL.Do(func() {
		Logger = loggerInit()
	})
	return Logger
}

func loggerInit() seelog.LoggerInterface {
	// 自定义输出格式
	err := seelog.RegisterCustomFormatter("Sscf", func(param string) seelog.FormatterFunc {
		return func(message string, level seelog.LogLevel, context seelog.LogContextInterface) interface{} {
			//获取匹配数据
			matchRe, _ := regexp.CompilePOSIX(`([0-9]|X|x)+`)
			mobileRe, _ := regexp.Compile(`^(1(3|4|5|6|7|8|9)\d{1})\d{5}(\d{3})$`)

			for _, Old := range matchRe.FindAllString(message, -1) {
				oldLen := len(Old)
				if oldLen != 11 && oldLen != 15 && oldLen != 18 {
					continue
				}

				xIndex := strings.Index(Old, "x") //身份证x的位置
				XIndex := strings.Index(Old, "X") //身份证X的位置
				if xIndex != -1 && XIndex != -1 || (xIndex != -1 || XIndex != -1) && oldLen == 11 || xIndex != -1 && xIndex != oldLen-1 || XIndex != -1 && XIndex != oldLen-1 {
					//不能同时存在xX
					//不能存在x|X时获取的是11位手机号
					//存在x|X时必须是结尾
					continue
				}

				var New string
				if oldLen == 11 {
					//匹配手机号码
					New = mobileRe.ReplaceAllString(Old, `${1}*****${3}`)
				} else {
					//匹配身份证
					New = Old[:6] + "********" + Old[14:]
				}

				message = strings.Replace(message, Old, New, -1)
			}

			prefix := "sscf"
			messageMap := make(map[string]string)
			messageMap[prefix+"_message"] = message
			messageMap[prefix+"_level"] = strconv.Itoa(int(level))
			messageMap[prefix+"_level_name"] = level.String()
			messageMap[prefix+"_datetime"] = time.Now().Format("2006-01-02 15:04:05")

			extra, _ := json.Marshal(map[string]interface{}{
				"trace": "func:" + context.Func() + ",path:" + context.FullPath() + ",line:" + strconv.Itoa(context.Line()),
			})
			messageMap[prefix+"_extra"] = string(extra)

			messageJson, _ := json.Marshal(messageMap)
			return string(messageJson) + "\n"
		}
	})

	if err != nil {
		panic("日志创建自定义格式函数失败:" + err.Error())
	}

	tmpl, err := template.ParseFiles(
		filepath.Join(RootPath, "config", "log.xml"),
	)

	if err != nil {
		panic("日志创建模板失败:" + err.Error())
	}

	var logByte bytes.Buffer
	err = tmpl.Execute(
		&logByte,
		struct {
			RootPath string
		}{
			RootPath: RootPath,
		},
	)

	if err != nil {
		panic("日志模板写入参数失败:" + err.Error())
	}

	// 使用模版替换，匹配对应的路径
	Logger, err = seelog.LoggerFromConfigAsBytes(logByte.Bytes())
	if err != nil {
		panic("日志读取配置文件失败:" + err.Error())
	}

	return Logger
}
