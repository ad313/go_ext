package ext

import "time"

type timeExt struct {
}

var Time timeExt

func (a timeExt) ToDateString(unix int64, format ...string) string {
	if unix <= 0 {
		return ""
	}

	//毫秒
	if unix > 1000000000000 {
		unix = unix / 1000
	}

	var f = "2006-01-02"
	if len(format) > 0 {
		f = format[0]
	}

	return time.Unix(unix, 0).Format(f)
}

func (a timeExt) ToDateTimeString(unix int64, format ...string) string {

	if unix <= 0 {
		return ""
	}

	//毫秒
	if unix > 1000000000000 {
		unix = unix / 1000
	}

	var f = "2006-01-02 15:04:05"
	if len(format) > 0 {
		f = format[0]
	}

	return time.Unix(unix, 0).Format(f)
}

func (a timeExt) ToUnix(str string) {

}
