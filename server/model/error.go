package model

import "errors"

var (
	ErrorConflictSaveUser = errors.New("user already exist")
	ErrorItemNotFound     = errors.New("item not found")
	ErrorWrongAuthData    = errors.New("incorrect login or password")

	ErrorVersionToLow  = errors.New("version of data to low")
	ErrorItemIsDeleted = errors.New("element is deleted")
	ErrorParamNotValid = errors.New("incoming parameter not valid")

	ErrAddingUser         = errors.New("ошибка добавления пользователя")
	ErrAuthenticatingUser = errors.New("ошибка авторизации пользователя")
	ErrGeneratingToken    = errors.New("ошибка генерации токена для пользователя")
)
