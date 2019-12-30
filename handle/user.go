package handle

import (
	"his/db"
	"his/util"
	"io/ioutil"
	"net/http"
)

const (
	pwd_salt = "*#890"
)
//处理用户注册请求
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("invalid parameter"))
		return
	}

	enc_passwd := util.Sha1([]byte(passwd + pwd_salt))
	suc := db.UserSignUp(username, enc_passwd)
	if suc {
		w.Write([]byte("SUCCESS"))
	}else {
		w.Write([]byte("FAIL"))
	}
}
