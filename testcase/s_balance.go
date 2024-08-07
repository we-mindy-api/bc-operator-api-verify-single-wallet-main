package testcase

import (
	"github.com/PGITAb/bc-operator-api-verify/api"
	"github.com/PGITAb/bc-operator-api-verify/config"
)

type SingleGetBalanceResponse struct {
	Balance  int    `json:"balance"`
	Currency string `json:"currency"`
	Time     int    `json:"time"`
}

func SingleGetBalance_Success() api.TestCaseResultGroup {
	result := api.
		Post("/balance", map[string]string{
			"token":      config.GetToken(),
			"operatorID": config.GetOperatorID(),
			"appSecret":  config.GetAppSecret(),
			"playerID":   config.GetPlayerID(),
		}, "Balance: Success").
		Expect([]api.TestCase{
			IsStatusCode(200),
			ContainFields([]string{"balance", "currency", "time"}),
			DecodeToStruct(&SingleGetBalanceResponse{}),
		})
	return api.TestCaseResultGroup{
		Name:     "Balance: Success",
		Endpoint: "/balance",
		Results:  result,
	}
}

func SingleGetBalance_InvalidParameter() api.TestCaseResultGroup {
	result := api.
		Post("/balance", map[string]string{
			"token":      config.GetToken(),
			"operatorID": config.GetOperatorID(),
			"appSecret":  config.GetAppSecret(),
			// "playerID":   config.GetPlayerID(),
		}, "Balance: Invalid parameter").
		Expect([]api.TestCase{
			IsStatusCode(400),
			ContainFields([]string{"error"}),
		})
	return api.TestCaseResultGroup{
		Name:     "Balance: Invalid parameter",
		Endpoint: "/balance",
		Results:  result,
	}
}

func SingleGetBalance_InvalidToken() api.TestCaseResultGroup {
	result := api.
		Post("/balance", map[string]string{
			"token":      "xxxxxxxxxxxxxxxxxxxx",
			"operatorID": config.GetOperatorID(),
			"appSecret":  config.GetAppSecret(),
			"playerID":   config.GetPlayerID(),
		}, "Balance: Invalid token").
		Expect([]api.TestCase{
			IsStatusCode(404),
			ContainFields([]string{"error"}),
		})
	return api.TestCaseResultGroup{
		Name:     "Balance: Invalid token",
		Endpoint: "/balance",
		Results:  result,
	}
}

func SingleGetBalance_InvalidSecret() api.TestCaseResultGroup {
	result := api.
		Post("/balance", map[string]string{
			"token":      config.GetToken(),
			"operatorID": config.GetOperatorID(),
			"appSecret":  "xxxxxxxxxxxxxxxxxxxx",
			"playerID":   config.GetPlayerID(),
		}, "Balance: Invalid secret").
		Expect([]api.TestCase{
			IsStatusCode(401),
			ContainFields([]string{"error"}),
		})
	return api.TestCaseResultGroup{
		Name:     "Balance: Invalid secret",
		Endpoint: "/balance",
		Results:  result,
	}
}
