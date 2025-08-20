package main

import "net/http"

func (c *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if c.platform == "dev" {
		c.db.DeleteUsers(r.Context())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Deleted all users records"))
		return
	}
	w.WriteHeader(403)
	w.Write([]byte("403 Forbiden"))
}
