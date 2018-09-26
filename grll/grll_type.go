package grll 

/* TODO:
 - overall status
 */

type GrllTestRun struct {
	Run         string          `xml:"run"      json:"run"`
	Timestamp   string          `xml:"ts"       json:"ts"`
	Results     []GrllResult    `xml:"results"  json:"results"`
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

// TODO: too artificial
type GrllHistorical struct {
	Test  string               `xml:"test"`
	Items []GrllHistoricalItem `xml:"items"`
}

type GrllHistoricalItem struct {
	Run       string `xml:"run"`
	Status    string `xml:"status"`
	Message   string `xml:"msg"`
	Timestamp string `xml:"ts"`
}

func NewGrllTestRun() *GrllTestRun {
	return &GrllTestRun{
		Run:      "Unnamed Test Run",
		Results:  make([]GrllResult,0),
		Tags:     make([]string, 0),
		Id:       -1,
	}
}

