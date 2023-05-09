package sqllte

import (
	"database/sql"
	"errors"
	"github.com/Xrefullx/YanDip/client/model"
	"github.com/Xrefullx/YanDip/client/pkg"
	"github.com/Xrefullx/YanDip/client/storage"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var _ storage.Storage = (*Storage)(nil)

const secretsTbl string = `
CREATE TABLE IF NOT EXISTS secrets (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	status_id INT NOT NULL,
    type_id INT,
	title TEXT NOT NULL,
	description TEXT,
	secret_data TEXT NOT NULL,
	secret_id UUID,
	secret_ver INT,
	time_stamp INTEGER NOT NULL
  );`

type Storage struct {
	db *sql.DB
}

// NewStorage inits new connection to psql storage.
func NewStorage(file string) (*Storage, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	if _, err = db.Exec(secretsTbl); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

// AddSecret adds new secret to storage
func (s *Storage) AddSecret(v model.Secret) (int64, error) {
	stmt, err := s.db.Prepare("INSERT INTO secrets(status_id, type_id, title, description, secret_id, secret_ver, secret_data, time_stamp) VALUES(?,?,?,?,?,?,?,?)")
	if err != nil {
		return 0, err
	}

	r, err := stmt.Exec(v.StatusID, v.TypeID, v.Title, v.Description, v.SecretID, v.SecretVer, v.SecretData, pkg.MakeTimestamp())
	if err != nil {
		return 0, err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateSecret adds new secret to storage
func (s *Storage) UpdateSecret(v model.Secret) error {
	query := `
		UPDATE secrets
		SET status_id = ?, type_id = ?, title=?, description=?, secret_id=?, secret_ver=?, secret_data=?,time_stamp=?
		WHERE id = ? AND time_stamp = ?;
`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(v.StatusID, v.TypeID, v.Title, v.Description, v.SecretID, v.SecretVer, v.SecretData, pkg.MakeTimestamp(), v.ID, v.TimeStamp)
	if err != nil {
		return err
	}

	exists, err := res.RowsAffected()

	if exists == 0 {
		return model.ErrorItemNotFound
	}

	return nil
}

// UpdateSecret adds new secret to storage
func (s *Storage) DeleteSecret(id int64) error {
	query := `DELETE FROM secrets WHERE id = ?`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	exists, err := res.RowsAffected()

	if exists == 0 {
		return model.ErrorItemNotFound
	}

	return nil
}

// GetSecret returns secret from storage
func (s *Storage) GetSecret(id int64) (model.Secret, error) {
	res := model.Secret{Info: model.Info{}}
	if err := s.db.QueryRow(
		"SELECT id, status_id, type_id, title, description, secret_id, secret_ver, secret_data, time_stamp FROM secrets WHERE id=@id",
		sql.Named("id", id),
	).Scan(
		&res.ID,
		&res.StatusID,
		&res.TypeID,
		&res.Title,
		&res.Description,
		&res.SecretID,
		&res.SecretVer,
		&res.SecretData,
		&res.TimeStamp); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return model.Secret{}, model.ErrorItemNotFound
		}
		return model.Secret{}, err
	}

	return res, nil
}

// GetSecret returns secret from storage
func (s *Storage) GetSecretByExtID(extID uuid.UUID) (model.Secret, error) {
	res := model.Secret{}
	if err := s.db.QueryRow(
		"SELECT id, status_id, type_id, title, description, secret_id, secret_ver, secret_data, time_stamp FROM secrets WHERE secret_id=@secret_id",
		sql.Named("secret_id", extID),
	).Scan(
		&res.ID,
		&res.StatusID,
		&res.TypeID,
		&res.Title,
		&res.Description,
		&res.SecretID,
		&res.SecretVer,
		&res.SecretData,
		&res.TimeStamp); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return model.Secret{}, model.ErrorItemNotFound
		}
		return model.Secret{}, err
	}

	return res, nil
}

// GetMetaList returns array of meta info secrets
func (s *Storage) GetMetaList() ([]model.SecretMeta, error) {
	var list []model.SecretMeta

	rows, err := s.db.Query(
		"SELECT id, status_id, secret_id, secret_ver, time_stamp FROM secrets")

	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err.Error())
		}
	}()

	for rows.Next() {
		var el model.SecretMeta
		err = rows.Scan(&el.ID, &el.StatusID, &el.SecretID, &el.SecretVer, &el.TimeStamp)
		if err != nil {
			return nil, err
		}

		list = append(list, el)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	if len(list) == 0 {
		list = make([]model.SecretMeta, 0)
	}

	return list, nil
}

// Close  closes database connection.
func (s *Storage) Close() {
	if s.db == nil {
		return
	}

	if err := s.db.Close(); err != nil {
		log.Println(err.Error())
	}
	s.db = nil
}
