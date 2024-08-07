package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	// Setup viper
	viper.SetConfigName("config")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("bc")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}
}

func GetWalletType() string {
	return viper.GetString("wallettype")
}

func GetDomain() string {
	return viper.GetString("server")
}

func GetToken() string {
	return viper.GetString("token")
}

func GetOperatorID() string {
	return viper.GetString("operatorid")
}

func GetAppSecret() string {
	return viper.GetString("secret")
}

func GetPlayerID() string {
	return viper.GetString("playerid")
}

func GetNickname() string {
	return viper.GetString("nickname")
}

func GetAmount() int {
	return viper.GetInt("amount")
}

func GetGameID() string {
	return viper.GetString("gameid")
}

func GetBetType() string {
	return viper.GetString("bettype")
}

func GetGameStatusWin() string {
	return viper.GetString("gamestatuswin")
}

func GetGameStatusLoss() string {
	return viper.GetString("gamestatusloss")
}

func GetGameResult() string {
	return viper.GetString("gameresult")
}

func GetAmount_2() int {
	return viper.GetInt("amount_2")
}

func GetAmount_3() int {
	return viper.GetInt("amount_3")
}

func GetCurrency() string {
	return viper.GetString("currency")
}

func GetTableID() string {
	return viper.GetString("tableid")
}

func GetPrefix() string {
	return viper.GetString("prefix")
}

func GetColorMode() bool {
	return viper.GetBool("colormode")
}

func GetFileOutput() bool {
	return viper.GetBool("fileoutput")
}

func GetIP() string {
	return viper.GetString("ip")
}

func GetOdd() string {
	return viper.GetString("odds")
}

func SetApiServer(url string) {
	viper.Set("server", url)
}

func SetToken(token string) {
	viper.Set("token", token)
}

func SetOperatorID(operatorID string) {
	viper.Set("operatorid", operatorID)
}

func SetAppSecret(appSecret string) {
	viper.Set("secret", appSecret)
}

func SetPlayerID(playerID string) {
	viper.Set("playerid", playerID)
}
