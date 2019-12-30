package handle

import (
	"fmt"
	"his/db"
	"his/util"
	"io/ioutil"
	"net/http"
	"time"
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
//SignInHandler 登录接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	pwd := r.Form.Get("password")
	encPasswd := util.Sha1([]byte(pwd + pwd_salt))
	//1.校验用户名及密码
	pwdChecked := db.UserSignUp(username, encPasswd)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
		return
	}
	//2.生成访问凭证(token)
	token := GenToken(username)
	upRes := db.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("FAILED"))
		return
	}
	//3.登录成功后重定向到首页
	w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
}

func GenToken(username string) string {
	//md5(username + timestamp + token_salt) + timestamp[:8]
	ts := fmt.Sprintf("%x",time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return  tokenPrefix + ts[:8]
}
