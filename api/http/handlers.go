package http

import (
	"fmt"
	"github.com/clearcodecn/wetalk/api/model"
	"github.com/clearcodecn/wetalk/pkg/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"regexp"
	"time"
)

const paramError = "参数错误"

var (
	emailRegexp  = regexp.MustCompile(`^[A-Za-z0-9]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`)
	mobileRegexp = regexp.MustCompile(`^1[3|5|6|7|8|9]\d{9}$`)
)

const (
	smsFormat = "欢迎注册wetalk,您的验证码是: %s"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// login user to login the system
func (s *Server) login(ctx *gin.Context) {
	req := new(LoginRequest)
	if err := ctx.BindJSON(req); err != nil {
		ctx.JSON(422, fail(paramError, nil))
		return
	}

	if req.Username == "" || req.Password == "" {
		ctx.JSON(422, fail(paramError, nil))
		return
	}

	var (
		user *model.User
		err  error
	)
	// if it is email
	if emailRegexp.MatchString(req.Username) {
		user, err = s.model.GetUserByEmail(req.Username)
	} else {
		user, err = s.model.GetUserByMobile(req.Username)
	}
	if err != nil {
		if err == model.ErrNotFound {
			ctx.JSON(422, fail("用户不存在", nil))
		} else {
			ctx.JSON(422, fail("登录失败", nil))
		}
		return
	}

	if user.Password != util.Md5(req.Password) {
		ctx.JSON(422, fail("密码错误", nil))
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.MapClaims{
		"id": user.ID,
	}
	if tokenString, err := token.SignedString(s.config.HttpConfig.JwtKey); err != nil {
		ctx.JSON(422, fail("密码错误", nil))
		return
	} else {
		if s.LoginHook != nil {
			for _, h := range s.LoginHook {
				h(user)
			}
		}
		ctx.JSON(200, successObject("", tokenString))
	}
}

// RegisterRequest is the params to register
type RegisterRequest struct {
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Code     string `json:"code"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// register register a new user
func (s *Server) register(ctx *gin.Context) {
	if !s.config.HttpConfig.EnableRegister {
		ctx.JSON(200, fail("服务器禁止注册", nil))
		return
	}

	req := new(RegisterRequest)
	if err := ctx.BindJSON(req); err != nil {
		ctx.JSON(422, fail(paramError, nil))
		return
	}

	var info string
	if req.Mobile != "" {
		info = req.Mobile
	}
	if req.Email != "" {
		info = req.Email
	}
	if info == "" || req.Avatar == "" || req.Nickname == "" || req.Password == "" {
		ctx.JSON(422, fail(paramError, nil))
		return
	}

	if s.config.HttpConfig.EnableVerify {
		if req.Code == "" {
			ctx.JSON(422, fail(paramError, nil))
			return
		}
		var vc = model.VerifyCode{
			Code:     req.Code,
			Info:     info,
			Verified: false,
			Type:     model.CodeRegister,
		}
		if !s.model.VerifyCode(&vc) {
			ctx.JSON(422, fail("验证码错误", nil))
			return
		}
	}

	user := model.User{
		Mobile:    req.Mobile,
		Avatar:    req.Avatar,
		Password:  util.Md5(req.Password),
		AddVerify: false,
		CreateAt:  time.Now(),
		DeleteAt:  time.Time{},
	}
	if err := s.model.CreateUser(&user); err != nil {
		ctx.JSON(500, fail("注册失败", err))
		return
	}
	ctx.JSON(200, success("注册成功"))
}

// userUpdate is update the user info
func (s *Server) userUpdate(ctx *gin.Context) {

}

type SendEmailVerifyCodeRequest struct {
	Email string `json:"email"`
}

// sendEmailVerifyCode send a email verify code
func (s *Server) sendEmailVerifyCode(ctx *gin.Context) {
	req := new(SendEmailVerifyCodeRequest)
	if err := ctx.BindJSON(req); err != nil || req.Email == "" {
		ctx.JSON(422, fail(paramError, nil))
		return
	}

	if !emailRegexp.MatchString(req.Email) {
		ctx.JSON(422, fail("邮箱格式错误", nil))
		return
	}
	code := &model.VerifyCode{
		Code: util.RandNumN(6),
		Info: req.Email,
		Type: model.CodeRegister,
	}
	if s.mailChan != nil {
		s.mailChan <- &MailInfo{
			VerifyCode: code,
			Title:      "welcome to register wetalk",
			Content:    fmt.Sprintf("your verify code is <br /><strong>%s</strong>", code.Code),
		}
	}
	ctx.JSON(200, success("success"))
}

// Upload upload a new file
func (s *Server) Upload(ctx *gin.Context) {
	mh, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(422, fail("获取文件失败", err))
		return
	}

	fi, err := s.uploader.UploadHTTP(mh)
	if err != nil {
		ctx.JSON(422, fail("上传文件失败", err))
		return
	}
	ctx.JSON(200, fi)
}

// SendSmsRequest send sms request.
type SendSmsRequest struct {
	Mobile string `json:"mobile"`
}

// sendSmsVerifyCode send sms code
func (s *Server) sendSmsVerifyCode(ctx *gin.Context) {
	req := new(SendSmsRequest)
	if err := ctx.BindJSON(req); err != nil {
		ctx.JSON(422, fail(paramError, nil))
		return
	}
	if !mobileRegexp.MatchString(req.Mobile) {
		ctx.JSON(422, fail("手机号格式错误", nil))
		return
	}
	code := &model.VerifyCode{
		Code:     util.RandN(6),
		Info:     req.Mobile,
		Type:     model.CodeRegister,
		Verified: false,
	}
	if s.smsChan != nil {
		s.smsChan <- &SmsInfo{
			VerifyCode: code,
			Content:    fmt.Sprintf(smsFormat, code.Code),
		}
	}
	ctx.JSON(200, success("success"))
}

// userinfo
func (s *Server) getUserInfo(ctx *gin.Context) {
	user, _, err := s.GetUser(ctx.Request.Context())
	if err != nil {
		ctx.JSON(422, fail("unexcept error", err))
		return
	}
	ctx.JSON(200, successObject("success", user))
}
