package handler

import (
	"bytes"
	"context"
	"douyin/gateway/auth"
	"douyin/gateway/response"
	pbvideo "douyin/grpc/video"
	"douyin/package/constant"
	"io"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type publishRequest struct {
	// 用户鉴权token
	Token string `form:"token"`
	// 视频标题
	Title string `form:"title"`
	// 新增 topic
	Topic string `form:"topic"`
}

type listRequest struct {
	// 用户鉴权token
	Token  string `query:"token"`
	UserID uint64 `query:"user_id"`
}

func PublishAction(c *fiber.Ctx) error {
	var req publishRequest
	err := c.BodyParser(&req)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
		}
		return c.JSON(res)
	}
	userClaim, err := auth.ParseToken(req.Token)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.WrongToken,
		}
		return c.JSON(res)
	}
	fileHeader, err := c.FormFile("data")
	if err != nil {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.FileFormatError,
		}
		return c.JSON(res)
	}
	// 检查文件后缀是不是mp4 大小在上传的时候会限制30MB
	if !strings.HasSuffix(fileHeader.Filename, constant.MP4Suffix) {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.FileFormatError,
		}
		return c.JSON(res)
	}
	zap.L().Info("PublishAction Filename:" + fileHeader.Filename)
	file, err := fileHeader.Open()
	if err != nil {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.FileFormatError,
		}
		return c.JSON(res)
	}
	defer file.Close()
	// 将文件转化为字节流
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.FileFormatError,
		}
		return c.JSON(res)
	}
	res, err := VideoClient.Publish(context.Background(), &pbvideo.PublishRequest{
		Title:  req.Title,
		Topic:  req.Topic,
		UserID: userClaim.UserID,
		Data:   buf.Bytes(),
	})
	if err != nil {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		return c.JSON(res)
	}
	return c.JSON(res)
}

func ListPublishedVideo(c *fiber.Ctx) error {
	var req listRequest
	err := c.QueryParser(&req)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.VideoListResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
		}
		return c.JSON(res)
	}
	var userID uint64
	if req.Token == "" {
		userID = 0
	} else {
		claims, err := auth.ParseToken(req.Token)
		if err != nil {
			res := response.UserRegisterOrLogin{
				StatusCode: constant.Failed,
				StatusMsg:  constant.WrongToken,
			}
			c.Status(fiber.StatusOK)
			return c.JSON(res)
		}
		userID = claims.UserID
	}
	resp, err := VideoClient.List(context.Background(), &pbvideo.ListRequest{
		UserID:      req.UserID,
		LoginUserID: userID,
	})
	if err != nil {
		res := response.UserRegisterOrLogin{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		return c.JSON(res)
	}
	return c.JSON(resp)
}