package handlers

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"src/mock.server/common"

	"github.com/golib/httprouter"
	redis "gopkg.in/redis.v5"
)

// MockDemoHandler router for mock demo handlers.
func MockDemoHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	if r.Method == "GET" {
		switch id {
		case "1":
			mockDemo01(w, r)
		case "2":
			mockDemo02(w, r)
		case "5":
			mockDemo05(w, r)
		default:
			common.ErrHandler(w, fmt.Errorf("GET for invalid path: %s", r.URL.Path))
			return
		}
	} else if r.Method == "POST" {
		switch id {
		case "2-1":
			mockDemo0201(w, r)
		case "3":
			mockDemo03(w, r)
		case "4":
			mockDemo04(w, r)
		default:
			common.ErrHandler(w, fmt.Errorf("POST for invalid path: %s", r.URL.Path))
		}
	} else {
		common.ErrHandler(w, fmt.Errorf("Method not support: %s", r.Method))
	}
}

// demo, parse get request.
// Get /demo/1?userid=xxx&username=xxx
func mockDemo01(w http.ResponseWriter, r *http.Request) {
	var userID, userName string

	r.ParseForm()
	log.Println("Request Method:", r.Method)

	log.Println("Form Data:")
	if val, ok := r.Form["userid"]; ok {
		userID = val[0]
	} else {
		userID = "null"
	}
	log.Println("userid:", userID)

	if val, ok := r.Form["username"]; ok {
		userName = val[0]
	} else {
		userName = "null"
	}
	log.Println("username:", userName)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hi, thanks for access %s", html.EscapeString(r.URL.Path[1:]))
}

// demo, parse get request.
// Get /demo/2?userid=xxx&username=xxx&key=val1&key=val2
func mockDemo02(w http.ResponseWriter, r *http.Request) {
	var userID, userName string

	log.Println("Request Query:")
	values := r.URL.Query()
	if val, ok := values["userid"]; ok {
		userID = val[0]
	} else {
		userID = "nil"
	}
	log.Println("userid:", userID)

	if val, ok := values["username"]; ok {
		userName = val[0]
	} else {
		userName = "nil"
	}
	log.Println("username:", userName)

	if val, ok := values["key"]; ok {
		for _, v := range val {
			log.Println("key:", v)
		}
	}

	b := []byte(fmt.Sprintf("hi, thanks for access %s", html.EscapeString(r.URL.Path[1:])))
	common.WriteOKHTMLResp(w, b)
}

// demo, parse input data with keyword, and return templated json body
// Post /demo/2-1?userid=xxx&username=xxx&age=19&key1=val1&key2=val2"
// Post /demo/2-2?userid=xxx&username=xxx&age=randint(27)&key1=val1&key2=randstr(8)
func mockDemo0201(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		common.ErrHandler(w, err)
		return
	}
	defer r.Body.Close()

	tmpl, err := template.New("mock").Parse(string(body))
	if err != nil {
		common.ErrHandler(w, err)
		return
	}

	params, err := common.ParseParamsForTempl(r.URL.Query())
	if err != nil {
		common.ErrHandler(w, err)
		return
	}

	w.Header().Set(common.TextContentType, common.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	if err := tmpl.Execute(w, params); err != nil {
		log.Fatalln(err)
	}
}

// demo, parse post json body.
type server struct {
	ServerName string `json:"server_name"`
	ServerIP   string `json:"server_ip"`
}

type serverInfo struct {
	SvrList  []server `json:"server_list"`
	SvrGrpID string   `json:"server_group_id"`
}

func mockDemo03(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		common.ErrHandler(w, err)
		return
	}
	defer r.Body.Close()

	var s serverInfo
	if err := json.Unmarshal(body, &s); err != nil {
		common.ErrHandler(w, err)
		return
	}
	log.Println("Server Info:")
	log.Println("server group id:", s.SvrGrpID)
	for _, svr := range s.SvrList {
		log.Printf("server name: %s, server ip: %s\n", svr.ServerName, svr.ServerIP)
	}

	common.WriteOKJSONResp(w, s)
	// fmt.Fprintf(w, "hi, thanks for access %s", html.EscapeString(r.URL.Path[1:]))
}

// demo, parse post form with cookie
// POST /demo/4
func mockDemo04(w http.ResponseWriter, r *http.Request) {
	log.Printf("Content type: %s\n", r.Header.Get(common.TextContentType))

	log.Println("Form data:")
	log.Printf("kv: key1=%s\n", r.PostFormValue("key1"))
	log.Printf("kv: key2=%s\n", r.PostFormValue("key2"))

	log.Println("Cookie data:")
	for _, v := range []string{"user", "pwd"} {
		cookie, err := r.Cookie(v)
		if err != nil {
			common.ErrHandler(w, err)
			return
		}
		printCookie(cookie)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hi, thanks for access %s", html.EscapeString(r.URL.Path[1:]))
}

func printCookie(c *http.Cookie) {
	log.Println("name:", c.Name)
	log.Println("value:", c.Value)
	log.Println("domain:", c.Domain)
	log.Println("expires:", c.Expires)
}

// demo, get total access count from redis.
// redis env: docker run --name redis -p 6379:6379 --rm -d redis:4.0
// GET /demo/5
func mockDemo05(w http.ResponseWriter, r *http.Request) {
	if len(common.RunConfigs.Server.RedisURI) == 0 {
		common.ErrHandler(w, fmt.Errorf("config redis uri is empty"))
		return
	}

	log.Println("Redis server:", common.RunConfigs.Server.RedisURI)
	accessCount, err := getAccessCountFromRedis()
	if err != nil {
		common.ErrHandler(w, err)
		return
	}
	log.Println("*Total Access:", accessCount)
	common.WriteOKHTMLResp(w, []byte(fmt.Sprintln("Total Access: "+accessCount)))
}

func getAccessCountFromRedis() (string, error) {
	errNum := "-1"
	options := redis.Options{
		Addr:     common.RunConfigs.Server.RedisURI,
		Password: "",
	}

	client := redis.NewClient(&options)
	defer client.Close()

	pong, err := client.Ping().Result()
	if err != nil {
		return errNum, err
	}
	log.Println(pong)

	const key = "server_total_access"
	total, err := client.Get(key).Result()
	if err != nil {
		if strings.Contains(err.Error(), "nil") {
			if err := client.Set(key, 1, 0).Err(); err != nil {
				return errNum, err
			}
			total = "1"
		} else {
			return errNum, err
		}
	}

	err = client.Incr(key).Err()
	if err != nil {
		return "-1", err
	}

	return total, nil
}
