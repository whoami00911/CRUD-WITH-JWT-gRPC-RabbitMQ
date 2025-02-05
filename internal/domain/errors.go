package domain

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrObsoleteToken = errors.New("Obsolete token")
var ErrTokenClaims = errors.New("token claims are not of type *tokenClaims")
var ErrSignInMethod = errors.New("invalid signing method")
var ErrTokenGen = errors.New("ошибка генерации токена")
