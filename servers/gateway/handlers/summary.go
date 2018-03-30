package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

const headerAccessControlAllowOrigin = "Access-Control-Allow-Origin"
const contentTypeHTML = "text/html"

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	/*TODO: add code and additional functions to do the following:
	- Add an HTTP header to the response with the name
	 `Access-Control-Allow-Origin` and a value of `*`. This will
	  allow cross-origin AJAX requests to your server.
	- Get the `url` query string parameter value from the request.
	  If not supplied, respond with an http.StatusBadRequest error.
	- Call fetchHTML() to fetch the requested URL. See comments in that
	  function for more details.
	- Call extractSummary() to extract the page summary meta-data,
	  as directed in the assignment. See comments in that function
	  for more details
	- Close the response HTML stream so that you don't leak resources.
	- Finally, respond with a JSON-encoded version of the PageSummary
	  struct. That way the client can easily parse the JSON back into
	  an object. Remember to tell the client that the response content
	  type is JSON.

	Helpful Links:
	https://golang.org/pkg/net/http/#Request.FormValue
	https://golang.org/pkg/net/http/#Error
	https://golang.org/pkg/encoding/json/#NewEncoder
	*/

	w.Header().Add(headerAccessControlAllowOrigin, "*")
	pageURL := r.URL.Query().Get("url")

	//STOP PROGRAM EXECUTION ????
	if len(pageURL) == 0 {
		http.Error(w, "Missing url query string parameter", http.StatusBadRequest)
		log.Fatalf("Missing url query string parameter: %v", http.StatusBadRequest)
	}

	html, err := fetchHTML(pageURL)
	if err != nil {
		log.Fatalf("Error in fetching URL: %v", err)
	}

	summary, err := extractSummary(pageURL, html)

	if err != nil {
		log.Fatalf("Error in extracting summary: %v", err)
	}

	json.NewEncoder(w).Encode(summary)
	//Need it here too???
	html.Close()
}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	/*TODO: Do an HTTP GET for the page URL. If the response status
	code is >= 400, return a nil stream and an error. If the response
	content type does not indicate that the content is a web page, return
	a nil stream and an error. Otherwise return the response body and
	no (nil) error.

	To test your implementation of this function, run the TestFetchHTML
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestFetchHTML

	Helpful Links:
	https://golang.org/pkg/net/http/#Get
	*/

	response, err := http.Get(pageURL)
	code := response.StatusCode
	contentType := response.Header.Get("Content-type")

	if !strings.HasPrefix(contentType, contentTypeHTML) {
		return nil, fmt.Errorf("Content type of response is not a web page, it is: %v", contentType)
	}

	if code >= 400 {
		return nil, fmt.Errorf("Bad Request %v", code)
	}

	if err != nil {
		return nil, fmt.Errorf("Error while getting url: %v", err)
	}

	return response.Body, nil
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	/*TODO: tokenize the `htmlStream` and extract the page summary meta-data
	according to the assignment description.

	To test your implementation of this function, run the TestExtractSummary
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestExtractSummary

	Helpful Links:
	https://drstearns.github.io/tutorials/tokenizing/
	http://ogp.me/
	https://developers.facebook.com/docs/reference/opengraph/
	https://golang.org/pkg/net/url/#URL.ResolveReference
	*/

	tokenizer := html.NewTokenizer(htmlStream)
	extracted := map[string]string{}
	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}
			log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
		}

		if tokenType == html.StartTagToken {
			token := tokenizer.Token()

			if token.Data == "head" {
				for {
					tokenType = tokenizer.Next()
					headToken := tokenizer.Token()

					if tokenType == html.EndTagToken && headToken.Data == "head" {
						break
					}

					if headToken.Data == "meta" {
						attr := headToken.Attr
						attrInfo := map[string]string{}
						//make function
						for _, att := range attr {
							log.Printf("Attribute: %s", att.Key)
							if att.Key == "property" {
								attrInfo["property"] = att.Val
							} else if att.Key == "content" {
								attrInfo["content"] = att.Val
							} else if att.Key == "name" {
								if att.Val == "description" {
									attrInfo["property"] = "og:description"
								}
							}
						}
						extracted[attrInfo["property"]] = attrInfo["content"]
					} else if headToken.Data == "title" && tokenType == html.StartTagToken {
						tokenType = tokenizer.Next()
						extracted["og:title"] = tokenizer.Token().Data
					}
					// else if headToken.Data == "link" {
					// 	attr := headToken.Attr

					// 	for _, att := range attr {

					// 	}

					// }

				}
			}
		}
	}
	p, err := constructSummary(extracted)

	if err != nil {
		log.Fatalf("Something wrong with construct summary: %v", err)
	}

	return p, nil
}

//constructSummary constructs pagesummary struct
func constructSummary(expected map[string]string) (*PageSummary, error) {
	images := make([]*PreviewImage, 0, 0)
	img := &PreviewImage{
		URL: expected["og:image"],
	}
	images = append(images, img)
	p := &PageSummary{
		Type:        expected["og:type"],
		URL:         expected["og:url"],
		Title:       expected["og:title"],
		SiteName:    expected["og:site_name"],
		Description: expected["og:description"],
		Author:      expected["og:author"],
		Images:      images,
	}
	return p, nil
}
