package fichier

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/errgo.v1"
)

func GetUploadHost() (string, error) {
	doc, err := goquery.NewDocument("https://1fichier.com/")
	if err != nil {
		return "", errgo.Notef(err, "fichier: could not retreive document")
	}

	action, ex := doc.Find("#files").Attr("action")
	if !ex {
		return "", errgo.Newf("fichier: could not find action attribute")
	}

	return action, nil
}

// UploadFile takes a filename and a reader and uploads the file to 1fichier.
// It returns the download and the delete link as seperate strings,
// or an error..
func UploadFile(fname string, r io.Reader) (string, string, error) {
	ulHost, err := GetUploadHost()
	if err != nil {
		return "", "", errgo.Notef(err, "fichier: UploadFile() failed to prepare upload host")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file[]", filepath.Base(fname))
	if err != nil {
		return "", "", errgo.Notef(err, "fichier: could not create multipart FormFile field")
	}

	if _, err := io.Copy(part, r); err != nil {
		return "", "", errgo.Notef(err, "fichier: could not copy data to multipart field")
	}

	if err := writer.WriteField("send_ssl", "on"); err != nil {
		return "", "", errgo.Notef(err, "fichier: could not copy data to multipart field")
	}

	if err := writer.Close(); err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("POST", ulHost, body)
	if err != nil {
		return "", "", errgo.Notef(err, "fichier: NewReq() failed")
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", errgo.Notef(err, "fichier: http request failed")
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", errgo.Newf("fichier: upload status not OK: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", "", errgo.Notef(err, "fichier: http request failed")
	}

	tr := doc.Find("table.premium tr")

	if l := tr.Length(); l != 2 {
		return "", "", errgo.Newf("fichier: table rows - wrong length: %d", l)
	}
	td := tr.Last()
	lnk, ok := td.Find("a").Attr("href")
	if !ok {
		return "", "", errgo.Newf("fichier: did not find download href - %q", td.Text())
	}

	urlDL, err := url.Parse(lnk)
	if err != nil {
		return "", "", errgo.Notef(err, "fichier: could not parse download link: %q", lnk)
	}

	lnkRm := td.Find("td").Last().Text()
	urlRM, err := url.Parse(lnkRm)
	if err != nil {
		return "", "", errgo.Notef(err, "fichier: could not parse remove link: %q", lnkRm)
	}

	return urlDL.String(), urlRM.String(), nil
}
