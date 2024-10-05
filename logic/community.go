package logic

import (
	"minireddit/dao/mysql"
	"minireddit/models"
)

// GetCommunityList 获取社区列表
func GetCommunityList() ([]*models.Community, error) {
	// 查数据库，查询所有社区并返回
	return mysql.GetCommunityList()
}

// GetCommunityDetail 获取社区详情
func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	// 查数据库，查询指定社区并返回
	return mysql.GetCommunityDetailByID(id)
}
