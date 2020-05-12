package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (app *CoreApplication) sendGet(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return "http.Get is wrong"
	}
	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "ioutil.ReadAll is wrong"
	}
	return string(s)
}

func (app *CoreApplication) DeliverTx_Msg(msg MsgTx) (int, string) {
	str := strings.Split(msg.Msg, "https://")
	if len(str) > 1 {
		url := "http://127.0.0.1:5000/gitcache/system/mirror/" + str[1]
		fmt.Println(url)
		go app.sendGet(url)
	}
	return 0, ""
}

func (app *CoreApplication) DeliverTx_Token(tokenObj TokenTx) (int, string) {
	_to := tokenObj.To
	_from := tokenObj.From
	_token := tokenObj.Token
	_amount := tokenObj.Amount
	if _to == "" || _to == _from {
		// create new token
		_, err := app.MongoDB_Query_CodeName(string(_token))
		if err == nil {
			return 1, "DeliverTx CodeName has existed"
		} else {
			if _, err := app.MongoDB_Add_CodeName(string(_token)); err != nil {
				return 1, "DeliverTx MongoDB_Add_CodeName failed"
			}
		}
		// add asset in assets
		assetNew := Asset{Publickey: _from, Token: _token, Amount: _amount}
		if _, err := app.MongoDB_Update_Assets(_from, _token, assetNew); err != nil {
			return 1, "DeliverTx MongoDB_Update_Assets failed"
		}
		return 0, ""
	}
	if _to != _from {
		fromPublic, err := app.MongoDB_Query_Assets(_from, _token)
		if err != nil {
			info := "you have any code of " + _token
			return 1, info
		}
		if fromPublic.Amount < _amount {
			return 1, "your amount is not enough"
		}
		fromAssets := Asset{Publickey: _from, Token: _token, Amount: fromPublic.Amount - _amount}
		toPublic, err := app.MongoDB_Query_Assets(_to, _token)
		var toAssets Asset
		if err != nil {
			toAssets = Asset{Publickey: _to, Token: _token, Amount: _amount}
		} else {
			toAssets = Asset{Publickey: _to, Token: _token, Amount: toPublic.Amount + _amount}
		}
		app.MongoDB_Update_Assets(_from, _token, fromAssets)
		app.MongoDB_Update_Assets(_to, _token, toAssets)
	}
	return 0, ""
}
