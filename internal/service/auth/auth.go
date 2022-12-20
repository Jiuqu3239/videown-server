package auth

import (
	"errors"
	"fmt"
	"time"
	"videown-server/internal/dto"
	"videown-server/pkg/setting"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JwtHelper struct {
	JwtKey []byte
	// by second
	ValidDuration int64
}

type CustomClaims struct {
	Name    string `json:"nam,omitempty"`
	IsAdmin bool   `json:"adm,omitempty"`
	jwt.StandardClaims
}

type UsernamePasswordVerifier interface {
	Verify(username, password string) error
}

type DefaultAdminVerifyProvider struct {
	Settings *setting.Settings
}

type UserVerifyProvider struct{}

func (UserVerifyProvider) Verify(username, password string) error {
	return nil
}

func (t DefaultAdminVerifyProvider) Verify(username, password string) error {
	if username == t.Settings.AppSetting.Username && password == t.Settings.AppSetting.Password {
		return nil
	}
	return fmt.Errorf("invalid username or password")
}

var (
	ErrTokenExpired     error = errors.New("token is expired")
	ErrTokenNotValidYet error = errors.New("token not active yet")
	ErrTokenMalformed   error = errors.New("that's not even a token")
	ErrTokenInvalid     error = errors.New("couldn't handle this token")
)

var jwtHelper *JwtHelper
var adminVerifyProvider UsernamePasswordVerifier

func GetAdminVerifier() UsernamePasswordVerifier {
	return adminVerifyProvider
}

func SetupAuth(jwtKey string, tokenValidDuration int64, adminVerifier UsernamePasswordVerifier) {
	if jwtHelper != nil {
		return
	}
	jwtHelper = &JwtHelper{
		[]byte(jwtKey),
		tokenValidDuration,
	}
	adminVerifyProvider = adminVerifier
}

func Jwth() *JwtHelper {
	return jwtHelper
}

func (j *JwtHelper) GenerateToken(name string, isAdmin bool) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		name,
		isAdmin,
		jwt.StandardClaims{
			NotBefore: now.Unix() - 30,
			ExpiresAt: now.Unix() + j.ValidDuration,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.JwtKey)
}

func (j *JwtHelper) GenerateTokenByClaims(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.JwtKey)
}

func (j *JwtHelper) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.JwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenInvalid
}

func (j *JwtHelper) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.JwtKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.GenerateTokenByClaims(*claims)
	}
	return "", ErrTokenInvalid
}

type AuthResult struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

func abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

type WalletSignHolder interface {
	Address() string
	Nonce() int64
	Sign() string
}

func VerifyWalletSignature(wsh WalletSignHolder) (error, time.Time, *common.Address) {
	now := time.Now()
	if abs(now.Unix()-wsh.Nonce()) > 3 { //TODO: 3 is hardcode now, to extract to config
		return fmt.Errorf("invalid nonce value"), now, nil
	}

	//msg := fmt.Sprintf("%s%d", wsh.Address(), wsh.Nonce())
	var err error
	var addr *common.Address
	//addr, err := ethsign.EthVerifySignAddress(wsh.Address(), msg, wsh.Sign())
	if err != nil {
		return err, now, addr
	}
	return nil, now, addr
}

func RequestWalletAuth(req dto.WalletAttachReq) (*AuthResult, error) {
	var err error
	var now time.Time
	var address *common.Address
	if err, now, address = VerifyWalletSignature(req); err != nil {
		return nil, err
	}

	token, err := jwtHelper.GenerateToken(address.Hex(), false)
	if err != nil {
		return nil, fmt.Errorf("permission authentication failure")
	}

	r := AuthResult{
		Token:     token,
		ExpiresAt: now.Unix() + jwtHelper.ValidDuration,
	}
	return &r, nil
}

type LoginInfo struct {
	Claims    *CustomClaims
	ClientIp  string
	LoginTime time.Time
	IsAdmin   bool
}

func RequestAuth(username, password string, verifier UsernamePasswordVerifier) (*AuthResult, error) {
	if err := verifier.Verify(username, password); err != nil {
		return nil, err
	}
	token, err := jwtHelper.GenerateToken(username, true)
	if err != nil {
		return nil, fmt.Errorf("permission authentication failure")
	}

	r := AuthResult{
		Token:     token,
		ExpiresAt: time.Now().Unix() + jwtHelper.ValidDuration,
	}
	return &r, nil
}

const ctxKeyLoginInfo = "loginInfo"

func StoreAuth(c *gin.Context, claims *CustomClaims, isAdmin bool) {
	loginInfo := LoginInfo{
		Claims:    claims,
		ClientIp:  c.ClientIP(),
		LoginTime: time.Now(),
		IsAdmin:   isAdmin,
	}
	c.Set(ctxKeyLoginInfo, &loginInfo)
}

func GetLoginInfo(c *gin.Context) *LoginInfo {
	if value, ok := c.Get(ctxKeyLoginInfo); ok {
		if loginInfo, ok := value.(*LoginInfo); ok {
			return loginInfo
		}
	}
	return nil
}

func ConfirmPassword(c *gin.Context, password string) error {
	loginInfo := GetLoginInfo(c)
	if loginInfo == nil {
		return fmt.Errorf("must login first")
	}
	if !loginInfo.IsAdmin {
		return fmt.Errorf("invalid access")
	}
	if err := adminVerifyProvider.Verify(loginInfo.Claims.Name, password); err != nil {
		return fmt.Errorf("password incorrect")
	}
	return nil
}
