package main

import (
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"os"
	"strings"
	"encoding/json"
	"encoding/xml"
	"io"
	"strconv"
	"flag"
. "rdb"
    . "gorilla/grll"
)

const (

	xslHeaderRdb    = `<?xml-stylesheet type="text/xsl" href="/rdb.xsl"?>` + "\n"
	xslHeaderGringo = `<?xml-stylesheet type="text/xsl" href="/grll.xsl"?>` + "\n"

	sqlGetTagsByTestRunId =
		`select tag.name from testrun
        left join testrun_tag
                on testrun.id = testrun_tag.testrun
         left join tag
                on tag.id = testrun_tag.tag
		where testrun.id = $1`

	sqlGetResultsByTestRunId =
		`SELECT test.name, status.name, COALESCE(result.message, '') as message
        FROM testrun
        INNER JOIN result
                ON result.testrun = testrun.id
        LEFT JOIN test
                ON test.id = result.test
        LEFT JOIN status
                ON result.status = status.id
        WHERE testrun.id = $1`

	sqlGetTestRuns      = `SELECT id, name, ts FROM testrun ORDER BY ts`

	sqlGetTestRunById   = `SELECT name FROM testrun WHERE id = $1`

	sqlPutTestRun       = "SELECT * FROM new_testrun($1, $2)"

	sqlPutResult        = "SELECT * FROM new_result($1, $2, NULL, $3)"

	sqlPutResultMessage = "SELECT * FROM new_result($1, $2, $4, $3)"

	sqlSetTag           = "SELECT set_tag($1, $2)"

	sqlGetTags          = "SELECT id, name FROM tag ORDER BY name"

	sqlGetTestRunsByTag =
		`SELECT id, name, ts FROM get_testruns_by_tag($1) ORDER BY ts`
)

var db *sql.DB

func init() {
	var err error
	//db, err := sql.Open("postgres", "user=postgres password=sparc1e dbname=gorilla sslmode=disable")
	//db, err = sql.Open("postgres", "postgres://postgres:sparc1e@localhost/gorilla")
	db, err = sql.Open("postgres", "postgres://postgres:sparc1e@localhost/gorilla?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func debug(msg string) {
	log.Println(msg)
}

func rdbToGrll(rdb *RDBTestSuite) *GrllTestRun {
	g := NewGrllTestRun()
	if "" != rdb.Suite.Name {
		g.Run = rdb.Suite.Name + " " + rdb.Config.Name
	} else {
		g.Run = rdb.Config.Name
	}
	g.Timestamp = rdb.Timestamp
	for _, t := range rdb.TestList.TestCases {
		g.Results = append(g.Results,
				GrllResult{Test:t.Name,
						Status:t.Status,
						Message:t.Message})
	}
	for _, v := range rdb.Config.Props.Properties {
		g.Tags = append(g.Tags, fmt.Sprintf("%s=%s", v.Key, v.Value))
	}
	return g
}


func getTestRuns() []GrllTestRun {
	list := make([]GrllTestRun,0)
	rows, err := db.Query(sqlGetTestRuns)
	if err != nil {
		return list
	}
	defer rows.Close()
	for rows.Next() {
		i := NewGrllTestRun()
		err := rows.Scan(&i.Id, &i.Run, &i.Timestamp)
		if err != nil {
			return list
		}
		list = append(list, *i)
	}
	return list
}

func newTestRun(n string, ts string) int {
	var id int
	// this is specific to Postgres (not having LastInsertId)
	err := db.QueryRow(sqlPutTestRun, n, ts).Scan(&id)
	if err!=nil {
		log.Println(err.Error())
	}
	return id
}

func setTag(runId int, tag string) {
	_, err := db.Exec(sqlSetTag, runId, tag)
	if err!=nil {
		log.Println(err.Error())
	}
}

func newResult(result *GrllResult, runId int) {
	var err error
	// TODO: there should be a better way of mapping '' -> NULL
	if "" != result.Message {
		_, err = db.Exec(sqlPutResultMessage, result.Test, result.Status, runId, result.Message)
	} else {
		_, err = db.Exec(sqlPutResult, result.Test, result.Status, runId)
	}
	if err != nil {
		log.Println(err.Error())
	}
}

func getTagsByRunId(runId int) []string {
	tags := make([]string,0)
	rows, err := db.Query(sqlGetTagsByTestRunId, runId)
	if err != nil {
		log.Println(err.Error())
		return tags
	}
	defer rows.Close()
	for rows.Next() {
		var t string
		err := rows.Scan(&t)
		if err != nil {
			log.Println("Error: ", err.Error())
			return tags
		}
		tags = append(tags, t)
	}
	return tags
}

func getResultsByTestRunId(runId int) []GrllResult {
	list := make([]GrllResult,0)
	rows, err := db.Query(sqlGetResultsByTestRunId, runId)
	if err != nil {
		log.Println("Error: ", err.Error())
		return list
	}
	defer rows.Close()
	for rows.Next() {
		i := new(GrllResult)
		err := rows.Scan(&i.Test, &i.Status, &i.Message)
		if err != nil {
			log.Println("Error: ", err.Error())
			return list
		}
		list = append(list, *i)
	}
	return list
}

func getTestRunById(runId int) *GrllTestRun {
	g := NewGrllTestRun()
	err := db.QueryRow(sqlGetTestRunById, runId).Scan(&g.Run)
	if err!=nil {
		log.Println(err.Error())
	}
	g.Results = getResultsByTestRunId(runId)
	g.Tags = getTagsByRunId(runId)
	return g
}

func getTags() []string {
	tags := make([]string,0)
	rows, err := db.Query(sqlGetTags)
	if err != nil {
		log.Println(err.Error())
		return tags
	}
	defer rows.Close()
	for rows.Next() {
		var t string
		var i int // unused so far
		err := rows.Scan(&i, &t)
		if err != nil {
			log.Println("Error: ", err.Error())
			return tags
		}
		tags = append(tags, t)
	}
	return tags
}

func getTestRunsByTag(tag string) []GrllTestRun {
	list := make([]GrllTestRun,0)
	rows, err := db.Query(sqlGetTestRunsByTag, tag)
	if err != nil {
		log.Println("Error: ", err.Error())
		return list
	}
	defer rows.Close()
	for rows.Next() {
		i := NewGrllTestRun()
		err := rows.Scan(&i.Id, &i.Run, &i.Timestamp)
		if err != nil {
			log.Println("Error: ", err.Error())
			return list
		}
		list = append(list, *i)
	}
	return list
}

func loadRDBSuiteToDB(g *GrllTestRun) {
	runId := newTestRun(g.Run, g.Timestamp)
	for _, t := range g.Tags {
		setTag(runId, t)
	}
	for _, r := range g.Results {
		//debug("        " + r.Test)
		newResult(&r, runId)
	}
}

func loadDir(dir string) {
	log.Println("Uploading reports from directory: ", dir)
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if ! info.IsDir() {
				log.Println("    file: ", path)
				rdbReport, _ := os.Open(path)
				importRDBXml(rdbReport)
			}
			return nil
		})
	if err != nil {
		fmt.Println(err)
	}
}

