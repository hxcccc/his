package handle

import (
	"fmt"
	"his/db"
	"his/util"
	"io/ioutil"
	"net/http"
	"strconv"
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
	pwdChecked := db.UserSignIn(username, encPasswd)
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
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	//http.Redirect(w, r, "/user/home", http.StatusFound)
	resp := util.RespMsg{
		Code:0,
		Msg:"OK",
		Data: struct {
			Location string
			Username string
			Token string
		}{
			Location:"http://" + r.Host + "/static/view/home.html",
			Username:username,
			Token:token,
		},
	}
	w.Write(resp.JSONBytes())
}
//UserInfoHandler 查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	//解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")
	//验证token
	isValidToken := IsTokenValid(username, token)
	if !isValidToken {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//查询用户信息
	user, err := db.GetUserInfo(username)
	if err != nil{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//组装并响应用户数据
	resp := util.RespMsg{
		Code:0,
		Msg:"OK",
		Data:user,
	}
	w.Write(resp.JSONBytes())
}

func GenToken(username string) string {
	//md5(username + timestamp + token_salt) + timestamp[:8]
	ts := fmt.Sprintf("%x",time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return  tokenPrefix + ts[:8]
}

func IsTokenValid(username string, token string) bool {
	if len(token) != 40 {
		return false
	}
	//判断token的时效性是否过期
	tokenTs, err := strconv.Atoi(string(token[len(token)-8:]))
	if err != nil {
		fmt.Println("ts:string to int failed")
		return false
	}
	nowTs, err := strconv.Atoi(fmt.Sprintf("%x", time.Now().Unix())[:8])
	if err != nil {
		fmt.Println("ts:string to int failed")
		return false
	}
	if keepTs := nowTs - tokenTs;keepTs > 86400 {
		return false
	}
	//从数据库表中查询username对应的token信息//对比两个token是否一致
	res := db.VerifyToken(username, token)
	if !res {
		return false
	}
	return true
}
