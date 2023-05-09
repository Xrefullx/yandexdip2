package sqllte

import (
	"errors"
	"github.com/Xrefullx/YanDip/client/model"
	"github.com/Xrefullx/YanDip/client/pkg"
	"github.com/google/uuid"
	"github.com/icrowley/fake"
)

func (s *TestSuite) TestStorage_UpdateSecret() {
	s.runDropSecrets("Update exist", func() {

		toAdd := getMockSecret()

		// add to storage
		id, err := s.storage.AddSecret(toAdd)
		s.Require().NoError(err)

		// read added
		toAdd, err = s.storage.GetSecret(id)
		s.Require().NoError(err)

		//  update field, except id and timestamp
		toUpdate := getMockSecret()
		toUpdate.ID = id
		toUpdate.TimeStamp = toAdd.TimeStamp

		//  update
		err = s.storage.UpdateSecret(toUpdate)
		s.Require().NoError(err)

		// read updated
		dbUpdated, err := s.storage.GetSecret(id)
		s.Require().NoError(err)

		s.Assert().EqualValues(dbUpdated.ID, toUpdate.ID)
		s.Assert().EqualValues(dbUpdated.TypeID, toUpdate.TypeID)
		s.Assert().EqualValues(dbUpdated.Title, toUpdate.Title)
		s.Assert().EqualValues(dbUpdated.Description, toUpdate.Description)
		s.Assert().EqualValues(dbUpdated.SecretID, toUpdate.SecretID)
		s.Assert().EqualValues(dbUpdated.SecretVer, toUpdate.SecretVer)
		s.Assert().EqualValues(dbUpdated.StatusID, toUpdate.StatusID)
		s.Require().True(dbUpdated.TimeStamp > toUpdate.TimeStamp)
	})

	s.runDropSecrets("Update wrong timestamp", func() {

		toAdd := getMockSecret()

		// add to storage
		id, err := s.storage.AddSecret(toAdd)
		s.Require().NoError(err)

		// read added
		toAdd, err = s.storage.GetSecret(id)
		s.Require().NoError(err)

		//  update field, except id and timestamp
		toUpdate := getMockSecret()
		toUpdate.ID = id
		toUpdate.TimeStamp = pkg.MakeTimestamp()

		//  update
		err = s.storage.UpdateSecret(toUpdate)
		s.Require().Error(err)
	})

	s.runDropSecrets("Update not exist", func() {
		secret1 := getMockSecret()
		secret1.ID = 201

		err := s.storage.UpdateSecret(secret1)
		s.Require().Error(err)
	})
}

func (s *TestSuite) TestStorage_GetInfoList() {
	s.runDropSecrets("Add list and get list of info", func() {
		count := 10
		arrSecrets := make([]model.Secret, count)

		for i := 0; i < count; i++ {
			secret := getMockSecret()

			id, err := s.storage.AddSecret(secret)
			s.Require().NoError(err)

			secret.ID = id

			arrSecrets[i] = secret
		}

		infoArr, err := s.storage.GetMetaList()
		s.Require().NoError(err)

		for i, el := range infoArr {
			s.Assert().EqualValues(arrSecrets[i].ID, el.ID)
			s.Assert().EqualValues(arrSecrets[i].SecretID, el.SecretID)
			s.Assert().EqualValues(arrSecrets[i].SecretVer, el.SecretVer)
			s.Assert().EqualValues(arrSecrets[i].StatusID, el.StatusID)
			s.Assert().NotEmpty(el.TimeStamp)
		}
	})
}

func (s *TestSuite) TestStorage_Add_Get() {
	s.runDropSecrets("Add and read secret", func() {
		testSecret := getMockSecret()

		id, err := s.storage.AddSecret(testSecret)
		s.Require().NoError(err)

		dbSecret, err := s.storage.GetSecret(id)

		s.Assert().EqualValues(testSecret.TypeID, dbSecret.TypeID)
		s.Assert().EqualValues(testSecret.Title, dbSecret.Title)
		s.Assert().EqualValues(testSecret.Description, dbSecret.Description)
		s.Assert().EqualValues(testSecret.SecretID, dbSecret.SecretID)
		s.Assert().EqualValues(testSecret.SecretVer, dbSecret.SecretVer)
		s.Assert().EqualValues(testSecret.StatusID, dbSecret.StatusID)
		s.Assert().EqualValues(testSecret.SecretData, dbSecret.SecretData)

		s.Assert().NotEmpty(dbSecret.TimeStamp)
	})

	s.runDropSecrets("Get not exist", func() {
		_, err := s.storage.GetSecret(564)
		s.Require().Error(err)
		s.Require().True(errors.Is(err, model.ErrorItemNotFound))
	})
}

func (s *TestSuite) runDropSecrets(name string, subtest func()) bool {
	defer s.dropSecretsTable()
	return s.Run(name, subtest)
}

func (s *TestSuite) dropSecretsTable() {
	_, err := s.storage.db.Exec("DELETE FROM Secrets")
	s.Require().NoError(err)
}

func getMockSecret() model.Secret {
	return model.Secret{
		Info: model.Info{
			TypeID:      model.SecretTypes["CARD"],
			Title:       fake.Company(),
			Description: fake.CharactersN(200),
		},
		SecretID:   uuid.New(),
		SecretVer:  1,
		StatusID:   3,
		SecretData: fake.CharactersN(2000),
	}
}
