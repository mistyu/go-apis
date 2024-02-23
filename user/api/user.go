package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go-apis/user/forms"
	"go-apis/user/global"
	"go-apis/user/global/response"
	"go-apis/user/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleValidatorError(ctx *gin.Context, err error) {
	// 返回错误信息
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
}

// HandleGrpcErrorToHttp 将grpc的code转换成http的状态码
func HandleGrpcErrorToHttp(err error, c *gin.Context) {

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})
			}
		}
	}
}

// GetUserList 获取用户列表
func GetUserList(ctx *gin.Context) {
	// 拨号连接用户grpc服务器
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvConfig.Host, global.ServerConfig.UserSrvConfig.Port))
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接【用户服务失败】", "msg", err.Error())
	}
	// 调用接口
	userSrvClient := proto.NewUserClient(userConn)

	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("pSize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{Pn: uint32(pnInt), PSize: uint32(pSizeInt)})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询【用户列表】失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			Birthday: response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			//Birthday: time.Time(time.Unix(int64(value.BirthDay), 0)).Format("2024-02-09"),
			Gender: value.Gender,

			Mobile: value.Mobile,
		}

		result = append(result, user)
	}

	ctx.JSON(http.StatusOK, result)
}

func PasswordLogin(ctx *gin.Context) {
	// 表单验证
	passwordLoginForm := forms.PassWordLoginForm{}
	if err := ctx.ShouldBindJSON(&passwordLoginForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvConfig.Host, global.ServerConfig.UserSrvConfig.Port))
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接【用户服务失败】", "msg", err.Error())
	}
	// 调用接口
	userSrvClient := proto.NewUserClient(userConn)

	// 登录
	if rsp, err := userSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败",
				})
			}
			return

		}
	} else {
		// 检查到了用户
		if passRsp, passErr := userSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          passwordLoginForm.PassWord,
			EncryptedPassword: rsp.Password,
		}); passErr != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]string{
				"mobile": "登录失败",
			})
		} else {
			if passRsp.Success {
				ctx.JSON(http.StatusOK, map[string]string{
					"mobile": "登录成功",
				})
			} else {
				ctx.JSON(http.StatusOK, map[string]string{
					"mobile": "登录失败",
				})
			}
		}
	}
}
