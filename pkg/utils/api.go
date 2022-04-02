package utils

const (
	MuseumService     = "http://museum:8080"
	ExhibitionService = "http://exhibition:8080"
	PictureService    = "http://picture:8080"
	ImageService      = "http://95.163.213.222/images/"
	VideoService      = "http://95.163.213.222/videos/"
)

const (
	GatewayApiMain         = "/api/v1/"
	GatewayApiPictureID    = "/api/v1/pictures/:id"
	GatewayApiExhibitionID = "/api/v1/exhibitions/:id"
	GatewayApiMuseumID     = "/api/v1/museums/:id"
	GatewayApiMuseums      = "/api/v1/museums"
	GatewayApiExhibitions  = "/api/v1/exhibitions"
	GatewayApiSearch       = "/api/v1/search"
	GatewayApiExhPictures  = "/api/v1/pictures"
)

const (
	BaseMuseumApi    = "/api/v1/museums"
	MuseumTop        = BaseMuseumApi + "/top"
	MuseumID         = BaseMuseumApi + "/:id"
	BaseMuseumSearch = BaseMuseumApi + "/search"
	MuseumSearch     = BaseMuseumSearch + "?name="
)

const (
	BaseExhibitionApi    = "/api/v1/exhibitions"
	ExhibitionTop        = BaseExhibitionApi + "/top"
	ExhibitionID         = BaseExhibitionApi + "/:id"
	ExhibitionByMuseum   = BaseExhibitionApi + "?museumID="
	BaseExhibitionSearch = BaseExhibitionApi + "/search"
	ExhibitionSearch     = BaseExhibitionSearch + "?name="
	ExhibitionSearchID   = BaseExhibitionSearch + "?name=%s&id=%d"
)

const (
	BasePictureApi      = "/api/v1/pictures"
	PictureID           = BasePictureApi + "/:id"
	PictureByExhibition = BasePictureApi + "?exhibitionID="
	BasePictureSearch   = BasePictureApi + "/search"
	PictureSearch       = BasePictureSearch + "?name="
	PictureSearchID     = BasePictureSearch + "?name=%s&id=%d"
)
