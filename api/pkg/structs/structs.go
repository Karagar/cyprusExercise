package structs

// Config - aggregate Routes props
type Config struct {
	Routes []Route `json:"routes"`
}

// Route - declaration of props to handle URIs
type Route struct {
	Path struct {
		Method string `json:"method"`
		URL    string `json:"url"`
	} `json:"path"`
	Handler    string `json:"handler"`
	IsNeedAuth bool   `json:"isNeedAuth"`
	IsCheckIP  bool   `json:"isCheckIP"`
	IsUseQueue bool   `json:"isUseQueue"`
	IsApi      bool   `json:"isApi"`
	TimeoutSec int    `json:"timeoutSec"`
}

// Company - main object of this exercise
type Company struct {
	Id          []byte `json:"-"`
	Uuid        string
	CompanyName string
	Code        string
	Country     string
	Website     string
	Phone       string
	Archive     bool
	DTCreated   string
	DTUpdated   string
	DTArchived  string
}

type CompanyResponse struct {
	Data  []*Company `json:"data"`
	Count int        `json:"count"`
}
