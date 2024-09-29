package logic

import (
	"minireddit/dao/mysql"
	"minireddit/models"
	"minireddit/pkg/snowflake"
)

// 存放业务逻辑的代码

func SignUp(p *models.ParamSignUp) {
	// 1. 判断用户是否存在
	mysql.QueryUserByUsername()
	// 2. 生成UID
	snowflake.GenID()
	// 3. 保存进MySQL
	mysql.InsertUser()
	// 4. 保存进Redis
}
