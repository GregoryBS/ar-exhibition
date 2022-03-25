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
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	// Country     string            `json:"country,omitempty"`
	// City        string            `json:"city,omitempty"`
	// Year        int               `json:"year,omitempty"`
	Description string   `json:"descr,omitempty"`
	Info        []*Param `json:"info,omitempty"`
	// Address     string            `json:"address,omitempty"`
	// Director    string            `json:"director,omitempty"`
	Image       string        `json:"picture,omitempty"`
	Sizes       *ImageSize    `json:"pictureSize,omitempty"`
	Exhibitions []*Exhibition `json:"exhibitions,omitempty"`
}

type MuseumPage struct {
	Page  int       `json:"page,omitempty"`
	Size  int       `json:"pageSize,omitempty"`
	Total int       `json:"totalElements,omitempty"`
	Items []*Museum `json:"items,omitempty"`
}

type Picture struct {
	ID          int        `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Image       string     `json:"picture,omitempty"`
	Sizes       *ImageSize `json:"pictureSize,omitempty"`
	Description string     `json:"descr,omitempty"`
	Info        []*Param   `json:"info,omitempty"`
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

type ErrorResponse struct {
	Message string `json:"message"`
}
