package rocket

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"fmt"
	"strings"

	"github.com/pquerna/ffjson/ffjson"
)

/* 公共调用的Handel */
type MsgHandel struct {
	Httpclient   *http.Client
	RocketUname  string
	RocketUpass  string
	RocketUID    string
	RocketUtoken string
	PData        PushData
	PDomain      string
	PMessage     string
}

/* 包初始化 */
func init() {
	fmt.Println("rocket packge import ok")
}

/* 提取git 用户的 push 行为 */
func (h *MsgHandel) GetPushData(m map[string]interface{}) {

	f := m["commits"].([]interface{})
	commits := f[0].(map[string]interface{})
	gituser := commits["author"].(map[string]interface{})

	repository := m["repository"].(map[string]interface{})

	pushData := PushData{
		Repository:    repository["description"].(string),
		RepositoryUrl: repository["url"].(string),
		Message:       commits["message"].(string),
		CommitUrl:     commits["url"].(string),
		UserName:      gituser["name"].(string),
	}

	h.PData = pushData

}

/* 登陆rocket.chat 获取token */
func (h *MsgHandel) getRocketUserToken() {
	postValues := url.Values{}
	postValues.Add("user", h.RocketUname)
	postValues.Add("password", h.RocketUpass)

	furl := h.PDomain + "/api/login"
	resp, _ := h.Httpclient.PostForm(furl, postValues)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var f interface{}
	ffjson.Unmarshal(body, &f)
	m := f.(map[string]interface{})
	//status := m["status"].(string)
	data := m["data"].(map[string]interface{})
	authToken := data["authToken"].(string)
	userID := data["userId"].(string)

	h.RocketUID = userID
	h.RocketUtoken = authToken
}

/* 给Team IM 推送消息 */
func (h *MsgHandel) PushRocketChat() {

	// 取用户token和id
	h.getRocketUserToken()

	postValues := url.Values{}
	postValues.Add("msg", h.PMessage)

	furl := h.PDomain + "/api/rooms/GENERAL/send"
	reqest, err := http.NewRequest("POST", furl, strings.NewReader(postValues.Encode()))

	if err != nil {
		log.Println("Fatal error ", err.Error())
	}
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	reqest.Header.Set("Accept", "application/json")
	reqest.Header.Add("X-Auth-Token", h.RocketUtoken)
	reqest.Header.Add("X-User-Id", h.RocketUID)

	response, err := h.Httpclient.Do(reqest)
	defer response.Body.Close()

	if err != nil {
		log.Println("Fatal error ", err.Error())
	}
}
