package getpdfs

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"golang.org/x/net/html"
)

// GetURL returns a document from a URL, retrying in case of error.
func GetURL(u string) (*http.Response, error) {
	var resp *http.Response
	var err error
	retries := 5
	for retries > 0 {
		retries--
		resp, err = http.Get(u)
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	return resp, err
}

// A Link stores the URL and text for a link in an HTML document.
type Link struct {
	Text string
	URL  string
}

// GetLinks returns all links in the given document.
func GetLinks(r io.Reader) []Link {
	var link Link
	var links []Link
	z := html.NewTokenizer(r)
	anchor := false
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		token := z.Token()
		if token.Data == "a" {
			switch token.Type {
			case html.StartTagToken:
				anchor = true
				link = Link{}
			case html.EndTagToken:
				anchor = false
				links = append(links, link)
			}
		}
		if anchor {
			switch token.Type {
			case html.StartTagToken:
				u, err := GetAttribute(token, "href")
				if err == nil {
					link.URL = u
				}
			case html.TextToken:
				link.Text = token.Data
			}
		}
	}
	return links
}

// GetAttribute returns the value for the given attribute key, or an error if none.
func GetAttribute(t html.Token, key string) (string, error) {
	for _, attr := range t.Attr {
		if attr.Key == key {
			return attr.Val, nil
		}
	}
	return "", errors.New("attribute not found")
}

// Filter returns an array of links where f(link) is true.
func Filter(l []Link, f func(Link) bool) []Link {
	r := []Link{}
	for _, link := range l {
		if f(link) {
			r = append(r, link)
		}
	}
	return (r)
}

// AlreadyProcessed returns whether a file with the given name, not including
// file extension, has been successfully processed.
func AlreadyProcessed(name string, project string) (bool, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, project)
	if err != nil {
		return false, err
	}
	defer client.Close()
	fileNoExt := strings.Trim(name, path.Ext(name))
	fileRef := client.Collection("dcgov_files").Doc(fileNoExt)
	snap, err := fileRef.Get(ctx)
	if !snap.Exists() {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	data := snap.Data()
	if ok, _ := data["ok"]; ok == "true" {
		return true, nil
	}
	return false, nil
}

// SaveToBucket saves the contents of file to the given bucket.
func SaveToBucket(file io.Reader, name string, bucket string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	wc := client.Bucket(bucket).Object(name).NewWriter(ctx)
	if _, err = io.Copy(wc, file); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}

// GetPDFs saves all PDFs linked to from the given URL in Google Cloud Storage.
func GetPDFs(u string, bucket string, project string) error {
	resp, err := GetURL(u)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	links := GetLinks(resp.Body)
	pdfs := Filter(links, func(l Link) bool {
		return strings.HasSuffix(l.URL, "pdf")
	})

	for _, link := range pdfs {
		name, _ := url.PathUnescape(path.Base(link.URL))
		if processed, _ := AlreadyProcessed(name, project); !processed {
			log.Printf("Fetching %s", name)

			file, err := GetURL(link.URL)
			if err != nil {
				return err
			}
			err = SaveToBucket(file.Body, name, bucket)
			if err != nil {
				return err
			}
		} else {
			log.Printf("Skipping %s", name)
		}
	}
	return nil
}
