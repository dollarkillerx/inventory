package simple

import (
	"github.com/dollarkillerx/inventory/internal/pkg/models"
	"github.com/pkg/errors"
)

func (s *Simple) GetUserCenter(account string) (*models.UserCenter, error) {
	var uc models.UserCenter
	err := s.DB().Model(&models.UserCenter{}).
		Where("account = ?", account).First(&uc).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &uc, nil
}
