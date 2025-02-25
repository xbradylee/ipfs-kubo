package main

import (
	"net/http"

	"github.com/xbradylee/ipfs-kubo/profile"
)

func init() {
	http.HandleFunc("/debug/stack",
		func(w http.ResponseWriter, _ *http.Request) {
			_ = profile.WriteAllGoroutineStacks(w)
		},
	)
}
