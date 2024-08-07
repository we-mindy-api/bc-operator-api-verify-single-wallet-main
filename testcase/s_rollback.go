package testcase

import (
	"errors"
	"fmt"
	"time"

	"github.com/PGITAb/bc-operator-api-verify/api"
	"github.com/PGITAb/bc-operator-api-verify/config"
	"github.com/rs/xid"
)

type SingleRollbackResponse struct {
	Balance  int    `json:"balance"`
	Currency string `json:"currency"`
	Time     int    `json:"time"`
	RefID    string `json:"refID"`
}

func SingleRollbackDebited_Success() api.TestCaseResultGroup {
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
		}, "Debit: Success-Get player balance").Scan(&bal1)
	if err != nil {
		return api.TestCaseResultGroup{
			Name:     "Debit: Success",
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
		}, "Rollback Debit: Success call deposit").
		Expect([]api.TestCase{
			IsStatusCode(200),
			ContainFields([]string{"balance", "currency", "time", "refID"}),
			DecodeToStruct(&SingleDebitResponse{}),
		})

	// Call rollback
	result = append(result, api.
		Post("/rollback", map[string]string{
			"operatorID": config.GetOperatorID(),
			"appSecret":  config.GetAppSecret(),
			"playerID":   config.GetPlayerID(),
			"gameID":     config.GetGameID(),
			"betID":      betID,
			"amount":     fmt.Sprint(config.GetAmount()),
			"currency":   config.GetCurrency(),
			"type":       "cancel",
			"time":       fmt.Sprint(time.Now().Unix()),
		}, "Rollback Debit: Success").
		Expect([]api.TestCase{
			IsStatusCode(200),
			ContainFields([]string{"balance", "currency", "time", "refID"}),
			DecodeToStruct(&SingleRollbackResponse{}),
		})...)

	return api.TestCaseResultGroup{
		Name:     "Rollback Debit: Success",
		Endpoint: "/rollback",
		Results:  result,
	}
}

func SingleRollbackCreditedPlayer_Success() api.TestCaseResultGroup {

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
		}, "Rollback Credited Player Win: Success get balabce").Scan(&bal1)
	if err != nil {
		return api.TestCaseResultGroup{
			Name:     "Debit: Success",
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
		}, "Rollback Credited Player Win: Success call deposit").
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
				"\"betID\": \"" + betID + "\"," +
				"\"amount\": \"" + fmt.Sprint(0) + "\"," +
				"\"validBetAmount\": \"" + fmt.Sprint(config.GetAmount_3()) + "\"," +
				"\"gameStatus\": \"" + config.GetGameStatusLoss() + "\"," +
				"\"gameResult\": \"" + config.GetGameResult() + "\"," +
				"\"currency\": \"" + config.GetCurrency() + "\"," +
				"\"type\": \"game\"," +
				"\"time\": \"" + fmt.Sprint(time.Now().Unix()) + "\"," +
				"\"odds\": \"" + config.GetOdd() + "\"" +
				"}]",
		}, "Rollback Credited Player Win: Success call credit").
		Expect([]api.TestCase{
			IsStatusCode(200),
			ContainFields([]string{"balance", "currency", "time", "refID"}),
			DecodeToStruct(&SingleCreditResponse{}),
		})...)

	// 已派彩玩家輸錢時, 金額=下注-派彩

	// Call rollback
	result = append(result, api.
		Post("/rollback", map[string]string{
			"operatorID": config.GetOperatorID(),
			"appSecret":  config.GetAppSecret(),
			"playerID":   config.GetPlayerID(),
			"gameID":     config.GetGameID(),
			"betID":      betID,
			"amount":     fmt.Sprint(config.GetAmount()),
			"currency":   config.GetCurrency(),
			"type":       "cancel",
			"time":       fmt.Sprint(time.Now().Unix()),
		}, "Rollback Credited Player Loss: Success").
		Expect([]api.TestCase{
			IsStatusCode(200),
			ContainFields([]string{"balance", "currency", "time", "refID"}),
			DecodeToStruct(&SingleRollbackResponse{}),
		})...)

	return api.TestCaseResultGroup{
		Name:     "Rollback Credited Player Loss: Success",
		Endpoint: "/rollback",
		Results:  result,
	}
}

