package constant

const (
	DatabaseMySQL         = "DBMySQL"
	DatabasePostgreSQL    = "DBPostgreSQL"
	DatabaseIMSPostgreSQL = "DBIMSPostgreSQL"
	RedisDefault          = "redis"

	NotFound = "not_found"

	//http status strings
	StatusUnauthorized        = "status_unauthorized"
	StatusBadRequest          = "bad_request"
	StatusNotFound            = "not_found"
	StatusRequestTimeout      = "timeout"
	StatusExpectationFailed   = "failed"
	StatusInternalServerError = "internalError"
	StatusTokenExpired        = "token_expired"
	StatusForbidden           = "forbidden"

	StatusUnauthorizedMessage = "status unauthorized"
	StatusBadRequestMessage   = "bad request"
	StatusNotFoundMessage     = "data not found"
	StatusOKMessage           = "ok"
	StatusTokenExpiredMessage = "token expired"
	StatusUnprocessableEntity = "status unprocessable entity"
	StatusOK                  = "status_ok"

	//GIN
	Release        = "release"
	REQUEST_ID     = contextKey("Request-ID")
	AUTH_API_KEY   = contextKey("X-AUTH-API-KEY")
	FORWARDED_FOR  = contextKey("X-FORWADED-FOR")
	TRUE_CLIENT_IP = contextKey("True-Client-Ip")
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}
