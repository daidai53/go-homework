// Copyright@daidai53 2023
package web

import (
	"errors"
	"fmt"
	"github.com/daidai53/go-homework/week2/webook/internal/domain"
	"github.com/daidai53/go-homework/week2/webook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = "(?=.*([a-zA-Z].*))(?=.*[0-9].*)[a-zA-Z0-9-*/+.~!@#$%^&*()]{8,72}$"
	birthdayRegexPattern = "^(?:(?!0000)[0-9]{4}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1[0-9]|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)" +
		"|(?:0[13578]|1[02])-31)|(?:[0-9]{2}(?:0[48]|[2468][048]|[13579][26])|(?:0[48]|[2468][048]|[13579][26])00)-02-29)$"
)

type UserHandler struct {
	emailRegExp    *regexp.Regexp
	passwordRegExp *regexp.Regexp
	birthdayRegExp *regexp.Regexp
	svc            *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRegExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		birthdayRegExp: regexp.MustCompile(birthdayRegexPattern, regexp.None),
		svc:            svc,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
}

func (u *UserHandler) SignUp(context *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var signUpReq SignUpReq
	if err := context.Bind(&signUpReq); err != nil {
		fmt.Println(err)
		return
	}

	isMail, err := u.emailRegExp.MatchString(signUpReq.Email)
	if err != nil {
		context.String(http.StatusOK, "系统错误")
		return
	}
	if !isMail {
		context.String(http.StatusOK, "邮箱输入非法")
		return
	}

	if signUpReq.Password != signUpReq.ConfirmPassword {
		context.String(http.StatusOK, "两次输入的密码不一致")
		return
	}

	isPassword, err := u.passwordRegExp.MatchString(signUpReq.Password)
	if err != nil {
		context.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		context.String(http.StatusOK, "密码必须包含字母、数字、特殊字符，并且长度不能小于8位")
		return
	}

	err = u.svc.SignUp(context, domain.User{
		Email:    signUpReq.Email,
		Password: signUpReq.Password,
	})
	switch {
	case err == nil:
		context.String(http.StatusOK, "hello 你在注册%v", isMail)
	case errors.Is(err, service.ErrDuplicateEmail):
		context.String(http.StatusOK, "邮箱冲突，请换一个")
	default:
		context.String(http.StatusOK, "系统错误:%v", err)
	}

}

func (u *UserHandler) Login(context *gin.Context) {
	type Req struct {
		Email    string
		Password string
	}
	var req Req
	err := context.Bind(&req)
	if err != nil {
		return
	}
	usr, err := u.svc.Login(context, req.Email, req.Password)
	switch {
	case err == nil:
		sess := sessions.Default(context)
		sess.Set("userId", usr.Id)
		sess.Options(sessions.Options{
			MaxAge: 900,
		})
		err = sess.Save()
		if err != nil {
			context.String(http.StatusOK, "系统错误")
			return
		}
		context.String(http.StatusOK, "登录成功")
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		context.String(http.StatusOK, "用户名或者密码不对")
	default:
		context.String(http.StatusOK, "系统错误")
	}

}

func (u *UserHandler) Edit(context *gin.Context) {
	type Req struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	sess := sessions.Default(context)
	userId := sess.Get("userId")
	var userIdInt64 int64
	userIdInt64, ok := userId.(int64)
	if !ok {
		context.String(http.StatusOK, "系统错误")
		return
	}
	var req Req
	err := context.Bind(&req)
	if err != nil {
		context.String(http.StatusOK, "系统错误")
		return
	}
	if runes := []rune(req.Nickname); len(runes) > 30 {
		context.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "昵称的长度不超过30个字符",
		})
		return
	}
	if runes := []rune(req.AboutMe); len(runes) > 300 {
		context.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "简介的长度不超过300个字符",
		})
		return
	}
	legalBirthday, err := u.birthdayRegExp.MatchString(req.Birthday)
	if err != nil {
		context.String(http.StatusOK, "系统错误")
		return
	}
	if !legalBirthday {
		context.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "生日的输入格式不合法",
		})
		return
	}
	err = u.svc.Edit(context, userIdInt64, req.Nickname, req.Birthday, req.AboutMe)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  err,
		})
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}

func (u *UserHandler) Profile(context *gin.Context) {
	sess := sessions.Default(context)
	userId := sess.Get("userId")
	var userIdInt64 int64
	userIdInt64, ok := userId.(int64)
	if !ok {
		context.String(http.StatusOK, "系统错误")
		return
	}
	err := u.svc.Profile(context, userIdInt64)
	switch {
	case err == nil:
	case errors.Is(err, service.ErrNoUserProfile):
		context.String(http.StatusOK, "系统错误")
	default:
		context.String(http.StatusOK, "系统错误")
	}
}
