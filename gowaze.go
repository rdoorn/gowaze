package gowaze

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	username string
	password string
	apiURL   string
	cookie   []*http.Cookie
}

var (
	apiURL = "https://www.waze.com"
)

const (
	urlLogin = "%s/login/get"
	//urlGetRoute = "%s/ow-RoutingManager/routingRequest?at=0&clientVersion=4.0.0&from=x%3A4.7515934%20y%3A52.2774235&nPaths=3&options=AVOID_TRAILS%3At%2CALLOW_UTURNS%3At&returnGeometries=true&returnInstructions=true&returnJSON=true&timeout=60000&to=x%3A4.834963%20y%3A52.645738"
	urlGetRoute = "%s/row-RoutingManager/routingRequest?at=0&clientVersion=4.0.0&from=x%%3A%.7f%%20y%%3A%.7f&nPaths=3&options=AVOID_TRAILS%%3At%%2CALLOW_UTURNS%%3At&returnGeometries=true&returnInstructions=true&returnJSON=true&timeout=60000&to=x%%3A%.7f%%20y%%3A%.7f"
)

func New() *Handler {
	return &Handler{
		apiURL: apiURL,
	}
}

func (h *Handler) Get(method, url string, args ...interface{}) ([]byte, error) {

	req, err := http.NewRequest(method, fmt.Sprintf(url, args...), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Referer", "https://www.waze.com/livemap?utm_source=waze_website&utm_campaign=waze_website")

	if url == urlLogin {
		//req.Header.Add("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("CPE/%s:%s", h.username, h.password)))))
	}
	// add cookies if any
	for _, c := range h.cookie {
		req.AddCookie(c)
	}

	client := &http.Client{}

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		for key, val := range via[0].Header {
			req.Header[key] = val
		}
		return nil
	}

	log.Printf("Request: %+v", req)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	for _, cookie := range res.Cookies() {
		delete := -1
		for id, ec := range h.cookie {
			if cookie.Name == ec.Name {
				delete = id
			}
		}
		if delete >= 0 {
			h.cookie = append(h.cookie[:delete], h.cookie[delete+1:]...)
		}
		h.cookie = append(h.cookie, cookie)
	}
	body, _ := ioutil.ReadAll(res.Body)
	log.Printf("body: %s", body)

	return body, nil
}

func (h *Handler) Login() error {
	// get a cookie if we have none
	if h.cookie == nil {
		if _, err := h.Get("POST", urlLogin, h.apiURL); err != nil {
			return err
		}
	}
	// get new cookie if it expired
	for id, c := range h.cookie {
		if c.Name == "_web_session" {
			if h.cookie[id].Expires.Before(time.Now()) {
				if _, err := h.Get("POST", urlLogin, h.apiURL); err != nil {
					return err
				}
			}

		}
	}
	// we already have a non-expired cookie
	return nil
}
