package udeskauth

import (
	"strings"
)

//获取token字符串

func (UdeskToken) getTokenString() string {
	token := GetUdeskAuthToken().Open_api_auth_token
	return token
}

//计算sign并返回authobj

func (UdeskToken) getAuthobj() Authobj {
	token := u.getTokenString()
	timestamp := GetTimeStamp()
	nonce := GetNonce()
	var builder strings.Builder
	builder.WriteString(Email + "&")
	builder.WriteString(token + "&")
	builder.WriteString(timestamp + "&")
	builder.WriteString(nonce + "&")
	builder.WriteString(Sign_version)
	str2sha256 := builder.String()
	hashResult := calculateSHA256(str2sha256)
	//fmt.Println(str2sha256)
	return Authobj{
		Email:        Email,
		Timestamp:    timestamp,
		Sign:         hashResult,
		Nonce:        nonce,
		Sign_version: Sign_version,
	}

}
