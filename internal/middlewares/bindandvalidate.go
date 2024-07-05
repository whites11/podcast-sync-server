package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"

	"github.com/whites11/podcast-sync-server/internal/errorrenderer"
	"github.com/whites11/podcast-sync-server/internal/key"
)

func BindAndValidate[T any](next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req T
		// parse json body
		ct := r.Header.Get("Content-Type")
		mediatype, _, err := mime.ParseMediaType(ct)
		if err == nil && mediatype == "application/json" {
			if r.Body != http.NoBody {
				// Look for a field in T with tag "jsonbodyfield"
				adjustedBody, err := io.ReadAll(r.Body)
				if err != nil {
					errorrenderer.Render(w, r, err, http.StatusBadRequest)
					return
				}
				t := reflect.TypeOf(req)
				for i := 0; i < t.NumField(); i++ {
					field := t.Field(i)
					tag := field.Tag.Get("jsonbodyfield")
					if tag != "" {
						adjustedBody = []byte(fmt.Sprintf("{ \"%s\": %s }", tag, string(adjustedBody)))
						break
					}
				}
				err = json.Unmarshal(adjustedBody, &req)
				if err != nil {
					errorrenderer.Render(w, r, err, http.StatusBadRequest)
					return
				}
			}
		}

		if reflect.TypeOf(req).Kind().String() != "slice" {
			// read path values
			{
				t := reflect.TypeOf(req)
				for i := 0; i < t.NumField(); i++ {
					field := t.Field(i)
					tag := field.Tag.Get("pathvalue")
					if tag != "" {
						reflect.ValueOf(&req).Elem().FieldByName(field.Name).SetString(r.PathValue(tag))
					}
				}
			}

			// read query string values
			{
				t := reflect.TypeOf(req)
				for i := 0; i < t.NumField(); i++ {
					field := t.Field(i)
					tag := field.Tag.Get("querystring")
					if tag != "" {
						f := reflect.ValueOf(&req).Elem().FieldByName(field.Name)
						switch f.Type().Name() {
						case "uint":
							v, err := strconv.ParseUint(r.URL.Query().Get(tag), 10, 64)
							if err != nil {
								errorrenderer.Render(w, r, fmt.Errorf("expected query string parameter %s to be uint", tag), http.StatusBadRequest)
								return
							}
							f.SetUint(v)
						case "string":
							f.SetString(r.URL.Query().Get(tag))
						case "int64":
							v, err := strconv.ParseInt(r.URL.Query().Get(tag), 10, 64)
							if err != nil {
								fmt.Println(err)
								errorrenderer.Render(w, r, fmt.Errorf("expected query string parameter %s to be int64", tag), http.StatusBadRequest)
								return
							}
							f.SetInt(v)
						case "bool":
							v := r.URL.Query().Get(tag) == "true"
							f.SetBool(v)
						default:
							errorrenderer.Render(w, r, fmt.Errorf("unhandled field type %s", f.Type().Name()), http.StatusInternalServerError)
							return
						}
					}
				}
			}

			// validate all fields
			{
				validate := validator.New(validator.WithRequiredStructEnabled())

				if err := validate.Struct(req); err != nil {
					errorrenderer.Render(w, r, err, http.StatusBadRequest)
					return
				}
			}
		}

		ctxWithReq := context.WithValue(r.Context(), key.CONTEXT_VALUE_REQUEST_KEY, req)

		next.ServeHTTP(w, r.WithContext(ctxWithReq))
	})
}
