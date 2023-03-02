package biz

import (
	"context"
	//"dms-project/graph/model"
	"errors"

	"github.com/jinzhu/copier"
)

func (r *Repository) Create_TABLENAME_(ctx context.Context, input model._TABLENAME_Input) (*model._TABLENAME_, error) {
	var bean model._TABLENAME_
	err := copier.Copy(&bean, input)
	if err != nil {
		return nil, err
	}
	err = r.DB.Create(&bean).Error
	if err != nil {
		return nil, err
	}
	return &bean, nil
}

func (r *Repository) Update_TABLENAME_(ctx context.Context, id int, input model._TABLENAME_Input) (*model._TABLENAME_, error) {
	var bean model._TABLENAME_
	err := copier.Copy(&bean, input)
	if err != nil {
		return nil, err
	}
	err = r.DB.Model(&bean).Where("id=?", id).Update(&bean).Error
	if err != nil {
		return nil, err
	}
	err = r.DB.Find(&bean, "id=?", id).Error
	if err != nil {
		return nil, err
	}
	return &bean, nil
}

func (r *Repository) Delete_TABLENAME_(ctx context.Context, id int) (bool, error) {
	err := r.DB.Where("id = ?", id).Delete(&model._TABLENAME_{}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Repository) Get_TABLENAME_(ctx context.Context, id *int, query *model._TABLENAME_Input) (*model._TABLENAME_, error) {
	var (
		queryModel model._TABLENAME_
		err        error
	)
	if id == nil && query == nil {
		return nil, errors.New("condition is nil")
	}
	if id != nil {
		err = r.DB.Where("id=?", id).Find(&queryModel).Error
		if err != nil {
			return nil, err
		}
		return &queryModel, nil
	} else {
		err = copier.Copy(&queryModel, *query)
		if err != nil {
			return nil, err
		}
		err = r.DB.Where(queryModel).Limit(1).Find(&queryModel).Error
		if err != nil {
			return nil, err
		}
		return &queryModel, nil
	}
}

func (r *Repository) Get_TABLENAME_sForPager(ctx context.Context, query *model._TABLENAME_Input, pager model.PagerInput, order *string) (*model._TABLENAME_PagerInfo, error) {
	var (
		reply      model._TABLENAME_PagerInfo
		items      []*model._TABLENAME_
		queryModel model._TABLENAME_
		total      int
		err        error
	)
	if query != nil {
		err = copier.Copy(&queryModel, *query)
		if err != nil {
			return nil, err
		}
	}
	items, err = r.Get_TABLENAME_s(ctx, query, pager, order)
	if err != nil {
		return nil, err
	}
	r.DB.Model(&queryModel).Where(queryModel).Count(&total)
	if err != nil {
		return nil, err
	}
	reply.Items = items
	reply.Total = total
	return &reply, nil
}

func (r *Repository) Get_TABLENAME_s(ctx context.Context, query *model._TABLENAME_Input, pager model.PagerInput, order *string) ([]*model._TABLENAME_, error) {
	var (
		items      []*model._TABLENAME_
		queryModel model._TABLENAME_
		err        error
		orderBy    string
	)
	if query != nil {
		err = copier.Copy(&queryModel, *query)
		if err != nil {
			return nil, err
		}
	}
	if pager.Limit > 1000 {
		pager.Limit = 1000
	}
	orderBy = TrimSafeSQL(order)
	err = r.DB.Offset(pager.Offset).Limit(pager.Limit).Where(queryModel).Order(orderBy).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
