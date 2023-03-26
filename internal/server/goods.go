package server

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dollarkillerx/inventory/internal/pkg/errs"
	"github.com/dollarkillerx/inventory/internal/pkg/models"
	"github.com/dollarkillerx/inventory/internal/pkg/request"
	"github.com/dollarkillerx/inventory/internal/pkg/response"
	"github.com/dollarkillerx/inventory/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func (s *Server) goods(ctx *gin.Context) {
	//model := utils.GetAuthModel(ctx)
}

func (s *Server) search(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)

	keyword := strings.TrimSpace(ctx.Query("keyword"))
	search, err := s.storage.Search(keyword, model.Account)
	if err != nil {
		response.Return(ctx, errs.BadRequest)
		return
	}

	for i, v := range search {
		if strings.TrimSpace(v.Img) == "" {
			search[i].Img = "https://jkl-1253341723.cos.ap-chengdu.myqcloud.com/default.png"
		}
	}

	response.Return(ctx, search)
}

func (s *Server) good(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)

	barcodes := ctx.Param("barcodes")
	good, err := s.storage.Good(barcodes, model.Account)
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.BadRequest)
		return
	}

	// https://jkl-1253341723.cos.ap-chengdu.myqcloud.com/default.png
	if strings.TrimSpace(good.Img) == "" {
		good.Img = "https://jkl-1253341723.cos.ap-chengdu.myqcloud.com/default.png"
	}
	response.Return(ctx, good)
}

func (s *Server) deleteGood(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)

	goodID := ctx.Param("goodID")
	err := s.storage.DeleteGood(goodID, model.Account)
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.BadRequest)
		return
	}

	response.Return(ctx, gin.H{})
}

func (s *Server) addGood(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)

	var good request.AddGoods
	if err := ctx.ShouldBindJSON(&good); err != nil {
		log.Println(err)
		response.Return(ctx, errs.BadRequest)
		return
	}

	if len(good.Barcode) != 13 {
		response.Return(ctx, errs.NewError("40001", "条形码不合法 长度不为 13"))
		return
	}

	err := s.storage.DB().Model(&models.Goods{}).Create(&models.Goods{
		BasicModel: models.BasicModel{ID: xid.New().String()},
		Barcode:    good.Barcode,
		Name:       good.Name,
		Spec:       good.Spec,
		Cost:       good.Cost,
		Price:      good.Price,
		Brand:      good.Brand,
		MadeIn:     good.MadeIn,
		Img:        good.Img,
		ByAccount:  model.Account,
	}).Error
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			response.Return(ctx, errs.NewError("4002", "商品以存在"))
			return
		}
		log.Println(err)
		response.Return(ctx, errs.SqlSystemError)
		return
	}

	response.Return(ctx, gin.H{})
}

// 更新商品
func (s *Server) upGood(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)
	var good request.ModifyGoods
	if err := ctx.ShouldBindJSON(&good); err != nil {
		log.Println(err)
		response.Return(ctx, errs.BadRequest)
		return
	}

	err := s.storage.DB().Model(&models.Goods{}).
		Where("id = ?", good.ID).
		//Where("by_account = ?", model.Account).
		Updates(&models.Goods{
			Name:      good.Name,
			Spec:      good.Spec,
			Cost:      good.Cost,
			Price:     good.Price,
			Brand:     good.Brand,
			MadeIn:    good.MadeIn,
			Img:       good.Img,
			UpAccount: model.Account,
		}).Error
	if err != nil {
		//        if strings.Contains(err.Error(), "unique") {
		//            response.Return(ctx, errs.NewError("4002", "商品以存在"))
		//            return
		//        }
		log.Println(err)
		response.Return(ctx, errs.SqlSystemError)
		return
	}

	response.Return(ctx, gin.H{})
}

func (s *Server) export(ctx *gin.Context) {
	ctx.Writer.Header().Add("Content-type", "text/html")
	ctx.Writer.Write([]byte(cs))
}

