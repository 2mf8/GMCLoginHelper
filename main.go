package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/2mf8/GMCLoginHelper/dto"

	"strconv"

	"google.golang.org/protobuf/proto"
)

type Logins struct {
	Logins   []*dto.CreateBotReq
	GMCPath  string
	BindPort int
}

func init() {
	fmt.Println("初始化成功,该程序用于非首次登录")
}

func main() {
	lv, err := LoginJsonRead()
	if err != nil && StartsWith(err.Error(), "open login.json") {
		LoginJsonCreate()
		fmt.Println("已为您创建默认登录文件，请修改 login.json 文件后重启该程序")
	}

	port := strconv.Itoa(lv.BindPort)
	url := fmt.Sprintf("http://localhost:%v/bot/create/v1/", lv.BindPort)

	cmd := exec.Command(lv.GMCPath, "--port", port)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	time.Sleep(time.Second * 2)
	for _, i := range lv.Logins {
		v, _ := proto.Marshal(i)
		resp, err := http.Post(url, "application/x-protobuf", bytes.NewBuffer(v))
		if err != nil {
			os.Stdout.WriteString(fmt.Sprintf("[错误] %s", err.Error()))
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		respString := "[返回] " + string(body) + "\n"
		os.Stdout.WriteString(fmt.Sprintf("[状态] %s \n[方法] %s\n", resp.Status, resp.Request.Method))
		os.Stdout.WriteString(respString)
	}
}

func LoginJsonCreate() {
	login := &dto.CreateBotReq{
		BotId:          1234567890,
		Password:       "123456ab",
		DeviceSeed:     1234567890,
		SignServer:     "http://kequ5060.cn:8080",
		SignServerAuth: "114514",
		ClientProtocol: 1,
	}

	var logins Logins
	logins.Logins = append(logins.Logins, login)
	logins.GMCPath = ".\\gmc.exe"
	logins.BindPort = 9000

	output, err := json.MarshalIndent(&logins, "", "\t")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	err = os.WriteFile("login.json", output, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file", err)
	}
}

func LoginJsonRead() (json_data Logins, err error) {
	jsonFile, err := os.Open("login.json")
	if err != nil {
		fmt.Println("Error reading JSON File:", err)
		return
	}
	defer jsonFile.Close()
	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON data:", err)
		return
	}
	json.Unmarshal(jsonData, &json_data)
	return
}

func StartsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
