package database

import "context"

type CallChecker interface {
	MakeCallAndGetRecent(ctx context.Context, userId int64, n int) (int, error)
}
