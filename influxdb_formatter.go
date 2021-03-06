package logrus

import (
	"os"
    "bytes"
    "fmt"
    "time"
    "strings"
)

//default measurement name and host name.
var (
	DefaultMeasurement = "default_measure"
)

//InfluxdbFormat influxdb format
type InfluxdbFormat struct {
	Measurement string
	Host        string
	Precision   string // ns, us, ms, s
}

var _ Formatter = (*InfluxdbFormat)(nil)

//NewInfluxdbFormat new a influxdb format.
func NewInfluxdbFormat(m, h, p string) *InfluxdbFormat{
    if m == ""{
        m = DefaultMeasurement
    }
    if h == "" {
        h,_ = os.Hostname()
    }
    if p == "" {
        p = "ms"
    }
    f := &InfluxdbFormat{
        Measurement:m,
        Host: h,
        Precision:p,
    }
    return f
}

//Format implement Formatter interface.
func (f *InfluxdbFormat) Format(entry *Entry) ([]byte, error) {
	var b *bytes.Buffer
	var keys = make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	prefixFieldClashes(entry.Data)
    msg := strings.Replace(strings.Replace(entry.Message, "\n", " ", -1), `"`, `\"`, -1)
    //sl.logger.Printf(`user=%s,db=%s usetime=%d,node="%s",ip="%s",sql="%s" %d`, user, db, t, node, ip, sql, time.Now().UnixNano())
    b.WriteString(fmt.Sprintf(`%s,host=%s,`, f.Measurement, f.Host))
	b.WriteString(fmt.Sprintf(`user=%s,db=%s,node=%s,ip=%s usetime=%d,sql="%s" %d`,
         entry.Data["user"],
         entry.Data["db"],
         entry.Data["node"],
         entry.Data["ip"],
         entry.Data["usetime"],
         msg,
         time.Now().UnixNano(),
    ))
    b.WriteByte('\n')
	return b.Bytes(), nil
}
