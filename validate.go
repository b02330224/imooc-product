package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"imooc-product/common"
	"imooc-product/datamodels"
	"imooc-product/encrypt"
	"imooc-product/rabbitmq"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var hostArray = []string{"127.0.0.1", "127.0.0.1"}
var localHost = ""
var port = "8081"
var hashConsistent *common.Consistent

//数量控制接口服务器内网IP，或者getOne的SLB内网IP
var GetOneIp  = "127.0.0.1"
var GetOnePort = "12345"

//rabbitmq
var rabbitMqValidate *rabbitmq.RabbitMQ

var interval = 20

//用来存放控制信息
type AccessControl struct {
	//用来存放用户想要存放的信息
	sourceArray map[int]time.Time
	sync.RWMutex
}

type BlackList struct {
	listArray map[int]bool
	sync.RWMutex
}

var blackList = &BlackList{listArray:make(map[int]bool)}
func (m *BlackList) GetBlackById(uid int) bool {
	m.RLock()
	defer m.RLock()

	return m.listArray[uid]
}

func (m *BlackList) SetBlackListById(uid int) bool {
	m.Lock()
	defer m.Unlock()
	m.listArray[uid] = true
	return true
}


var accessControl = &AccessControl{sourceArray:make(map[int]time.Time)}

func (m *AccessControl) GetNewRecord(uid int) time.Time {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()

	return m.sourceArray[uid]
}

func (m *AccessControl) SetNewRecord(uid int) {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	m.sourceArray[uid] = time.Now()
}

func (m *AccessControl) GetDistributedRight(req *http.Request) bool {
	uid , err := req.Cookie("uid")
	if err != nil {
		return false
	}

	hostRequest , err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}

	if hostRequest == localHost {
		//执行本机数据读取和校验
		return m.GetDataFromMap(uid.Value)
	} else {
		//不是本机充当代理访问数据返回结果
		return GetDataFromOtherMap(hostRequest, req)
	}
}

//获取本机map，并且处理业务逻辑，返回的结果类型为bool类型
func (m *AccessControl) GetDataFromMap(uid string) (isOk bool) {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return
	}

	if blackList.GetBlackById(uidInt) {
		return false
	}

	dataRecord := m.GetNewRecord(uidInt)
	if !dataRecord.IsZero() {
		if dataRecord.Add(time.Duration(interval) * time.Second).After(time.Now()) {
			return false
		}
	}

	m.SetNewRecord(uidInt)

	return true
}


func GetDataFromOtherMap(host string, request *http.Request) bool {
	hostUrl := "http://" + host + ":" + port + "/checkRight"
	response, body, err := GetCurl(hostUrl, request)
	if err != nil {
		return false
	}

	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false
}

//统一验证拦截器，每个接口都需要提前验证
func Auth(w http.ResponseWriter, r *http.Request) error {

	err := CheckUserInfo(r)
	if err != nil {
		return errors.New("验证失败")
	}

	return nil
}

func CheckUserInfo (r *http.Request) error {
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		return errors.New("用户uid cookie 获取失败！")
	}

	//获取用户加密串
	signCookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("用户加密串 cookie 获取失败！")
	}

	signStr, err := url.QueryUnescape(signCookie.Value)
	if err != nil {
		return errors.New("url decode 失败")
	}
	signByte, err := encrypt.DePwdCode(signStr)
	if err != nil {
		return errors.New("加密串被篡改")
	}
	if checkInfo(uidCookie.Value , string(signByte)) {
		return nil
	}
	return errors.New("身份校验失败")
}

func checkInfo(checkStr string, signStr string) bool {
	if checkStr == signStr {
		return true
	}
	return false
}

func GetCurl(hostUrl string, request *http.Request) (response *http.Response, body []byte, err error) {
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return
	}

	uidSign, err := request.Cookie("sign")
	if err != nil {
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", hostUrl, nil)
	if err != nil {
		return
	}

	cookieUid := &http.Cookie{Name:"uid", Value:uidPre.Value, Path:"/"}
	cookieSign := &http.Cookie{Name:"sign", Value:uidSign.Value, Path:"/"}
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)

	response, err = client.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	return
}

func CheckRight(w http.ResponseWriter, r *http.Request) {
	right := accessControl.GetDistributedRight(r)
	if right == false {
		w.Write([]byte("false"))
	}

	w.Write([]byte("true"))
	return
}

func Check(w http.ResponseWriter, r *http.Request) {
	queryForm, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil || len(queryForm["productId"]) <= 0 {
		w.Write([]byte("false"))
		return
	}

	productString := queryForm["productId"][0]

	userCookie, err := r.Cookie("uid")
	if err != nil {
		w.Write([]byte("false"))
		return
	}

	right := accessControl.GetDistributedRight(r)
	if right == false {
		w.Write([]byte("false"))
		return
	}

	hostUrl := "http://" + GetOneIp + ":" + GetOnePort + "/getOne"
	responseValidate, validateBody, err := GetCurl(hostUrl, r)
	if err != nil {
		fmt.Println("test", err)
		w.Write([]byte("false"))
		return
	}

	if responseValidate.StatusCode == 200 {
		if string(validateBody) == "true" {
			productId , err := strconv.ParseInt(productString, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			userId, err := strconv.ParseInt(userCookie.Value, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			message := datamodels.NewMessage(userId, productId)
			byteMessage, err := json.Marshal(message)
			if err != nil {
				w.Write([]byte("false"))
				return
			}


			err = rabbitMqValidate.PublishSimple(string(byteMessage))
			if err != nil {
				w.Write([]byte("false"))
				return
			}

			w.Write([]byte("true"))
			return
		}
	}

	w.Write([]byte("false"))
	return
}

func main() {
	//负载均衡器设置
	//采用一致性哈希算法
	hashConsistent = common.NewConsistent()
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}

	localIp, err := common.GetIntranceIp()
	if err != nil {
		fmt.Println(err)
	}
	localHost = localIp

	rabbitMqValidate = rabbitmq.NewRabbitMQSimple("imoocProduct")
	defer rabbitMqValidate.Destory()

	//设置静态文件目录
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("./fronted/web/htmlProductShow"))))
	//设置静态资源目录
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./fronted/web/public"))))

	//过滤器
	filter := common.NewFilter()
	filter.RegisterFilterUri("/check", Auth)
	filter.RegisterFilterUri("/checkRight", Auth)

	http.HandleFunc("/check", filter.Handle(Check))
	http.HandleFunc("/checkRight", filter.Handle(CheckRight))

	//test
	http.HandleFunc("/testPublishSimple", testPublishSimple)


	http.ListenAndServe(":8081", nil)
}


//测试发送导队列的性能 qps在550左右
func testPublishSimple(w http.ResponseWriter, r *http.Request){

	message := datamodels.NewMessage(1, 1)
	byteMessage, err := json.Marshal(message)
	if err != nil {
		w.Write([]byte("false"))
		return
	}


	err = rabbitMqValidate.PublishSimple(string(byteMessage))
	if err != nil {
		w.Write([]byte("false"))
		return
	}

	w.Write([]byte("true"))
	return
}