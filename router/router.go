package router

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PGITAb/bc-operator-api-verify/api"
	"github.com/PGITAb/bc-operator-api-verify/config"
	"github.com/PGITAb/bc-operator-api-verify/helper"
	"github.com/PGITAb/bc-operator-api-verify/testlist"
)

type SummaryItem struct {
	Name   string
	Total  int
	Failed int
}

type TemplateData struct {
	Results         []api.TestCaseResultGroup
	Summary         []SummaryItem
	HtmlcurrentTime string
}

func ResultHtml() {
	questions := [5]string{
		"Please enter the api url：",
		"Please enter the token：",
		"Please enter the operatorID：",
		"Please enter the appSecret：",
		"Please enter the playerID：",
	}
	htmlcurrentTime := time.Now().Format("Monday, 02 January 2006 15:04:05")
	var vals []string
	var val string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println(questions[0])
	for i := 1; scanner.Scan(); i++ {
		val = scanner.Text()
		if strings.TrimSpace(val) == "" {
			i--
		} else {
			fmt.Println(val)
			vals = append(vals, strings.TrimSpace(val))
		}

		if i == 5 {
			break
		}
		fmt.Println(questions[i])
	}

	config.SetApiServer(vals[0])
	config.SetToken(vals[1])
	config.SetOperatorID(vals[2])
	config.SetAppSecret(vals[3])
	config.SetPlayerID(vals[4])

	// 文件
	var f *os.File
	if config.GetFileOutput() {
		filename := "output/" + time.Now().Format("2006-01-02-150405-op-api-verify-"+vals[2]+".html")
		var err error
		f, err = os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 创建结果收集器
	results := make([]api.TestCaseResultGroup, 0)
	collector := api.ResultCollector{
		Func: func(gp api.TestCaseResultGroup) {
			results = append(results, gp)
			fmt.Println(gp.Print(config.GetColorMode()))
		},
	}

	// 运行任务
	testlist.SingleWallet(collector)

	// 准备数据
	totalRequests := api.GetRequestCount()
	totalFailedRequests := api.GetRequestFailedCount()
	totalAssertions := 0
	totalFailedAssertions := 0

	for _, gp := range results {
		for _, r := range gp.Results {
			if r.Error != nil {
				totalFailedAssertions++
			}
			totalAssertions++
		}
	}

	summaryItems := []SummaryItem{
		{"Requests", totalRequests, totalFailedRequests},
		{"Assertions", totalAssertions, totalFailedAssertions},
	}

	// 加载模板文件
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	// 准备数据
	data := TemplateData{
		Results:         results,
		Summary:         summaryItems,
		HtmlcurrentTime: htmlcurrentTime,
	}

	// 执行模板并将结果写入文件
	if config.GetFileOutput() {
		err = tmpl.Execute(f, data)
		if err != nil {
			log.Fatal("Error executing template:", err)
		}
	} else {
		err = tmpl.Execute(os.Stdout, data)
		if err != nil {
			log.Fatal("Error executing template:", err)
		}
	}
}

//---------------------------------------------------------------------//
//---------------------------------------------------------------------//
//---------------------------------------------------------------------//

func ResultText() {
	questions := [5]string{
		"Please enter the api url：",
		"Please enter the token：",
		"Please enter the operatorID：",
		"Please enter the appSecret：",
		"Please enter the playerID：",
	}

	var vals []string
	var val string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println(questions[0])
	for i := 1; scanner.Scan(); i++ {
		val = scanner.Text()
		if strings.TrimSpace(val) == "" {
			i--
		} else {
			fmt.Println(val)
			vals = append(vals, strings.TrimSpace(val))
		}

		if i == 5 {
			break
		}
		fmt.Println(questions[i])
	}

	config.SetApiServer(vals[0])
	config.SetToken(vals[1])
	config.SetOperatorID(vals[2])
	config.SetAppSecret(vals[3])
	config.SetPlayerID(vals[4])

	// file
	var f *os.File
	if config.GetFileOutput() {
		filename := "output/" + time.Now().Format("2006-01-02-150405-op-api-verify-"+vals[2]+".txt")
		var err error
		f, err = os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 创建结果收集器
	results := make([]api.TestCaseResultGroup, 0) // 创建一个空的 TestCaseResultGroup 切片，用于收集测试结果
	collector := api.ResultCollector{
		Func: func(gp api.TestCaseResultGroup) {
			// 将新的测试结果组添加到 results 切片中
			results = append(results, gp)
			// 打印测试结果组，打印时根据配置决定是否使用颜色
			fmt.Println(gp.Print(config.GetColorMode()))
			// 如果配置了文件输出，则将结果写入文件
			if config.GetFileOutput() {
				f.WriteString(gp.Print(false) + "\n")
			}
		},
	}

	// 运行任务
	testlist.SingleWallet(collector) // 运行单钱包测试，并将 collector 传递给它以收集结果

	// 打印总结
	fmt.Println(printSummary(results, config.GetColorMode())) // 打印所有测试结果的总结，是否使用颜色根据配置决定
	if config.GetFileOutput() {
		f.WriteString(printSummary(results, false)) // 如果配置了文件输出，则将总结写入文件
	}
}

func printSummary(gps []api.TestCaseResultGroup, color bool) string {
	strBu := strings.Builder{}                                                                      // 创建一个字符串构建器，用于高效地拼接字符串
	strBu.WriteString("Task Name\t\t\t\t\tPass\tFail\n")                                            // 写入表头
	strBu.WriteString("========================================================================\n") // 写入分隔线

	totalRequests := api.GetRequestCount()             // 获取总请求数
	totalFailedRequests := api.GetRequestFailedCount() // 获取总失败请求数
	totalAssertions := 0                               // 初始化总断言数
	totalFailedAssertions := 0                         // 初始化总失败断言数

	// 遍历每个测试结果组
	for _, gp := range gps {
		pass := 0 // 初始化通过的测试计数
		fail := 0 // 初始化失败的测试计数
		// 遍历每个测试结果
		for _, r := range gp.Results {
			if r.Error != nil {
				fail++                  // 如果有错误，增加失败计数
				totalFailedAssertions++ // 增加总失败断言数
			} else {
				pass++ // 没有错误，增加通过计数
			}
			totalAssertions++ // 每个测试用例计为一个断言
		}

		str := "" // 用于存储格式化的字符串
		if color {
			// 使用颜色输出
			str = helper.RightPad(gp.Name, 40) + "\t\033[0;32m%d\t" // 右对齐任务名称，绿色显示通过的数量
			if fail == 0 {
				str += "\033[0m%d\n" // 如果没有失败，恢复默认颜色显示通过的数量
			} else {
				str += "\033[0;31m%d\033[0m\n" // 如果有失败，红色显示失败的数量
			}
		} else {
			// 不使用颜色输出
			str = helper.RightPad(gp.Name, 40) + "\t%d\t%d\n" // 右对齐任务名称，显示通过和失败的数量
		}
		strBu.WriteString(fmt.Sprintf(str, pass, fail)) // 将格式化的字符串写入字符串构建器
	}

	// 添加总结表格
	strBu.WriteString("\nSUMMARY ITEM\tTOTAL\tFAILED\n")                                           // 写入总结表头
	strBu.WriteString(fmt.Sprintf("Requests\t%d\t%d\n", totalRequests, totalFailedRequests))       // 写入总请求数和总失败数
	strBu.WriteString(fmt.Sprintf("Assertions\t%d\t%d\n", totalAssertions, totalFailedAssertions)) // 写入总断言数和总失败数

	return strBu.String() // 返回最终的总结字符串
}
