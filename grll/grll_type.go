package grll 

/* TODO:
 - timestamp
 - overall status
 */

type GrllTestRun struct {
	Run         string          `xml:"run"      json:"run"`
	Timestamp   string          `xml:"ts"       json:"ts"`
	Results     []GrllResult  `xml:"results"  json:"results"`
	Tags        []string        `xml:"tags"     json:"tags"` // ??? ,omitempty
	// ??
	Id          int             `xml:"id"       json:"id"` // ??? ,omitempty
}

type GrllResult struct {
	//XMLName  xml.Name  `xml:"result"`
	Test     string    `xml:"test"   json:"test"`
	Status   string    `xml:"status" json:"status"`
	Message  string    `xml:"msg"    json:"msg,omitempty"`
}

func NewGrllTestRun() *GrllTestRun {
	return &GrllTestRun{
		Run:      "Unnamed Test Run",
		Results:  make([]GrllResult,0),
		Tags:     make([]string, 0),
		Id:       -1,
	}
}