func SingleRollback_InvalidParameter() api.TestCaseResultGroup {
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
		}, "Rollback: Invalid parameter call debit").
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
				"\"betID\": \"" + betID + "\"," +
				"\"amount\": \"" + fmt.Sprint(config.GetAmount_2()) + "\"," +
				"\"validBetAmount\": \"" + fmt.Sprint(config.GetAmount_3()) + "\"," +
				"\"gameStatus\": \"" + config.GetGameStatusWin() + "\"," +
				"\"gameResult\": \"" + config.GetGameResult() + "\"," +
				"\"currency\": \"" + config.GetCurrency() + "\"," +
				"\"type\": \"game\"," +
				"\"time\": \"" + fmt.Sprint(time.Now().Unix()) + "\"," +
				"\"odds\": \"" + config.GetOdd() + "\"" +
				"}]",
		}, "Rollback: Invalid parameter call credit").
		Expect([]api.TestCase{
			IsStatusCode(200),
			ContainFields([]string{"balance", "currency", "time", "refID"}),
			DecodeToStruct(&SingleCreditResponse{}),
		})...)

	// Call rollback
	result = append(result, api.
		Post("/rollback", map[string]string{
			"operatorID": config.GetOperatorID(),
			"appSecret":  config.GetAppSecret(),
			// "playerID":   config.GetPlayerID(),
			"gameID": config.GetGameID(),
			"betID":  betID,
			// "amount":   fmt.Sprint(config.GetAmount() - config.GetAmount_2()),
			"currency": config.GetCurrency(),
			"type":     "cancel",
			"time":     fmt.Sprint(time.Now().Unix()),
		}, "Rollback: Invalid parameter").
		Expect([]api.TestCase{
			IsStatusCode(400),
			ContainFields([]string{"error"}),
		})...)
	return api.TestCaseResultGroup{
		Name:     "Rollback: Invalid parameter",
		Endpoint: "/rollback",
		Results:  result,
	}
}

func SingleRollback_InvalidSecret() api.TestCaseResultGroup {
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
		}, "Rollback: Invalid secret call debit").
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
				"\"betID\": \"" + betID + "\"," +
				"\"amount\": \"" + fmt.Sprint(config.GetAmount_2()) + "\"," +
				"\"validBetAmount\": \"" + fmt.Sprint(config.GetAmount_3()) + "\"," +
				"\"gameStatus\": \"" + config.GetGameStatusWin() + "\"," +
				"\"gameResult\": \"" + config.GetGameResult() + "\"," +
				"\"currency\": \"" + config.GetCurrency() + "\"," +
				"\"type\": \"game\"," +
				"\"time\": \"" + fmt.Sprint(time.Now().Unix()) + "\"," + "\"odds\": \"" + config.GetOdd() + "\"" +
				"}]",
		}, "Rollback: Invalid secret call credit").
		Expect([]api.TestCase{
			IsStatusCode(200),
			ContainFields([]string{"balance", "currency", "time", "refID"}),
			DecodeToStruct(&SingleCreditResponse{}),
		})...)

	// Call rollback
	result = append(result, api.
		Post("/rollback", map[string]string{
			"operatorID": config.GetOperatorID(),
			"appSecret":  "xxxxxxxxxxxxxxxxxxxx",
			"playerID":   config.GetPlayerID(),
			"gameID":     config.GetGameID(),
			"betID":      betID,
			"amount":     fmt.Sprint(config.GetAmount() - config.GetAmount_2()),
			"currency":   config.GetCurrency(),
			"type":       "cancel",
			"time":       fmt.Sprint(time.Now().Unix()),
		}, "Rollback: Invalid secret").
		Expect([]api.TestCase{
			IsStatusCode(401),
			ContainFields([]string{"error"}),
		})...)
	return api.TestCaseResultGroup{
		Name:     "Rollback: Invalid secret",
		Endpoint: "/rollback",
		Results:  result,
	}
}

