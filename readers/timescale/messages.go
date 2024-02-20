// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package timescale

import (
	"encoding/json"
	"fmt"

	"github.com/absmach/magistrala/pkg/errors"
	"github.com/absmach/magistrala/pkg/transformers/senml"
	"github.com/absmach/magistrala/readers"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx" // required for DB access
)

var _ readers.MessageRepository = (*timescaleRepository)(nil)

type timescaleRepository struct {
	db *sqlx.DB
}

// New returns new TimescaleSQL writer.
func New(db *sqlx.DB) readers.MessageRepository {
	return &timescaleRepository{
		db: db,
	}
}

func (tr timescaleRepository) ReadAll(chanID string, rpm readers.PageMetadata) (readers.MessagesPage, error) {
	order := "time"
	format := defTable

	if rpm.Format != "" && rpm.Format != defTable {
		order = "created"
		format = rpm.Format
	}

	q := fmt.Sprintf(`SELECT * FROM %s WHERE %s ORDER BY %s DESC LIMIT :limit OFFSET :offset;`, format, fmtCondition(chanID, rpm), order)

	params := map[string]interface{}{
		"channel":      chanID,
		"limit":        rpm.Limit,
		"offset":       rpm.Offset,
		"subtopic":     rpm.Subtopic,
		"publisher":    rpm.Publisher,
		"name":         rpm.Name,
		"protocol":     rpm.Protocol,
		"value":        rpm.Value,
		"bool_value":   rpm.BoolValue,
		"string_value": rpm.StringValue,
		"data_value":   rpm.DataValue,
		"from":         rpm.From,
		"to":           rpm.To,
	}

	rows, err := tr.db.NamedQuery(q, params)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.UndefinedTable {
				return readers.MessagesPage{}, nil
			}
		}
		return readers.MessagesPage{}, errors.Wrap(readers.ErrReadMessages, err)
	}
	defer rows.Close()

	page := readers.MessagesPage{
		PageMetadata: rpm,
		Messages:     []readers.Message{},
	}
	switch format {
	case defTable:
		for rows.Next() {
			msg := senmlMessage{Message: senml.Message{}}
			if err := rows.StructScan(&msg); err != nil {
				return readers.MessagesPage{}, errors.Wrap(readers.ErrReadMessages, err)
			}

			page.Messages = append(page.Messages, msg.Message)
		}
	default:
		for rows.Next() {
			msg := jsonMessage{}
			if err := rows.StructScan(&msg); err != nil {
				return readers.MessagesPage{}, errors.Wrap(readers.ErrReadMessages, err)
			}
			m, err := msg.toMap()
			if err != nil {
				return readers.MessagesPage{}, errors.Wrap(readers.ErrReadMessages, err)
			}
			page.Messages = append(page.Messages, m)
		}
	}

	q = fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s;`, format, fmtCondition(chanID, rpm))
	rows, err = tr.db.NamedQuery(q, params)
	if err != nil {
		return readers.MessagesPage{}, errors.Wrap(readers.ErrReadMessages, err)
	}
	defer rows.Close()

	total := uint64(0)
	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return page, err
		}
	}
	page.Total = total

	q = fmt.Sprintf(`SELECT time_bucket('1 hour', time) AS bucket, SUM(*) FROM %s WHERE %s GROUP BY bucket ORDER BY bucket ASC;`, format, fmtCondition(chanID, rpm))
	rows, err = tr.db.NamedQuery(q, params)
	if err != nil {
		return readers.MessagesPage{}, errors.Wrap(readers.ErrReadMessages, err)
	}
	defer rows.Close()

	var buckets []readers.TimeBucket
	for rows.Next() {
		var bucket readers.TimeBucket
		if err := rows.Scan(&bucket.Time, &bucket.Total); err != nil {
			return page,err
		}
		buckets = append(buckets, bucket)
	}
	page.Buckets = buckets

	sum := float64(0)
	if rows.Next() {
		if err := rows.Scan(&sum); err != nil {
			return page, err
		}
	}
	page.Sum = sum

	q = fmt.Sprintf(`SELECT time_bucket('1 hour', time) AS bucket, AVG(*) FROM %s WHERE %s GROUP BY bucket ORDER BY bucket ASC;`, format, fmtCondition(chanID, rpm))
	rows, err = tr.db.NamedQuery(q, params)
	if err != nil {
		return readers.MessagesPage{}, errors.Wrap(readers.ErrReadMessages, err)
	}
	defer rows.Close()

	avg := float64(0)
	if rows.Next() {
		if err := rows.Scan(&avg); err != nil {
			return page, err
		}
	}
	page.Avg = avg

	for rows.Next() {
		var bucket readers.TimeBucket
		if err := rows.Scan(&bucket.Time, &bucket.Total); err != nil {
			return page,err
		}
		buckets = append(buckets, bucket)
	}
	page.Buckets = buckets


	q = fmt.Sprintf(`SELECT time_bucket('1 hour', time) AS bucket, MAX(*) FROM %s WHERE %s GROUP BY bucket ORDER BY bucket ASC;`, format, fmtCondition(chanID, rpm))
	rows, err = tr.db.NamedQuery(q, params)
	if err != nil {
		return readers.MessagesPage{}, errors.Wrap(readers.ErrReadMessages, err)
	}
	defer rows.Close()

	max := float64(0)
	if rows.Next() {
		if err := rows.Scan(&max); err != nil {
			return page, err
		}
	}
	page.Max = max

	for rows.Next() {
		var bucket readers.TimeBucket
		if err := rows.Scan(&bucket.Time, &bucket.Total); err != nil {
			return page,err
		}
		buckets = append(buckets, bucket)
	}
	page.Buckets = buckets

	q = fmt.Sprintf(`SELECT time_bucket('1 hour', time) AS bucket, MIN(*) FROM %s WHERE %s GROUP BY bucket ORDER BY bucket ASC;`, format, fmtCondition(chanID, rpm))
	rows, err = tr.db.NamedQuery(q, params)
	if err != nil {
		return readers.MessagesPage{}, errors.Wrap(readers.ErrReadMessages, err)
	}
	defer rows.Close()

	min := float64(0)
	if rows.Next() {
		if err := rows.Scan(&min); err != nil {
			return page, err
		}
	}
	page.Min = min

	for rows.Next() {
		var bucket readers.TimeBucket
		if err := rows.Scan(&bucket.Time, &bucket.Total); err != nil {
			return page,err
		}
		buckets = append(buckets, bucket)
	}
	page.Buckets = buckets

	return page, nil
}

func fmtCondition(chanID string, rpm readers.PageMetadata) string {
	condition := `channel = :channel`

	var query map[string]interface{}
	meta, err := json.Marshal(rpm)
	if err != nil {
		return condition
	}
	if err := json.Unmarshal(meta, &query); err != nil {
		return condition
	}

	for name := range query {
		switch name {
		case
			"subtopic",
			"publisher",
			"name",
			"protocol":
			condition = fmt.Sprintf(`%s AND %s = :%s`, condition, name, name)
		case "v":
			comparator := readers.ParseValueComparator(query)
			condition = fmt.Sprintf(`%s AND value %s :value`, condition, comparator)
		case "vb":
			condition = fmt.Sprintf(`%s AND bool_value = :bool_value`, condition)
		case "vs":
			comparator := readers.ParseValueComparator(query)
			switch comparator {
			case "=":
				condition = fmt.Sprintf("%s AND string_value = :string_value ", condition)
			case ">":
				condition = fmt.Sprintf("%s AND string_value LIKE '%%' || :string_value || '%%' AND string_value <> :string_value", condition)
			case ">=":
				condition = fmt.Sprintf("%s AND string_value LIKE '%%' || :string_value || '%%'", condition)
			case "<=":
				condition = fmt.Sprintf("%s AND :string_value LIKE '%%' || string_value || '%%'", condition)
			case "<":
				condition = fmt.Sprintf("%s AND :string_value LIKE '%%' || string_value || '%%' AND string_value <> :string_value", condition)
			}
		case "vd":
			comparator := readers.ParseValueComparator(query)
			condition = fmt.Sprintf(`%s AND data_value %s :data_value`, condition, comparator)
		case "from":
			condition = fmt.Sprintf(`%s AND time >= :from`, condition)
		case "to":
			condition = fmt.Sprintf(`%s AND time < :to`, condition)
		}
	}
	return condition
}

type senmlMessage struct {
	ID string `db:"id"`
	senml.Message
}

type jsonMessage struct {
	Channel   string `db:"channel"`
	Created   int64  `db:"created"`
	Subtopic  string `db:"subtopic"`
	Publisher string `db:"publisher"`
	Protocol  string `db:"protocol"`
	Payload   []byte `db:"payload"`
}

func (msg jsonMessage) toMap() (map[string]interface{}, error) {
	ret := map[string]interface{}{
		"channel":   msg.Channel,
		"created":   msg.Created,
		"subtopic":  msg.Subtopic,
		"publisher": msg.Publisher,
		"protocol":  msg.Protocol,
		"payload":   map[string]interface{}{},
	}
	pld := make(map[string]interface{})
	if err := json.Unmarshal(msg.Payload, &pld); err != nil {
		return nil, err
	}
	ret["payload"] = pld
	return ret, nil
}

