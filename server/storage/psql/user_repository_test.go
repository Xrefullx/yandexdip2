package psql

import (
	"errors"
	"github.com/Xrefullx/YanDip/server/model"
	"github.com/google/uuid"
	"github.com/icrowley/fake"
)

func (s *TestSuite) TestUsers_CreateUser() {
	user := model.User{
		Login:        fake.CharactersN(8),
		PasswordHash: fake.CharactersN(60),
		MasterHash:   fake.CharactersN(60),
	}

	s.Run("Create non-existing user", func() {
		res, err := s.storage.userRepo.Create(s.ctx, user)
		s.Require().NoError(err)

		s.Assert().EqualValues(user.Login, res.Login)
		s.Assert().EqualValues(user.PasswordHash, res.PasswordHash)
		s.Assert().EqualValues(user.MasterHash, res.MasterHash)
		s.Assert().NotEqual(uuid.Nil, res.ID)
	})

	s.Run("Try to create existing user", func() {
		_, err := s.storage.userRepo.Create(s.ctx, user)
		s.Require().Error(err)
		s.Require().True(errors.Is(err, model.ErrorConflictSaveUser))
	})
}

func (s *TestSuite) TestUsers_Exist() {
	s.Run("Check non exist user", func() {
		exist, err := s.storage.userRepo.Exist(s.ctx, uuid.New())
		s.Require().NoError(err)
		s.Require().False(exist)
	})

	s.Run("Check exist user", func() {
		user := model.User{
			Login:        fake.CharactersN(8),
			PasswordHash: fake.CharactersN(60),
		}

		res, err := s.storage.userRepo.Create(s.ctx, user)
		s.Require().NoError(err)

		exist, err := s.storage.userRepo.Exist(s.ctx, res.ID)
		s.Require().NoError(err)
		s.Require().True(exist)
	})
}

func (s *TestSuite) TestUsers_GetByLogin() {
	s.Run("Get not exist", func() {
		_, err := s.storage.userRepo.GetByLogin(s.ctx, "login_not_exist")
		s.Require().Error(err)
		s.Require().True(errors.Is(err, model.ErrorItemNotFound))
	})

	s.Run("Get with empty login", func() {
		_, err := s.storage.userRepo.GetByLogin(s.ctx, "login_not_exist")
		s.Require().Error(err)
	})

	s.Run("Get exist", func() {
		user := model.User{
			Login:        fake.CharactersN(8),
			PasswordHash: fake.CharactersN(60),
		}

		_, err := s.storage.userRepo.Create(s.ctx, user)
		s.Require().NoError(err)

		res, err := s.storage.userRepo.GetByLogin(s.ctx, user.Login)
		s.Require().NoError(err)
		s.Require().EqualValues(res.Login, user.Login)
	})
}
