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
	Tags		map[string]string
}

type StartSessionRequestBody struct {
	DeviceId    string `json:"deviceId"`
	AppId       string `json:"appId"`
	AutoSync    bool   `json:"autoSync"`
	Screenshots bool   `json:"screenshots"`
	Tags		map[string]string   `json:"tags"`
}

type StopSessionOptions struct {
	IncludeSessionJsonInResponse bool
	OutputDir string
}

type StopSessionRequestBody struct {
	IncludeSessionJsonInResponse bool `json:"includeSessionJsonInResponse"`
	OutputDir *string `json:"outputDir"`
}

type ExecuteShellCommandRequestBody struct {
	Command string `json:"command"`
}

type ExecuteShellCommandResponseBody struct {
	Output string `json:"output"`
}

type ServerVersionInfo struct {
	MajorVersion string `json:"majorVersion"`
	BuildNumber int `json:"buildNumber"`
	CommitHash string `json:"commitHash"`
}

func New(config *Config) *GbaClient {
	if os.Getenv("GBA_BASE_URL") != "" {
		config.BaseUrl = os.Getenv("GBA_BASE_URL")
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

func (c *GbaClient) ListSessions() ([]Session, error) {
	sessions := make([]Session, 0)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/sessions", c.Config.BaseUrl), nil)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(body, &sessions)

	return sessions, nil
}

func (c *GbaClient) StartSession(deviceId string, appId string, options *StartSessionOptions) (*Session, error) {
	var session *Session

	requestBody := &StartSessionRequestBody{
		DeviceId:    deviceId,
		AppId:       appId,
	}

	if options != nil {
		requestBody.AutoSync = options.AutoSync
		requestBody.Screenshots = options.Screenshots
		if options.Tags != nil {
			requestBody.Tags = options.Tags
		}
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

func (c *GbaClient) StopSession(sessionId string, options *StopSessionOptions) (*string, error) {
	requestBody := &StopSessionRequestBody{}

	if options != nil {
		requestBody.IncludeSessionJsonInResponse = options.IncludeSessionJsonInResponse;

		if options.OutputDir != "" {
			requestBody.OutputDir = &options.OutputDir
		}
	}

	encodedRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sessions/%s/stop", c.Config.BaseUrl, sessionId), bytes.NewBuffer(encodedRequestBody))
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

	stringBody := string(body)

	return &stringBody, nil
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

func (c *GbaClient) GetProperties() (map[string]interface{}, error) {
	properties := make(map[string]interface{})

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/properties", c.Config.BaseUrl), nil)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &properties)
	if err != nil {
		return nil, err
	}

	return properties, nil
}

func (c *GbaClient) SetProperties(requestBody map[string]interface{}) error {
	encodedRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/properties", c.Config.BaseUrl), bytes.NewBuffer(encodedRequestBody))
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == 400 || resp.StatusCode == 404 || resp.StatusCode == 500 {
		var errorResponse map[string]string
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return err
		}
		return errors.New(errorResponse["error"])
	}

	return nil
}

func (c *GbaClient) GenerateSessionJson(sessionPath string, targetPath string) error {
	requestBody := make(map[string]interface{})
	requestBody["sessionPath"] = sessionPath
	requestBody["targetPath"] = targetPath

	encodedRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/generate-json", c.Config.BaseUrl), bytes.NewBuffer(encodedRequestBody))
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == 400 || resp.StatusCode == 500 {
		var errorResponse map[string]string
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return err
		}
		return errors.New(errorResponse["error"])
	}

	return nil
}

func (c *GbaClient) EnableWifiProf(deviceId string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/devices/%s/enable-wifi-prof", c.Config.BaseUrl, deviceId), nil)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == 400 || resp.StatusCode == 500 {
		var errorResponse map[string]string
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return err
		}
		return errors.New(errorResponse["error"])
	}

	return nil
}

func (c *GbaClient) DisableWifiProf(deviceId string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/devices/%s/disable-wifi-prof", c.Config.BaseUrl, deviceId), nil)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == 400 || resp.StatusCode == 500 {
		var errorResponse map[string]string
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return err
		}
		return errors.New(errorResponse["error"])
	}

	return nil
}

func (c *GbaClient) GetServerVersionInfo() (*ServerVersionInfo, error) {
	var serverVersionInfo ServerVersionInfo

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/version", c.Config.BaseUrl), nil)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 400 || resp.StatusCode == 500 {
		var errorResponse map[string]string
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errorResponse["error"])
	}

	err = json.Unmarshal(body, &serverVersionInfo)
	if err != nil {
		return nil, err
	}

	return &serverVersionInfo, nil
}

func (c *GbaClient) ExecuteShellCommandOnDevice(deviceId, command string) (*string, error) {
	requestBody := &ExecuteShellCommandRequestBody{
		Command: command,
	}

	encodedRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/devices/%s/shell", c.Config.BaseUrl, deviceId), bytes.NewBuffer(encodedRequestBody))
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

	if resp.StatusCode == 400 || resp.StatusCode == 500 {
		var errorResponse map[string]string
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errorResponse["error"])
	}

	responseBody := &ExecuteShellCommandResponseBody{}

	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody.Output, nil
}
