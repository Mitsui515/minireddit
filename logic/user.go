package logic

import (
	"minireddit/dao/mysql"
	"minireddit/models"
	"minireddit/pkg/jwt"
	"minireddit/pkg/snowflake"
)

// 存放业务逻辑的代码

func SignUp(p *models.ParamSignUp) (err error) {
	// 判断用户是否存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		// 数据库查询出错
		return err
	}
	// 生成UID
	userID := snowflake.GenID()
	// 构造一个User实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	// 保存进MySQL
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 传递的是指针，就能拿到UserID
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	// 生成JWT
	aToken, rToken, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}
	user.AccessToken = aToken
	user.RefreshToken = rToken
	return
}
