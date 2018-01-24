package json

import (
	"fmt"
	"gitlab.niceprivate.com/golang/cook/orm"
	"reflect"
	"testing"
)

type M_Live struct {
	*orm.M            `json:"-"`
	Id                uint64  `orm:"pk" json:"id"`
	Uid               string  `json:"uid"`
	Pic               string  `json:"pic"`
	Content           string  `json:"content"`
	Status            string  `json:"status"`
	AddTime           int     `json:"add_time"`
	StartTime         int     `json:"start_time"`
	EndTime           int     `json:"end_time"`
	ForeshowTime      int     `json:"foreshow_time"`
	TemplateId        int     `json:"template_id"`
	Latitude          string  `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	LikeNum           int     `json:"like_num"`
	AudienceNum       int     `json:"audience_num"`
	AudienceOlTopNum  int     `json:"audience_ol_top_num"`
	ServiceType       string  `json:"service_type"`
	ServiceLiveid     string  `json:"service_liveid"`
	Group             int     `json:"group"`
	Visibility        int     `json:"visibility"`
	Location          string  `json:"location"`
	LocationId        int     `json:"location_id"`
	Province          string  `json:"province"`
	SupportFeedReplay int     `json:"support_feed_replay"`
	EncodeType        string  `json:"encode_type"`
}

func Test_marshal_obj(t *testing.T) {
	AutoDequoteOn()
	HumpToUnderlineOn()
	defer AutoDequoteOff()
	defer HumpToUnderlineOff()
	m := M_Live{
		Id:                196182630970425847,
		Uid:               "3496138",
		Pic:               "/upload/show/2017/12/18/ef4048af7c6a9b731786548496818422",
		Content:           "E-Pa-Ra的直播",
		Status:            "noreplay",
		AddTime:           1513588768,
		StartTime:         1513588773,
		EndTime:           1513589119,
		ForeshowTime:      0,
		TemplateId:        0,
		Latitude:          "22578967.6581",
		Longitude:         108231062.2061,
		LikeNum:           41,
		AudienceNum:       2,
		AudienceOlTopNum:  2,
		ServiceType:       "ws",
		ServiceLiveid:     "196182630970425847",
		Visibility:        1,
		Group:             8,
		SupportFeedReplay: 1,
		Location:          "南宁",
		LocationId:        186,
		Province:          "广西",
		EncodeType:        "h264",
	}
	ref_j := `{"id":"196182630970425847","uid":"3496138","pic":"\/upload\/show\/2017\/12\/18\/ef4048af7c6a9b731786548496818422","content":"E-Pa-Ra\u7684\u76f4\u64ad","status":"noreplay","add_time":"1513588768","start_time":"1513588773","end_time":"1513589119","foreshow_time":"0","template_id":"0","latitude":22578967.6581,"longitude":"108231062.2061","like_num":"41","audience_num":"2","audience_ol_top_num":"2","service_type":"ws","service_liveid":"196182630970425847","visibility":"1","group":"8","support_feed_replay":"1","location":"\u5357\u5b81","location_id":"186","province":"\u5e7f\u897f","encode_type":"h264"}`
	if j, err := Marshal_string(&m); err != nil {
		t.Logf("marshale error: %s", err)
		t.Logf("ref_json: %s", ref_j)
		t.Logf("marshal result: %s", j)
		t.Fail()
	}
}

func Test_unmarshal_into_obj(t *testing.T) {
	AutoDequoteOn()
	HumpToUnderlineOn()
	defer AutoDequoteOff()
	defer HumpToUnderlineOff()
	j := `{"id":"196182630970425847","uid":3496138,"pic":"\/upload\/show\/2017\/12\/18\/ef4048af7c6a9b731786548496818422","content":"E-Pa-Ra\u7684\u76f4\u64ad","status":"noreplay","add_time":"1513588768","start_time":"1513588773","end_time":"1513589119","foreshow_time":"0","template_id":"0","latitude":22578967.6581,"longitude":"108231062.2061","like_num":"41","audience_num":"2","audience_ol_top_num":"2","service_type":"ws","service_liveid":"196182630970425847","visibility":"1","group":"8","support_feed_replay":"1","location":"\u5357\u5b81","location_id":"186","province":"\u5e7f\u897f","encode_type":"h264"}`
	ref_m := M_Live{
		Id:                196182630970425847,
		Uid:               "3496138",
		Pic:               "/upload/show/2017/12/18/ef4048af7c6a9b731786548496818422",
		Content:           "E-Pa-Ra的直播",
		Status:            "noreplay",
		AddTime:           1513588768,
		StartTime:         1513588773,
		EndTime:           1513589119,
		ForeshowTime:      0,
		TemplateId:        0,
		Latitude:          "22578967.6581",
		Longitude:         108231062.2061,
		LikeNum:           41,
		AudienceNum:       2,
		AudienceOlTopNum:  2,
		ServiceType:       "ws",
		ServiceLiveid:     "196182630970425847",
		Visibility:        1,
		Group:             8,
		SupportFeedReplay: 1,
		Location:          "南宁",
		LocationId:        186,
		Province:          "广西",
		EncodeType:        "h264",
	}
	m := M_Live{}
	err := Unmarshal_string(j, &m)
	if err != nil {
		t.Logf("unmarshal error: %s", err)
		t.Fail()
	}
	if m.Id != 196182630970425847 || m.Uid != "3496138" || m.Pic != "/upload/show/2017/12/18/ef4048af7c6a9b731786548496818422" || m.Latitude != "22578967.6581" || m.Province != "广西" {
		err_data := ""
		for i := 0; i < reflect.TypeOf(m).NumField(); i++ {
			rt_f := reflect.TypeOf(m).Field(i)
			r_f := reflect.ValueOf(m).Field(i)
			err_data = fmt.Sprintf("%s\n\t%20s: %-60v\t(except: %v)", err_data, rt_f.Name+"("+rt_f.Type.Kind().String()+")", r_f.Interface(), reflect.ValueOf(ref_m).Field(i).Interface())
		}
		t.Logf("unmarshal data error:%s", err_data)
		t.Fail()
	}
}
