package service

import (
	"awesome-go/internal/models"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func (s *Service) generateToken(ip, userAgent string, expires time.Time) string {
	now := time.Now()
	nowStamp := strconv.Itoa(int(now.Unix()))
	expiresStamp := strconv.Itoa(int(expires.Unix()))
	uniqKey := userAgent + ip + expiresStamp + nowStamp
	token := md5.Sum([]byte(uniqKey))
	return hex.EncodeToString(token[:])
}

func (s *Service) CreateSession(user models.User, userAgent, IP string, expires time.Time) (models.Session, error) {
	session := models.Session{
		User:      user,
		Token:     s.generateToken(IP, userAgent, expires),
		UserAgent: userAgent,
		IP:        IP,
		ExpiresOn: expires,
	}
	err := gorm.G[models.Session](s.db).Create(s.context(), &session)
	return session, err
}

func (s *Service) GetSessionByToken(token string) (models.Session, error) {
	session, err := gorm.G[models.Session](s.db).Preload("User", nil).Where("token = ?", token).First(s.context())
	if err != nil {
		return models.Session{}, err
	}
	return session, nil
}
