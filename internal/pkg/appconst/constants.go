package appconst

const (
	MsgRegisterCreated             = "Successfully registered user"
	MsgVerifiedOk                  = "Successfully verified user"
	MsgListUserOk                  = "Successfully retrieved user list"
	MsgLoginOk                     = "Successfully logged in"
	MsgResetPasswordOk             = "Successfully send reset password token"
	MsgConfirmResetOk              = "Successfully create new password"
	MsgLogoutOk                    = "Successfully logged out"
	MsgGetCartItemsOk              = "Successfully retrieved cart item(s)"
	MsgAddCartItemCreated          = "Successfully added cart item"
	MsgUpdateCartItemOk            = "Successfully updated cart item"
	MsgRemoveCartItemOk            = "Successfully removed cart item"
	MsgAddPharmacistCreated        = "Successfully added new pharmacist"
	MsgListPartnersOk              = "Successfully retrieved partner list"
	MsgListPharmacistsOk           = "Successfully retrieved pharmacist list"
	MsgGetPharmacistOk             = "Successfully retrieved pharmacist"
	MsgUpdatePharmacistOk          = "Successfully updated pharmacist"
	MsgRemovePharmacistOk          = "Successfully removed pharmacist"
	MsgAddPartnerCreated           = "Successfully added new partner"
	MsgUpdatePartnerOk             = "Successfully updated partner"
	MsgGetPartnerOk                = "Successfully retrieved partner"
	MsgAddPharmacyCreated          = "Successfully added new pharmacy"
	MsgListPharmaciesOk            = "Successfully retrieved pharmacy list"
	MsgListAvailPharmaciesOk       = "Successfully retrieved available pharmacy list"
	MsgListPharmacyLogisticsOk     = "Successfully retrieved pharmacy logistic list"
	MsgAddAddressCreated           = "Successfully added new address"
	MsgListAddressesOk             = "Successfully retrieved user address list"
	MsgListProvinceOk              = "Successfully retrieved province list"
	MsgListCityOk                  = "Successfully retrieved city list"
	MsgListProductCategoriesOk     = "Successfully retrieved categories list"
	MsgAddProductCategoryCreated   = "Successfully added new product category"
	MsgUpdateProductCategoryOk     = "Successfully updated product category"
	MsgRemoveProductCategoryOk     = "Successfully removed product category"
	MsgListProductOk               = "Successfully retrieved product list"
	MsgListManufacturerOk          = "Successfully retrieved manufacturer list"
	MsgListLogisticOk              = "Successfully retrieved logistic list"
	MsgListProductClassificationOk = "Successfully retrieved product classification list"
	MsgListProductFormOk           = "Successfully retrieved product form list"
	MsgListPopularProductOk        = "Successfully retrieved popular product(s)"
	MsgGetProductDetailsOk         = "Successfully retrieved product details"
	MsgListPharmacyProductsOk      = "Successfully retrieved added product(s)"
	MsgGetPharmacyProductOk        = "Successfully retrieved pharmacy product"
	MsgAddPharmacyProductCreated   = "Successfully added as pharmacy product"
	MsgUpdatePharmacyProductOk     = "Successfully updated pharmacy product"
	MsgDeletePharmacyProductOk     = "Successfully removed pharmacy product"
	MsgAddProductOk                = "Successfully added new product"
	MsgCreateOrderCreated          = "Successfully created order"
	MsgListUnpaidOrderOk           = "Successfully retrieved unpaid orders"
	MsgListOrderOk                 = "Successfully retrieved orders"
	MsgUploadPaymentProofOk        = "Successfully uploaded payment proof"
	MsgUpdateOrderStatusOk         = "Successfully updated order status"
)

const (
	BindingRequired              = "required"
	BindingEq                    = "eq"
	BindingEmail                 = "email"
	BindingMin                   = "min"
	BindingGte                   = "gte"
	BindingNumber                = "number"
	BindingNumeric               = "numeric"
	BindingRequiredWithout       = "required_without"
	BindingRequiredWithoutAll    = "required_without_all"
	BindingOneOf                 = "oneof"
	DefaultFieldErrorTranslation = "this field contains invalid input"
)

const (
	RoleAdmin      = "admin"
	RoleCustomer   = "customer"
	RolePharmacist = "pharmacist"
)

const (
	KeyUserID     = "user_id"
	KeyRole       = "role"
	KeyIsVerified = "is_verified"
)

const (
	MinDayId = 0
	MaxDayId = 6
)

const (
	QueueDefault  = "default"
	QueueCritical = "critical"
	QueueLow      = "low"
)

const (
	WIB = "WIB"
	UTC = "UTC"
)

const (
	ContentTypeFormUrlEncoded = "application/x-www-form-urlencoded"
)

const (
	LogisticNameOfficial = "Official"
	LogisticServiceYES   = "YES"
)

const (
	CacheBestSeller = "BestSellerProduct"
)
