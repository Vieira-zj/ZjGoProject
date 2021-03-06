package control

import (
	"net/http"

	"src/tools.app/apps/webmvc/args"
	"src/tools.app/apps/webmvc/service"
	"src/tools.app/apps/webmvc/util"
)

// UserCtrl user控制器对象
type UserCtrl struct {
}

// Router 将url和处理函数绑定
func (ctrl *UserCtrl) Router() {
	RouterPost("/user/login", ctrl.authWithID)
	// RouterGet("/user/find",ctrl.Find)
	// RouterGet("/user/quit",ctrl.quit)
	RegExRouter("/d/.*", ctrl.regText)
}

var userService service.UserService

func (ctrl *UserCtrl) authWithID(w http.ResponseWriter, req *http.Request) {
	var reqData args.AuthArg
	util.Bind(req, &reqData)

	if u, err := userService.LoginWithOpenID(reqData.User.ID); err != nil {
		util.RespFail(w, http.StatusInternalServerError, err.Error())
	} else {
		util.RespOK(w, u)
	}
}

func (ctrl *UserCtrl) regText(w http.ResponseWriter, req *http.Request) {
	util.RespOK(w, req.RequestURI)
}
