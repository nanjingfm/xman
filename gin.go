package xman

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/unrolled/secure"
)

type Engine struct {
	*gin.Engine
}

func NewEngine() *Engine {
	var r = gin.Default()
	// Router.Use(middleware.LoadTls())  // 打开就能玩https了

	r.Use(Cors())            // 跨域
	r.Use(I18n())            // 多语言
	RegisterCaptchaRouter(r) // 注册验证码相关路由

	return &Engine{Engine: r}
}

// AddRoute 注册路由
func (e *Engine) AddRoute(r func(r *gin.Engine)) {
	r(e.Engine)
}

// Run 启动
func (e *Engine) Run() {
	address := fmt.Sprintf(":%d", sysConf().System.Addr)
	LogDebug("server run success on ", address)
	LogError(e.Engine.Run(address))
}

func LoadTls() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:443",
		})
		err := middleware.Process(c.Writer, c.Request)
		if err != nil {
			// 如果出现错误，请不要继续
			fmt.Println(err)
			return
		}
		// 继续往下处理
		c.Next()
	}
}

// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		// 处理请求
		c.Next()
	}
}

// Custom claims structure
type CustomClaims struct {
	UUID        uuid.UUID
	ID          uint
	NickName    string
	AuthorityId string // 用户角色
	jwt.StandardClaims
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localSstorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get("x-token")
		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				Return(c, ECodeSystemErr, gin.H{
					"reload": true,
				}, "授权已过期")
				c.Abort()
				return
			}
			Return(c, ECodeSystemErr, gin.H{
				"reload": true,
			}, err.Error())
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("can't handle this token")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(sysConf().System.SigningKey),
	}
}

// 创建一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid

	}

}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}

// 登录以后签发jwt
func GetJWT(clams CustomClaims) (token string, err error) {
	j := &JWT{
		SigningKey: []byte(sysConf().System.SigningKey), // 唯一签名
	}
	token, err = j.CreateToken(clams)
	return
	//if !sysConf().System.UseMultipoint {
	//	return
	//}
	//var loginJwt mmysql.JwtBlacklist
	//loginJwt.Jwt = token
	//err, jwtStr := service.GetRedisJWT(user.Username)
	//if err == redis.Nil {
	//	if err := service.SetRedisJWT(loginJwt, user.Username); err != nil {
	//		Return(c, ERROR, nil, "设置登录状态失败")
	//		return
	//	}
	//	Return(c, SUCCESS, resp.LoginResponse{
	//		User:      user,
	//		Token:     token,
	//		ExpiresAt: clams.StandardClaims.ExpiresAt * 1000,
	//	})
	//} else if err != nil {
	//	Return(c, ERROR, nil, fmt.Sprintf("%v", err))
	//} else {
	//	var blackJWT mmysql.JwtBlacklist
	//	blackJWT.Jwt = jwtStr
	//	if err := service.AddBlacklist(blackJWT); err != nil {
	//		Return(c, ERROR, nil, "jwt作废失败")
	//		return
	//	}
	//	if err := service.SetRedisJWT(loginJwt, user.Username); err != nil {
	//		Return(c, ERROR, nil, "设置登录状态失败")
	//		return
	//	}
	//	Return(c, SUCCESS, resp.LoginResponse{
	//		User:      user,
	//		Token:     token,
	//		ExpiresAt: clams.StandardClaims.ExpiresAt * 1000,
	//	})
	//}
}
