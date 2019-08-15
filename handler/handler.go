package handler

import (
    "fmt"
	"reflect"
	"strings"
    "time"
	"log"
	"net/http"
	"github.com/hprose/hprose-golang/rpc"
	"github.com/shioyim/samarit-my/constant"
    "github.com/shioyim/samarit-my/config"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)




type response struct {
	Success bool
	Message string
	Data    interface{}
}

type sendHeaderEvent2 struct{}


 func (e sendHeaderEvent2) OnSendHeader(ctx *rpc.HTTPContext) {
// 	ctx.Response.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8081");
// 	ctx.Response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE");
 	ctx.Response.Header().Set("Access-Control-Allow-Headers","Authorization")//Origin , X-Requested-With, Content-Type, Accept,
 	// header.Set("Access-Control-Allow-Credentials", "true") 
// //     // ctx.Response.Header().Set("Content-Type", "application/json;charset=utf-8");
 }


func Server() {
	domain  := config.String("domain")
	port    := config.String("port")
	service := rpc.NewHTTPService()
	handler := struct {
		//config logic will be wrote here.
		User      user
		// Log       logger
	}{}


	service.Event = &sendHeaderEvent2{}

	service.AddBeforeFilterHandler(func(request []byte, ctx rpc.Context, next rpc.NextFilterHandler) (response []byte, err error) {
		ctx.SetInt64("start", time.Now().UnixNano())
		httpContext := ctx.(*rpc.HTTPContext)
		if httpContext != nil {
			ctx.SetString("username", parseToken(httpContext.Request.Header.Get("authorization")))
		}
		return next(request, ctx)
	})
	service.AddInvokeHandler(func(name string, args []reflect.Value, ctx rpc.Context, next rpc.NextInvokeHandler) (results []reflect.Value, err error) {
		name = strings.Replace(name, "_", ".", 1)
		results, err = next(name, args, ctx)
		spend := (time.Now().UnixNano() - ctx.GetInt64("start")) / 1000000
		spendInfo := ""
		if spend > 1000 {
			spendInfo = fmt.Sprintf("%vs", spend/1000)
		} else {
			spendInfo = fmt.Sprintf("%vms", spend)
		}
		log.Printf("%16s() spend %s", name, spendInfo)
		return
	})


	service.AddAllMethods(handler)

	http.Handle("/api", service)
	// http.Handle("/", http.FileServer(http.Dir("web/dist")))
	log.Printf("%v v%v running at http://%v:%v\n",constant.ProjectName, constant.Version, domain, port)
	http.ListenAndServe(":"+port, nil)
}
