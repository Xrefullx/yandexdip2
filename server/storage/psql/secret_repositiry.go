package psql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Xrefullx/YanDip/server/model"
	"github.com/Xrefullx/YanDip/server/storage"
	"github.com/google/uuid"
	"log"
)

var _ storage.SecretRepository = (*secretRepository)(nil)

// secretRepository implements SecretRepository interface, provides actions with order records in psql storage.
type secretRepository struct {
	db *sql.DB
}

// newOrderRepository inits new order repository.
func newSecretRepository(db *sql.DB) *secretRepository {
	return &secretRepository{
		db: db,
	}
}

func (r *secretRepository) Add(ctx context.Context, secret model.Secret) (uuid.UUID, error) {
	err := r.db.QueryRowContext(
		ctx,
		"INSERT INTO secrets(ver,user_id,data,is_deleted) VALUES($1,$2,$3,$4) "+
			"RETURNING id",
		secret.Ver,
		secret.UserID,
		secret.Data,
		secret.IsDeleted,
	).Scan(
		&secret.ID,
	)

	if err != nil {
		return uuid.Nil, err
	}

	return secret.ID, nil
}

func (r *secretRepository) Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (model.Secret, error) {
	res := model.Secret{}
	if err := r.db.QueryRowContext(ctx,
		"SELECT id, ver,user_id,data,is_deleted FROM secrets WHERE id=$1 AND user_id=$2",
		id, userID,
	).Scan(
		&res.ID,
		&res.Ver,
		&res.UserID,
		&res.Data,
		&res.IsDeleted,
	); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return model.Secret{}, model.ErrorItemNotFound
		}

		return model.Secret{}, err
	}

	return res, nil
}

func (r *secretRepository) Update(ctx context.Context, el model.Secret) error {
	query := `
		UPDATE secrets
		SET ver = $2, user_id = $3, data=$4, is_deleted = $5
		WHERE id = $1 AND user_id = $3;
`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, el.ID, el.Ver, el.UserID, el.Data, el.IsDeleted)
	if err != nil {
		return err
	}

	exists, err := res.RowsAffected()
	if exists == 0 {
		return model.ErrorItemNotFound
	}

	return nil
}

func (r *secretRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `
		UPDATE secrets
		SET is_deleted = $3
		WHERE id = $1 AND user_id = $2;
`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, id, userID, true)
	if err != nil {
		return err
	}

	exists, err := res.RowsAffected()
	if exists == 0 {
		return model.ErrorItemNotFound
	}

	return nil
}

func (r *secretRepository) GetUserVersionList(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]int, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, ver from secrets WHERE user_id = $1 AND is_deleted = $2", userID, false)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err.Error())
		}
	}()

	res := make(map[uuid.UUID]int)

	for rows.Next() {
		var key uuid.UUID
		var val int

		err = rows.Scan(&key, &val)
		if err != nil {
			return nil, err
		}

		res[key] = val
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}
