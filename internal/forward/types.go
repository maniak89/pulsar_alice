package forward

type accountInfo struct {
	webAccountID string
	accountID    string
	accountPID   string
	dbID         string
}

type meterResponse struct {
	MeterID         string `json:"@MeterId"`
	ServiceNameID   string `json:"@ServiceNameId"`
	ServiceName     string `json:"@ServiceName"`
	RateNumbers     string `json:"@RateNumbers"`
	RemarkID        string `json:"@RemarkId"`
	RemarkIDText    string `json:"@RemarkIdText"`
	RateNumber1     string `json:"@RateNumber1"`
	RateNumber1Text string `json:"@RateNumber1Text"`
	Value1          string `json:"@Value1"`
	RateNumber2     string `json:"@RateNumber2"`
	RateNumber2Text string `json:"@RateNumber2Text"`
	Value2          string `json:"@Value2"`
	EDate           string `json:"@EDate"`
}

type meterServiceData struct {
	MeterID         string `json:"MeterId"`
	RateNumbers     string `json:"RateNumbers"`
	EDate           string `json:"EDate"`
	RateNumber1     string `json:"RateNumber1"`
	RateNumber1Text string `json:"RateNumber1Text"`
	Value1          string `json:"Value1"`
}
type meterService struct {
	ServiceNameID string             `json:"ServiceNameId"`
	ServiceName   string             `json:"ServiceName"`
	Data          []meterServiceData `json:"data"`
}

type responseService struct {
	RemarkID     string               `json:"RemarkId"`
	RemarkIDText string               `json:"RemarkIdText"`
	Data         map[int]meterService `json:"data"`
}

type metersResponse struct {
	Answer struct {
		Rcode  string `json:"rcode"`
		Rmsg   string `json:"rmsg"`
		Rsol   string `json:"rsol"`
		Trans  string `json:"trans"`
		Meters struct {
			Row []meterResponse `json:"row"`
		} `json:"meters"`
		Data map[int]responseService `json:"data"`
	} `json:"answer"`
}
