package http

import (
	"fmt"
	"github.com/clearcodecn/wetalk/api/model"
	"github.com/clearcodecn/wetalk/pkg/util"
	"github.com/gin-gonic/gin"
	"regexp"
	"time"
)

const paramError = "param error"

var (
	emailRegexp = regexp.MustCompile(`^[A-Za-z0-9]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`)
)

func (s *Server) login(ctx *gin.Context) {

}

type RegisterRequest struct {
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Code     string `json:"code"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

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

	_ = user
}

func (s *Server) userUpdate(ctx *gin.Context) {

}

type SendEmailVerifyCodeRequest struct {
	Email string `json:"email"`
}

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
		User: req.Email,
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