func SingleRollback_DebitNotFound() api.TestCaseResultGroup {
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
		}, "Rollback: Bet ID not found call debit").
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
				"\"betID\": \"" + betID + "\"," +
				"\"amount\": \"" + fmt.Sprint(config.GetAmount_2()) + "\"," +
				"\"validBetAmount\": \"" + fmt.Sprint(config.GetAmount_3()) + "\"," +
				"\"gameStatus\": \"" + config.GetGameStatusWin() + "\"," +
				"\"gameResult\": \"" + config.GetGameResult() + "\"," +
				"\"currency\": \"" + config.GetCurrency() + "\"," +
				"\"type\": \"game\"," +
				"\"time\": \"" + fmt.Sprint(time.Now().Unix()) + "\"," + "\"odds\": \"" + config.GetOdd() + "\"" +
				"}]",
		}, "Rollback: Bet ID not found call credit").
		Expect([]api.TestCase{
			IsStatusCode(200),
			ContainFields([]string{"balance", "currency", "time", "refID"}),
			DecodeToStruct(&SingleCreditResponse{}),
		})...)

	// Call rollback
	result = append(result, api.
		Post("/rollback", map[string]string{
			"operatorID": config.GetOperatorID(),
			"appSecret":  config.GetAppSecret(),
			"playerID":   config.GetPlayerID(),
			"gameID":     config.GetGameID(),
			"betID":      "betIDNotExist",
			"amount":     fmt.Sprint(config.GetAmount() - config.GetAmount_2()),
			"currency":   config.GetCurrency(),
			"type":       "cancel",
			"time":       fmt.Sprint(time.Now().Unix()),
		}, "Rollback: Bet ID not found").
		Expect([]api.TestCase{
			IsStatusCode(410),
			ContainFields([]string{"error"}),
		})...)
	return api.TestCaseResultGroup{
		Name:     "Rollback: Bet ID not found",
		Endpoint: "/rollback",
		Results:  result,
	}
}

func SingleRollback_DuplicateTransaction() api.TestCaseResultGroup {
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
		}, "Rollback: Duplicate call debit").
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
				"\"betID\": \"" + betID + "\"," +
				"\"amount\": \"" + fmt.Sprint(config.GetAmount_2()) + "\"," +
				"\"validBetAmount\": \"" + fmt.Sprint(config.GetAmount_3()) + "\"," +
				"\"gameStatus\": \"" + config.GetGameStatusWin() + "\"," +
				"\"gameResult\": \"" + config.GetGameResult() + "\"," +
				"\"currency\": \"" + config.GetCurrency() + "\"," +
				"\"type\": \"game\"," +
				"\"time\": \"" + fmt.Sprint(time.Now().Unix()) + "\"," + "\"odds\": \"" + config.GetOdd() + "\"" +
				"}]",
		}, "Rollback: Duplicate call credit").
		Expect([]api.TestCase{
			IsStatusCode(200),
			ContainFields([]string{"balance", "currency", "time", "refID"}),
			DecodeToStruct(&SingleCreditResponse{}),
		})...)

	formData := map[string]string{
		"operatorID": config.GetOperatorID(),
		"appSecret":  config.GetAppSecret(),
		"playerID":   config.GetPlayerID(),
		"gameID":     config.GetGameID(),
		"betID":      betID,
		"amount":     fmt.Sprint(0),
		"currency":   config.GetCurrency(),
		"type":       "cancel",
		"time":       fmt.Sprint(time.Now().Unix()),
	}

	// Call rollback
	result = append(result, api.
		Post("/rollback", formData, "Rollback: Duplicate call rollback").
		Expect([]api.TestCase{
			IsStatusCode(200),
			ContainFields([]string{"balance", "currency", "time", "refID"}),
			DecodeToStruct(&SingleRollbackResponse{}),
		})...)

	// Call duplicate rollback
	result = append(result, api.
		Post("/rollback", formData, "Rollback: Duplicate").
		Expect([]api.TestCase{
			IsStatusCode(409),
			ContainFields([]string{"error"}),
		})...)

	return api.TestCaseResultGroup{
		Name:     "Rollback: Duplicate rollback transaction",
		Endpoint: "/rollback",
		Results:  result,
	}
}
