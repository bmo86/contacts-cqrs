package means

import (
	"encoding/json"
	"net/http"
)

func MessageErr(code int, message string, w http.ResponseWriter) {
	fields := make(map[string]interface{})
	fields["status"] = "error"
	fields["error"] = message

	msg, err := json.Marshal(fields)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("an error acurred internaly"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(msg)
}
