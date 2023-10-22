package parser

import (
	"context"
	"tenderhack-parser/internal/models"
	"time"
)

type IntervalLogs struct {
	Count  int           `json:"count"`
	Logs   []*models.Log `json:"logs"`
	TSUnix int64         `json:"-"`
	TS     string        `json:"ts"`
}

func (s *Service) GetLogsByDate(ctx context.Context, from time.Time, to time.Time, resolved bool, interval uint32) ([]*IntervalLogs, error) {
	var (
		logs []*models.Log
		err  error
	)

	if resolved {
		logs, err = s.logs.GetResolved(ctx, from, to)
	} else {
		logs, err = s.logs.GetUnresolved(ctx, from, to)
	}

	if err != nil {
		return nil, err
	}

	res := []*IntervalLogs{}
	for _, l := range logs {
		rounded := l.TS.Truncate(time.Second * time.Duration(interval))
		roundedUnix := rounded.Unix()

		found := false
		for i := range res {
			if res[i].TSUnix == roundedUnix {
				found = true
				res[i].Count++
				res[i].Logs = append(res[i].Logs, l)
				break
			}
		}

		if !found {
			res = append(res, &IntervalLogs{
				Count: 1,
				Logs: []*models.Log{
					l,
				},
				TS:     rounded.String(),
				TSUnix: roundedUnix,
			})
		}
	}

	return res, nil
}
