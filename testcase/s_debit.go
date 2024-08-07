package testcase

import (
	"errors"
	"fmt"
	"time"

	"github.com/PGITAb/bc-operator-api-verify/api"
	"github.com/PGITAb/bc-operator-api-verify/config"
	"github.com/rs/xid"
)

type SingleDebitResponse struct {
	Balance  int    `json:"balance"`
	Currency string `json:"currency"`
	Time     int    `json:"time"`
	RefID    string `json:"refID"`
}

// 成功案例
func SingleDebit_Success() api.TestCaseResultGroup {
	betID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	gameRoundID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	parentBetID := betID
	// Get original balance
	bal1 := SingleGetBalanceResponse{}
	err := api.
		Post("/balance", map[string]string{
			"token":      config.GetToken(),
			"operatorID": config.GetOperatorID(),
			"appSecret":  config.GetAppSecret(),
			"playerID":   config.GetPlayerID(),
		}, "Debit: Success").Scan(&bal1)
	if err != nil {
		return api.TestCaseResultGroup{
			Name:     "Debit: Success call balance",
			Endpoint: "/debit",
			Results: []api.TestCaseResult{
				{
					Name:  "Get player balance",
					Error: errors.New("Get player balance failed: " + err.Error()),
				},
			},
		}
	}

	// Call deposit
	result := api.
		Post("/debit", map[string]string{
			"token":       config.GetToken(),
			"operatorID":  config.GetOperatorID(),
			"appSecret":   config.GetAppSecret(),
			"playerID":    config.GetPlayerID(),
			"gameID":      config.GetGameID(),
			"betID":       betID,
			"gameRoundID": gameRoundID,
			"parentBetID": parentBetID,
			"betType":     config.GetBetType(),
			"amount":      fmt.Sprint(config.GetAmount()),
			"currency":    config.GetCurrency(),
			"type":        "bet",
			"time":        fmt.Sprint(time.Now().Unix()),
			"ip":          config.GetIP(),
		}, "Debit: Success").
		Expect([]api.TestCase{
			IsStatusCode(200),
			ContainFields([]string{"balance", "currency", "time", "refID"}),
			DecodeToStruct(&SingleDebitResponse{}),
		})

	return api.TestCaseResultGroup{
		Name:     "Debit: Success",
		Endpoint: "/debit",
		Results:  result,
	}
}

// 無效參數
func SingleDebit_InvalidParameter() api.TestCaseResultGroup {
	betID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	gameRoundID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	parentBetID := betID
	result := api.
		Post("/debit", map[string]string{
			"token":      config.GetToken(),
			"operatorID": config.GetOperatorID(),
			"appSecret":  config.GetAppSecret(),
			// "playerID":    config.GetPlayerID(),
			"gameID":      config.GetGameID(),
			"betID":       betID,
			"gameRoundID": gameRoundID,
			"parentBetID": parentBetID,
			"betType":     config.GetBetType(),
			"amount":      fmt.Sprint(config.GetAmount()),
			"currency":    config.GetCurrency(),
			"type":        "bet",
			"time":        fmt.Sprint(time.Now().Unix()),
			"ip":          config.GetIP(),
		}, "Debit: Invalid parameter").
		Expect([]api.TestCase{
			IsStatusCode(400),
			ContainFields([]string{"error"}),
		})
	return api.TestCaseResultGroup{
		Name:     "Debit: Invalid parameter",
		Endpoint: "/debit",
		Results:  result,
	}
}

// TOKEN失效
func SingleDebit_InvalidToken() api.TestCaseResultGroup {
	betID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	gameRoundID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	parentBetID := betID
	result := api.
		Post("/debit", map[string]string{
			"token":       "xxxxxxxxxxxxxxxxxxxx",
			"operatorID":  config.GetOperatorID(),
			"appSecret":   config.GetAppSecret(),
			"playerID":    config.GetPlayerID(),
			"gameID":      config.GetGameID(),
			"betID":       betID,
			"gameRoundID": gameRoundID,
			"parentBetID": parentBetID,
			"betType":     config.GetBetType(),
			"amount":      fmt.Sprint(config.GetAmount()),
			"currency":    config.GetCurrency(),
			"type":        "bet",
			"time":        fmt.Sprint(time.Now().Unix()),
			"ip":          config.GetIP(),
		}, "Debit: Invalid token").
		Expect([]api.TestCase{
			IsStatusCode(404),
			ContainFields([]string{"error"}),
		})
	return api.TestCaseResultGroup{
		Name:     "Debit: Invalid token",
		Endpoint: "/debit",
		Results:  result,
	}
}

