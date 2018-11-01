package util

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
)

type ExposedError struct {
	PublicMessage  string
	PrivateMessage string
}

func (err *ExposedError) Error() string {
	return err.PublicMessage
}

type ErrorWriter func(w http.ResponseWriter, r *http.Request, data interface{}, code int)

func HTTPRecoveryCustom(errWriter ErrorWriter) Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					var stacktrace string
					for i := 1; ; i++ {
						_, f, l, got := runtime.Caller(i)
						if !got {
							break
						}

						stacktrace += fmt.Sprintf("%s:%d\n", f, l)
					}

					// when stack finishes
					var publicMessage, logMessage string
					logMessage = fmt.Sprintln("Recovered from failed handler")
					switch errmsg := err.(type) {
					case ExposedError:
						logMessage += fmt.Sprintf("Pub Message: %s\n", errmsg.PublicMessage)
						logMessage += fmt.Sprintf("Prv Message: %s\n", errmsg.PrivateMessage)
						// TODO use i18n to translate public messages
						publicMessage = errmsg.PublicMessage
					case error:
						logMessage += fmt.Sprintf("Message: %s\n", errmsg)
						publicMessage = "An internal system error has occurred."
					default:
						logMessage += fmt.Sprint("Unknown error type\n")
						publicMessage = "An unknown error has occurred."
					}
					logMessage += fmt.Sprintf("Message: %s\n", err)
					logMessage += fmt.Sprintf("Trace:\n%s\n", stacktrace)
					log.Print(logMessage)

					errWriter(w, r, publicMessage, http.StatusInternalServerError)
				}
			}()
			handler.ServeHTTP(w, r)
		})
	}
}

func HTTPRecovery(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				var stacktrace string
				for i := 1; ; i++ {
					_, f, l, got := runtime.Caller(i)
					if !got {
						break
					}

					stacktrace += fmt.Sprintf("%s:%d\n", f, l)
				}

				// when stack finishes
				var publicMessage, logMessage string
				logMessage = fmt.Sprintln("Recovered from failed handler")
				switch errmsg := err.(type) {
				case ExposedError:
					logMessage += fmt.Sprintf("Pub Message: %s\n", errmsg.PublicMessage)
					logMessage += fmt.Sprintf("Prv Message: %s\n", errmsg.PrivateMessage)
					// TODO use i18n to translate public messages
					publicMessage = errmsg.PublicMessage
				case error:
					logMessage += fmt.Sprintf("Message: %s\n", errmsg)
					publicMessage = "An internal system error has occurred."
				default:
					logMessage += fmt.Sprint("Unknown error type\n")
					publicMessage = "An unknown error has occurred."
				}
				logMessage += fmt.Sprintf("Message: %s\n", err)
				logMessage += fmt.Sprintf("Trace:\n%s\n", stacktrace)
				log.Print(logMessage)

				http.Error(w, publicMessage, http.StatusInternalServerError)
			}
		}()
		handler.ServeHTTP(w, r)
	})
}
