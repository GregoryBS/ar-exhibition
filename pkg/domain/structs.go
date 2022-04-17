package domain

type ImageSize struct {
	Height int `json:"height,omitempty"`
	Width  int `json:"width,omitempty"`
}

type Param struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Museum struct {
	ID          int           `json:"id,omitempty"`
	Name        string        `json:"name,omitempty"`
	Description string        `json:"descr,omitempty"`
	Info        []*Param      `json:"info,omitempty"`
	Image       string        `json:"picture,omitempty"`
	Sizes       *ImageSize    `json:"pictureSize,omitempty"`
	Exhibitions []*Exhibition `json:"exhibitions,omitempty"`
}

type Page struct {
	Number int           `json:"page"`
	Size   int           `json:"pageSize"`
	Total  int           `json:"totalElements"`
	Items  []interface{} `json:"items,omitempty"`
}

type Picture struct {
	ID          int        `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Image       string     `json:"picture,omitempty"`
	Sizes       *ImageSize `json:"pictureSize,omitempty"`
	Description string     `json:"descr,omitempty"`
	Info        []*Param   `json:"info,omitempty"`
	Video       string     `json:"video,omitempty"`
	VideoSize   string     `json:"videoSize,omitempty"`
}

type Exhibition struct {
	ID          int        `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Info        []*Param   `json:"info,omitempty"`
	Image       string     `json:"picture,omitempty"`
	Sizes       *ImageSize `json:"pictureSize,omitempty"`
	Pictures    []*Picture `json:"content,omitempty"`
}

type MainPage struct {
	Museums     []*Museum     `json:"topMuseum"`
	Exhibitions []*Exhibition `json:"topExhibition"`
	Pictures    []*Picture    `json:"recommendation"`
}

type SearchPage struct {
	Museums     []*Museum     `json:"museums,omitempty"`
	Exhibitions []*Exhibition `json:"exhibitions,omitempty"`
	Pictures    []*Picture    `json:"pictures,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type User struct {
	ID       int    `json:"id,omitempty"`
	Login    string `json:"login"`
	PassIn   string `json:"password,omitempty"`
	Password []byte `json:"-"`
	Museum   string `json:"museum,omitempty"`
	Token    string `json:"token,omitempty"`
}
