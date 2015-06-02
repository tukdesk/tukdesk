package helpers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tukdesk/tukdesk/backend/models"

	"github.com/astaxie/beego/validation"
	"github.com/tukdesk/httputils/tools"
)

const (
	TokenType                 = "token"
	TokenDefaultExpirationSec = 3600 * 24 * 30
	TokenDefaultExpiration    = TokenDefaultExpirationSec * time.Second

	errInvalidTokenPrefix = "helpers: invalid token: "
)

func newErrInvalidToken(reason string) error {
	return fmt.Errorf("%s%s", errInvalidTokenPrefix, reason)
}

var (
	ErrTokenNotFound = fmt.Errorf("token not found")

	ErrInvalidTokenExpired             = newErrInvalidToken("token expired")
	ErrInvalidTokenParseFailed         = newErrInvalidToken("failed to parse")
	ErrInvalidTokenInvalidType         = newErrInvalidToken("invalid type")
	ErrInvalidTokenInvalidUserId       = newErrInvalidToken("invalid userId")
	ErrInvalidTokenUserNotFound        = newErrInvalidToken("user not found")
	ErrInvalidTokenChannelInfoRequired = newErrInvalidToken("channel info required")
	ErrInvalidTokenInvalidChannelName  = newErrInvalidToken("invalid channel name")
	ErrInvalidTokenInvalidEmail        = newErrInvalidToken("invalid email")
)

func IsInvalidToken(err error) bool {
	return strings.HasPrefix(err.Error(), errInvalidTokenPrefix)
}

func TokenForUser(user *models.User, key string) string {
	data := map[string]interface{}{
		"type": TokenType,
		"uid":  user.Id.Hex(),
	}

	return tools.GenerateToken(data, TokenDefaultExpiration, []byte(key))
}

func UserFromRequest(r *http.Request, key string) (*models.User, bool, error) {
	t := r.Header.Get("Authorization")
	if t == "" {
		return nil, false, ErrTokenNotFound
	}
	return UserFromToken(t, []byte(key))
}

func UserFromToken(t string, key []byte) (*models.User, bool, error) {
	token, err := tools.ParseToken(t, key)
	if err == tools.ErrTokenExpired {
		return nil, false, ErrInvalidTokenExpired
	}

	if err != nil || !token.Valid {
		return nil, false, ErrInvalidTokenParseFailed
	}

	if typ, _ := token.Claims["type"]; typ != TokenType {
		return nil, false, ErrInvalidTokenInvalidType
	}

	// uid exists
	if uid, _ := token.Claims["uid"].(string); uid != "" {
		userId, ok := IdFromString(uid)
		if !ok {
			return nil, false, ErrInvalidTokenInvalidUserId
		}

		user, err := UserFindById(userId)
		// user not found
		if IsNotFound(err) {
			return nil, false, ErrInvalidTokenUserNotFound
		}

		if err != nil {
			return nil, false, err
		}

		return user, false, nil
	}

	// uid not found
	chName, _ := token.Claims["chName"].(string)
	chId, _ := token.Claims["chId"].(string)

	// chName and chId required
	if chName == "" || chId == "" {
		return nil, false, ErrInvalidTokenChannelInfoRequired
	}

	// channel agent not allowed
	if chName == UserChannelAgent {
		return nil, false, ErrInvalidTokenInvalidChannelName
	}

	// user email
	var email string

	// if is channel email , chId should be the email address
	if chName == UserChannelEmail {
		email = chId
	} else {
		email, _ = token.Claims["email"].(string)
	}

	if email != "" {
		v := &validation.Validation{}
		v.Email(email, "email")
		if v.HasErrors() {
			if chName == UserChannelEmail {
				return nil, false, ErrInvalidTokenInvalidEmail
			} else {
				email = ""
			}
		}
	}

	name, _ := token.Claims["name"].(string)
	if name == "" {
		if email != "" {
			name = UserGetValidNameFromEmail(email)
		}
	} else {
		name = UserGetValidName(name)
	}

	return UserMustByChannel(chName, chId, email, name)

}
