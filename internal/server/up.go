package server

import (
	"context"
	"fmt"
	"github.com/rs/xid"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/dollarkillerx/inventory/internal/conf"
	"github.com/dollarkillerx/inventory/internal/pkg/errs"
	"github.com/dollarkillerx/inventory/internal/pkg/response"
	"github.com/dollarkillerx/inventory/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func (s *Server) UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("imgFile")
	if err != nil {
		response.Return(ctx, errs.BadRequest)
		return
	}
	ctType, pfx := utils.GetContentTypeByFileName(file.Filename)
	if ctType == "" {
		response.Return(ctx, errs.BadRequest)
		return
	}

	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", conf.CONF.OSSConf.Bucket, conf.CONF.OSSConf.Region))
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		//设置超时时间
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			//如实填写账号和密钥，也可以设置为环境变量
			SecretID:  conf.CONF.OSSConf.SecretID,
			SecretKey: conf.CONF.OSSConf.SecretKey,
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   false,
			},
		},
	})

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: ctType,
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			XCosACL: "public-read",
		},
	}

	open, err := file.Open()
	if err != nil {
		response.Return(ctx, errs.BadRequest)
		return
	}
	defer open.Close()
	filename := fmt.Sprintf("inventory/%s%s", xid.New().String(), pfx)
	_, err = c.Object.Put(context.Background(), filename, open, opt)
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.BadRequest)
		return
	}

	response.Return(ctx, gin.H{"url": fmt.Sprintf("%s/%s", conf.CONF.OSSConf.Url, filename)})
}
