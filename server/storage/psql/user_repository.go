package psql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Xrefullx/YanDip/server/model"
	"github.com/Xrefullx/YanDip/server/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

var _ storage.UserRepository = (*userRepository)(nil)

// userRepository implements UserRepository interface, provides actions with user records in psql storage.
type userRepository struct {
	db *sql.DB
}

// newUserRepository inits new user repository.
func newUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

// Save saves user to database.
// If login exist return ErrorConflictSaveUser
func (u *userRepository) Create(ctx context.Context, user model.User) (model.User, error) {
	err := u.db.QueryRowContext(
		ctx,
		"INSERT INTO users (login, pass_hash, master_hash) VALUES ($1, $2, $3) RETURNING id, login, pass_hash, master_hash",
		user.Login,
		user.PasswordHash,
		user.MasterHash,
	).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.MasterHash)

	if err != nil {
		//  if exist return ErrorConflictSaveUser
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == pgerrcode.UniqueViolation && pqErr.Constraint == "users_login_key" {
			return model.User{}, model.ErrorConflictSaveUser
		}

		return model.User{}, err
	}

	return user, nil
}

//	 GetByLogin selects user by login
//		if not found, returns ErrorItemNotFound
func (u *userRepository) GetByLogin(ctx context.Context, login string) (model.User, error) {
	var user model.User
	if err := user.ValidateLogin(login); err != nil {
		return model.User{}, err
	}

	if err := u.db.QueryRowContext(ctx,
		`SELECT id, login, pass_hash, master_hash FROM users WHERE login = $1`,
		login,
	).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.MasterHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, model.ErrorItemNotFound
		}
		return model.User{}, err
	}

	return user, nil
}

// Exist checks that user is exist in database.
func (u *userRepository) Exist(ctx context.Context, userID uuid.UUID) (bool, error) {
	count := 0
	err := u.db.QueryRowContext(ctx,
		"SELECT  COUNT(*) as count FROM users WHERE id = $1", userID).Scan(&count)

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
