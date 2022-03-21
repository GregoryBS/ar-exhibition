package domain

type Museum struct {
	ID          int          `json:"id,omitempty"`
	Name        string       `json:"name,omitempty"`
	Country     string       `json:"country,omitempty"`
	City        string       `json:"city,omitempty"`
	Year        int          `json:"year,omitempty"`
	Description string       `json:"descr,omitempty"`
	Address     string       `json:"address,omitempty"`
	Director    string       `json:"director,omitempty"`
	Image       string       `json:"image,omitempty"`
	Exhibitions []Exhibition `json:"exhibitions,omitempty"`
}

type MuseumPage struct {
	Page  int       `json:"page,omitempty"`
	Size  int       `json:"pageSize,omitempty"`
	Total int       `json:"totalElements,omitempty"`
	Items []*Museum `json:"items,omitempty"`
}

type Picture struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Year      int    `json:"year,omitempty"`
	Author    string `json:"author,omitempty"`
	Technique string `json:"technique,omitempty"`
	Image     string `json:"image,omitempty"`
	Height    int    `json:"height,omitempty"`
	Width     int    `json:"width,omitempty"`
}

type Exhibition struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	From        string `json:"dateFrom,omitempty"`
	To          string `json:"dateTo,omitempty"`
	Image       string `json:"image,omitempty"`
}

type MainPage struct {
	Museums     []*Museum     `json:"topMuseum"`
	Exhibitions []*Exhibition `json:"topExhibition"`
	Pictures    []*Picture    `json:"recommendation"`
}
