package constant

import "time"

const PatternStr = "/[`~!@#$%^&*()_\\-+=<>?:\"{}|,.\\/;'\\\\[\\]·~！@#￥%……&*（）——\\-+={}|《》？：“”【】、；‘'，。、]/"

const (
	UsernameMinLen = 6
	PasswordMinLen = 8
)

const (
	BlankString = ""
)

// Jwt

const KeyJwtUserId = "jwtUserId"
const TokenPrefixStr = "login:token:"
const TokenExpireTime = time.Hour * 24 * 7

// Redis Key

const KeyUserId = "user_id"
const KeyUserRole = "user_role"
const KeyUsername = "username"
const KeyAvatarUrl = "avatar_url"

// CORS

const (
	AllOrigins = "*"
)
