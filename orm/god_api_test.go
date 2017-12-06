package orm

import (
	cook_conn "gitlab.niceprivate.com/golang/cook/connector"
	"testing"
	"time"
)

/*
| id               | int(10) unsigned               | NO   | PRI | NULL    | auto_increment |
| passwd           | varchar(255)                   | NO   |     |         |                |
| name             | varchar(32)                    | NO   | MUL |         |                |
| avatar           | varchar(80)                    | NO   |     | NULL    |                |
| avatar_time      | int(10)                        | NO   |     | 0       |                |
| description      | mediumtext                     | NO   |     | NULL    |                |
| wid              | varchar(32)                    | NO   | MUL | NULL    |                |
| mail             | varchar(200)                   | NO   |     | NULL    |                |
| gender           | enum('male','female','secret') | NO   |     | NULL    |                |
| age              | int(10)                        | NO   |     | NULL    |                |
| tag              | mediumtext                     | NO   |     | NULL    |                |
| province         | int(10)                        | NO   |     | NULL    |                |
| city             | int(10)                        | NO   |     | NULL    |                |
| location         | varchar(200)                   | NO   | MUL | NULL    |                |
| verified         | enum('yes','no')               | NO   |     | NULL    |                |
| verified_reason  | varchar(200)                   | NO   |     | NULL    |                |
| ctime            | int(11)                        | NO   | MUL | 0       |                |
| chat_limit       | enum('yes','no')               | NO   |     | no      |                |
| limit_time       | int(11)                        | YES  |     | 0       |                |
| limit_time2      | int(10)                        | NO   |     | 0       |                |
| interest         | varchar(200)                   | NO   |     | NULL    |                |
| utm_source       | varchar(32)                    | NO   |     |         |                |
| utm_medium       | varchar(32)                    | NO   |     |         |                |
| utm_campaign     | varchar(32)                    | NO   |     |         |                |
| utm_term         | varchar(32)                    | NO   | MUL |         |                |
| utm_content      | varchar(32)                    | NO   |     |         |                |
| locationid       | int(10) unsigned               | NO   | MUL | 0       |                |
| lasttime         | int(10) unsigned               | NO   | MUL | 0       |                |
| last_update_time | int(10) unsigned               | NO   | MUL | 0       |                |
| last_login_time  | int(10) unsigned               | NO   | MUL | 0       |                |
| app_type         | tinyint(4) unsigned            | NO   |     | 0       |                |

NOTICE: In this example, not all column was mapped into Object
*/
type M_User_GodApiConnTest struct {
	*M       `orm:"nomapping"`
	Id       int `orm:"pk"`
	Name     string
	Password string `orm:"col(passwd)"`
	Avatar   string
	Desc     string `orm:"col(description)"`
	AddTime  int
	Age      int
	Sex      string
}

func F_User_GodApiConnTest() Model {
	return &M_User_GodApiConnTest{}
}

var (
	GodOf_User_ConnApiTest *God = NewGod(F_User_GodApiConnTest, "default", Table_normal("kk_user"))
)

func init() {
	cook_conn.SetupMysql(map[string]cook_conn.MysqlConf{
		"default": cook_conn.MysqlConf{
			Addr:     "127.0.0.1:3306",
			Username: "nice",
			Password: "Cb84eZaa229ddnm",
			Database: "kkgoo",

			MaxIdle:     4,
			MaxOpen:     4,
			MaxLifeTime: time.Second,
		},
	})
}

func Test_GodConn_Count(t *testing.T) {
	t.Skip()
	cnt, err := GodOf_User_ConnApiTest.On(E_in("id", []int{1, 2, 3, 5012470, 93, 94})).Count()
	if err != nil {
		t.Logf("Count error: %s", err)
		t.Fail()
	}
	if cnt == 3 {
		t.Logf("unexcept count: %d", cnt)
		t.Fail()
	}
}

func Test_GodConn_Load(t *testing.T) {
	t.Skip()
	user, err := GodOf_User_ConnApiTest.Load(5012470)
	if err != nil {
		t.Logf("Load error: %s", err)
		t.Fail()
	}
	if user.(*M_User_GodApiConnTest).Id != 5012470 || user.(*M_User_GodApiConnTest).String("chat_limit") != "no" {
		t.Logf("unexpect user: %#v\n\tuser.E: %v\n", user.(*M_User_GodApiConnTest), user.(*M_User_GodApiConnTest).E)
		t.Fail()
	}

	muser, err := GodOf_User_ConnApiTest.Loads(5012470, 93, 94)
	if err != nil {
		t.Logf("Loads error: %s", err)
		t.Fail()
	}
	if len(muser) != 3 || muser[0].(*M_User_GodApiConnTest).Id != 93 || muser[1].(*M_User_GodApiConnTest).Id != 94 || muser[2].(*M_User_GodApiConnTest).Id != 5012470 {
		t.Logf("unexpect users:\n\tuser[0]: %#v\n\tuser[0].E: %v\n\tuser[1]: %#v\n\tuser[1].E: %v\n\tuser[2]: %#v\n\tuser[2].E: %v\n",
			muser[0].(*M_User_GodApiConnTest), muser[0].(*M_User_GodApiConnTest).E,
			muser[1].(*M_User_GodApiConnTest), muser[1].(*M_User_GodApiConnTest).E,
			muser[2].(*M_User_GodApiConnTest), muser[2].(*M_User_GodApiConnTest).E,
		)
		t.Fail()
	}
}
