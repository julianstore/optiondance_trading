package codegen

import (
	"bytes"
	"errors"
	"html/template"
	"option-dance/pkg/util"
	"os"
	"path"
	"strings"
	"syscall"
)

const (
	service_tmpl = `package service

import (
	"youmai/common/util"
	"youmai/core"
	"youmai/dao"
)

func {{.upName}}Page(current int64, size int64, qs string) (list []dao.{{.upName}}, total int64, pages int64) {
	dao.{{.upName}}Mgr(core.Db).Count(&total)
	pages = util.GetPages(total, size)
	dao.{{.upName}}Mgr(core.Db).Offset(int((current - 1) * size)).Limit(int(size)).Order("created_at desc").Find(&list)
	return
}

func {{.upName}}Create({{.name}} dao.{{.upName}}) error {
	if err := dao.{{.upName}}Mgr(core.Db).Create(&{{.name}}).Error; err != nil {
		return err
	}
	return nil
}

func {{.upName}}Detail(id string) ({{.name}} dao.{{.upName}},err error) {
	if err := dao.{{.upName}}Mgr(core.Db).Where("id=?",id).First(&{{.name}}).Error; err != nil {
		return {{.name}},err
	}
	return {{.name}},nil
}

func {{.upName}}Update({{.name}} dao.{{.upName}}) error {
	if err := dao.{{.upName}}Mgr(core.Db).Where("id=?",{{.name}}.ID).Updates({{.name}}).Error; err != nil {
		return err
	}
	return nil
}

func {{.upName}}Delete(id int64) error {
	if err := dao.{{.upName}}Mgr(core.Db).Where("id=?",id).Update("status",9).Error; err != nil {
		return err
	}
	return nil
}
`

	api_tplt = `package api

import (
	"github.com/gin-gonic/gin"
	"youmai/common/util"
	"youmai/dao"
	"youmai/service"
)

func {{.upName}}Page(c *gin.Context) {
	current, size, qs, err := util.PageInfo(c)
	if err != nil {
		http.FailWithErr(err,c)
  	    return
	}
	list, total, pages := service.{{.upName}}Page(current, size, qs)
	http.OkWithData(gin.H{
		"code":    0,
		"records": list,
		"total":   total,
		"pages":   pages,
	},c)
}

func {{.upName}}Create(c *gin.Context) {
	var {{.name}} dao.{{.upName}}
	err := c.BindJSON(&{{.name}})
	if err != nil {
		http.FailWithErr(err,c)
		return
	}
	err = service.{{.upName}}Create({{.name}})
	if err != nil {
		http.FailWithErr(err,c)
		return
	}
	http.Ok(c)
	return
}


func {{.upName}}Detail(c *gin.Context) {
	id := c.Param("id")
	detail, err := service.{{.upName}}Detail(id)
	if err != nil {
		http.FailWithErr(err,c)
		return
	}
	http.OkWithData(detail,c)
	return
}

func {{.upName}}Update(c *gin.Context) {
	id, _ := util.Int64Param(c, "id")
	var {{.name}} dao.{{.upName}}
	err := c.BindJSON(&{{.name}})
	if err != nil {
		http.FailWithErr(err,c)
		return
	}
	{{.name}}.ID = id
	err = service.{{.upName}}Update({{.name}})
	if err != nil {
		http.FailWithErr(err,c)
		return
	}
	http.Ok(c)
	return
}

func {{.upName}}Delete(c *gin.Context) {
	id, _ := util.Int64Param(c, "id")
	err := service.{{.upName}}Delete(id)
	if err != nil {
		http.FailWithErr(err,c)
		return
	}
	http.Ok(c)
	return
}
`

	router_tplt = `package router

	app.POST("/api/v1/{{.name}}", api.{{.upName}}Create)
	app.GET("/api/v1/{{.name}}/:id", api.{{.upName}}Detail)
	app.GET("/api/v1/{{.name}}s", api.{{.upName}}Page)
	app.POST("/api/v1/{{.name}}/:id", api.{{.upName}}Update)
	app.DELETE("/api/v1/{{.name}}/:id", api.{{.upName}}Delete)
`
)

func CodeGen(tableName string) {
	FirstUpCaseName := strings.ToUpper(tableName[:1]) + tableName[1:]
	p := util.AbsPath()
	root := path.Join(p, "../../")
	//create service
	file, err := os.Create(path.Join(root, "service", tableName+"_service.go"))
	service, err := template.New("service").Parse(service_tmpl)
	buf := new(bytes.Buffer)
	err = service.Execute(buf, map[string]interface{}{
		"name":   tableName,
		"upName": FirstUpCaseName,
	})
	_, err = file.Write([]byte(buf.String()))
	if err != nil {
		println(err.Error())
	}
	//create api
	file, err = os.Create(path.Join(root, "api", tableName+".go"))
	api, err := template.New("api").Parse(api_tplt)
	buf = new(bytes.Buffer)
	err = api.Execute(buf, map[string]interface{}{
		"name":   tableName,
		"upName": FirstUpCaseName,
	})
	_, err = file.Write([]byte(buf.String()))
	if err != nil {
		println(err.Error())
	}

	//create router
	routerPath := path.Join(root, "router", tableName+".go")
	f, err := os.OpenFile(routerPath, os.O_WRONLY|os.O_APPEND, 0666)
	if errors.Is(err, syscall.ERROR_FILE_NOT_FOUND) {
		f, err = os.Create(routerPath)
	}
	router, err := template.New("router").Parse(router_tplt)
	buf = new(bytes.Buffer)
	err = router.Execute(buf, map[string]interface{}{
		"name":   tableName,
		"upName": FirstUpCaseName,
	})
	_, err = f.Write([]byte(buf.String()))
	if err != nil {
		println(err.Error())
	}

}
