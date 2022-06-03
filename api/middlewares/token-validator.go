package middlewares

import (
	"context"
	"encoding/json"
	json_web_token "github.com/mecamon/shoppingify-server/core/json-web-token"
	"net/http"
)

func TokenValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customClaims, err := json_web_token.Validate(r.Header.Get("Authorization"))
		if err != nil {
			errorMap := map[string]interface{}{"unauthorized": "You have not access to this area."}
			output, _ := json.MarshalIndent(errorMap, "", "    ")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(output)
			return
		}

		ctx := context.WithValue(r.Context(), "ID", customClaims.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
