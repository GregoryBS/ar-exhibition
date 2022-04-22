package utils

const (
	MuseumService     = "http://museum:8081"
	ExhibitionService = "http://exhibition:8082"
	PictureService    = "http://picture:8083"
	UserService       = "http://user:8000"
	GatewayService    = "http://gateway:8080"
	ImageService      = "http://95.163.213.222/images/"
	VideoService      = "http://95.163.213.222/videos/"
)

const (
	GatewayApiMain         = "/api/v1/"
	GatewayApiPictureID    = "/api/v1/pictures/:id"
	GatewayApiPictureImage = "/api/v1/pictures/:id/images"
	GatewayApiExhibitionID = "/api/v1/exhibitions/:id"
	GatewayApiMuseumID     = "/api/v1/museums/:id"
	GatewayApiMuseumImage  = "/api/v1/museums/:id/images"
	GatewayApiMuseums      = "/api/v1/museums"
	GatewayApiExhibitions  = "/api/v1/exhibitions"
	GatewayApiSearch       = "/api/v1/search"
	GatewayApiPictures     = "/api/v1/pictures"
	GatewayApiMuseumShow   = "/api/v1/museums/:id/public"
)

const (
	BaseMuseumApi    = "/api/v1/museums"
	MuseumTop        = BaseMuseumApi + "/top"
	MuseumID         = BaseMuseumApi + "/:id"
	BaseMuseumSearch = BaseMuseumApi + "/search"
	MuseumSearch     = BaseMuseumSearch + "?name="
	MuseumImage      = MuseumID + "/images"
	MuseumShow       = MuseumID + "/public"
)

const (
	BaseExhibitionApi    = "/api/v1/exhibitions"
	ExhibitionTop        = BaseExhibitionApi + "/top"
	ExhibitionID         = BaseExhibitionApi + "/:id"
	ExhibitionByMuseum   = BaseExhibitionApi + "?museumID="
	BaseExhibitionSearch = BaseExhibitionApi + "/search"
	ExhibitionSearch     = BaseExhibitionSearch + "?name="
	ExhibitionSearchID   = BaseExhibitionSearch + "?name=%s&id=%d"
	ExhibitionShow       = BaseExhibitionApi + "/public"
	ExhibitionShowID     = ExhibitionID + "/public"
)

const (
	BasePictureApi      = "/api/v1/pictures"
	PictureID           = BasePictureApi + "/:id"
	PictureImage        = PictureID + "/images"
	PictureByExhibition = BasePictureApi + "?exhibitionID="
	PictureByIDs        = BasePictureApi + "?id="
	BasePictureSearch   = BasePictureApi + "/search"
	PictureSearch       = BasePictureSearch + "?name="
	PictureSearchID     = BasePictureSearch + "?name=%s&id=%d"
	PictureShow         = BasePictureApi + "/public"
	PictureShowExh      = PictureShow + "?exhibitionID="
	PictureShowID       = PictureID + "/public"
)

const (
	ExhibitionStart = "Начало"
	ExhibitionEnd   = "Конец"
)

const (
	BaseUserApi = "/api/v1/users"
	UserSignup  = BaseUserApi + "/signup"
	UserLogin   = BaseUserApi + "/login"
	UserID      = BaseUserApi + "/id"
)

const UserHeader = "X-User"
