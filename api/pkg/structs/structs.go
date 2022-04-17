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
	ID          []byte  `gorm:"default:newid()" json:"-"`
	Uuid        string  `gorm:"-"`
	CompanyName *string `json:"Name"`
	Code        *string
	Country     *string
	Website     *string
	Phone       *string
	Archive     *bool   `gorm:"<-:update" json:"-"`
	DTCreated   *string `gorm:"-" json:"-"`
	DTUpdated   *string `gorm:"-" json:"-"`
	DTArchived  *string `gorm:"<-:update" json:"-"`
}

// CompanyResponse - structure that server send to client on request
type CompanyResponse struct {
	Data  []*Company `json:"data"`
	Count int        `json:"count"`
}
