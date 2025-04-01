package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
	"log"
	"net/http"
)

// 自定义go的中间件 拦截器
func myHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("token", "123456")
		c.Next() //放行
		//c.Abort() //阻止
	}
}
func main() {
	//创建一个服务
	ginServer := gin.Default()
	ginServer.Use(myHandler())
	ginServer.Use(favicon.New("./favicon.ico"))

	//加载静态页面
	ginServer.LoadHTMLGlob("templates/*")
	//加载资源文件
	ginServer.Static("/static", "./static")

	//响应一个页面给前端
	ginServer.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"msg": "这是go后台传递过来的数据",
		})
	})

	//1.接受前端传递过来的参数
	ginServer.GET("/user/info", myHandler(), func(c *gin.Context) {
		//获取拦截器的值
		token := c.MustGet("token").(string)
		log.Println("token is ", token)

		userid := c.Query("userid")
		username := c.Query("username")
		c.JSON(http.StatusOK, gin.H{
			"userid":   userid,
			"username": username,
		})
	})

	//2.
	ginServer.GET("/user/info/:userid/:username", func(c *gin.Context) {
		userid := c.Param("userid")
		username := c.Param("username")
		c.JSON(http.StatusOK, gin.H{
			"userid":   userid,
			"username": username,
		})
	})

	//前端给后端传json
	ginServer.POST("/json", func(c *gin.Context) {
		data, _ := c.GetRawData()
		var m map[string]interface{}
		_ = json.Unmarshal(data, &m)
		c.JSON(http.StatusOK, m)
	})

	//重定向
	ginServer.GET("/test", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://www.baidu.com")
	})

	//404
	ginServer.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
	})

	//路由组
	userGroup := ginServer.Group("/user")
	{
		userGroup.GET("/add")
		userGroup.DELETE("/delete")
		userGroup.PUT("/put")
	}
	//端口
	ginServer.Run(":8082")
}
