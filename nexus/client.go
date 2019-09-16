package nexus

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/hashicorp/go-getter"
)

// client is used to interact with Nexus API
type client struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
}

// Put executes a put request to upload a file
func (c client) Put(repository, assetPath string, r io.Reader) error {
	u := path.Join("/repository", repository, assetPath)
	u = fmt.Sprintf("%s%s", c.baseURL, u)
	fmt.Println("Upload URL: ", u)

	req, err := http.NewRequest("PUT", u, r)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.username, c.password)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return nil
	case http.StatusUnauthorized:
		return fmt.Errorf("Invalid credentials")
	default:
		return fmt.Errorf("Invalid status code: %v", resp.Status)
	}
}

func (c client) Delete(repository, assetPath string) error {
	u := path.Join("/repository", repository, assetPath)
	u = fmt.Sprintf("%s%s", c.baseURL, u)
	fmt.Println("Delete URL: ", u)

	req, err := http.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.username, c.password)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	default:
		return fmt.Errorf("Error status code %v", resp.Status)
	}
}

func (c client) FileExists(repository, assetPath string) (bool, error) {
	u := path.Join("/repository", repository, assetPath)
	u = fmt.Sprintf("%s%s", c.baseURL, u)

	req, err := http.NewRequest("HEAD", u, nil)
	if err != nil {
		return false, err
	}
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("Error status code %v", resp.Status)
	}
}

func getFileContents(src string) ([]byte, error) {
	u, err := url.Parse(src)
	if err != nil {
		return nil, err
	}

	// Add flag to prevent extract of contents. We're interested in the single
	// file
	q := u.Query()
	q.Set("archive", "false")
	u.RawQuery = q.Encode()

	// Create temp file to store downloaded artifact
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}
	if err := getter.GetFile(f.Name(), u.String()); err != nil {
		return nil, err
	}

	defer os.Remove(f.Name())
	return ioutil.ReadFile(f.Name())
}

func newClient(baseURL, username, password string) *client {
	return &client{
		baseURL:  baseURL,
		username: username,
		password: password,
	}
}
