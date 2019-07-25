package apiClient

import (
	"encoding/json"
	"fmt"
	"github.com/aleph-icabod/json-placeholder/entity"
	"io/ioutil"
	"net/http"
	"strings"
)

// PhotosClient struct for implement specific HTTP behavior
type PhotosClient struct {
	BaseURL string
	client  *http.Client
}

func (c *PhotosClient) GetPhotos() ([]*entity.Photo, error) {
	request, err := http.NewRequest(http.MethodGet, c.BaseURL, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the API responses with a %v status", response.StatusCode)
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var photos []*entity.Photo
	err = json.Unmarshal(data, &photos)
	if err != nil {
		return nil, err
	}
	return photos, nil
}

func (c *PhotosClient) GetPhoto(id int) (*entity.Photo, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%d", c.BaseURL, id), nil)
	if err != nil {
		return nil, err
	}
	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("photo with the id %v does not exist", id)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the API responses with a %v status", response.StatusCode)
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var photo entity.Photo
	err = json.Unmarshal(data, &photo)
	if err != nil {
		return nil, err
	}
	return &photo, nil
}

func (c *PhotosClient) CreatePhoto(photo *entity.Photo) error {
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s", c.BaseURL),
		strings.NewReader(photo.ToJsonString()))
	if err != nil {
		return err
	}
	response, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("the API responses with a %v status", response.StatusCode)
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, photo)
	if err != nil {
		return err
	}
	return nil
}

func (c *PhotosClient) DeletePhoto(id int) error {
	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%d", c.BaseURL, id), nil)
	if err != nil {
		return err
	}
	response, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNotFound {
		return fmt.Errorf("the photo with id %d has not found", id)
	}
	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return fmt.Errorf("the API responses with a %v status", response.StatusCode)
	}
	return nil
}

func (c *PhotosClient) UpdatePhoto(photo *entity.Photo) error {
	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%d", c.BaseURL, photo.ID),
		strings.NewReader(photo.ToJsonString()))
	if err != nil {
		return err
	}
	response, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNotFound {
		return fmt.Errorf("the photo with id %d has not found", photo.ID)
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("the API responses with a %v status", response.StatusCode)
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, photo)
	if err != nil {
		return err
	}
	return nil
}

// NewClient constructor to create a new JsonPlaceholderClient
func NewClient(url string) *PhotosClient {
	return &PhotosClient{
		BaseURL: url,
		client:  http.DefaultClient,
	}
}
