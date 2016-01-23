package httpcfg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Get the current lookupd addresses from config endpoint
func Get(cfgURL string) ([]string, error) {
	res, err := http.Get(cfgURL)
	if err != nil {
		return []string{}, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []string{}, err
	}

	addrs := []string{}
	json.Unmarshal(body, &addrs)
	return addrs, nil
}

// Set lookupd addresses on config endpoint
func Set(cfgURL string, addrs []string) error {
	body, err := json.Marshal(addrs)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", cfgURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("type=error msg=unable to set config status=%d", res.StatusCode)
	}

	return nil
}
