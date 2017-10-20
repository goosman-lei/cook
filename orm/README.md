# simple orm support

## No sharding demo

```
// user.go
package user

import (
    cook_sql "gitlab.niceprivate.com/golang/cook/sql/mysql"
    cook_orm "gitlab.niceprivate.com/golang/cook/orm"
)

var (
    GodOf_User *cook_orm.God = cook_orm.NewGod_shard_none(F_User, "default", "kk_user") // no sharding
)

type M_User struct {
    Id int
    Password string `orm:"col(passwd)"` // current, only support col command in orm tag
    Name string
    Avatar string
    AvatarTime string
    Description string
}

func F_User() interface{} {
    return &M_User{}
}

func Load(uid int) (*M_User, error) {
    user, err := GodOf_User.One(Select(
        E_field("*"),
    ).From(
        GodOf_User.Table(),
    ).Where(
        E_eq("id", uid),
    ).Limit(1))

    return user.(*M_User), err
}

func MultiLoad(uids []int) ([]*M_User, error) {
    users, err := GodOf_User.Multi(Select(
        E_field("*"),
    ).From(
        GodOf_User.Table(),
    ).Where(
        E_in("id", uids),
    ).Limit(1))

    return users.([]*M_User), err
}

```

## Sharding demo

```
// user.go
package user

import (
    cook_sql "gitlab.niceprivate.com/golang/cook/sql/mysql"
    cook_orm "gitlab.niceprivate.com/golang/cook/orm"
)

var (
    GodOf_User *cook_orm.God = cook_orm.NewGod_shard_mod_int(F_User, "default", "kk_user_%d", 128) // sharding with one int value
)

type M_User struct {
    Id int
    Password string `orm:"col(passwd)"` // current, only support col command in orm tag
    Name string
    Avatar string
    AvatarTime string
    Description string
}

func F_User() interface{} {
    return &M_User{}
}

func Load(uid int) (*M_User, error) {
    user, err := GodOf_User.One(Select(
        E_field("*"),
    ).From(
        GodOf_User.Table(uid),
    ).Where(
        E_eq("id", uid),
    ).Limit(1))

    return user.(*M_User), err
}

func MultiLoad(uids []int) ([]*M_User, error) {
    var (
        users []*M_User
        tables []*cook_sql.Expr
        table_args [][][]interface{}
        err error
    )

    if tables, table_args, err = GodOf_User.Tables(uids); err != nil {
        return nil, err
    }

    for i, table := range tables {
        if tmp_users, tmp_err := GodOf_User.Multi(Select(
            E_field("*"),
        ).From(table).Where(
            E_in("id", table_args[i][0]),
        ).Limit(1)); tmp_err != nil {
            return nil, tmp_err
        } else {
            users = append(users, tmp_users.([]*M_User)...)
        }

    }

    return users, nil
}

```
