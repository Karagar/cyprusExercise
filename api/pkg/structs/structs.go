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
	Uuid        []byte `json:"-"`
	Guid        string `json:"Guid"`
	CompanyName string `json:"CompanyName"`
	Code        string `json:"Code"`
	Country     string `json:"Country"`
	Website     string `json:"Website"`
	Phone       string `json:"Phone"`
	Archive     bool   `json:"Archive"`
	DTCreated   string `json:"DTCreated"`
	DTUpdated   string `json:"DTUpdated"`
	DTArchived  string `json:"DTArchived"`
}
