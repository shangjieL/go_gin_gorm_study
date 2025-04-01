package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/http"
	"time"
)

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//解决表名复数问题 例如user表会自动转换为users表
			SingularTable: true,
		},
	})

	//数据库连接池
	sqlDB, _ := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量.
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间.
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	// 自动迁移 (这是GORM自动创建表的一种方式)
	type User struct {
		gorm.Model
		//binding:"required" 表示该字段不能为空
		Name     string    `gorm:"type:varchar(20);not null" json:"name" binding:"required"`
		Age      int       `gorm:"type:int" json:"age"`
		Birthday time.Time `gorm:"type:datetime" json:"birthday"`
	}
	db.AutoMigrate(&User{})

	ginServer := gin.Default()
	//增
	ginServer.POST("/user/add", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(user); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg": "参数错误",
			})
		} else {
			//插入数据
			db.Create(&user)
			c.JSON(http.StatusOK, gin.H{
				"msg": "添加成功",
			})
		}
	})

	//也支持写sql

	// INSERT INTO `users` (`birthday`,`updated_at`) VALUES ("2020-01-01 00:00:00.000", "2020-07-04 11:05:21.775")
	user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}
	db.Omit("Name", "Age", "CreatedAt").Create(&user)

	//DELETE FROM users WHERE id IN (1,2,3);
	db.Delete(&user, []int{1, 2, 3})

	// 根据条件更新
	// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE active=true;
	db.Model(&User{}).Where("active = ?", true).Update("name", "hello")

	// SELECT * FROM users WHERE id = "1b74413f-f3b8-409f-ac47-e8c062e3472a";
	db.First(&user, "id = ?", "1b74413f-f3b8-409f-ac47-e8c062e3472a")

	// LIKE
	// SELECT * FROM users WHERE name LIKE '%jin%';
	db.Where("name LIKE ?", "%jin%").Find(&user)

}
