package mysql

import (
	"testing"
)

func Benchmark_parse_select(b *testing.B) {
	q := Select(
		E_field("u.id").Alias("uid"),
		E_field("u.name").Alias("uname"),
		E_field("s.id").Alias("sid"),
		E_field("s.pic").Alias("pic_url"),
		E_field("s.add_time").Alias("publish_time"),
	).From(
		E_table("kk_user").Alias("u").Join(
			E_table("kk_user_show").Alias("s"),
			E_literal("u.id = s.uid"),
		),
	).Where(
		E_in("u.id", []int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
		E_in("s.status", []string{"wait", "hide"}),
	).Orderby(E_field("s.id").Desc()).Limit(10, 0)

	for i := 0; i < b.N; i++ {
		q.Parse()
	}
}

func Benchmark_parse_update(b *testing.B) {
	q := Update(E_table("kk_user")).Set(
		E_literal("limit_time = limit_time + 1"),
		E_literal("last_update_time = unix_timestamp(now())"),
	).Where(
		E_lt("id", 100),
	).Orderby(
		E_field("id").Desc(),
	).Limit(10)

	for i := 0; i < b.N; i++ {
		q.Parse()
	}
}
