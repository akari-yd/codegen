package biz

import (
	"context"
	"github.com/jinzhu/gorm"
	"regexp"
)

var regSafeSQL = regexp.MustCompile(`(?i)(update|delete|update)\s+|['()";\<\>\[\]]`) // todo 两个update？
var sqlColumnNullErr = regexp.MustCompile(`(?i)Column\s+'([\w]+)'\s+cannot\s+be\s+null`)
var sqlRecordNotFoundErr = regexp.MustCompile(`record not found`)

func TrimSafeSQL(s *string) string {
	if s == nil {
		return ""
	}
	return regSafeSQL.ReplaceAllString(*s, "")
}

type Pager struct {
	Offset int
	Limit  int
}

func (r *Repository) Query(ctx context.Context, result interface{}, where string, args []interface{}, order string, pager Pager) error {
	return r.DB.Model(result).Where(where, args...).Order(order).Offset(pager.Offset).Limit(pager.Limit).Find(result).Error
}

type Repository struct {
	DB *gorm.DB
	//Redis *redis.Client
}
