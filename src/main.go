package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type FriendType struct {
	Uin    string `json:"uin"`
	Remark string `json:"remark"`
	Name   string `json:"name"`
}

type FriendTaotaoType struct {
	friend     FriendType
	taotaoList []interface{}
}

// use to get friends
var (
	cookie     = ""
	referer    = ""
	qzonetoken = ""
	g_tk       = ""
	myUin      = ""
)

// use to get taotao
var (
	cookie2  = ""
	referer2 = ""
	g_tk2    = ""
)

var logger *log.Logger

func init() {
	// log output file
	file, err := os.OpenFile("qqtaotao.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Fail to open error logger file:", err)
	}
	// custom log format
	logger = log.New(io.MultiWriter(file, os.Stderr), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	// 1.get friend list
	var friends = getAllFriends(myUin, cookie, referer)

	var ch = make(chan FriendTaotaoType, 10)

	// 2.produce
	for _, friend := range friends {
		go func(f FriendType) {
			taotaoList := process(f)
			var ft FriendTaotaoType
			if len(taotaoList) != 0 {
				ft = FriendTaotaoType{
					friend:     f,
					taotaoList: taotaoList,
				}
			} else {
				ft = FriendTaotaoType{
					friend: f,
				}
			}
			ch <- ft
		}(friend)
	}

	// 3.consume
	var outgoingFriends []FriendType
	var friendTaotaos []FriendTaotaoType
	count := 0
	for ft := range ch {
		friendTaotaos = append(friendTaotaos, ft)
		if len(ft.taotaoList) > 1 {
			outgoingFriends = append(outgoingFriends, ft.friend)
		}
		count++
		if count >= len(friends) {
			break
		}
	}

	// 4.print
	bytes, _ := json.Marshal(outgoingFriends)
	logger.Println("当前有可见说说的好友共：" + strconv.Itoa(len(outgoingFriends)) + "位")
	logger.Println("当前有可见说说的好友：" + string(bytes))
	logger.Printf("所有好友的说说信息：%v", friendTaotaos)
}

func process(f FriendType) []interface{} {
	content := doRequest(f.Uin, cookie2, referer2)
	msglist := content["msglist"]
	if msglist == nil {
		return nil
	}
	list := msglist.([]interface{})
	//fmt.Println(f.remark)
	//for i := range list {
	//	oneTaotao := list[i].(map[string]interface{})
	//	fmt.Println(oneTaotao["content"].(string))
	//	if i > 5 {
	//		break
	//	}
	//}
	//fmt.Println("\n\n\n\n")
	return list
}

func getAllFriends(uin string, cookie string, referer string) []FriendType {
	var client = &http.Client{}
	request, _ := http.NewRequest("GET", "https://user.qzone.qq.com/proxy/domain/r.qzone.qq.com/cgi-bin/tfriend/friend_show_qqfriends.cgi?uin="+
		uin+
		"&follow_flag=0&groupface_flag=0&fupdate=1&g_tk="+g_tk+"&qzonetoken="+qzonetoken+"&g_tk="+g_tk, nil)
	request.Header.Add("Cookie", cookie)
	request.Header.Add("referer", referer)

	response, _ := client.Do(request)
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	s := string(body)
	s = strings.ReplaceAll(s, "_Callback(", "")
	s = strings.ReplaceAll(s, ");", "")
	m := make(map[string]interface{})
	_ = json.Unmarshal([]byte(s), &m)
	m = m["data"].(map[string]interface{})
	items := m["items"].([]interface{})
	var result []FriendType
	for _, item := range items {
		element := item.(map[string]interface{})
		f := &FriendType{
			Uin:    strconv.Itoa(int(element["uin"].(float64))),
			Remark: element["remark"].(string),
			Name:   element["name"].(string),
		}
		result = append(result, *f)
	}
	return result
}

func doRequest(uin string, cookie string, referer string) map[string]interface{} {
	var client = &http.Client{}
	request, _ := http.NewRequest("GET", "https://user.qzone.qq.com/proxy/domain/taotao.qq.com/cgi-bin/emotion_cgi_msglist_v6?uin="+
		uin+
		"&ftype=0&sort=0&pos=0&num=20&replynum=100&g_tk="+g_tk2+
		"&callback=_preloadCallback&code_version=1&format=jsonp&need_private_comment=1&g_tk="+g_tk2, nil)
	request.Header.Add("Cookie", cookie)
	request.Header.Add("referer", referer)

	response, _ := client.Do(request)
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	s := string(body)
	s = strings.ReplaceAll(s, "_preloadCallback(", "")
	s = strings.ReplaceAll(s, ");", "")
	m := make(map[string]interface{})
	_ = json.Unmarshal([]byte(s), &m)
	return m
}
