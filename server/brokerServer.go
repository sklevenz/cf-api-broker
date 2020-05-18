package server

import (
	"fmt"
	"net/http"
)

func BrokerServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "20")
}
