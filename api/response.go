package api

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Response struct {
	Body     []byte
	Code     int
	Error    error
	WriteCSV bool  // 是否成功写入 CSV
	CSVError error // CSV 写入时的错误信息
	Name     string
}

type TestCase struct {
	Name string
	Func func(Response) error
}

type TestCaseResult struct {
	Name  string
	Error error
}

type TestCaseResultGroup struct {
	Name     string
	Endpoint string
	Results  []TestCaseResult
}

type ResultCollector struct {
	Func    func(TestCaseResultGroup)
	Results []TestCaseResultGroup
}

func (rsp Response) Expect(testcases []TestCase) []TestCaseResult {
	result := make([]TestCaseResult, 0)
	for _, tc := range testcases {
		err := tc.Func(rsp)
		result = append(result, TestCaseResult{
			Name:  tc.Name,
			Error: err,
		})
	}
	return result
}

func (rsp Response) Scan(pointer interface{}) error {
	err := json.Unmarshal(rsp.Body, pointer)
	if err != nil {
		return errors.New("parse json to struct failed")
	}
	return nil
}

func (rsp Response) Print() Response {
	fmt.Println(string(rsp.Body))
	return rsp
}

func (tcr TestCaseResult) Print(color bool) string {
	if tcr.Error == nil {
		if color {
			return fmt.Sprintf("\033[0;32m[PASS] %s\033[0m", tcr.Name)
		}
		return fmt.Sprintf("[PASS] %s", tcr.Name)
	}
	if color {
		return fmt.Sprintf("\033[0;31m[FAILURE] %s: %s\033[0m", tcr.Name, tcr.Error.Error())
	}
	return fmt.Sprintf("[FAILURE] %s: %s", tcr.Name, tcr.Error.Error())
}

func (gp TestCaseResultGroup) Print(color bool) string {
	str := fmt.Sprintf("# %s - %s\n", gp.Name, gp.Endpoint)
	for _, r := range gp.Results {
		str += r.Print(color) + "\n"
	}
	return str
}
