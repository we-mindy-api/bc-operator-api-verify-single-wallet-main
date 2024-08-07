package testcase

import (
	"errors"
	"fmt"
	"time"

	"github.com/PGITAb/bc-operator-api-verify/api"
	"github.com/PGITAb/bc-operator-api-verify/config"
	"github.com/rs/xid"
)

type SingleCreditResponse struct {
	Balance  int    `json:"balance"`
	Currency string `json:"currency"`
	Time     int    `json:"time"`
	RefID    string `json:"refID"`
}

// 成功案例
func SingleCredit_Success() api.TestCaseResultGroup {
	betID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	gameRoundID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	parentBetID := betID
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
		}, "Credit: Success-Call deposit").
		Expect([]api.TestCase{
			IsStatusCode(200),
			ContainFields([]string{"balance", "currency", "time", "refID"}),
			DecodeToStruct(&SingleDebitResponse{}),
		})

	// Get original balance
	bal1 := SingleGetBalanceResponse{}
	err := api.
		Post("/balance", map[string]string{
			"token":      config.GetToken(),
			"operatorID": config.GetOperatorID(),
			"appSecret":  config.GetAppSecret(),
			"playerID":   config.GetPlayerID(),
		}, "Credit: Success-Get original balance").Scan(&bal1)
	if err != nil {
		return api.TestCaseResultGroup{
			Name:     "Credit: Success",
			Endpoint: "/credit",
			Results: []api.TestCaseResult{
				{
					Name:  "Get player balance",
					Error: errors.New("get player balance failed"),
				},
			},
		}
	}

	// Call credit
	result = append(result, api.
		Post("/credit", map[string]string{
			"data": "[{" +
				"\"operatorID\": \"" + config.GetOperatorID() + "\"," +
				"\"appSecret\": \"" + config.GetAppSecret() + "\"," +
				"\"playerID\": \"" + config.GetPlayerID() + "\"," +
				"\"gameID\": \"" + config.GetGameID() + "\"," +
				"\"betID\": \"" + betID + "\"," +
				"\"amount\": \"" + fmt.Sprint(config.GetAmount_2()) + "\"," +
				"\"validBetAmount\": \"" + fmt.Sprint(config.GetAmount_3()) + "\"," +
				"\"gameStatus\": \"" + config.GetGameStatusWin() + "\"," +
				"\"gameResult\": \"" + config.GetGameResult() + "\"," +
				"\"currency\": \"" + config.GetCurrency() + "\"," +
				"\"type\": \"game\"," +
				"\"time\": \"" + fmt.Sprint(time.Now().Unix()) + "\"," + "\"odds\": \"" + config.GetOdd() + "\"" +
				"}]",
		}, "Credit: Success-call credit").
		Expect([]api.TestCase{
			IsStatusCode(200),
			ContainFields([]string{"balance", "currency", "time", "refID"}),
			DecodeToStruct(&SingleCreditResponse{}),
		})...)

	return api.TestCaseResultGroup{
		Name:     "Credit: Success",
		Endpoint: "/credit",
		Results:  result,
	}
}

// 無效參數
func SingleCredit_InvalidParameter() api.TestCaseResultGroup {
	betID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	gameRoundID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	parentBetID := betID
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
		}, "Credit: Invalid parameter-call deposit").
		Expect([]api.TestCase{
			IsStatusCode(200),
			ContainFields([]string{"balance", "currency", "time", "refID"}),
			DecodeToStruct(&SingleDebitResponse{}),
		})

	result = append(result, api.
		Post("/credit", map[string]string{
			"data": "[{" +
				"\"operatorID\": \"" + config.GetOperatorID() + "\"," +
				"\"appSecret\": \"" + config.GetAppSecret() + "\"," +
				// "\"playerID\": \"" + config.GetPlayerID() + "\"," +
				"\"gameID\": \"" + config.GetGameID() + "\"," +
				"\"betID\": \"" + betID + "\"," +
				"\"amount\": \"" + fmt.Sprint(config.GetAmount()) + "\"," +
				"\"validBetAmount\": \"" + fmt.Sprint(config.GetAmount_3()) + "\"," +
				"\"gameStatus\": \"" + config.GetGameStatusWin() + "\"," +
				"\"gameResult\": \"" + config.GetGameResult() + "\"," +
				"\"currency\": \"" + config.GetCurrency() + "\"," +
				"\"type\": \"game\"," +
				"\"time\": \"" + fmt.Sprint(time.Now().Unix()) + "\"," + "\"odds\": \"" + config.GetOdd() + "\"" +
				"}]",
		}, "Credit: Invalid parameter").
		Expect([]api.TestCase{
			IsStatusCode(400),
			ContainFields([]string{"error"}),
		})...)
	return api.TestCaseResultGroup{
		Name:     "Credit: Invalid parameter",
		Endpoint: "/credit",
		Results:  result,
	}
}

