package mysql

import (
	"errors"
)

var (
	empty_sql_args SqlArgs = make(SqlArgs, 0)

	ErrParseEmpty = errors.New("parsed empty clause")
)

func (s *Statement) parse_select() (string, SqlArgs, error) {
	var (
		err  error
		args SqlArgs = make(SqlArgs, 0)

		buf []byte = make([]byte, Max_sql_len)
		off int    = 0
	)

	off += copy(buf[off:], "SELECT ")
	if err = multi_expr_to_string(s.SelectExprs, ", ", &args, &buf, &off); err != nil {
		return "", empty_sql_args, err
	}

	if s.FromExprs != nil {
		off += copy(buf[off:], " FROM ")
		if err = multi_expr_to_string(s.FromExprs, ", ", &args, &buf, &off); err != nil {
			return "", empty_sql_args, err
		}
	}

	if s.WhereExprs != nil {
		off += copy(buf[off:], " WHERE ")
		if err = multi_expr_to_string(s.WhereExprs, " AND ", &args, &buf, &off); err != nil {
			return "", empty_sql_args, err
		}
	}

	if s.GroupbyExprs != nil {
		off += copy(buf[off:], " GROUP BY ")
		if err = multi_expr_to_string(s.GroupbyExprs, ", ", &args, &buf, &off); err != nil {
			return "", empty_sql_args, err
		}
	}

	if s.HavingExprs != nil {
		off += copy(buf[off:], " HAVING ")
		if err = multi_expr_to_string(s.HavingExprs, ", ", &args, &buf, &off); err != nil {
			return "", empty_sql_args, err
		}
	}

	if s.OrderbyExprs != nil {
		off += copy(buf[off:], " ORDER BY ")
		if err = multi_expr_to_string(s.OrderbyExprs, ", ", &args, &buf, &off); err != nil {
			return "", empty_sql_args, err
		}
	}

	if s.LimitExpr != nil {
		off += copy(buf[off:], " LIMIT ")
		if err = expr_to_string(s.LimitExpr, &args, &buf, &off); err != nil {
			return "", empty_sql_args, err
		}
	}

	return string(buf[:off]), args, nil
}

func (s *Statement) parse_update() (string, SqlArgs, error) {
	var (
		err  error
		args SqlArgs = make(SqlArgs, 0)

		buf []byte = make([]byte, 16*1024)
		off int    = 0
	)

	off += copy(buf[off:], "UPDATE ")
	if err = multi_expr_to_string(s.UpdateExprs, ", ", &args, &buf, &off); err != nil {
		return "", empty_sql_args, err
	}

	if s.SetExprs != nil {
		off += copy(buf[off:], " SET ")
		if err = multi_expr_to_string(s.SetExprs, ", ", &args, &buf, &off); err != nil {
			return "", empty_sql_args, err
		}
	}

	if s.WhereExprs != nil {
		off += copy(buf[off:], " WHERE ")
		if err = multi_expr_to_string(s.WhereExprs, " AND ", &args, &buf, &off); err != nil {
			return "", empty_sql_args, err
		}
	}

	if s.OrderbyExprs != nil {
		off += copy(buf[off:], " ORDER BY ")
		if err = multi_expr_to_string(s.OrderbyExprs, ", ", &args, &buf, &off); err != nil {
			return "", empty_sql_args, err
		}
	}

	if s.LimitExpr != nil {
		off += copy(buf[off:], " LIMIT ")
		if err = expr_to_string(s.LimitExpr, &args, &buf, &off); err != nil {
			return "", empty_sql_args, err
		}
	}

	return string(buf[:off]), args, nil
}

func (s *Statement) parse_insert() (string, SqlArgs, error) {
	var (
		err  error
		args SqlArgs = make(SqlArgs, 0)

		buf []byte = make([]byte, 16*1024)
		off int    = 0
	)

	off += copy(buf[off:], "INSERT INTO ")
	if err = multi_expr_to_string(s.InsertExprs, ", ", &args, &buf, &off); err != nil {
		return "", empty_sql_args, err
	}

	if s.SetExprs != nil {
		off += copy(buf[off:], " SET ")
		if err = multi_expr_to_string(s.SetExprs, ", ", &args, &buf, &off); err != nil {
			return "", empty_sql_args, err
		}
	}

	if s.OndupExprs != nil {
		off += copy(buf[off:], " ON DUPLICATE KEY UPDATE ")
		if err = multi_expr_to_string(s.OndupExprs, ", ", &args, &buf, &off); err != nil {
			return "", empty_sql_args, err
		}
	}

	return string(buf[:off]), args, nil
}

func (s *Statement) parse_delete() (string, SqlArgs, error) {
	var (
		err  error
		args SqlArgs = make(SqlArgs, 0)

		buf []byte = make([]byte, 16*1024)
		off int    = 0
	)

	off += copy(buf[off:], "DELETE FROM ")
	if err = multi_expr_to_string(s.DeleteExprs, ", ", &args, &buf, &off); err != nil {
		return "", empty_sql_args, err
	}

	if s.WhereExprs != nil {
		off += copy(buf[off:], " WHERE ")
		if err = multi_expr_to_string(s.WhereExprs, " AND ", &args, &buf, &off); err != nil {
			return "", empty_sql_args, err
		}
	}

	if s.OrderbyExprs != nil {
		off += copy(buf[off:], " ORDER BY ")
		if err = multi_expr_to_string(s.OrderbyExprs, ", ", &args, &buf, &off); err != nil {
			return "", empty_sql_args, err
		}
	}

	if s.LimitExpr != nil {
		off += copy(buf[off:], " LIMIT ")
		if err = expr_to_string(s.LimitExpr, &args, &buf, &off); err != nil {
			return "", empty_sql_args, err
		}
	}

	return string(buf[:off]), args, nil
}
