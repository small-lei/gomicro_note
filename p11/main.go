package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
)

// 商品服务
func main() {

	etcdReg := etcd.NewRegistry(
		registry.Addrs("127.0.0.1:2379"))

	r := gin.Default()
	// 路由分组
	v1Group := r.Group("/v1")
	{
		v1Group.Handle("POST", "/prods", func(c *gin.Context) {
			var pr ProdsRequest
			// 给默认值
			err := c.Bind(&pr)
			if err != nil || pr.Size <= 0 {
				log.Println(err)
				pr = ProdsRequest{Size: 2}
			}
			c.JSON(http.StatusOK, gin.H{
				"data": NewProdList(pr.Size),
			})
		})
	}

	service := web.NewService(
		web.Name("ProdSrv"),
		web.Handler(r),
		web.Registry(etcdReg),
	)

	// 通过命令行参数启动
	// --server_address 指定地址端口，或者环境变量$MICRO_SERVER_ADDRESS]
	// 运行2个服务
	// go run main.go prodModels.go --server_address  127.0.0.1:8000
	// go run main.go prodModels.go --server_address  127.0.0.1:8001
	service.Init()
	service.Run()
}
