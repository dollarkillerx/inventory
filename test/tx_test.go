package test

import (
	"testing"
)

func TestTx(t *testing.T) {

	////将<bucket>和<region>修改为真实的信息
	////bucket的命名规则为{name}-{appid} ，此处填写的存储桶名称必须为此格式
	//u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", conf.CONF.OSSConf.Bucket, conf.CONF.OSSConf.Region))
	//b := &cos.BaseURL{BucketURL: u}
	//c := cos.NewClient(b, &http.Client{
	//	//设置超时时间
	//	Timeout: 100 * time.Second,
	//	Transport: &cos.AuthorizationTransport{
	//		//如实填写账号和密钥，也可以设置为环境变量
	//		SecretID:  conf.CONF.OSSConf.SecretID,
	//		SecretKey: conf.CONF.OSSConf.SecretKey,
	//		Transport: &debug.DebugRequestTransport{
	//			RequestHeader:  true,
	//			RequestBody:    false,
	//			ResponseHeader: true,
	//			ResponseBody:   false,
	//		},
	//	},
	//})
	//
	//// Case2 使用options上传对象
	//file, err := ioutil.ReadFile("start.png")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//opt := &cos.ObjectPutOptions{
	//	ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
	//		ContentType: utils.GetContentTypeByFileName("start.png"),
	//	},
	//	ACLHeaderOptions: &cos.ACLHeaderOptions{
	//		XCosACL: "public-read",
	//	},
	//}
	//_, err = c.Object.Put(context.Background(), "inventory/f.png", bytes.NewReader(file), opt)
	//if err != nil {
	//	log.Fatalln(err)
	//}

}
