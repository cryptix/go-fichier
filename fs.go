package fichier

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"

	"gopkg.in/errgo.v1"
)

type Client struct {
	http http.Client
}

func NewClient(u, p string) (*Client, error) {
	c := Client{}

	j, err := cookiejar.New(nil)
	if err != nil {
		return nil, errgo.Notef(err, "fichier: could not create cookiejar")
	}
	c.http.Jar = j

	// TODO: options for purge
	lvals := make(url.Values)
	lvals.Set("mail", u)
	lvals.Set("pass", p)
	lvals.Set("lt", "on")
	lvals.Set("restrict", "on")

	resp, err := c.http.PostForm("https://1fichier.com/login.pl", lvals)
	if err != nil {
		return nil, errgo.Notef(err, "fichier: login PostForm failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errgo.Newf("fichier: login response: %s", resp.Status)
	}

	return &c, nil
}

type GetInfoResp struct {
	NumberOfFiles int
	UsedSpace     string
	TotalAccess   int
}

// GetInfo ...
func (c *Client) GetInfo() (*GetInfoResp, error) {
	req, err := http.NewRequest("GET", "https://1fichier.com/console/infog.pl", nil)
	if err != nil {
		return nil, errgo.Notef(err, "fichier: could not create GetInfo request")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, errgo.Notef(err, "fichier: GetInfo request failed")
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, errgo.Notef(err, "fichier: failed to create goquery document from response")
	}

	sel := doc.Find("table.premium td:nth-child(2)")

	if sel.Length() != 3 {
		return nil, errgo.Notef(err, "fichier: failed to find the expected number of fields")
	}

	var r GetInfoResp
	r.NumberOfFiles, err = strconv.Atoi(sel.Eq(0).Text())
	if err != nil {
		return nil, errgo.Notef(err, "fichier: could not convert NumOfFiles HTML response to int: %q", sel.Eq(0).Text())
	}

	r.UsedSpace = sel.Eq(1).Text()

	r.TotalAccess, err = strconv.Atoi(sel.Eq(2).Text())
	if err != nil {
		return nil, errgo.Notef(err, "fichier: could not convert TotalAccess HTML response to int: %q", sel.Eq(2).Text())
	}

	return &r, nil
}

type Dir struct {
	Id   uint
	Name string
}

func (c *Client) Dirs(parent int) ([]Dir, error) {
	u, err := url.Parse("https://1fichier.com/console/dirs.pl")
	if err != nil {
		return nil, errgo.Notef(err, "fichier: could not parse dirs url")
	}

	v := u.Query()
	v.Set("dir_id", strconv.Itoa(parent))
	v.Set("map", "0") // TODO: figure out what this is for
	u.RawQuery = v.Encode()()

	req, err := http.NewRequest("GET", u.Encode(), nil)
	if err != nil {
		return nil, errgo.Notef(err, "fichier: could not create Dirs() request")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, errgo.Notef(err, "fichier: Dirs request failed")
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, errgo.Notef(err, "fichier: failed to create goquery document from response")
	}

	return nil, errgo.New("TODO")
}
