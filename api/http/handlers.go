package http

import "github.com/gin-gonic/gin"

func (s *Server) login(ctx *gin.Context) {

}

func (s *Server) register(ctx *gin.Context) {

}

func (s *Server) userUpdate(ctx *gin.Context) {

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
	}
	ctx.JSON(200, fi)
}


