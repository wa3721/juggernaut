package auth

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strconv"
	"time"
)

//新建请求体对象

func newReqUdeskToken() RequestUdeskBody {
	return RequestUdeskBody{
		Email:    Email,
		Password: Password,
	}
}

// 新建返回体对象
func newRespUdeskBody() UdeskToken {
	return UdeskToken{
		Code:                "",
		Open_api_auth_token: "",
	}
}

//获取Unix时间戳

func GetTimeStamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

//获取nonce

func GetNonce() string {
	randomUUID := uuid.New()
	nonce := randomUUID.String()
	return nonce
}

// sha256转换函数
func calculateSHA256(input string) string {
	// 将字符串转换为字节数组
	inputBytes := []byte(input)

	// 创建SHA-256哈希对象
	hasher := sha256.New()

	// 将字节数组写入哈希对象
	hasher.Write(inputBytes)

	// 计算哈希值并返回
	hashInBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	return hashString
}

// 获取鉴权token对象

func GetUdeskAuthToken() UdeskToken {
	reqBody := newReqUdeskToken()
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("获取udesk管理员token过程中转换请求体json失败，错误是：%v", err)
	}
	payload := bytes.NewBufferString(string(jsonData))
	resp, err := http.Post(Auth_token_url, "application/json", payload)
	if err != nil {
		fmt.Printf("获取udesk管理员token过程中请求获取token接口失败，错误是：%v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("获取udesk管理员token过程中获取token接口响应失败，错误是：%v", err)
	}
	json.Unmarshal(body, &u)
	return u
}
