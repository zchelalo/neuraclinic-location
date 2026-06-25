package postgresutil

import "github.com/jackc/pgx/v5/pgtype"

func TextPtr(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	v := value.String
	return &v
}
