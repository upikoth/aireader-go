// Code generated by ogen, DO NOT EDIT.

package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/ogen-go/ogen/uri"
)

func (s *Server) cutPrefix(path string) (string, bool) {
	prefix := s.cfg.Prefix
	if prefix == "" {
		return path, true
	}
	if !strings.HasPrefix(path, prefix) {
		// Prefix doesn't match.
		return "", false
	}
	// Cut prefix from the path.
	return strings.TrimPrefix(path, prefix), true
}

// ServeHTTP serves http request as defined by OpenAPI v3 specification,
// calling handler that matches the path or returning not found error.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	elem := r.URL.Path
	elemIsEscaped := false
	if rawPath := r.URL.RawPath; rawPath != "" {
		if normalized, ok := uri.NormalizeEscapedPath(rawPath); ok {
			elem = normalized
			elemIsEscaped = strings.ContainsRune(elem, '%')
		}
	}

	elem, ok := s.cutPrefix(elem)
	if !ok || len(elem) == 0 {
		s.notFound(w, r)
		return
	}
	args := [1]string{}

	// Static code generated router with unwrapped path search.
	switch {
	default:
		if len(elem) == 0 {
			break
		}
		switch elem[0] {
		case '/': // Prefix: "/api/v1/"
			origElem := elem
			if l := len("/api/v1/"); len(elem) >= l && elem[0:l] == "/api/v1/" {
				elem = elem[l:]
			} else {
				break
			}

			if len(elem) == 0 {
				break
			}
			switch elem[0] {
			case 'h': // Prefix: "health"
				origElem := elem
				if l := len("health"); len(elem) >= l && elem[0:l] == "health" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					// Leaf node.
					switch r.Method {
					case "GET":
						s.handleV1CheckHealthRequest([0]string{}, elemIsEscaped, w, r)
					default:
						s.notAllowed(w, r, "GET")
					}

					return
				}

				elem = origElem
			case 'o': // Prefix: "oauth"
				origElem := elem
				if l := len("oauth"); len(elem) >= l && elem[0:l] == "oauth" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					switch r.Method {
					case "POST":
						s.handleV1AuthorizeUsingOauthRequest([0]string{}, elemIsEscaped, w, r)
					default:
						s.notAllowed(w, r, "POST")
					}

					return
				}
				switch elem[0] {
				case 'R': // Prefix: "Redirect/"
					origElem := elem
					if l := len("Redirect/"); len(elem) >= l && elem[0:l] == "Redirect/" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'm': // Prefix: "mail"
						origElem := elem
						if l := len("mail"); len(elem) >= l && elem[0:l] == "mail" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "GET":
								s.handleV1AuthorizeUsingOauthHandleMailRedirectRequest([0]string{}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "GET")
							}

							return
						}

						elem = origElem
					case 'v': // Prefix: "vk"
						origElem := elem
						if l := len("vk"); len(elem) >= l && elem[0:l] == "vk" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "GET":
								s.handleV1AuthorizeUsingOauthHandleVkRedirectRequest([0]string{}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "GET")
							}

							return
						}

						elem = origElem
					case 'y': // Prefix: "yandex"
						origElem := elem
						if l := len("yandex"); len(elem) >= l && elem[0:l] == "yandex" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "GET":
								s.handleV1AuthorizeUsingOauthHandleYandexRedirectRequest([0]string{}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "GET")
							}

							return
						}

						elem = origElem
					}

					elem = origElem
				}

				elem = origElem
			case 'p': // Prefix: "passwordRecoveryRequests"
				origElem := elem
				if l := len("passwordRecoveryRequests"); len(elem) >= l && elem[0:l] == "passwordRecoveryRequests" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					// Leaf node.
					switch r.Method {
					case "PATCH":
						s.handleV1ConfirmPasswordRecoveryRequestRequest([0]string{}, elemIsEscaped, w, r)
					case "POST":
						s.handleV1CreatePasswordRecoveryRequestRequest([0]string{}, elemIsEscaped, w, r)
					default:
						s.notAllowed(w, r, "PATCH,POST")
					}

					return
				}

				elem = origElem
			case 'r': // Prefix: "registrations"
				origElem := elem
				if l := len("registrations"); len(elem) >= l && elem[0:l] == "registrations" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					// Leaf node.
					switch r.Method {
					case "PATCH":
						s.handleV1ConfirmRegistrationRequest([0]string{}, elemIsEscaped, w, r)
					case "POST":
						s.handleV1CreateRegistrationRequest([0]string{}, elemIsEscaped, w, r)
					default:
						s.notAllowed(w, r, "PATCH,POST")
					}

					return
				}

				elem = origElem
			case 's': // Prefix: "session"
				origElem := elem
				if l := len("session"); len(elem) >= l && elem[0:l] == "session" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					switch r.Method {
					case "GET":
						s.handleV1CheckCurrentSessionRequest([0]string{}, elemIsEscaped, w, r)
					default:
						s.notAllowed(w, r, "GET")
					}

					return
				}
				switch elem[0] {
				case 's': // Prefix: "s"
					origElem := elem
					if l := len("s"); len(elem) >= l && elem[0:l] == "s" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch r.Method {
						case "POST":
							s.handleV1CreateSessionRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "POST")
						}

						return
					}
					switch elem[0] {
					case '/': // Prefix: "/"
						origElem := elem
						if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
							elem = elem[l:]
						} else {
							break
						}

						// Param: "id"
						// Leaf parameter
						args[0] = elem
						elem = ""

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "DELETE":
								s.handleV1DeleteSessionRequest([1]string{
									args[0],
								}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "DELETE")
							}

							return
						}

						elem = origElem
					}

					elem = origElem
				}

				elem = origElem
			case 'u': // Prefix: "user"
				origElem := elem
				if l := len("user"); len(elem) >= l && elem[0:l] == "user" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					switch r.Method {
					case "GET":
						s.handleV1GetCurrentUserRequest([0]string{}, elemIsEscaped, w, r)
					default:
						s.notAllowed(w, r, "GET")
					}

					return
				}
				switch elem[0] {
				case 's': // Prefix: "s"
					origElem := elem
					if l := len("s"); len(elem) >= l && elem[0:l] == "s" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "GET":
							s.handleV1GetUsersRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "GET")
						}

						return
					}

					elem = origElem
				}

				elem = origElem
			case 'v': // Prefix: "voices"
				origElem := elem
				if l := len("voices"); len(elem) >= l && elem[0:l] == "voices" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					// Leaf node.
					switch r.Method {
					case "POST":
						s.handleV1CreateVoiceRequest([0]string{}, elemIsEscaped, w, r)
					default:
						s.notAllowed(w, r, "POST")
					}

					return
				}

				elem = origElem
			}

			elem = origElem
		}
	}
	s.notFound(w, r)
}

