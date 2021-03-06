package handlers

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"src/mock.server/common"
	myutils "src/tools.app/utils"

	"github.com/golib/httprouter"
)

const (
	uriName              = "uri"
	queryFilePathPattern = "%s/%s_query.txt"
	bodyFilePathPattern  = "%s/%s_body.txt"
)

var (
	dataDirPath = filepath.Join(myutils.GetCurPath(), "data")
)

// MockAPIRegisterHandler register a uri with params and template body.
// Post /mock/register/:uri
func MockAPIRegisterHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// TODO: use db instead of text files to save uri:params:template_body.
	if err := myutils.MakeDir(dataDirPath); err != nil {
		common.ErrHandler(w, err)
		return
	}

	uri := params.ByName(uriName)
	filePath := fmt.Sprintf(queryFilePathPattern, dataDirPath, uri)
	if err := myutils.WriteContentToFile(filePath, r.URL.RawQuery, true); err != nil {
		common.ErrHandler(w, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		common.ErrHandler(w, err)
		return
	}
	defer r.Body.Close()

	filePath = fmt.Sprintf(bodyFilePathPattern, dataDirPath, uri)
	if err := myutils.WriteContentToFile(filePath, string(body), true); err != nil {
		common.ErrHandler(w, err)
		return
	}

	respJSON := CmdRespJSON{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("register uri success: %s", uri),
		Results: string(body),
	}
	common.WriteOKJSONResp(w, respJSON)
}

// MockAPIHandler sends templated json response by register params and body.
// Post /mock/:uri
func MockAPIHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set(common.TextContentType, common.ContentTypeJSON)
	if err := common.MockReturnCode(r, w); err != nil {
		common.ErrHandler(w, err)
		return
	}

	uri := params.ByName(uriName)
	filePath := fmt.Sprintf(bodyFilePathPattern, dataDirPath, uri)
	body, err := myutils.ReadFileContentBuf(filePath)
	if err != nil {
		common.ErrHandler(w, err)
		return
	}

	filePath = fmt.Sprintf(queryFilePathPattern, dataDirPath, uri)
	query, err := myutils.ReadFileContent(filePath)
	if err != nil {
		common.ErrHandler(w, err)
		return
	}

	// 优先级：当前请求的参数 覆盖 注册参数
	queryMap := common.QueryToMap(query)
	for k, v := range r.URL.Query() {
		queryMap[k] = v
	}
	if len(queryMap) == 0 {
		if _, err := w.Write([]byte(body)); err != nil {
			common.ErrHandler(w, err)
		}
		return
	}

	// template 处理
	tmplParams, err := common.ParseParamsForTempl(queryMap)
	if err != nil {
		common.ErrHandler(w, err)
		return
	}

	tmpl, err := template.New("mockapi").Parse(string(body))
	if err != nil {
		common.ErrHandler(w, err)
		return
	}
	if err := tmpl.Execute(w, tmplParams); err != nil {
		log.Fatalln(err)
	}
}
