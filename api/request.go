package api

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/PGITAb/bc-operator-api-verify/config"
	"github.com/go-resty/resty/v2"
)

func Get(path string, query map[string]string) Response {
	httpClient := resty.New()
	rsp, err := httpClient.R().EnableTrace().
		SetQueryParams(query).
		Get(config.GetDomain() + path)
	fmt.Println("wrong typo")
	return Response{
		Body:  rsp.Body(),
		Code:  rsp.StatusCode(),
		Error: err,
	}
}

// 初始化部分
var writer *csv.Writer

func init() {
	// 获取当前日期和时间
	currentTime := time.Now()

	logTime := currentTime.Format("2006-01-02-150405")

	// 构建日志文件名
	logFileName := "log/op-api-log-" + logTime + ".csv"
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("无法打开CSV:", err)
	}
	writer = csv.NewWriter(file)
}

// 修改 Post 函数

var (
	requestCount       int
	requestFailedCount int // 新增一个计数器用于记录请求失败的次数
	mutex              sync.Mutex
)

// 增加请求计数器
func incrementRequestCount() {
	mutex.Lock()
	defer mutex.Unlock()
	requestCount++
}

// 增加请求失败计数器
func incrementRequestFailedCount() {
	mutex.Lock()
	defer mutex.Unlock()
	requestFailedCount++
}

// 增加一个计数器来记录请求次数
// GetRequestCount returns the total number of POST requests made
func GetRequestCount() int {
	mutex.Lock()
	defer mutex.Unlock()
	return requestCount
}

// 获取请求失败计数
func GetRequestFailedCount() int {
	mutex.Lock()
	defer mutex.Unlock()
	return requestFailedCount
}

// 增加请求计数并返回响应
func Post(path string, formData map[string]string, name string) Response {
	incrementRequestCount() // 增加请求计数器

	httpClient := resty.New().
		SetTimeout(10 * time.Second) // 设置超时时间为10秒

	jsonRequest, err := json.Marshal(formData)
	if err != nil {
		fmt.Println("JSON编码失败:", err)
	}
	sadsada := string(jsonRequest)

	rsp, err := httpClient.R().
		EnableTrace().
		SetFormData(formData).
		Post(config.GetDomain() + path)

	if err != nil {
		// 检查是否是超时错误
		if os.IsTimeout(err) {
			fmt.Println("HTTP请求超时:", err)
			// 记录超时错误到 CSV
			record := []string{time.Now().Format("15:04:05"), name, config.GetDomain() + path, sadsada, "OP response timeout"}
			if err := writer.Write(record); err != nil {
				fmt.Println("写入CSV失败:", err)
			}
			writer.Flush()
			incrementRequestFailedCount() // 增加请求失败计数器
			return Response{
				Body:     nil,
				Code:     0,
				Error:    fmt.Errorf("request timeout"),
				WriteCSV: false,
				CSVError: err,
				Name:     name,
			}
		}

		fmt.Println("HTTP请求失败:", err)
		incrementRequestFailedCount() // 增加请求失败计数器
		return Response{
			Body:     nil,
			Code:     0,
			Error:    err,
			WriteCSV: false,
			CSVError: err,
			Name:     name,
		}
	}

	// 获取当前日期和时间
	currentTime := time.Now()
	netwie := currentTime.Format("15:04:05")

	apiaa := config.GetDomain() + path

	// 打印原始的响应内容
	responseBody := string(rsp.Body())
	fmt.Println("响应内容:", responseBody)

	// 解码 JSON 响应
	var decodedResponse map[string]interface{}
	if err := json.Unmarshal(rsp.Body(), &decodedResponse); err != nil {
		fmt.Println("解码JSON失败:", err)
		fmt.Println("完整的响应内容:", responseBody)
		// 即使解码失败，仍将原始的响应内容写入 CSV
		decodedResponse = map[string]interface{}{} // 空的 map
	}

	// 将解码后的内容转换为 JSON 字符串，以便记录到 CSV
	jsonResponse, err := json.Marshal(decodedResponse)
	if err != nil {
		fmt.Println("JSON编码失败:", err)
		incrementRequestFailedCount() // 增加请求失败计数器
		return Response{
			Body:     rsp.Body(),
			Code:     rsp.StatusCode(),
			Error:    err,
			WriteCSV: false,
			CSVError: err,
			Name:     name,
		}
	}

	// 创建 CSV 记录
	record := []string{netwie, name, apiaa, sadsada, string(jsonResponse)}
	fmt.Println(record)

	// 将记录写入 CSV 文件
	if err := writer.Write(record); err != nil {
		fmt.Println("写入CSV失败:", err)
		// 如果写入 CSV 失败，返回错误，但仍然返回原始的响应内容
		incrementRequestFailedCount() // 增加请求失败计数器
		return Response{
			Body:     rsp.Body(),
			Code:     rsp.StatusCode(),
			Error:    err,
			WriteCSV: false,
			CSVError: err,
			Name:     name,
		}
	}
	// 刷新缓冲区
	writer.Flush()

	// 如果没有发生错误，返回正常的 Response
	return Response{
		Body:     rsp.Body(),
		Code:     rsp.StatusCode(),
		Error:    nil,
		WriteCSV: true,
		CSVError: nil,
		Name:     name,
	}
}
func PostWithSignature(path string, formData map[string]string, signStr string) Response {
	httpClient := resty.New()
	has := md5.Sum([]byte(signStr))
	signature := fmt.Sprintf("%x", has)
	rsp, err := httpClient.R().EnableTrace().
		SetHeader("signature", signature).
		SetFormData(formData).
		Post(config.GetDomain() + path)
	fmt.Println("wrong typo")
	return Response{
		Body:  rsp.Body(),
		Code:  rsp.StatusCode(),
		Error: err,
	}
}

func PostWithWrongSignature(path string, formData map[string]string) Response {
	httpClient := resty.New()
	signature := "wrongbase64input1234567890"
	rsp, err := httpClient.R().EnableTrace().
		SetHeader("signature", signature).
		SetFormData(formData).
		Post(config.GetDomain() + path)
	fmt.Println("wrong typo")

	return Response{
		Body:  rsp.Body(),
		Code:  rsp.StatusCode(),
		Error: err,
	}
}