// 憑證錯誤
func SingleCredit_InvalidSecret() api.TestCaseResultGroup {
	betID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	gameRoundID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	parentBetID := betID
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
		}, "Credit: Invalid secret-call debit").
		Expect([]api.TestCase{
			IsStatusCode(200),
			ContainFields([]string{"balance", "currency", "time", "refID"}),
			DecodeToStruct(&SingleDebitResponse{}),
		})
	result = append(result, api.
		Post("/credit", map[string]string{
			"data": "[{" +
				"\"operatorID\": \"" + config.GetOperatorID() + "\"," +
				"\"appSecret\": \"xxxxxxxxxxxxxxxxxxxx\"," +
				"\"playerID\": \"" + config.GetPlayerID() + "\"," +
				"\"gameID\": \"" + config.GetGameID() + "\"," +
				"\"betID\": \"" + xid.New().String() + "\"," +
				"\"amount\": \"" + fmt.Sprint(config.GetAmount()) + "\"," +
				"\"validBetAmount\": \"" + fmt.Sprint(config.GetAmount_3()) + "\"," +
				"\"gameStatus\": \"" + config.GetGameStatusWin() + "\"," +
				"\"gameResult\": \"" + config.GetGameResult() + "\"," +
				"\"currency\": \"" + config.GetCurrency() + "\"," +
				"\"type\": \"game\"," +
				"\"time\": \"" + fmt.Sprint(time.Now().Unix()) + "\"," + "\"odds\": \"" + config.GetOdd() + "\"" +
				"}]",
		}, "Credit: Invalid secret").
		Expect([]api.TestCase{
			IsStatusCode(401),
			ContainFields([]string{"error"}),
		})...)
	return api.TestCaseResultGroup{
		Name:     "Credit: Invalid secret",
		Endpoint: "/credit",
		Results:  result,
	}
}

// 查無投注注單
func SingleCredit_DebitNotFound() api.TestCaseResultGroup {
	betID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	gameRoundID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	parentBetID := betID
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
		}, "Credit: Bet ID not found call debit").
		Expect([]api.TestCase{
			IsStatusCode(200),
			ContainFields([]string{"balance", "currency", "time", "refID"}),
			DecodeToStruct(&SingleDebitResponse{}),
		})

	// Call credit
	result = append(result, api.
		Post("/credit", map[string]string{
			"data": "[{" +
				"\"operatorID\": \"" + config.GetOperatorID() + "\"," +
				"\"appSecret\": \"" + config.GetAppSecret() + "\"," +
				"\"playerID\": \"" + config.GetPlayerID() + "\"," +
				"\"gameID\": \"" + config.GetGameID() + "\"," +
				"\"betID\": \"" + "betIDNotExist" + "\"," +
				"\"amount\": \"" + fmt.Sprint(config.GetAmount_2()) + "\"," +
				"\"validBetAmount\": \"" + fmt.Sprint(config.GetAmount_3()) + "\"," +
				"\"gameStatus\": \"" + config.GetGameStatusWin() + "\"," +
				"\"gameResult\": \"" + config.GetGameResult() + "\"," +
				"\"currency\": \"" + config.GetCurrency() + "\"," +
				"\"type\": \"game\"," +
				"\"time\": \"" + fmt.Sprint(time.Now().Unix()) + "\"," + "\"odds\": \"" + config.GetOdd() + "\"" +
				"}]",
		}, "Credit: Bet ID not found").
		Expect([]api.TestCase{
			IsStatusCode(410),
			ContainFields([]string{"error"}),
		})...)
	return api.TestCaseResultGroup{
		Name:     "Credit: Bet ID not found",
		Endpoint: "/credit",
		Results:  result,
	}
}

// 重複交易
func SingleCredit_DuplicateTransaction() api.TestCaseResultGroup {
	betID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	gameRoundID := time.Now().Format("2006-01-02") + "-" + xid.New().String()
	parentBetID := betID
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
		}, "Credit: Duplicate credit transaction call deposit").
		Expect([]api.TestCase{
			IsStatusCode(200),
			ContainFields([]string{"balance", "currency", "time", "refID"}),
			DecodeToStruct(&SingleDebitResponse{}),
		})

	// Credit data
	formData := map[string]string{
		"data": "[{" +
			"\"operatorID\": \"" + config.GetOperatorID() + "\"," +
			"\"appSecret\": \"" + config.GetAppSecret() + "\"," +
			"\"playerID\": \"" + config.GetPlayerID() + "\"," +
			"\"gameID\": \"" + config.GetGameID() + "\"," +
			"\"betID\": \"" + betID + "\"," +
			"\"parentBetID\": \"" + parentBetID + "\"," +
			"\"amount\": \"" + fmt.Sprint(config.GetAmount_2()) + "\"," +
			"\"validBetAmount\": \"" + fmt.Sprint(config.GetAmount_3()) + "\"," +
			"\"gameStatus\": \"" + config.GetGameStatusWin() + "\"," +
			"\"gameResult\": \"" + config.GetGameResult() + "\"," +
			"\"currency\": \"" + config.GetCurrency() + "\"," +
			"\"type\": \"game\"," +
			"\"time\": \"" + fmt.Sprint(time.Now().Unix()) + "\"," + "\"odds\": \"" + config.GetOdd() + "\"" +
			"}]",
	}

	// Call credit
	result = append(result, api.
		Post("/credit", formData, "Credit: Duplicate credit transaction call credit ").
		Expect([]api.TestCase{
			IsStatusCode(200),
			ContainFields([]string{"balance", "currency", "time", "refID"}),
			DecodeToStruct(&SingleCreditResponse{}),
		})...)

	// Call duplicate credit
	result = append(result, api.
		Post("/credit", formData, "Credit: Duplicate credit transaction call deposit- Call duplicate credit").
		Expect([]api.TestCase{
			IsStatusCode(409),
			ContainFields([]string{"error"}),
		})...)
	return api.TestCaseResultGroup{
		Name:     "Credit: Duplicate credit transaction",
		Endpoint: "/credit",
		Results:  result,
	}
}