// Route is route object.
type Route struct {
	name        string
	summary     string
	operationID string
	pathPattern string
	count       int
	args        [1]string
}

// Name returns ogen operation name.
//
// It is guaranteed to be unique and not empty.
func (r Route) Name() string {
	return r.name
}

// Summary returns OpenAPI summary.
func (r Route) Summary() string {
	return r.summary
}

// OperationID returns OpenAPI operationId.
func (r Route) OperationID() string {
	return r.operationID
}

// PathPattern returns OpenAPI path.
func (r Route) PathPattern() string {
	return r.pathPattern
}

// Args returns parsed arguments.
func (r Route) Args() []string {
	return r.args[:r.count]
}

// FindRoute finds Route for given method and path.
//
// Note: this method does not unescape path or handle reserved characters in path properly. Use FindPath instead.
func (s *Server) FindRoute(method, path string) (Route, bool) {
	return s.FindPath(method, &url.URL{Path: path})
}

// FindPath finds Route for given method and URL.
func (s *Server) FindPath(method string, u *url.URL) (r Route, _ bool) {
	var (
		elem = u.Path
		args = r.args
	)
	if rawPath := u.RawPath; rawPath != "" {
		if normalized, ok := uri.NormalizeEscapedPath(rawPath); ok {
			elem = normalized
		}
		defer func() {
			for i, arg := range r.args[:r.count] {
				if unescaped, err := url.PathUnescape(arg); err == nil {
					r.args[i] = unescaped
				}
			}
		}()
	}

	elem, ok := s.cutPrefix(elem)
	if !ok {
		return r, false
	}

	// Static code generated router with unwrapped path search.
	switch {
	default:
		if len(elem) == 0 {
			break
		}
		switch elem[0] {
		case '/': // Prefix: "/api/v1/"
			origElem := elem
			if l := len("/api/v1/"); len(elem) >= l && elem[0:l] == "/api/v1/" {
				elem = elem[l:]
			} else {
				break
			}

			if len(elem) == 0 {
				break
			}
			switch elem[0] {
			case 'h': // Prefix: "health"
				origElem := elem
				if l := len("health"); len(elem) >= l && elem[0:l] == "health" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					// Leaf node.
					switch method {
					case "GET":
						r.name = V1CheckHealthOperation
						r.summary = ""
						r.operationID = "V1CheckHealth"
						r.pathPattern = "/api/v1/health"
						r.args = args
						r.count = 0
						return r, true
					default:
						return
					}
				}

				elem = origElem
			case 'o': // Prefix: "oauth"
				origElem := elem
				if l := len("oauth"); len(elem) >= l && elem[0:l] == "oauth" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					switch method {
					case "POST":
						r.name = V1AuthorizeUsingOauthOperation
						r.summary = ""
						r.operationID = "V1AuthorizeUsingOauth"
						r.pathPattern = "/api/v1/oauth"
						r.args = args
						r.count = 0
						return r, true
					default:
						return
					}
				}
				switch elem[0] {
				case 'R': // Prefix: "Redirect/"
					origElem := elem
					if l := len("Redirect/"); len(elem) >= l && elem[0:l] == "Redirect/" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'm': // Prefix: "mail"
						origElem := elem
						if l := len("mail"); len(elem) >= l && elem[0:l] == "mail" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch method {
							case "GET":
								r.name = V1AuthorizeUsingOauthHandleMailRedirectOperation
								r.summary = ""
								r.operationID = "V1AuthorizeUsingOauthHandleMailRedirect"
								r.pathPattern = "/api/v1/oauthRedirect/mail"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}

						elem = origElem
					case 'v': // Prefix: "vk"
						origElem := elem
						if l := len("vk"); len(elem) >= l && elem[0:l] == "vk" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch method {
							case "GET":
								r.name = V1AuthorizeUsingOauthHandleVkRedirectOperation
								r.summary = ""
								r.operationID = "V1AuthorizeUsingOauthHandleVkRedirect"
								r.pathPattern = "/api/v1/oauthRedirect/vk"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}

						elem = origElem
					case 'y': // Prefix: "yandex"
						origElem := elem
						if l := len("yandex"); len(elem) >= l && elem[0:l] == "yandex" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch method {
							case "GET":
								r.name = V1AuthorizeUsingOauthHandleYandexRedirectOperation
								r.summary = ""
								r.operationID = "V1AuthorizeUsingOauthHandleYandexRedirect"
								r.pathPattern = "/api/v1/oauthRedirect/yandex"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}

						elem = origElem
					}

					elem = origElem
				}

				elem = origElem
			case 'p': // Prefix: "passwordRecoveryRequests"
				origElem := elem
				if l := len("passwordRecoveryRequests"); len(elem) >= l && elem[0:l] == "passwordRecoveryRequests" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					// Leaf node.
					switch method {
					case "PATCH":
						r.name = V1ConfirmPasswordRecoveryRequestOperation
						r.summary = ""
						r.operationID = "V1ConfirmPasswordRecoveryRequest"
						r.pathPattern = "/api/v1/passwordRecoveryRequests"
						r.args = args
						r.count = 0
						return r, true
					case "POST":
						r.name = V1CreatePasswordRecoveryRequestOperation
						r.summary = ""
						r.operationID = "V1CreatePasswordRecoveryRequest"
						r.pathPattern = "/api/v1/passwordRecoveryRequests"
						r.args = args
						r.count = 0
						return r, true
					default:
						return
					}
				}

				elem = origElem
			case 'r': // Prefix: "registrations"
				origElem := elem
				if l := len("registrations"); len(elem) >= l && elem[0:l] == "registrations" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					// Leaf node.
					switch method {
					case "PATCH":
						r.name = V1ConfirmRegistrationOperation
						r.summary = ""
						r.operationID = "V1ConfirmRegistration"
						r.pathPattern = "/api/v1/registrations"
						r.args = args
						r.count = 0
						return r, true
					case "POST":
						r.name = V1CreateRegistrationOperation
						r.summary = ""
						r.operationID = "V1CreateRegistration"
						r.pathPattern = "/api/v1/registrations"
						r.args = args
						r.count = 0
						return r, true
					default:
						return
					}
				}

				elem = origElem
			case 's': // Prefix: "session"
				origElem := elem
				if l := len("session"); len(elem) >= l && elem[0:l] == "session" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					switch method {
					case "GET":
						r.name = V1CheckCurrentSessionOperation
						r.summary = ""
						r.operationID = "V1CheckCurrentSession"
						r.pathPattern = "/api/v1/session"
						r.args = args
						r.count = 0
						return r, true
					default:
						return
					}
				}
				switch elem[0] {
				case 's': // Prefix: "s"
					origElem := elem
					if l := len("s"); len(elem) >= l && elem[0:l] == "s" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "POST":
							r.name = V1CreateSessionOperation
							r.summary = ""
							r.operationID = "V1CreateSession"
							r.pathPattern = "/api/v1/sessions"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}
					switch elem[0] {
					case '/': // Prefix: "/"
						origElem := elem
						if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
							elem = elem[l:]
						} else {
							break
						}

						// Param: "id"
						// Leaf parameter
						args[0] = elem
						elem = ""

						if len(elem) == 0 {
							// Leaf node.
							switch method {
							case "DELETE":
								r.name = V1DeleteSessionOperation
								r.summary = ""
								r.operationID = "V1DeleteSession"
								r.pathPattern = "/api/v1/sessions/{id}"
								r.args = args
								r.count = 1
								return r, true
							default:
								return
							}
						}

						elem = origElem
					}

					elem = origElem
				}

				elem = origElem
			case 'u': // Prefix: "user"
				origElem := elem
				if l := len("user"); len(elem) >= l && elem[0:l] == "user" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					switch method {
					case "GET":
						r.name = V1GetCurrentUserOperation
						r.summary = ""
						r.operationID = "V1GetCurrentUser"
						r.pathPattern = "/api/v1/user"
						r.args = args
						r.count = 0
						return r, true
					default:
						return
					}
				}
				switch elem[0] {
				case 's': // Prefix: "s"
					origElem := elem
					if l := len("s"); len(elem) >= l && elem[0:l] == "s" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch method {
						case "GET":
							r.name = V1GetUsersOperation
							r.summary = ""
							r.operationID = "V1GetUsers"
							r.pathPattern = "/api/v1/users"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}

					elem = origElem
				}

				elem = origElem
			case 'v': // Prefix: "voices"
				origElem := elem
				if l := len("voices"); len(elem) >= l && elem[0:l] == "voices" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					// Leaf node.
					switch method {
					case "POST":
						r.name = V1CreateVoiceOperation
						r.summary = ""
						r.operationID = "V1CreateVoice"
						r.pathPattern = "/api/v1/voices"
						r.args = args
						r.count = 0
						return r, true
					default:
						return
					}
				}

				elem = origElem
			}

			elem = origElem
		}
	}
	return r, false
}
