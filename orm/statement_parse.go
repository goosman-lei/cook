package orm

import (
	"errors"
)

const (
	MAX_SQL_LEN = 16 * 1024
)

func (s *Statement) parse_select() error {
	var (
		err  error
		args SqlArgs = make(SqlArgs, 0)

		buf []byte = make([]byte, MAX_SQL_LEN)
		off int    = 0
	)

	off += copy(buf[off:], "SELECT ")
	if err = multi_expr_to_string(s.SelectClause, ", ", &args, &buf, &off); err != nil {
		goto Failure
	}

	if s.TableClause != nil {
		off += copy(buf[off:], " FROM ")
		if err = multi_expr_to_string(s.TableClause, ", ", &args, &buf, &off); err != nil {
			goto Failure
		}
	}

	if s.OnClause != nil {
		off += copy(buf[off:], " WHERE ")
		if err = multi_expr_to_string(s.OnClause, " AND ", &args, &buf, &off); err != nil {
			goto Failure
		}
	}

	if s.GroupbyClause != nil {
		off += copy(buf[off:], " GROUP BY ")
		if err = multi_expr_to_string(s.GroupbyClause, ", ", &args, &buf, &off); err != nil {
			goto Failure
		}
	}

	if s.HavingClause != nil {
		off += copy(buf[off:], " HAVING ")
		if err = multi_expr_to_string(s.HavingClause, ", ", &args, &buf, &off); err != nil {
			goto Failure
		}
	}

	if s.OrderbyClause != nil {
		off += copy(buf[off:], " ORDER BY ")
		if err = multi_expr_to_string(s.OrderbyClause, ", ", &args, &buf, &off); err != nil {
			goto Failure
		}
	}

	if s.LimitClause != nil {
		off += copy(buf[off:], " LIMIT ")
		if err = expr_to_string(s.LimitClause, &args, &buf, &off); err != nil {
			goto Failure
		}
	}

	s.Query = string(buf[:off])
	s.Args = args
	return nil
Failure:
	s.Error = err
	return err
}

func (s *Statement) parse_update() error {
	var (
		err  error
		args SqlArgs = make(SqlArgs, 0)

		buf []byte = make([]byte, MAX_SQL_LEN)
		off int    = 0
	)

	off += copy(buf[off:], "UPDATE ")
	if err = multi_expr_to_string(s.TableClause, ", ", &args, &buf, &off); err != nil {
		goto Failure
	}

	off += copy(buf[off:], " SET ")
	if err = multi_expr_to_string(s.UpdateClause, ", ", &args, &buf, &off); err != nil {
		goto Failure
	}

	if s.OnClause != nil {
		off += copy(buf[off:], " WHERE ")
		if err = multi_expr_to_string(s.OnClause, " AND ", &args, &buf, &off); err != nil {
			goto Failure
		}
	}

	if s.OrderbyClause != nil {
		off += copy(buf[off:], " ORDER BY ")
		if err = multi_expr_to_string(s.OrderbyClause, ", ", &args, &buf, &off); err != nil {
			goto Failure
		}
	}

	if s.LimitClause != nil {
		off += copy(buf[off:], " LIMIT ")
		if err = expr_to_string(s.LimitClause, &args, &buf, &off); err != nil {
			goto Failure
		}
	}

	s.Query = string(buf[:off])
	s.Args = args
	return nil
Failure:
	s.Error = err
	return err

}

func (s *Statement) parse_insert() error {
	var (
		err  error
		args SqlArgs = make(SqlArgs, 0)

		buf []byte = make([]byte, MAX_SQL_LEN)
		off int    = 0
	)

	off += copy(buf[off:], "INSERT INTO ")
	if err = multi_expr_to_string(s.TableClause, ", ", &args, &buf, &off); err != nil {
		goto Failure
	}

	if s.InsertClause.values != nil {
		if len(s.InsertClause.cols) > 0 {
			(buf)[off] = '('
			off++
			is_first := true
			for _, col := range s.InsertClause.cols {
				if is_first {
					is_first = false
					off += copy(buf[off:], col)
				} else {
					off += copy(buf[off:], ", "+col)
				}
			}
			(buf)[off] = ')'
			off++
		}

		off += copy(buf[off:], " VALUES ")
		is_first := true
		for _, expr := range s.InsertClause.values {
			if is_first {
				is_first = false
			} else {
				(buf)[off] = ','
				off++
				(buf)[off] = ' '
				off++
			}
			if err = expr_to_string(expr, &args, &buf, &off); err != nil {
				goto Failure
			}
		}
	} else if s.InsertClause.sets != nil {
		off += copy(buf[off:], " SET ")
		if err = multi_expr_to_string(s.InsertClause.sets, ", ", &args, &buf, &off); err != nil {
			goto Failure
		}

		if s.OndupClause != nil {
			off += copy(buf[off:], " ON DUPLICATE KEY UPDATE ")
			if err = multi_expr_to_string(s.OndupClause, ", ", &args, &buf, &off); err != nil {
				goto Failure
			}
		}
	}

	s.Query = string(buf[:off])
	s.Args = args
	return nil
Failure:
	s.Error = err
	return err
}

func (s *Statement) parse_delete() error {
	var (
		err  error
		args SqlArgs = make(SqlArgs, 0)

		buf []byte = make([]byte, MAX_SQL_LEN)
		off int    = 0
	)

	off += copy(buf[off:], "DELETE FROM ")
	if err = multi_expr_to_string(s.TableClause, ", ", &args, &buf, &off); err != nil {
		goto Failure
	}

	if s.DeleteClause != nil {
		off += copy(buf[off:], " WHERE ")
		if err = multi_expr_to_string(s.DeleteClause, " AND ", &args, &buf, &off); err != nil {
			goto Failure
		}
	}

	if s.OrderbyClause != nil {
		off += copy(buf[off:], " ORDER BY ")
		if err = multi_expr_to_string(s.OrderbyClause, ", ", &args, &buf, &off); err != nil {
			goto Failure
		}
	}

	if s.LimitClause != nil {
		off += copy(buf[off:], " LIMIT ")
		if err = expr_to_string(s.LimitClause, &args, &buf, &off); err != nil {
			goto Failure
		}
	}

	s.Query = string(buf[:off])
	s.Args = args
	return nil
Failure:
	s.Error = err
	return err
}
