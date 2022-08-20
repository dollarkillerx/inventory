package test

import (
	"fmt"
	"testing"
	"time"
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

/**
获取本周周一的日期
*/
func GetFirstDateOfWeek() (weekMonday time.Time) {
	now := time.Now()

	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekStartDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	return weekStartDate
}

/**
获取上周的周一日期
*/
func GetLastWeekFirstDate() (weekMonday time.Time) {
	thisWeekMonday := GetFirstDateOfWeek()
	lastWeekMonday := thisWeekMonday.AddDate(0, 0, -7)
	return lastWeekMonday
}

func TestTime2(tx *testing.T) {
	type timeFilter struct {
		Key   string    `json:"key"`
		Start time.Time `json:"time"`
		End   time.Time `json:"end"`
	}

	t := time.Now()
	ts := time.Now().AddDate(0, 0, -1)
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}

	moon := t.AddDate(0, 0, -t.Day()+1)
	lastMonthFirstDay := t.AddDate(0, -1, -t.Day()+1)
	lastMonthStart := time.Date(lastMonthFirstDay.Year(), lastMonthFirstDay.Month(), lastMonthFirstDay.Day(), 0, 0, 0, 0, t.Location())

	timeFilters := []timeFilter{
		{
			Key:   "今天",
			Start: time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()),
			End:   time.Now().Add(6 * time.Hour),
		},
		{
			Key:   "昨天",
			Start: time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location()),
			End:   time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()),
		},
		{
			Key:   "本周",
			Start: time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset),
			End:   time.Now().Add(6 * time.Hour),
		},
		{
			Key:   "上周",
			Start: GetLastWeekFirstDate(),
			End:   GetFirstDateOfWeek(),
		},
		{
			Key:   "本月",
			Start: time.Date(moon.Year(), moon.Month(), moon.Day(), 0, 0, 0, 0, moon.Location()),
			End:   time.Now().Add(6 * time.Hour),
		},
		{
			Key:   "上個月",
			Start: lastMonthStart,
			End:   time.Date(moon.Year(), moon.Month(), moon.Day(), 0, 0, 0, 0, moon.Location()),
		},
	}

	for _, v := range timeFilters {
		fmt.Println(v.Key, "   ", v.Start, "  ", v.End)
	}
}
