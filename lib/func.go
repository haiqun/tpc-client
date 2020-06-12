package lib

import (
	"path/filepath"
	"net/url"
	"net/http"
	"fmt"
	"io/ioutil"
	"time"
	"math/rand"
	"crypto/md5"
	"sort"
	"strings"
	"reflect"
	"strconv"
	"errors"
	"tcp_client/providers"
)

var areaParamError = errors.New("区域参数不正确")

/**
 * Db操作错误
 */
func DbError() error {
	return fmt.Errorf("%s","Mysql数据操作失败")
}

// http post方法
func HttpPostValues(uri string, values url.Values) ([]byte, error) {
	response, err := http.PostForm(uri, values)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK{
		return nil, fmt.Errorf("http post error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}

	return ioutil.ReadAll(response.Body)
}

// 生成签名
func genSign(signData map[string]interface{}, signKey string, keyVal string) (sign string) {
	signData[keyVal] = signKey

	//排序
	var keys []string
	for k := range signData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	//拼凑
	var tmpStr string
	for _, k := range keys {
		value := fmt.Sprintf("%v", signData[k])
		tmpStr += k + "=" + value + "&"
	}

	tmpStr = strings.Trim(tmpStr, "&")

	sign = Md5(tmpStr)

	return
}

// md5加密
func Md5(str string) string {
	data := []byte(str)
	return fmt.Sprintf("%x", md5.Sum(data))
}

// 随机字符
func RandString(n int) string {
	rand.Seed(time.Now().UnixNano())

	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// 获取时间戳
func GetNow() int64 {
	now := time.Now()
	return now.Unix()
}

// 产品中心 post签名请求
func ProductCentrePost(uri string, values url.Values) ([]byte, error) {
	// 产品中心配置
	productCentre := providers.Config.GetStringMapString("productCentre")
	uri = productCentre["domain"] + uri

	// 生成签名
	signData := make(map[string]interface{})
	for k, v := range values {
		// 服务器之间调用, 默认只允许一个key对应一个value传参
		if len(v) > 0 {
			signData[k] = v[0]
		}
	}

	values["sign"] = []string{ genSign(signData, productCentre["sign"], "sscfsalt", ) }

	return HttpPostValues(uri, values)
}

/**
 * 获取项目中的路径
 */
func GetProjectPath(pathStr string) string {
	return filepath.Join(providers.RootPath, pathStr)
}

/**
 * struct转换post参数
 */
func StructToUrlMap(i interface{}) (values url.Values) {
	values = url.Values{}
	iVal := reflect.ValueOf(i).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		// You ca use tags here...
		// tag := typ.Field(i).Tag.Get("json")
		// Convert each type into a string for the url.Values string map
		var v string
		switch f.Interface().(type) {
		case int, int8, int16, int32, int64:
			v = strconv.FormatInt(f.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			v = strconv.FormatUint(f.Uint(), 10)
		case float32:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 32)
		case float64:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 64)
		case []byte:
			v = string(f.Bytes())
		case string:
			v = f.String()
		}
		// values.Set(typ.Field(i).Name, v)
		jsonTag := typ.Field(i).Tag.Get("json")
		if jsonTag != "" {
			values.Set(jsonTag, v)
		}
	}
	return
}

/**
 * 获取[0, max] 的随机数
 */
func GetRandNum(max int) int {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rnd.Intn(max)
}

// map反转
func ReverseMap(m map[string]string) map[string]string {
	n := make(map[string]string)
	for k, v := range m {
		n[v] = k
	}
	return n
}

// 呼叫系统area参数 转换到 产品中心area
func AreaConvert(hjArea string) (area string, err error) {
	areaConfig := providers.Config.GetStringMapString("area")
	areaConfigRe := ReverseMap(areaConfig)

	// GZ/CS
	area = strings.ToUpper(areaConfigRe[hjArea])
	if area == "" {
		err = areaParamError
	}
	return
}

/**
 * 产品中心area 转换到 呼叫系统/客服系统area参数
 * defaultVal 是否返回默认值
 */
func AreaToCommArea(area string, defaultVal bool) (areaCall string)  {
	area = strings.ToLower(area)

	areaCall = providers.Config.GetString("area." + area)

	// 返回默认值
	if areaCall == "" && defaultVal == true {
		return "hj-gz"
	}

	return
}

/**
 * vip操盘 area转用户类型值
 */
func VipOperateUserType(area string) uint {
	area = strings.ToLower(area)

	userType := providers.Config.GetInt("vipUserType." + area)

	return uint(userType)
}

/**
 * 产品中心 允许更改用户的区域检测
 */
func ProCentreAllowChangeArea(area string, userArea string) bool {
	area = strings.ToLower(area)
	userArea = strings.ToLower(userArea)

	// 广州用户改为 广州
	if area == "gz" && userArea == "gz" {
		return true
	} else if area == "cs" && (userArea == "cs" || userArea == "hy") {
		// 长沙/华远用户改为 长沙, 他们共用一个公众号
		return true
	} else if area == "hy" && (userArea == "cs" || userArea == "hy") {
		// 长沙/华远用户改为 华远, 他们共用一个公众号
		return  true
	}

	return false
}

/**
 * 订单区域转用户区域
 */
func OrderAreaToUserArea(orderArea string, serveType string) string {
	orderArea = strings.ToUpper(orderArea)

	// 线上产品
	if serveType == "app" || serveType == "998price" {
		switch orderArea {
		case "GZ":
			return "APPGZ"
		case "BJ":
			return "APPGZ"
		case "CS":
			return "APPCS"
		case "HY":
			return "APPHY"
		}
	} else {
		switch orderArea {
		case "GZ":
			return "GZ"
		case "BJ":
			return "GZ"
		case "CS":
			return "CS"
		case "HY":
			return "HY"
		}
	}

	return ""
}

/**
 * 根据客服区域，转换对应的客服系统
 */
func AdminServiceSystemType(area string) string {
	area = strings.ToUpper(area)

	switch area {
	case "GZ":
		// 旧客服系统
		return "old"
	case "CS":
		return "old"
	case "BJ":
		// 新客服系统
		return "new"
	case "HY":
		return "new"
	}

	return ""
}

/**
 * 产品区域转crm对应的区域
 */
func ProductAreaToCrmArea(area string) string {
	area = strings.ToUpper(area)

	if area == "GZ" || area == "APPGZ" {
		return "GZ"
	} else if area == "CS" || area == "APPCS" {
		return "CS"
	} else if area == "HY" || area == "APPHY" {
		return "HY"
	}

	return area
}

/**
 * 判断是否线上产品
 */
func IsOnlineProduct(area string) bool {
	area = strings.ToUpper(area)

	// 线上产品
	if area == "APPGZ" || area == "APPCS" || area == "APPHY" {
		return true
	}

	return false
}

/**
 * crm区域对应线上用户区域
 */
func CrmAreaToOnlineArea(area string) string {
	area = strings.ToUpper(area)

	// 线上产品
	switch area {
	case "GZ":
		return "APPGZ"
	case "CS":
		return "APPCS"
	case "HY":
		return "APPHY"
	}

	return ""
}

/**
 * 隐藏字符串
 */
func HideString(string string, start int, length int, char string) string {
	strLength := len(string)

	if start == 0 && length == 0 {
		length = strLength
	} else if start == 0 {
		start = strLength - length
	} else if length == 0 {
		length = strLength - start
	}

	if start+length > strLength {
		return string
	}

	return string[0:start] + strings.Repeat(char, length) + string[start+length:strLength]
}
