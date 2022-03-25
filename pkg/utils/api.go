package utils

const (
	MuseumService     = "http://museum:8080"
	ExhibitionService = "http://exhibition:8080"
	PictureService    = "http://picture:8080"
	Service           = "http://95.163.213.222/"
)

const (
	GatewayApiMain         = "/api/v1/"
	GatewayApiPictureID    = "/api/v1/pictures/:id"
	GatewayApiExhibitionID = "/api/v1/exhibitions/:id"
	GatewayApiMuseumID     = "/api/v1/museums/:id"
	GatewayApiMuseums      = "/api/v1/museums"
	GatewayApiExhibitions  = "/api/v1/exhibitions"
)

const (
	BaseMuseumApi = "/api/v1/museums"
	MuseumTop     = BaseMuseumApi + "/top"
	MuseumID      = BaseMuseumApi + "/:id"
)

const (
	BaseExhibitionApi  = "/api/v1/exhibitions"
	ExhibitionTop      = BaseExhibitionApi + "/top"
	ExhibitionID       = BaseExhibitionApi + "/:id"
	ExhibitionByMuseum = BaseExhibitionApi + "?museumID="
)

const (
	BasePictureApi      = "/api/v1/pictures"
	PictureID           = BasePictureApi + "/:id"
	PictureByExhibition = BasePictureApi + "?exhibitionID="
)
