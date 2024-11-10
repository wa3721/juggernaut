package auth

const (
	Email          = "hbqi@alauda.io"
	Sign_version   = "v2"
	Password       = "Ye_qiu@123"
	Auth_token_url = "https://servicecenter-alauda.udesk.cn/open_api_v1/log_in"
)

//获取鉴权token的账号密码请求体对象

type RequestUdeskBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//udesk鉴权接口返回token对象

type UdeskToken struct {
	Code                string `json:"code"`
	Open_api_auth_token string `json:"open_api_auth_token"`
}

// 定义认证对象
type Authobj struct {
	Email,
	Timestamp,
	Sign,
	Nonce,
	Sign_version string
}

//全局声明返回体对象

var u = newRespUdeskBody()
