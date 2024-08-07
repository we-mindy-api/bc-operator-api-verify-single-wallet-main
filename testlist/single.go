package testlist

import (
	"github.com/PGITAb/bc-operator-api-verify/api"
	"github.com/PGITAb/bc-operator-api-verify/testcase"
)

func SingleWallet(collect api.ResultCollector) {
	collect.Func(testcase.SingleValidate_Success())
	collect.Func(testcase.SingleValidate_InvalidParameter())
	collect.Func(testcase.SingleValidate_InvalidToken())
	collect.Func(testcase.SingleValidate_InvalidSecret())

	collect.Func(testcase.SingleGetBalance_Success())
	collect.Func(testcase.SingleGetBalance_InvalidParameter())
	collect.Func(testcase.SingleGetBalance_InvalidToken())
	collect.Func(testcase.SingleGetBalance_InvalidSecret())

	collect.Func(testcase.SingleDebit_Success())
	collect.Func(testcase.SingleDebit_InvalidParameter())
	collect.Func(testcase.SingleDebit_InvalidToken())
	collect.Func(testcase.SingleDebit_InvalidSecret())
	collect.Func(testcase.SingleDebit_BalanceNotEnough())
	collect.Func(testcase.SingleDebit_DuplicateTransaction())

	collect.Func(testcase.SingleCredit_Success())
	collect.Func(testcase.SingleCredit_InvalidParameter())
	collect.Func(testcase.SingleCredit_InvalidSecret())
	collect.Func(testcase.SingleCredit_DebitNotFound())
	collect.Func(testcase.SingleCredit_DuplicateTransaction())

	collect.Func(testcase.SingleRollbackDebited_Success())
	collect.Func(testcase.SingleRollbackCreditedPlayer_Success())
	collect.Func(testcase.SingleRollback_InvalidParameter())
	collect.Func(testcase.SingleRollback_InvalidSecret())
	collect.Func(testcase.SingleRollback_DebitNotFound())
	collect.Func(testcase.SingleRollback_DuplicateTransaction())
}