func (s *Server) coreExport(ctx *gin.Context) {
	var payload request.CoreExport

	//err := ctx.ShouldBindJSON(&payload)
	//if err != nil {
	//	log.Println(err)
	//	response.Return(ctx, errs.BadRequest)
	//	return
	//}

	payload.Account = ctx.Param("account")
	payload.Password = ctx.Param("password")

	uc, err := s.storage.GetUserCenter(payload.Account)
	if err != nil {
		response.Return(ctx, errs.LoginFailed)
		return
	}

	if uc.Password != payload.Password {
		response.Return(ctx, errs.LoginFailed)
		return
	}

	var ins []models.Inventory
	err = s.storage.DB().Model(&models.Inventory{}).
		Where("account = ?", payload.Account).
		Where("quantity != 0").
		Order("created_at desc").Find(&ins).Error
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.SqlSystemError)
		return
	}

	var gs []models.Goods
	err = s.storage.DB().Model(&models.Goods{}).Find(&gs).Error
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.SqlSystemError)
		return
	}

	var vm = map[string]string{}
	for _, v := range gs {
		vm[v.ID] = v.Name
	}

	var css = bytes.NewBuffer([]byte(""))
	fw := csv.NewWriter(css)
	fw.Write([]string{"创建时间: ", fmt.Sprintf("%s", time.Now().Format("2006-01-02")), "店铺: ", fmt.Sprintf("%s", uc.Storehouse), fmt.Sprintf("%s", payload.Account)})
	fw.Write([]string{"Barcode", "商品名称", "库存数量", "总成本", "更新时期"})
	for _, v := range ins {
		fw.Write([]string{fmt.Sprintf("'%s", v.Barcode), vm[v.GoodsID], fmt.Sprintf("%d", v.Quantity), fmt.Sprintf("%.2f", v.Cost), v.UpdatedAt.Format("2006-01-02")})
	}

	fw.Flush()

	//返回文件流
	ctx.Writer.Header().Add("Content-type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+fmt.Sprintf("%s-%s-%s.csv", uc.Storehouse, payload.Account, time.Now().Format("2006-01-02")))
	ctx.Writer.Write(css.Bytes())
}

