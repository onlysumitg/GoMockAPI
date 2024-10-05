package httputils

var CodesList []int = []int{
	100,
	101,
	102,
	103,
	200,
	201,
	202,
	203,
	204,
	205,
	206,
	207,
	208,
	226,
	300,
	301,
	302,
	303,
	304,
	305,
	306,
	307,
	308,
	400,
	401,
	402,
	403,
	404,
	405,
	406,
	407,
	408,
	409,
	410,
	411,
	412,
	413,
	414,
	415,
	416,
	417,
	418,
	421,
	422,
	423,
	424,
	425,
	426,
	428,
	429,
	431,
	451,
	500,
	501,
	502,
	503,
	504,
	505,
	506,
	507,
	508,
	510,
	511,
}

var Codes map[int]string = map[int]string{
	100: "StatusContinue",           // RFC 9110, 15.2.1
	101: "StatusSwitchingProtocols", // RFC 9110, 15.2.2
	102: "StatusProcessing",         // RFC 2518, 10.1
	103: "StatusEarlyHints",         // RFC 8297

	200: "StatusOK",                   // RFC 9110, 15.3.1
	201: "StatusCreated",              // RFC 9110, 15.3.2
	202: "StatusAccepted",             // RFC 9110, 15.3.3
	203: "StatusNonAuthoritativeInfo", // RFC 9110, 15.3.4
	204: "StatusNoContent",            // RFC 9110, 15.3.5
	205: "StatusResetContent",         // RFC 9110, 15.3.6
	206: "StatusPartialContent",       // RFC 9110, 15.3.7
	207: "StatusMultiStatus",          // RFC 4918, 11.1
	208: "StatusAlreadyReported",      // RFC 5842, 7.1
	226: "StatusIMUsed",               // RFC 3229, 10.4.1

	300: "StatusMultipleChoices",   // RFC 9110, 15.4.1
	301: "StatusMovedPermanently",  // RFC 9110, 15.4.2
	302: "StatusFound",             // RFC 9110, 15.4.3
	303: "StatusSeeOther",          // RFC 9110, 15.4.4
	304: "StatusNotModified",       // RFC 9110, 15.4.5
	305: "StatusUseProxy",          // RFC 9110, 15.4.6
	306: "_",                       // RFC 9110, 15.4.7 (Unused)
	307: "StatusTemporaryRedirect", // RFC 9110, 15.4.8
	308: "StatusPermanentRedirect", // RFC 9110, 15.4.9

	400: "StatusBadRequest",                   // RFC 9110, 15.5.1
	401: "StatusUnauthorized",                 // RFC 9110, 15.5.2
	402: "StatusPaymentRequired",              // RFC 9110, 15.5.3
	403: "StatusForbidden",                    // RFC 9110, 15.5.4
	404: "StatusNotFound",                     // RFC 9110, 15.5.5
	405: "StatusMethodNotAllowed",             // RFC 9110, 15.5.6
	406: "StatusNotAcceptable",                // RFC 9110, 15.5.7
	407: "StatusProxyAuthRequired",            // RFC 9110, 15.5.8
	408: "StatusRequestTimeout",               // RFC 9110, 15.5.9
	409: "StatusConflict",                     // RFC 9110, 15.5.10
	410: "StatusGone",                         // RFC 9110, 15.5.11
	411: "StatusLengthRequired",               // RFC 9110, 15.5.12
	412: "StatusPreconditionFailed",           // RFC 9110, 15.5.13
	413: "StatusRequestEntityTooLarge",        // RFC 9110, 15.5.14
	414: "StatusRequestURITooLong",            // RFC 9110, 15.5.15
	415: "StatusUnsupportedMediaType",         // RFC 9110, 15.5.16
	416: "StatusRequestedRangeNotSatisfiable", // RFC 9110, 15.5.17
	417: "StatusExpectationFailed",            // RFC 9110, 15.5.18
	418: "StatusTeapot",                       // RFC 9110, 15.5.19 (Unused)
	421: "StatusMisdirectedRequest",           // RFC 9110, 15.5.20
	422: "StatusUnprocessableEntity",          // RFC 9110, 15.5.21
	423: "StatusLocked",                       // RFC 4918, 11.3
	424: "StatusFailedDependency",             // RFC 4918, 11.4
	425: "StatusTooEarly",                     // RFC 8470, 5.2.
	426: "StatusUpgradeRequired",              // RFC 9110, 15.5.22
	428: "StatusPreconditionRequired",         // RFC 6585, 3
	429: "StatusTooManyRequests",              // RFC 6585, 4
	431: "StatusRequestHeaderFieldsTooLarge",  // RFC 6585, 5
	451: "StatusUnavailableForLegalReasons",   // RFC 7725, 3

	500: "StatusInternalServerError",           // RFC 9110, 15.6.1
	501: "StatusNotImplemented",                // RFC 9110, 15.6.2
	502: "StatusBadGateway",                    // RFC 9110, 15.6.3
	503: "StatusServiceUnavailable",            // RFC 9110, 15.6.4
	504: "StatusGatewayTimeout",                // RFC 9110, 15.6.5
	505: "StatusHTTPVersionNotSupported",       // RFC 9110, 15.6.6
	506: "StatusVariantAlsoNegotiates",         // RFC 2295, 8.1
	507: "StatusInsufficientStorage",           // RFC 4918, 11.5
	508: "StatusLoopDetected",                  // RFC 5842, 7.2
	510: "StatusNotExtended",                   // RFC 2774, 7
	511: "StatusNetworkAuthenticationRequired", // RFC 6585, 6
}
