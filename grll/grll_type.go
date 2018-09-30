package grll 

/* TODO:
 - overall status
 */

type GrllTestRun struct {
	Run         string          `xml:"run"      json:"run"`
	Timestamp   string          `xml:"ts"       json:"ts"`
	Link        string          `xml:"link"     json:"link"`
	Results     []GrllResult    `xml:"results"  json:"results"`
	Values      []GrllValue     `xml:"values"   json:"values"`
	Tags        []string        `xml:"tags"     json:"tags"`
	// ??
	Id          int             `xml:"id"       json:"id"`
}

type GrllResult struct {
	//XMLName  xml.Name  `xml:"result"`
	Test     string    `xml:"test"   json:"test"`
	Status   string    `xml:"status" json:"status"`
	Message  string    `xml:"msg"    json:"msg,omitempty"`
	// TODO: duration?
}

type GrllValue struct {
	Test   string   `xml:"test"             json:"test"`
	Value  float64  `xml:"value,omitempty"  json:"value,omitempty"`
	Unit   string   `xml:"unit,omitempty"   json:"unit,omitempty"`
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

