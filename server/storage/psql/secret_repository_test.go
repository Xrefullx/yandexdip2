package psql

import (
	"context"
	"errors"
	"github.com/Xrefullx/YanDip/server/model"
	"github.com/google/uuid"
	"github.com/icrowley/fake"
	"math/rand"
)

var user = model.User{
	Login:        fake.CharactersN(60),
	PasswordHash: fake.CharactersN(60),
}

func (s *TestSuite) TestSecret_AddGet() {
	s.Run("Get not existed", func() {
		_, err := s.storage.Secret().Get(s.ctx, uuid.New(), uuid.New())
		s.Require().Error(err)
		s.Require().True(errors.Is(err, model.ErrorItemNotFound))
	})

	s.Run("Create and Get", func() {
		user, err := s.storage.User().Create(s.ctx, user)

		s.Require().NoError(err)

		secret := getMockSecret(user.ID)

		secret.ID, err = s.storage.Secret().Add(s.ctx, secret)
		s.Require().NoError(err)

		secretRes, err := s.storage.Secret().Get(s.ctx, secret.ID, user.ID)
		s.Require().NoError(err)

		s.Assert().EqualValues(secret.ID, secretRes.ID)

		s.Assert().EqualValues(secret.UserID, secretRes.UserID)
		s.Assert().EqualValues(secret.Ver, secretRes.Ver)
		s.Assert().EqualValues(secret.Data, secretRes.Data)
		s.Assert().EqualValues(secret.IsDeleted, secretRes.IsDeleted)

	})

	s.dropTables()
}

func (s *TestSuite) TestStorage_GetUserVersionList() {
	s.Run("Add list and get list of ver", func() {
		user, err := s.storage.User().Create(s.ctx, user)
		count := 10
		arrSecrets := make([]model.Secret, count)

		for i := 0; i < count; i++ {
			secret := getMockSecret(user.ID)

			id, err := s.storage.Secret().Add(context.Background(), secret)
			s.Require().NoError(err)

			secret.ID = id

			arrSecrets[i] = secret
		}

		info, err := s.storage.Secret().GetUserVersionList(context.Background(), user.ID)
		s.Require().NoError(err)

		for _, el := range arrSecrets {
			item, ok := info[el.ID]

			s.Assert().True(ok)

			s.Assert().EqualValues(item, el.Ver)
		}
	})

	s.dropTables()
}

func (s *TestSuite) dropTables() {
	_, err := s.storage.db.Exec("DELETE FROM Secrets")
	s.Require().NoError(err)
	_, err = s.storage.db.Exec("DELETE FROM Users")
	s.Require().NoError(err)
}

func getMockSecret(userID uuid.UUID) model.Secret {
	return model.Secret{
		UserID:    userID,
		Ver:       rand.Intn(20),
		IsDeleted: false,
		Data:      fake.CharactersN(2000),
	}
}
