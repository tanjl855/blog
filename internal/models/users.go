package models

import (
	"encoding/json"
	"errors"
	"fmt"
)

/**
 * 定义User表结构
 */

type User struct {
	ID           int    `json:"id" db:"id"`
	UserName     string `json:"user_name" db:"user_name"`
	Password     string `json:"password" db:"password"`
	Email        string `json:"email" db:"email"`
	Gender       int    `json:"gender" db:"gender"`
	AccessToken  string
	RefreshToken string
}

func (u *User) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Gender   int    `json:"gender"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		// 后续替换成自定义的日志系统打log
		fmt.Printf("Unmarshal JSON failed: %v", err)
		return
	}
	// check fields
	if len(required.UserName) == 0 {
		err = errors.New("缺少必填字段user_name")
		return
	}
	if len(required.Password) == 0 {
		err = errors.New("缺少必填字段password")
		return
	}
	u.UserName = required.UserName
	u.Password = required.Password
	u.Email = required.Email
	u.Gender = required.Gender
	return
}

/**
 * RegisterForm 注册请求参数
 */

type RegisterForm struct {
	UserName        string `json:"user_name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Gender          int    `json:"gender" binding:"oneof=0 1 2"` // 0:未知 1:男 2:女
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

/**
 * LoginForm 登录请求参数
 */

type LoginForm struct {
	UserName    string `json:"user_name" binding:"required"`
	Password    string `json:"password" binding:"required"`
	MessageCode string `json:"message_code"` // TODO: 后续加入sendSms校验登录
}

func (r *RegisterForm) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		UserName        string `json:"user_name"`
		Email           string `json:"email"`
		Gender          int    `json:"gender"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		fmt.Printf("Unmarshal JSON failed, err(%v)", err)
		return
	}
	// check
	if len(required.UserName) == 0 {
		err = errors.New("缺少必填字段user_name")
		return
	}
	if len(required.Password) == 0 {
		err = errors.New("缺少必填字段password")
		return
	}
	if len(required.Email) == 0 {
		err = errors.New("缺少必填字段email")
		return
	}
	if required.Password != required.ConfirmPassword {
		err = errors.New("两次密码不一致")
		return
	}
	r.UserName = required.UserName
	r.Email = required.Email
	r.Gender = required.Gender
	r.Password = required.Password
	r.ConfirmPassword = required.ConfirmPassword
	return
}

/**
 * VoteDataForm 投票数据
 */

type VoteDataForm struct {
	UserID    string `json:"user_id"`                          // 当前的用户
	PostID    string `json:"post_id" binding:"required"`       // 帖子id
	Direction int    `json:"direction" binding:"oneof=1 0 -1"` // 赞成票(1) 反对票(-1) 取消投票(0)
}

func (v *VoteDataForm) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		PostID    string `json:"post_id"`
		Direction int    `json:"direction"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		fmt.Printf("Unmarshal JSON failed, err(%v)", err)
		return
	}
	// check
	if len(required.PostID) == 0 {
		err = errors.New("缺少必填字段post_id")
		return
	}
	if required.Direction == 0 {
		err = errors.New("缺少必填字段direction")
	}
	v.PostID = required.PostID
	v.Direction = required.Direction

	return
}
