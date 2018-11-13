package rest

const (
    ResultOK                  = 200 // RFC 7231, 6.3.1
    ResultBadRequest          = 400 // RFC 7231, 6.5.1
    ResultUnauthorized        = 401 // RFC 7235, 3.1
    ResultPaymentRequired     = 402 // RFC 7231, 6.5.2
    ResultForbidden           = 403 // RFC 7231, 6.5.3
    ResultNotFound            = 404 // RFC 7231, 6.5.4
    ResultMethodNotAllowed    = 405 // RFC 7231, 6.5.5
    ResultNotAcceptable       = 406 // RFC 7231, 6.5.6
    ResultProxyAuthRequired   = 407 // RFC 7235, 3.2
    ResultRequestTimeout      = 408 // RFC 7231, 6.5.7
    ResultConflict            = 409 // RFC 7231, 6.5.8
    ResultGone                = 410 // RFC 7231, 6.5.9
    ResultInternalServerError = 500 // RFC 7231, 6.6.1
    ResultNotImplemented      = 501 // RFC 7231, 6.6.2
    ResultServiceUnavailable  = 503 // RFC 7231, 6.6.4
    ResultRequireUpdate       = 600
    ResultUserError           = 700
    ResultSystemError         = 800
    ResultFatalError          = 900

    PlatformNone   = 0
    PlatformApple  = 1
    PlatformGoogle = 2

    EnvironmentLocal       = 9
    EnvironmentDemo        = 8
    EnvironmentDevelop1    = 6
    EnvironmentDevelop2    = 5
    EnvironmentTest        = 4
    EnvironmentStaging     = 3
    EnvironmentAppleReview = 2
    EnvironmentProduction  = 1
)
