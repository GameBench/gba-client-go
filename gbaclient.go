package gba

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type GbaClient struct {
	HttpClient *http.Client
	Config     *Config
}

type Device struct {
	Name string
	Id   string
}

type Config struct {
	BaseUrl  string
	Username string
	Password string
	Token    string
}

func (c *Config) UseToken() bool {
	if c.Token != "" {
		return true
	}

	return false
}

func (c *Config) GetAuthPassword() string {
	if c.Password != "" {
		return c.Password
	}

	if c.Token != "" {
		return c.Token
	}

	return ""
}

type App struct {
	Identifier string
	Debuggable bool
}

type Session struct {
	Id string
}

type StartSessionOptions struct {
	AutoSync    bool
	Screenshots bool
}

type StartSessionRequestBody struct {
	DeviceId    string `json:"deviceId"`
	AppId       string `json:"appId"`
	Username    string `json:"username"`
	PassOrToken string `json:"passOrToken"`
	UseToken    bool   `json:"useToken"`
	AutoSync    bool   `json:"autoSync"`
	Screenshots bool   `json:"screenshots"`
}

func New(config *Config) *GbaClient {
	if os.Getenv("GBA_BASE_URL") != "" {
		config.BaseUrl = os.Getenv("GBA_BASE_URL")
	}

	if os.Getenv("GBA_USERNAME") != "" {
		config.Username = os.Getenv("GBA_USERNAME")
	}

	if os.Getenv("GBA_PASSWORD") != "" {
		config.Password = os.Getenv("GBA_PASSWORD")
	}

	if os.Getenv("GBA_TOKEN") != "" {
		config.Token = os.Getenv("GBA_TOKEN")
	}

	client := &http.Client{}
	return &GbaClient{HttpClient: client, Config: config}
}

func (c *GbaClient) ListDevices() ([]Device, error) {
	devices := make([]Device, 0)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/devices", c.Config.BaseUrl), nil)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(body, &devices)

	return devices, nil
}

func (c *GbaClient) GetDevice(deviceId string) (*Device, error) {
	var device *Device

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/devices/%s", c.Config.BaseUrl, deviceId), nil)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, errors.New("device not found")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &device)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (c *GbaClient) GetDeviceApps(deviceId string) ([]App, error) {
	var apps []App

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/devices/%s/apps", c.Config.BaseUrl, deviceId), nil)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, errors.New("device not found")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &apps)
	if err != nil {
		return nil, err
	}

	return apps, nil
}

func (c *GbaClient) StartSession(deviceId string, appId string, options *StartSessionOptions) (*Session, error) {
	var session *Session

	requestBody := &StartSessionRequestBody{
		DeviceId:    deviceId,
		AppId:       appId,
		Username:    c.Config.Username,
		PassOrToken: c.Config.GetAuthPassword(),
		UseToken:    c.Config.UseToken(),
	}

	if options != nil {
		requestBody.AutoSync = options.AutoSync
		requestBody.Screenshots = options.Screenshots
	}

	encodedRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sessions", c.Config.BaseUrl), bytes.NewBuffer(encodedRequestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 400 || resp.StatusCode == 401 {
		var errorResponse map[string]string
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errorResponse["error"])
	}

	err = json.Unmarshal(body, &session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (c *GbaClient) StopSession(sessionId string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sessions/%s/stop", c.Config.BaseUrl, sessionId), nil)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *GbaClient) Sync() error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sessions/sync", c.Config.BaseUrl), nil)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
