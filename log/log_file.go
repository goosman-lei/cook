package log

import (
	"fmt"
	cook_util "gitlab.niceprivate.com/golang/cook/util"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

var (
	fname_format_regexp = regexp.MustCompile("\\*(?:{\\d+})?")

	time_edges = []int{9999, 12, 31, 24, 60}
	time_fmts  = []string{"%04d", "%02d", "%02d", "%02d", "%02d"}

	deadline_fmts = []string{"2006 Z07:00", "2006-01 Z07:00", "2006-01-02 Z07:00", "2006-01-02 15 Z07:00", "2006-01-02 15:04 Z07:00"}
)

func open_log_file(fname string) (*os.File, error) {
	if fp, err := os.OpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755); err == nil {
		// open file success
		return fp, nil
	} else if !cook_util.Err_NoSuchFileOrDir(err) {
		// open failure, but is not have no file or dir error
		return nil, err
	} else if err = os.MkdirAll(filepath.Dir(fname), 0755); err != nil {
		// mkdir failure
		return nil, err
	}

	return os.OpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
}

func parse_log_fname(fname string) (string, time.Time) {
	var (
		formated_fname string

		start_time      time.Time
		start_time_eles []int

		item_n         int = 0 // index of current item
		item_val       int     // current value of item
		last_item_lack int     // distance to next roundtrip of last item's value
	)

	start_time = time.Now()

	// year, month and day, start from 1, 2, ...
	// hour and minute start from 0, 1, ...
	start_time_eles = []int{start_time.Year() - 1, int(start_time.Month()) - 1, start_time.Day() - 1, start_time.Hour(), start_time.Minute()}

	formated_fname = fname_format_regexp.ReplaceAllStringFunc(fname, func(item string) (replacement string) {
		defer func() {
			if r := recover(); r != nil {
				replacement = item
			}
		}()
		if item_n >= len(time_edges) {
			// only replacement first 5th, skip other
			return item
		} else if len(item) == 1 {
			// have no interval
			item_val = start_time_eles[item_n]
			last_item_lack = 1
		} else {
			if interval_i64, err := strconv.ParseUint(item[2:len(item)-1], 10, 64); err != nil {
				item_val = start_time_eles[item_n]
				last_item_lack = 1
			} else {
				interval := int(interval_i64)
				if interval >= time_edges[item_n] {
					interval = time_edges[item_n]
				}
				item_val = int(math.Floor(float64(start_time_eles[item_n])/float64(interval))) * interval
				last_item_lack = item_val + interval - start_time_eles[item_n]
			}
		}
		if item_n < 3 {
			// year, month and day, start from 1, 2, ...
			replacement = fmt.Sprintf(time_fmts[item_n], item_val+1)
		} else {
			// hour and minute start from 0, 1, ...
			replacement = fmt.Sprintf(time_fmts[item_n], item_val)
		}
		item_n++
		return
	})
	return formated_fname, calc_deadline(start_time, last_item_lack, item_n-1)
}

func calc_deadline(start_time time.Time, last_item_lack int, item_n int) (deadline time.Time) {
	switch item_n {
	case -1:
		deadline = start_time.AddDate(9999, 0, 0).Local()
	case 0:
		deadline, _ = time.Parse(
			deadline_fmts[item_n],
			time.Now().AddDate(last_item_lack, 0, 0).Format(deadline_fmts[item_n]),
		)
		deadline = deadline.Local()
	case 1:
		deadline, _ = time.Parse(
			deadline_fmts[item_n],
			time.Now().AddDate(0, last_item_lack, 0).Format(deadline_fmts[item_n]),
		)
		deadline = deadline.Local()
	case 2:
		deadline, _ = time.Parse(
			deadline_fmts[item_n],
			time.Now().AddDate(0, 0, last_item_lack).Format(deadline_fmts[item_n]),
		)
		deadline = deadline.Local()
	case 3:
		deadline, _ = time.Parse(
			deadline_fmts[item_n],
			time.Now().Add(time.Hour*time.Duration(last_item_lack)).Format(deadline_fmts[item_n]),
		)
		deadline = deadline.Local()
	case 4:
		deadline, _ = time.Parse(
			deadline_fmts[item_n],
			time.Now().Add(time.Minute*time.Duration(last_item_lack)).Format(deadline_fmts[item_n]),
		)
		deadline = deadline.Local()
	}
	return
}