var cs = "<!DOCTYPE html>\n<html lang=\"zh-cmn-Hans\">\n\n<head>\n    <title>EXPORT</title>\n    <meta charset=\"utf-8\">\n    <script src=\"https://cdn.staticfile.org/vue/2.6.10/vue.min.js\"></script>\n    <script src=\"https://cdn.jsdelivr.net/npm/axios@0.27.2/dist/axios.min.js\"></script>\n    <meta name=\"viewport\" content=\"width=device-width,initial-scale=1,shrink-to-fit=no\">\n    <link rel=\"stylesheet\" href=\"https://www.dollarkiller.com/lib/css/bootstrap.min.css\">\n    <link rel=\"stylesheet\" href=\"https://www.dollarkiller.com/lib/css/style.css\">\n    <link rel=\"stylesheet\" href=\"https://www.dollarkiller.com/lib/css/themecolor.css\">\n</head>\n\n<body>\n<header>\n    <div class=\"collapse bg-light\" id=\"navbarHeader\">\n        <div class=\"container\">\n            <div class=\"row\">\n                <div class=\"col-sm-12 col-md-11 py-4\">\n                    <h4 class=\"text-muted\">关于</h4>\n                    <p class=\"text-muted\">DollarKiller是一个程序,他也具有一定的金融知识 (贪婪唤起赌博的天性，被弥漫的繁荣所驱使)。</p>\n                    <a href=\"https://blog.dollarkiller.com\" class=\"btn btn-dark d-sm-none\">Blog</a>\n                </div>\n            </div>\n        </div>\n    </div>\n    <div class=\"navbar navbar-dark bg-info box-shadow\">\n        <div class=\"container d-flex justify-content-between\">\n            <a href=\"/\" class=\"navbar-brand d-flex align-items-center\">\n                <svg xmlns=\"http://www.w3.org/2000/svg\" width=\"20\" height=\"20\" viewBox=\"0 0 24 24\" fill=\"none\"\n                     stroke=\"currentColor\"\n                     stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"mr-2\">\n                    <path d=\"M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z\"></path>\n                    <circle cx=\"12\" cy=\"13\" r=\"4\"></circle>\n                </svg>\n                <strong>DollarKiller</strong></a>\n            <ul class=\"navbar-nav flex-row ml-md-auto d-none d-md-flex\">\n                <li class=\"nav-item\">\n                    <a class=\"nav-link p-2\" href=\"https://blog.dollarkiller.com\" rel=\"noopener\">Blog</a></li>\n            </ul>\n            <button class=\"navbar-toggler\" type=\"button\" data-toggle=\"collapse\" data-target=\"#navbarHeader\"\n                    aria-controls=\"navbarHeader\"\n                    aria-expanded=\"false\" aria-label=\"Toggle navigation\">\n                <span class=\"navbar-toggler-icon\"></span>\n            </button>\n        </div>\n    </div>\n</header>\n\n<main role=\"main\">\n    <div class=\"album py-5 bg-light\">\n        <div class=\"container\">\n            <h2 class=\"text-success\">导出库存\n            </h2>\n            <div class=\"row\">\n                <div class=\"col-md-12\">\n                    <div class=\"card card_default mb-4 box-shadow\">\n                        <div class=\"card-body\">\n\n                            <div id=\"app\">\n                                <div class=\"row\">\n\n                                    <div class=\"col-sm\">\n\n                                        <div class=\"input-group input-group-sm mb-3\">\n                                            <div class=\"input-group-prepend\">\n                                                <span class=\"input-group-text\" id=\"principal\">账户</span>\n                                            </div>\n                                            <input type=\"text\" class=\"form-control\" aria-label=\"Small\"\n                                                   aria-describedby=\"principal\" v-model=\"account\">\n                                        </div>\n\n                                        <div class=\"input-group input-group-sm mb-3\">\n                                            <div class=\"input-group-prepend\">\n                                                <span class=\"input-group-text\" id=\"fixed_throw\">密码</span>\n                                            </div>\n                                            <input type=\"text\" class=\"form-control\" aria-label=\"Small\"\n                                                   aria-describedby=\"fixed_throw\" v-model=\"password\">\n                                        </div>\n\n\n\n                                        <button class=\"btn btn-primary\" @click=\"calculation\">导出</button>\n\n                                    </div>\n\n                                </div>\n                            </div>\n\n                        </div>\n                    </div>\n                </div>\n            </div>\n        </div>\n    </div>\n\n</main>\n\n<footer>\n    <div class=\"container\">\n        <p class=\"float-right\">\n            <a href=\"#\">回到顶部</a></p>\n        <p>快速进入我的博客\n            <a href=\"https://blog.dollarkiller.com\" target=\"_blank\" class=\"badge badge-primary my-2\">Blog</a>\n        <p>Copyright @DollarKiller dollarkiller.com. </p>\n    </div>\n</footer>\n<script>\n    let vue = new Vue({\n        el: \"#app\",\n        data: {\n            account: \"\",\n            password: \"\",\n\n        },\n        methods: {\n            calculation: function () {\n                if (this.account === \"\" || this.password === \"\") {\n                    alert(\"请认真填写 表单!!!\");\n                    return;\n                }\n\n\n                var url = `/export_core/${this.account}/${this.password}`;\n                window.open(url)\n                // window.location.href= url\n            }\n        }\n    })\n</script>\n<script src=\"https://www.dollarkiller.com/lib/js/jquery-3.2.1.slim.min.js\"></script>\n<script src=\"https://www.dollarkiller.com/lib/js/popper-1.12.9.min.js\"></script>\n<script src=\"https://www.dollarkiller.com/lib/js/bootstrap.min.js\"></script>\n</body>\n</html>\n"
