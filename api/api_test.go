package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const serverStartupMillis = 1000

const testUsername = "username"
const testPassword = "Super@35Secure"
const testCreateUserPassword = "somepassword"
const testBitcoinP2PKHAddress = "1DxBaADfhTSWsevbzDghrhKSqQwsBpuM5A"
const testDomain = "example.com"
const TestNanoAddress = "xrb_3xnpp3eh6fhnfztx46ypubizd5q1fgds3dbbkp5ektwut3tumrykyx6u5qpd"

func TestAPISuccess(t *testing.T) {
	server := Start()
	defer server.Shutdown(nil)

	time.Sleep(serverStartupMillis * time.Millisecond) //wait for server to start

	// Create a user
	url := "http://localhost:" + os.Getenv("PORT") + "/v1/users"
	params := []byte(`{
		"alias": "` + testUsername + `$` + testDomain + `",
		"password": "` + testPassword + `",
		"create_user_password": "` + testCreateUserPassword + `"
		}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(params))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	// Login to that user
	url = "http://localhost:" + os.Getenv("PORT") + "/v1/auth"
	params = []byte(`{
		"alias": "` + testUsername + "$" + testDomain + `",
		"password": "` + testPassword + `"
		}`)

	req, err = http.NewRequest("POST", url, bytes.NewBuffer(params))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)

	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)

	authResponse := postAuthResponse{}
	err = json.Unmarshal(body, &authResponse)
	assert.Nil(t, err)

	// Add a Bitcoin address
	url = "http://localhost:" + os.Getenv("PORT") + "/v1/addresses"
	params = []byte(`{
		"address_type": 100,
		"address": "` + testBitcoinP2PKHAddress + `"
		}`)

	req, err = http.NewRequest("PUT", url, bytes.NewBuffer(params))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResponse.Token)
	resp, err = client.Do(req)

	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	// Add a Nano address
	url = "http://localhost:" + os.Getenv("PORT") + "/v1/addresses"
	params = []byte(`{
		"address_type": 300,
		"address": "` + TestNanoAddress + `"
		}`)

	req, err = http.NewRequest("PUT", url, bytes.NewBuffer(params))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResponse.Token)
	resp, err = client.Do(req)

	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	// Get an address
	url = "http://localhost:" + os.Getenv("PORT") + "/v1/addresses?alias=" + testUsername + "$" + testDomain + "&address_type=100"
	req, err = http.NewRequest("GET", url, bytes.NewBuffer(params))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)

	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	// Get addresses
	url = "http://localhost:" + os.Getenv("PORT") + "/v1/addresses?alias=" + testUsername + "$" + testDomain
	req, err = http.NewRequest("GET", url, bytes.NewBuffer(params))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	assert.True(t, len(body) > 0)

	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, len(body) > 0)

	// Delete an address
	url = "http://localhost:" + os.Getenv("PORT") + "/v1/addresses/100"
	req, err = http.NewRequest("DELETE", url, bytes.NewBuffer(params))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResponse.Token)
	resp, err = client.Do(req)

	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	// Delete an address
	url = "http://localhost:" + os.Getenv("PORT") + "/v1/addresses/300"
	req, err = http.NewRequest("DELETE", url, bytes.NewBuffer(params))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResponse.Token)
	resp, err = client.Do(req)

	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	// Get address and fail
	url = "http://localhost:" + os.Getenv("PORT") + "/v1/addresses?alias=" + testUsername + "$" + testDomain + "&address_type=100"
	req, err = http.NewRequest("GET", url, bytes.NewBuffer(params))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 404, resp.StatusCode)

	// Delete the user
	url = "http://localhost:" + os.Getenv("PORT") + "/v1/users"
	req, err = http.NewRequest("DELETE", url, bytes.NewBuffer(params))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResponse.Token)
	resp, err = client.Do(req)

	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	// Fail the login
	url = "http://localhost:" + os.Getenv("PORT") + "/v1/auth"
	params = []byte(`{
		"alias": "` + testUsername + "$" + testDomain + `",
		"password": "` + testPassword + `"
		}`)

	req, err = http.NewRequest("POST", url, bytes.NewBuffer(params))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)

	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 400, resp.StatusCode)
}
