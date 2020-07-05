package im

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	//发送消息使用导的url
	sendurl = `https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=`
	//获取token使用导的url
	get_token = `https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=`
)

var requestError = errors.New("request error,check url or network")

type access_token struct {
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
}

//定义一个简单的文本消息格式
type send_msg struct {
	Touser  string            `json:"touser"`
	Toparty string            `json:"toparty"`
	Totag   string            `json:"totag"`
	Msgtype string            `json:"msgtype"`
	Agentid int               `json:"agentid"`
	Text    map[string]string `json:"text"`
	Safe    int               `json:"safe"`
}

type send_msg_error struct {
	Errcode int    `json:"errcode`
	Errmsg  string `json:"errmsg"`
}

func main() {
	touser := flag.String("t", "GaoXing", "-t user 直接接收消息的用户昵称")
	// 这里agentid填写自己的
	agentid := flag.Int("i", 111, "-i 0 指定agentid")
	content := flag.String("c", "微信报警", "-c 'Hello world' 指定要发送的内容")
	corpid := flag.String("p", "------（填写自己的）------", "-p corpid 必须指定")
	corpsecret := flag.String("s", "------(填写自己的)-------", "-s corpsecret 必须指定")
	flag.Parse()

	if *corpid == "" || *corpsecret == "" {
		flag.Usage()
		return
	}

	var m send_msg = send_msg{Touser: *touser, Toparty: "C_G_X", Msgtype: "text", Agentid: *agentid, Text: map[string]string{"content": *content}}

	token, err := Get_token(*corpid, *corpsecret)
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Println("获取token成功：", token)
	buf, err := json.Marshal(m)
	if err != nil {
		return
	}
	err = Send_msg(token.Access_token, buf)
	if err != nil {
		println(err.Error())
	} else {
		fmt.Println("发送消息成功", string(buf))
	}

}

//发送消息.msgbody 必须是 API支持的类型
func Send_msg(Access_token string, msgbody []byte) error {
	body := bytes.NewBuffer(msgbody)
	resp, err := http.Post(sendurl+Access_token, "application/json", body)
	if resp.StatusCode != 200 {
		return requestError
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var e send_msg_error
	err = json.Unmarshal(buf, &e)
	if err != nil {
		return err
	}
	if e.Errcode != 0 && e.Errmsg != "ok" {
		return errors.New(string(buf))
	}
	return nil
}

//通过corpid 和 corpsecret 获取token
func Get_token(corpid, corpsecret string) (at access_token, err error) {
	resp, err := http.Get(get_token + corpid + "&corpsecret=" + corpsecret)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = requestError
		return
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(buf, &at)
	if at.Access_token == "" {
		err = errors.New("corpid or corpsecret error.")
	}
	return
}
