package app

import (
	"context"
	"time"

	"github.com/g4web/otus_anti_brute_force/configs"
	"github.com/g4web/otus_anti_brute_force/internal/bucket"
)

type App struct {
	loginBuckets    *bucket.StringBuckets
	passwordBuckets *bucket.StringBuckets
	ipBuckets       *bucket.IPBuckets
}

func NewApp(ctx context.Context, cfg *configs.Config) *App {
	loginBuckets := bucket.NewStringBuckets(cfg.LoginTimeLimit, cfg.LoginMaxCountForTimeLimit)
	passwordBuckets := bucket.NewStringBuckets(cfg.PasswordTimeLimit, cfg.PasswordMaxCountForTimeLimit)
	ipBuckets := bucket.NewIPBuckets(cfg.IPTimeLimit, cfg.IPMaxCountForTimeLimit)

	garbageCleanerStart(ctx, cfg, loginBuckets, passwordBuckets, ipBuckets)

	return &App{
		loginBuckets:    loginBuckets,
		passwordBuckets: passwordBuckets,
		ipBuckets:       ipBuckets,
	}
}

func (u App) IsOk(ip string, login string, password string) (bool, error) {
	isBannedByIP, err := u.ipBuckets.IsBanned(ip)
	if err != nil {
		return false, nil
	}

	isBannedByLogin, err := u.loginBuckets.IsBanned(login)
	if err != nil {
		return false, nil
	}

	isBannedByPassword, err := u.passwordBuckets.IsBanned(password)
	if err != nil {
		return false, nil
	}

	return !isBannedByIP && !isBannedByLogin && !isBannedByPassword, nil
}

func (u *App) DeleteLoginStats(login string) {
	u.loginBuckets.Forget(login)
}

func (u *App) DeleteIPStats(rawIP string) {
	u.ipBuckets.Forget(rawIP)
}

func (u *App) AddNetworkToWhiteList(rawNetwork string) error {
	return u.ipBuckets.AddWhiteListNetwork(rawNetwork)
}

func (u *App) AddNetworkToBlackList(rawNetwork string) error {
	return u.ipBuckets.AddBlackListNetwork(rawNetwork)
}

func (u *App) RemoveNetworkFromWhiteList(rawNetwork string) error {
	return u.ipBuckets.RemoveWhiteListNetwork(rawNetwork)
}

func (u *App) RemoveNetworkFromBlackList(rawNetwork string) error {
	return u.ipBuckets.RemoveBlackListNetwork(rawNetwork)
}

func garbageCleanerStart(
	ctx context.Context,
	cfg *configs.Config,
	loginBuckets *bucket.StringBuckets,
	passwordBuckets *bucket.StringBuckets, ipBuckets *bucket.IPBuckets,
) {
	ticker := time.NewTicker(cfg.CleanUpPeriod)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				loginBuckets.DeleteGarbage()
				passwordBuckets.DeleteGarbage()
				ipBuckets.DeleteGarbage()
			}
		}
	}()
}