// TODO: return status
func importRDBXml(rdr io.Reader) {
	var s RDBTestSuite
	d := xml.NewDecoder(rdr)
	if err := d.Decode(&s); err != nil {
		fmt.Println("parsing config file", err.Error())
	}
	log.Println("        report: ", s.Config.Name)
	g := rdbToGrll(&s)
	loadRDBSuiteToDB(g)
	//if j, err := json.MarshalIndent(g, "", "    "); err == nil {
	//	fmt.Printf("%s\n", j)
	//}
}


func main() {

	port := flag.Int("port", 3000, "port to listen")
	www := flag.String("www", "www", "web dir to serve")
	importDir := flag.String("import", "", "file or directory with reports to import to db")
	flag.Parse()

	if *importDir != "" {
		loadDir(*importDir)
		fmt.Println("exiting...")
		return
	}

	err := filepath.Walk(*www,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if ! info.IsDir() {
			// TODO: path.Join
				to := "/" + strings.TrimPrefix(strings.TrimPrefix(path, *www), "/")
				http.HandleFunc(to,
					func (w http.ResponseWriter, r *http.Request) {
					http.ServeFile(w, r, "./" + path)
				})
			}
			return nil
		})
	if err != nil {
	fmt.Println(err.Error())
		return
	}
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "www/grll.html")
	})

	http.HandleFunc("/api/testrun", getTestRun)
	http.HandleFunc("/api/tag", getTag)

	http.HandleFunc("/api/upload", uploadReport)

	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)

return
}


func getTestRun(w http.ResponseWriter, r *http.Request) {
	debug(r.URL.String() + " from " + r.RemoteAddr)
	runId := r.URL.Query().Get("id")
	tag := r.URL.Query().Get("tag")
	if runId != "" {
		w.Header().Set("Content-Type", "application/xml")
		i, _ := strconv.Atoi(runId)
		if x, err := xml.MarshalIndent(getTestRunById(i), "", "  "); err == nil {
	x = []byte(xml.Header + xslHeaderGringo + string(x))
			fmt.Fprintf(w, "%s\n", x)
		}
	} else if tag != "" {
		log.Println("Requesting tag ", tag)
		w.Header().Set("Content-Type", "application/json")
		e := json.NewEncoder(w)
		e.Encode(getTestRunsByTag(tag))
	} else {
		w.Header().Set("Content-Type", "application/json")
		e := json.NewEncoder(w)
		e.Encode(getTestRuns())
	}
}

func getTag(w http.ResponseWriter, r *http.Request) {
	debug(r.URL.String() + " from " + r.RemoteAddr)
	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)
	e.Encode(getTags())
}

func uploadReport(w http.ResponseWriter, r *http.Request) {
	debug(r.URL.String() + " from " + r.RemoteAddr)
	w.Header().Set("Content-Type", "text/plain")
	url := r.URL.Query().Get("url")
	// upload by url
	// r.Method == GET
	if url != "" {
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Println("bad status: ", resp.Status)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		importRDBXml(resp.Body)
	}
	// TODO: else upload file
	// https://zupzup.org/go-http-file-upload-download/
	fmt.Fprintln(w, "ok")
}
