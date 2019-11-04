package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GbaClient struct {
	HttpClient *http.Client
	Config *Config
}

type Device struct {
	Name string
	Id string
}

type Config struct {
	BaseUrl string
}

type App struct {
	Identifier string
	Debuggable bool
}

type Session struct {
	Id string
}

func New(config *Config) *GbaClient {
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

func (c *GbaClient) StartSession(deviceId string, appId string) (*Session, error) {
	var session *Session

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sessions", c.Config.BaseUrl), nil)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (c *GbaClient) StopSession(sessionId string) error {
	var session *Session

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sessions/%s/stop", c.Config.BaseUrl, sessionId), nil)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, session)
	if err != nil {
		return err
	}

	return nil
}
