package flood_control

import (
	"context"
	"log"
	"task/internal/config"
	"task/internal/database"
	"task/internal/database/postgres"
	"time"
)

type FloodController struct {
	Db              database.CallChecker
	NSecondInterval int
	KCallLimit      int
}

func New(cfg *config.Config) FloodController {
	d, err := postgres.New(cfg.Name, cfg.User, cfg.Password, cfg.Host, cfg.Port)
	if err != nil {
		log.Fatalf("%s", err)
	}
	return FloodController{
		Db:              d,
		NSecondInterval: cfg.NSecondInterval,
		KCallLimit:      cfg.KCallLimit,
	}
}
func (f *FloodController) Check(ctx context.Context, userID int64) (bool, error) {
	ctxT, cancel := context.WithTimeout(ctx, time.Duration(f.NSecondInterval*1500000000))
	defer cancel()
	calls, err := f.Db.MakeCallAndGetRecent(ctxT, userID, f.NSecondInterval)
	if err != nil {
		return true, err
	}
	if calls >= f.KCallLimit {
		return true, nil
	}
	return false, nil

}
