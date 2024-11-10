package auth

import "strings"

//最终请求的url后缀拼接

func Geturlstring(url string) string {
	var authobj = u.getAuthobj() //必须要同步处理timestamp和sign的关系，否则会鉴权失败，所以从同一个结构体中取值
	var builder strings.Builder
	builder.WriteString(url)
	builder.WriteString("email=" + authobj.Email)
	builder.WriteString("&timestamp=" + authobj.Timestamp)
	builder.WriteString("&sign=" + authobj.Sign)
	builder.WriteString("&nonce=" + authobj.Nonce)
	builder.WriteString("&sign_version=" + authobj.Sign_version)
	return builder.String()
}
