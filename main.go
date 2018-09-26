/*
TODO:
	- set limits for queries (in settings?)
	- rdb perf reports
	- storing numbers, series
	- tests
	- move sql to server functions
 */

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
	"time"
	"path"
)

const (

	//xslHeaderRdb   = `<?xml-stylesheet type="text/xsl" href="/rdb.xsl"?>` + "\n"
	xslGrllTestRun = `<?xml-stylesheet type="text/xsl" href="/grll.xsl"?>` + "\n"
	xslGrllRusults = `<?xml-stylesheet type="text/xsl" href="/grll-results.xsl"?>` + "\n"

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

	sqlGetResultsByTest =
		`SELECT testrun.name, status.name,
		COALESCE(result.message, '') as message, testrun.ts FROM result
		INNER JOIN test ON test.id = result.test
		INNER JOIN testrun ON testrun.id = result.testrun
		INNER JOIN status ON status.id = result.status
		WHERE test.name = $1
		ORDER BY testrun.ts LIMIT 100`

	sqlGetTestRuns      = `SELECT id, name, ts FROM testrun ORDER BY ts LIMIT 50`

	sqlGetTestRunById   = `SELECT name FROM testrun WHERE id = $1`

	sqlPutTestRun       = "SELECT * FROM new_testrun($1, $2)"

	sqlPutResult        = "SELECT * FROM new_result($1, $2, NULL, $3)"

	sqlPutResultMessage = "SELECT * FROM new_result($1, $2, $4, $3)"

	sqlSetTag           = "SELECT set_tag($1, $2)"

	sqlGetTags          = "SELECT id, name FROM tag ORDER BY name LIMIT 50"

	sqlGetTestRunsByTag =
		`SELECT id, name, ts FROM get_testruns_by_tag($1) ORDER BY ts`

	sqlSearchTestRun =
		`select testrun.id, testrun.name, testrun.ts
		from testrun where testrun.name like $1
		order by testrun.ts limit 100`

	sqlSearchTest =
		`select test.name from test
		where test.name like $1
		order by test.name limit 100`

	sqlSearchTag =
		`select tag.name from tag
		where tag.name like $1
		order by tag.name limit 50`
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

func justPrint(e error) {
	if nil != e {
		fmt.Println(e.Error())
	}
}


func debug(msg string) {
	log.Println(msg)
}

/////////// convert report ////////////////////////////////

func rdbToGrll(rdb *RDBTestSuite) *GrllTestRun {
	g := NewGrllTestRun()
	if "" != rdb.Suite.Name {
		g.Run = rdb.Suite.Name + " " + rdb.Config.Name
	} else {
		g.Run = rdb.Config.Name
	}
	if "" != rdb.Timestamp {
		g.Timestamp = rdb.Timestamp
	} else { // set current time as last resort
		g.Timestamp = time.Now().Format(time.RFC3339)
	}
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

////////////  db queries ///////////////////////////////

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
	justPrint(err)
}

func newResult(result *GrllResult, runId int) {
	var err error
	// TODO: there should be a better way of mapping '' -> NULL
	if "" != result.Message {
		_, err = db.Exec(sqlPutResultMessage, result.Test, result.Status, runId, result.Message)
	} else {
		_, err = db.Exec(sqlPutResult, result.Test, result.Status, runId)
	}
	justPrint(err)
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

// TODO: merge
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
func getTagsLike(likeThis string) []string {
	tags := make([]string,0)
	rows, err := db.Query(sqlSearchTag, "%" + likeThis + "%")
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

// TODO: refactor this copy-paste stuff
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
func getTestRunsLike(likeThis string) []GrllTestRun {
	list := make([]GrllTestRun,0)
	// TODO: check if ends with %
	rows, err := db.Query(sqlSearchTestRun, "%" + likeThis + "%")
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

func getTestRunByTest(testName string) GrllHistorical {
	list := GrllHistorical{Test:testName,
							Items:make([]GrllHistoricalItem,0)}
	rows, err := db.Query(sqlGetResultsByTest, testName)
	if err != nil {
		log.Println("Error: ", err.Error())
		return list
	}
	defer rows.Close()
	for rows.Next() {
		i := GrllHistoricalItem{}
		err := rows.Scan(&i.Run, &i.Status, &i.Message, &i.Timestamp)
		if err != nil {
			log.Println("Error: ", err.Error())
			return list
		}
		list.Items = append(list.Items, i)
	}
	return list
}


////////// report imports ////////////////////////////////

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

///////// http handlers ////////////////////////////////

func getTestRun(w http.ResponseWriter, r *http.Request) {
	debug(r.URL.String() + " from " + r.RemoteAddr)
	runId := r.URL.Query().Get("id")
	tag := r.URL.Query().Get("tag")
	like := r.URL.Query().Get("like")
	test := r.URL.Query().Get("test")
	if runId != "" {
		w.Header().Set("Content-Type", "application/xml")
		i, _ := strconv.Atoi(runId)
		if x, err := xml.MarshalIndent(getTestRunById(i), "", "  "); err == nil {
			x = []byte(xml.Header + xslGrllTestRun + string(x))
			fmt.Fprintf(w, "%s\n", x)
		}
		return
	} else if test != "" {
		w.Header().Set("Content-Type", "application/xml")
		if x, err := xml.MarshalIndent(getTestRunByTest(test), "", "  "); err == nil {
			x = []byte(xml.Header + xslGrllRusults + string(x))
			fmt.Fprintf(w, "%s\n", x)
		}
		return
	} else if tag != "" {
		log.Println("Requesting tag ", tag)
		w.Header().Set("Content-Type", "application/json")
		e := json.NewEncoder(w)
		e.Encode(getTestRunsByTag(tag))
		return
	} else if like != "" {
		log.Println("Requesting like ", like)
		w.Header().Set("Content-Type", "application/json")
		e := json.NewEncoder(w)
		e.Encode(getTestRunsLike(like))
		return
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
	like := r.URL.Query().Get("like")
	if like != "" {
		log.Println("Requesting like ", like)
		e.Encode(getTagsLike(like))
		return
	} else {
		e.Encode(getTags())
	}
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

////////// the main //////////////////////////////////

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
		},
	)
	if err != nil {
	fmt.Println(err.Error())
		return
	}
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join(*www, "grll.html"))
	})

	http.HandleFunc("/api/testrun", getTestRun)
	http.HandleFunc("/api/tag", getTag)

	http.HandleFunc("/api/upload", uploadReport)

	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)

	return
}