// 憑證錯誤
func SingleDebit_InvalidSecret() api.TestCaseResultGroup {
	betID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	gameRoundID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	parentBetID := betID
	result := api.
		Post("/debit", map[string]string{
			"token":       config.GetToken(),
			"operatorID":  config.GetOperatorID(),
			"appSecret":   "xxxxxxxxxxxxxxxxxxxx",
			"playerID":    config.GetPlayerID(),
			"gameID":      config.GetGameID(),
			"betID":       betID,
			"gameRoundID": gameRoundID,
			"parentBetID": parentBetID,
			"betType":     config.GetBetType(),
			"amount":      fmt.Sprint(config.GetAmount()),
			"currency":    config.GetCurrency(),
			"type":        "bet",
			"time":        fmt.Sprint(time.Now().Unix()),
			"ip":          config.GetIP(),
		}, "Debit: Invalid secret").
		Expect([]api.TestCase{
			IsStatusCode(401),
			ContainFields([]string{"error"}),
		})
	return api.TestCaseResultGroup{
		Name:     "Debit: Invalid secret",
		Endpoint: "/debit",
		Results:  result,
	}
}

// 餘額不足
func SingleDebit_BalanceNotEnough() api.TestCaseResultGroup {
	betID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	gameRoundID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	parentBetID := betID
	result := api.
		Post("/debit", map[string]string{
			"token":       config.GetToken(),
			"operatorID":  config.GetOperatorID(),
			"appSecret":   config.GetAppSecret(),
			"playerID":    config.GetPlayerID(),
			"gameID":      config.GetGameID(),
			"betID":       betID,
			"gameRoundID": gameRoundID,
			"parentBetID": parentBetID,
			"betType":     config.GetBetType(),
			"amount":      fmt.Sprint(99999999999999),
			"currency":    config.GetCurrency(),
			"type":        "bet",
			"time":        fmt.Sprint(time.Now().Unix()),
			"ip":          config.GetIP(),
		}, "Debit: Balance not enough").
		Expect([]api.TestCase{
			IsStatusCode(402),
			ContainFields([]string{"error"}),
		})
	return api.TestCaseResultGroup{
		Name:     "Debit: Balance not enough",
		Endpoint: "/debit",
		Results:  result,
	}
}

// 重複交易
func SingleDebit_DuplicateTransaction() api.TestCaseResultGroup {
	betID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	gameRoundID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	parentBetID := betID
	// 定义字符串
	//Namecsv := "Debit: Duplicate debit transaction"
	formData := map[string]string{
		"token":       config.GetToken(),
		"operatorID":  config.GetOperatorID(),
		"appSecret":   config.GetAppSecret(),
		"playerID":    config.GetPlayerID(),
		"gameID":      config.GetGameID(),
		"betID":       betID,
		"gameRoundID": gameRoundID,
		"parentBetID": parentBetID,
		"betType":     config.GetBetType(),
		"amount":      fmt.Sprint(config.GetAmount()),
		"currency":    config.GetCurrency(),
		"type":        "bet",
		"time":        fmt.Sprint(time.Now().Unix()),
		"ip":          config.GetIP(),
	}
	// Call deposit
	result := api.
		Post("/debit", formData, "Duplicate debit call deposit").
		Expect([]api.TestCase{
			IsStatusCode(200),
			ContainFields([]string{"balance", "currency", "time", "refID"}),
			DecodeToStruct(&SingleDebitResponse{}),
		})

	// Call duplicate deposit
	result = append(result, api.
		Post("/debit", formData, "Debit: Duplicate debit transaction").
		Expect([]api.TestCase{
			IsStatusCode(409),
			ContainFields([]string{"error"}),
		})...)

	return api.TestCaseResultGroup{
		Name:     "Debit: Duplicate debit transaction",
		Endpoint: "/debit",
		Results:  result,
	}
}
