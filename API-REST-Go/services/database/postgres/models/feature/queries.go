package feature

import (
	"API-REST/services/database/postgres/predicates"
	"errors"
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
)

func (m *Model) GetAll(p *predicates.Predicates) ([]*Feature, error) {
	query := m.Db.Table("features").Select(
		"GeomFromEWKB(geom) as geom",
		"timestamp",
		"user_id",
	)

	if p != nil {
		query = predicates.Apply(query, p)
	}

	res, err := query.Get()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("features not found")
	}

	var features []*Feature
	for _, r := range res {
		point := postgis.PointS{}
		point.Scan([]byte(r["geom"].(string)))
		p := Feature{
			Geom:      &point,
			Timestamp: r["timestamp"].(time.Time),
			UserID:    uuid.MustParse(r["user_id"].(string)),
		}

		features = append(features, &p)
	}

	return features, nil
}
func (m *Model) GetAllMostRecent() ([]*Feature, error) {
	rows, err := m.Db.Sql().Query(`
	SELECT GeomFromEWKB(t.geom), t.timestamp, t.user_id
	FROM features AS t
	INNER JOIN (
		SELECT user_id, MAX(timestamp) AS max_timestamp
		FROM features
		GROUP BY user_id
	) AS t2 ON t.user_id = t2.user_id AND t.timestamp = t2.max_timestamp;
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var features []*Feature
	for rows.Next() {
		var feature Feature
		feature.Geom = &postgis.PointS{}
		err = rows.Scan(feature.Geom, &feature.Timestamp, &feature.UserID)
		if err != nil {
			return nil, err
		}
		features = append(features, &feature)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return features, nil
}
func (m *Model) GetMostRecentByUserID(userID uuid.UUID) (*Feature, error) {
	res, err := m.Db.Table("features").Select(
		"GeomFromEWKB(geom) as geom",
		"timestamp",
	).Where("user_id", "=", userID.String()).
		OrderBy("timestamp", "DESC").
		Limit(1).
		Get()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("feature not found")
	}

	r := res[0]
	point := postgis.PointS{}
	point.Scan([]byte(r["geom"].(string)))
	feature := Feature{
		Geom:      &point,
		Timestamp: r["timestamp"].(time.Time),
	}

	return &feature, nil
}
func (m *Model) Insert(f *Feature) error {
	_, err := m.Db.Sql().Exec(`
	INSERT INTO features(geom, timestamp, user_id)
	VALUES(GeomFromEWKB($1), $2, $3);
    `, *f.Geom, f.Timestamp, f.UserID.String())
	if err != nil {
		return err
	}
	return nil
}
