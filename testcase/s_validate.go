package testcase

import (
	"github.com/PGITAb/bc-operator-api-verify/api"
	"github.com/PGITAb/bc-operator-api-verify/config"
)

type SingleValidateResponse struct {
	PlayerID string `json:"playerID"`
	Nickname string `json:"nickname"`
	Currency string `json:"currency"`
	Time     int    `json:"time"`
}

func SingleValidate_Success() api.TestCaseResultGroup {
	result := api.
		Post("/validate", map[string]string{
			"token":      config.GetToken(),
			"operatorID": config.GetOperatorID(),
			"appSecret":  config.GetAppSecret(),
		}, "Validate: Success").
		Expect([]api.TestCase{
			IsStatusCode(200),
			ContainFields([]string{"playerID", "nickname", "currency", "time"}),
			DecodeToStruct(&SingleValidateResponse{}),
		})
	return api.TestCaseResultGroup{
		Name:     "Validate: Success",
		Endpoint: "/validate",
		Results:  result,
	}
}

func SingleValidate_InvalidParameter() api.TestCaseResultGroup {
	result := api.
		Post("/validate", map[string]string{
			"token":      config.GetToken(),
			"operatorID": "",
			"appSecret":  config.GetAppSecret(),
		}, "Validate: Invalid parameter").
		Expect([]api.TestCase{
			IsStatusCode(400),
			ContainFields([]string{"error"}),
		})
	return api.TestCaseResultGroup{
		Name:     "Validate: Invalid parameter",
		Endpoint: "/validate",
		Results:  result,
	}
}

func SingleValidate_InvalidToken() api.TestCaseResultGroup {
	result := api.
		Post("/validate", map[string]string{
			"token":      "xxxxxxxxxxxxxxxxxxxx",
			"operatorID": config.GetOperatorID(),
			"appSecret":  config.GetAppSecret(),
		}, "Validate: Invalid token").
		Expect([]api.TestCase{
			IsStatusCode(404),
			ContainFields([]string{"error"}),
		})
	return api.TestCaseResultGroup{
		Name:     "Validate: Invalid token",
		Endpoint: "/validate",
		Results:  result,
	}
}

func SingleValidate_InvalidSecret() api.TestCaseResultGroup {
	result := api.
		Post("/validate", map[string]string{
			"token":      config.GetToken(),
			"operatorID": config.GetOperatorID(),
			"appSecret":  "xxxxxxxxxxxxxxxxxxxx",
		}, "Validate: Invalid secret").
		Expect([]api.TestCase{
			IsStatusCode(401),
			ContainFields([]string{"error"}),
		})
	return api.TestCaseResultGroup{
		Name:     "Validate: Invalid secret",
		Endpoint: "/validate",
		Results:  result,
	}
}
