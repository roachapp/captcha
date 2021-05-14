
package main

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	log "github.com/sirupsen/logrus"

	"github.com/roachapp/captcha"
)

func serve(w http.ResponseWriter, r *http.Request, id string) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("captcha-ID", id)

	var content bytes.Buffer
	w.Header().Set("Content-Type", "image/png")
	if err := captcha.WriteImage(&content, id, 240, 80); err != nil {
		log.Fatal(err)
		return err
	}

	http.ServeContent(w, r, id, time.Time{}, bytes.NewReader(content.Bytes()))
	return nil
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	id := captcha.New()

	if r.FormValue("reload") != "" {
		captcha.Reload(id)
	}

	if serve(w, r, id) == captcha.ErrNotFound {
		http.NotFound(w, r)
	}
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	captchaID := r.FormValue("id")
	captchaSolution := r.FormValue("sol")

	if !captcha.VerifyString(captchaID, captchaSolution) {
		_, err := io.WriteString(w, "{\"message\": \"try again :(\", \"status\": 400}\n")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		_, err := io.WriteString(w, "{\"message\": \"that went smoothly :)\", \"status\": 200}\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	http.Handle("/", tollbooth.LimitFuncHandler(tollbooth.NewLimiter(0.01, nil), getHandler))
	http.Handle("/validate", tollbooth.LimitFuncHandler(tollbooth.NewLimiter(0.01, nil), validateHandler))
	log.Debugf("Captcha Server running on %s", "localhost:8666")

	if err := http.ListenAndServe("localhost:8666", nil); err != nil {
		log.Fatal(err)
	}
}
