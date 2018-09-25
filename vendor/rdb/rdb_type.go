package rdb 


import (
	"encoding/xml"
	//"encoding/json"
	"os"
)



type RDBTestSuite struct {
	XMLName   xml.Name        `xml:"report"    json:"report"`
	Config    RDBConfig       `xml:"config"    json:"config"`
	Timestamp string          `xml:"timestamp" json:"timestamp"` //date '+%Y-%m-%dT%H:%M:%S%:z'
	TestList  RDBTestList     `xml:"test_list" json:"test_list"`
	// the rest is garbage
	Kind      string          `xml:"type,attr" json:"type,attr"`
	Project   string          `xml:"project"   json:"project"`
	Log       string          `xml:"log"       json:"log"`
	Milestone string          `xml:"milestone" json:"milestone"`
	Platform  string          `xml:"platform"  json:"platform"`
	Owner     RDBOwner        `xml:"owner"     json:"owner"`
	Suite     RDBSuite        `xml:"suite"     json:"suite"`
	Target    RDBTarget       `xml:"target"    json:"target"`
}

type RDBOwner struct {
	Email    string    `xml:"email" json:"email"`
}

type RDBConfig struct {
	Name  string      `xml:"name"          json:"name"`
	Props RDBPropList `xml:"property_list" json:"property_list"`
}

type RDBPropList struct {
	Properties []RDBProp   `xml:"item" json:"item"`
}

type RDBProp struct {
	Key   string `xml:"key"   json:"key"`
	Type  string `xml:"type"  json:"type"`
	Value string `xml:"value" json:"value"`
}

type RDBSuite struct {
	Name    string `xml:"name"    json:"name"`
	Version string `xml:"version" json:"version"`
}

type RDBTarget struct {
	Version    string `xml:"version"    json:"version"`
	Competitor string `xml:"competitor" json:"competitor"`
}

type RDBTestList struct {
	TestCases []RDBTestCase   `xml:"item" json:"item"`
}

type RDBTestCase struct {
	Name           string    `xml:"name"               json:"name"`
	Status         string    `xml:"status"             json:"status"`
	// TODO: format
	Duration       float64   `xml:"duration,omitempty" json:"duration,omitempty"`
	Message        string    `xml:"message,omitempty"  json:"message,omitempty"`
}


func NewRDBTestSuite() *RDBTestSuite {
	return &RDBTestSuite{
			Config:     RDBConfig { Name: os.Getenv("TPP_CONFIG"),
									Props: RDBPropList{} },
			TestList:   RDBTestList {TestCases: []RDBTestCase{}},
			Timestamp:  "2018-08-31T02:20:10-07:00", // TODO
			Kind:       "auto",
			Project:    os.Getenv("TPP_PROJECT"),
			Log:        os.Getenv("TPP_LOG_BASE"),
			Milestone:  os.Getenv("TPP_MILESTONE"),
			Platform:   os.Getenv("TPP_PLATFORM"),
			Owner:      RDBOwner { Email: os.Getenv("TPP_EMAIL_OWNER") },
			Suite:      RDBSuite { Name: os.Getenv("TPP_SUITE"),
									Version: os.Getenv("TPP_SUITE_VERSION") },
			Target:     RDBTarget { Version: os.Getenv("TPP_TARGET"),
									Competitor: os.Getenv("TPP_TARGET_COMPETITOR") },
	}
}
