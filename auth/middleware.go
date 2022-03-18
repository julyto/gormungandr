package auth

import (
	"database/sql"
	"encoding/base64"
	"strings"
	"time"

	"github.com/CanalTP/gormungandr"
	"github.com/gin-gonic/gin"
	cache "github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

func getToken(c *gin.Context) string {
	username, _, ok := c.Request.BasicAuth()
	if ok {
		return username
	} else if authorization := c.GetHeader("Authorization"); len(authorization) > 0 {
		if strings.HasPrefix(strings.ToLower(authorization), "basic ") {
			//this is basic authentication with the missing ":" at the end
			//it isn't valid per the standard but jormungandr understand it, so...
			value, err := base64.StdEncoding.DecodeString(authorization[6:])
			if err != nil {
				return ""
			}
			return string(value)
		}
		return authorization
	} else {
		return c.Query("key")
	}
}

func AuthenticationMiddleware(db *sql.DB, authCache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware(c, db, authCache)
	}
}

func middleware(c *gin.Context, db *sql.DB, authCache *cache.Cache) {
	coverage := c.Param("coverage")
	token := getToken(c)
	logger := logrus.WithFields(logrus.Fields{
		"coverage": coverage,
		"token":    "######################",
	})
	logger.Debug("authentifying request")
	if token == "" {
		c.Header("WWW-Authenticate", "basic realm=\"Token Required\"")
		c.AbortWithStatusJSON(401, gin.H{"message": "no token"})
		return
	}
	user, err := CachedAuthenticate(token, time.Now(), db, authCache)
	if err == ErrAuthenticationFailed {
		c.Header("WWW-Authenticate", "basic realm=\"Token Required\"")
		c.AbortWithStatusJSON(401, gin.H{"message": "authentication failed"})
		return
	} else if err != nil {
		panic(err)
	}
	ok, err := CachedIsAuthorized(user, coverage, db, authCache)
	if err != nil {
		panic(err)
	}
	if !ok {
		c.Header("WWW-Authenticate", "basic realm=\"Token Required\"")
		c.AbortWithStatusJSON(403, gin.H{"message": "authentication failed"})
		return
	}
	gormungandr.SetUser(c, user)
	gormungandr.SetCoverage(c, coverage)
	c.Next()
}
