package psql

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"pickfighter.com/auth/pkg/model"
	logs "pickfighter.com/pkg/logger"
)

const (
	searchUsersQuery = `SELECT u.user_id, u.name, u.claim, u.rank, u.flags, u.created_at, u.updated_at
	FROM public.pf_users AS u`
)

// FindUser searches for a user based on the provided UserRequest.
// It uses the SearchUsers method with a UsersRequest containing a single user ID.
// If a matching user is found, it returns the user details; otherwise, it returns pgx.ErrNoRows.
// In case of an inconsistency where more than one result is found, it returns an error.
func (r *Repository) FindUser(ctx context.Context, req *model.UserRequest) (*model.User, error) {
	results, err := r.SearchUsers(ctx, &model.UsersRequest{
		UserIds: []int32{req.UserId},
	})
	if err != nil {
		return nil, err
	}

	if len(results) == 1 {
		return results[0], nil
	} else if len(results) > 1 {
		return nil, fmt.Errorf("inconsistent: more than one result: %d", len(results))
	}

	return nil, pgx.ErrNoRows
}

// SearchUsers performs a search for users based on the provided UsersRequest.
// It constructs a dynamic SQL query using the parameters from the UsersRequest.
// The query is customized based on the filtering criteria such as user IDs, flags, and limits.
// If an error occurs during the query execution, it is logged along with the SQL query.
// The method returns a slice of User pointers representing the search results.
func (r *Repository) SearchUsers(ctx context.Context, req *model.UsersRequest) ([]*model.User, error) {
	q := searchUsersQuery

	args := r.PerformUsersRequestQuery(req)

	if len(args) > 0 {
		q += ` WHERE `
		q += strings.Join(args, sep)
	}

	q += ` ORDER BY u.user_id DESC`

	if req.Limit > 0 {
		q += fmt.Sprintf(` LIMIT %d OFFSET %d`, req.Limit, req.Offset)
	}

	rows, err := r.GetPool().Query(ctx, q)
	if err != nil {
		if err != pgx.ErrNoRows {
			logs.Debugf("query: \n%s", q)
		}
		return nil, err
	}
	defer rows.Close()

	var res []*model.User

	for rows.Next() {
		var u model.User
		var flags, updatedAt pgtype.Int8
		var rootClaim, rank pgtype.Varchar

		if err := rows.Scan(&u.UserId, &u.Name, &rootClaim, &rank, &flags, &u.CreatedAt, &updatedAt); err != nil {
			return nil, r.DebugLogSqlErr(q, err)
		}
		u.Rank = rank.String
		u.Claim = rootClaim.String
		u.Flags = uint64(flags.Int)
		u.UpdatedAt = updatedAt.Int

		res = append(res, &u)
	}

	return res, nil
}

// PerformUsersRequestQuery constructs a list of SQL query conditions based on the provided UsersRequest.
// It checks various filtering criteria, such as user IDs, names, emails, and creation timestamps,
// and creates corresponding SQL conditions for each valid criterion.
// The resulting list of conditions is intended to be used in a WHERE clause when querying user data.
// The sanitized strings are used to prevent SQL injection.
// The method returns a slice of strings representing the SQL conditions.
func (r *Repository) PerformUsersRequestQuery(req *model.UsersRequest) []string {
	var args []string
	if req == nil {
		return args
	}

	if len(req.UserIds) > 0 {
		if len(req.UserIds) == 1 {
			if req.UserIds[0] > 0 {
				args = append(args, fmt.Sprintf(`u.user_id = %d`, req.UserIds[0]))
			}
		} else {
			stringIds := ""
			for i := range req.UserIds {
				if req.UserIds[i] > 0 {
					if len(stringIds) > 0 {
						stringIds += ","
					}
					stringIds += fmt.Sprintf("%d", req.UserIds[i])
				}
			}
			if len(stringIds) > 0 {
				args = append(args, fmt.Sprintf(`u.user_id IN (%s)`, stringIds))
			}
		}
	}

	if len(req.Name) > 0 {
		args = append(args, fmt.Sprintf(`u.name ILIKE '%%%s%%'`, r.SanitizeString(req.Name)))
	}

	if len(req.Email) > 0 {
		args = append(args, fmt.Sprintf(`u.public_email ILIKE '%%%s%%'`, r.SanitizeString(req.Email)))
	}

	if req.CreatedFrom > 0 {
		args = append(args, fmt.Sprintf(`u.created_at > '%d'`, req.CreatedFrom))
	}

	if req.CreatedUntil > 0 {
		args = append(args, fmt.Sprintf(`u.created_at < '%d'`, req.CreatedUntil))
	}

	return args
}
