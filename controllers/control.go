package controllers

import (
	"compress/gzip"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joncalhoun/viewcon"
	"github.com/julienschmidt/httprouter"
	"github.com/vadimtk/mview/views"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

type MetricsController struct {
	viewcon.Controller
}

var Metrics MetricsController

type ChartOverview struct {
	Id   string
	Sid  uint // service_id (1-system, 2-mysql)
	Name string
}

type MetricPoint struct {
	Ts  uint
	Avg float64
}

type Resp struct {
	Charts []*ChartOverview
}

func (self MetricsController) Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {

	// Pretend to lookup cats in some way.
	cats := []string{"heathcliff", "garfield"}

	// render the view
	return views.Metrics.Index.Render(w, cats)
}

func (self MetricsController) Browse(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {

	log.Println("Browse called")

	db, err := sql.Open("mysql", "root@/dev_o1")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	var (
		mid  int
		name string
		sid  uint
	)

	log.Println("Connected to MySQL")

	rows, err := db.Query("SELECT metric_id, name FROM metrics")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var charts []*ChartOverview

	for rows.Next() {
		err := rows.Scan(&mid, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(mid, name)
		if strings.HasPrefix(name, "mysql") {
			sid = 2
		} else {
			sid = 1
		}
		charts = append(charts, &ChartOverview{Id: strconv.Itoa(mid), Sid: sid, Name: name})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	p := Resp{Charts: charts}
	// render the view
	log.Println("Charts len:", len(p.Charts))

	t := template.New("browse.tmpl")
	t, err = t.ParseFiles("html/browse.tmpl")
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return t.Execute(w, p)
	//return views.Metrics.Browse.Render(w, p)
}

func (self MetricsController) Metric(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {

	log.Println("Metric called")

	db, err := sql.Open("mysql", "root@/dev_o1")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	var (
		ts  int64
		avg float64
	)

	log.Println("Connected to MySQL")

	rows, err := db.Query("SELECT UNIX_TIMESTAMP(ts) tts,avg FROM metric_data WHERE metric_id=" +
		ps.ByName("mid") + " AND service_id=" +
		ps.ByName("sid") + " AND instance_id=" +
		ps.ByName("hostid") + " AND ts > '2015-05-29 02:50:00'  ORDER BY ts ASC")
	//ps.ByName("hostid") + " AND ts > DATE_SUB(NOW(), INTERVAL 1 HOUR)  ORDER BY ts ASC")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	//var metrics []*MetricPoint
	str := ""

	for rows.Next() {
		err := rows.Scan(&ts, &avg)
		if err != nil {
			log.Fatal(err)
		}

		//log.Println(ts, avg)
		//metrics = append(metrics, &MetricPoint{Ts: ts, Avg: avg})
		if str != "" {
			str = str + ",\n"
		}
		str = str + fmt.Sprintf("[%s,%f]", strconv.FormatInt(ts*1000, 10), avg)
	}

	err = rows.Err()

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Encoding", "gzip")
	gz := gzip.NewWriter(w)
	defer gz.Close()
	_, err = gz.Write([]byte("[\n" + strings.TrimRight(str, ",") + "\n]"))
	return err
	//return views.Metrics.Browse.Render(w, p)
}

func (self MetricsController) Feed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	// render the view
	return views.Metrics.Feed.Render(w, nil)
}

func (self *MetricsController) ReqKey(a viewcon.Action) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if r.FormValue("key") != "topsecret" {
			http.Error(w, "Invalid key.", http.StatusUnauthorized)
		} else {
			self.Controller.Perform(a)(w, r, ps)
		}
	})
}
