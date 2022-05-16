package domain

import "time"

type ImageSize struct {
	Height float32 `json:"height,omitempty"`
	Width  float32 `json:"width,omitempty"`
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
	Show        int           `json:"show,omitempty"`
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
	Show        int        `json:"show,omitempty"`
}

type Exhibition struct {
	ID          int        `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Info        []*Param   `json:"info,omitempty"`
	Image       string     `json:"picture,omitempty"`
	Sizes       *ImageSize `json:"pictureSize,omitempty"`
	Pictures    []*Picture `json:"content,omitempty"`
	Show        int        `json:"show,omitempty"`
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
	Login    string `json:"login,omitempty"`
	PassIn   string `json:"password,omitempty"`
	Password []byte `json:"-"`
	Museum   string `json:"museum,omitempty"`
	Token    string `json:"token,omitempty"`
}

type MuseumExhibition struct {
	Mus *Museum     `json:"museum,omitempty"`
	Exh *Exhibition `json:"exhibition,omitempty"`
}

type Stats struct {
	Port     int       `json:"port"`
	Method   string    `json:"method"`
	Status   int       `json:"status"`
	URL      string    `json:"url,omitempty"`
	Duration int       `json:"duration"`
	When     time.Time `json:"when,omitempty"`
}
