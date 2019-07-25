package controller

import (
	"bufio"
	"fmt"
	"os"

	"github.com/aleph-icabod/json-placeholder/entity"
)

// JsonPlaceholderClient interface to abstract the behavior of the API
type JsonPlaceholderClient interface {
	GetPhotos() ([]*entity.Photo, error)
	GetPhoto(id int) (*entity.Photo, error)
	CreatePhoto(photo *entity.Photo) error
	DeletePhoto(id int) error
	UpdatePhoto(photo *entity.Photo) error
}
type SlackAPI interface {
	SendNotification(message string, photo *entity.Photo) error
}

type Controller struct {
	api   JsonPlaceholderClient
	slack SlackAPI
}

// NewController create a new Controller with JsonPlaceholderClient injected
func NewController(api JsonPlaceholderClient, slack SlackAPI) *Controller {
	return &Controller{
		api:   api,
		slack: slack,
	}
}

// ListPhotos list photos returned by the API
func (c *Controller) ListPhotos() error {
	photos, err := c.api.GetPhotos()
	if err != nil {
		return fmt.Errorf("Something bad: %v\n", err.Error())
	}
	for _, p := range photos {
		fmt.Println(p.ToJsonString())
	}
	return nil
}

func (c *Controller) GetPhoto() error {
	fmt.Println("Input the photo's id: ")
	var id int
	_, err := fmt.Scan(&id)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}
	photo, err := c.api.GetPhoto(id)
	if err != nil {
		return fmt.Errorf("Something bad getting the photo: %v", err)
	}
	fmt.Println("PHOTO found")
	fmt.Println(photo.ToJsonString())
	return nil
}

func (c *Controller) CreatePhoto() error {
	photo, err := c.capturePhoto()
	if err != nil {
		return err
	}
	err = c.api.CreatePhoto(photo)
	if err != nil {
		return fmt.Errorf("Something bad creating photo: %v", err)
	}
	err = c.slack.SendNotification("Someone has create a new photo", photo)
	if err != nil {
		fmt.Println("could not been posible sent notification to slack: ", err)
	}
	fmt.Println(photo.ToJsonString())

	return nil
}

func (c *Controller) UpdatePhoto() error {
	fmt.Println("Input the photo's id: ")
	var id int
	_, err := fmt.Scan(&id)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}
	photo, err := c.capturePhoto()
	if err != nil {
		return err
	}
	photo.ID = id
	err = c.api.UpdatePhoto(photo)
	if err != nil {
		return fmt.Errorf("Something bad creating photo: %v", err)
	}
	fmt.Println(photo.ToJsonString())
	return nil
}

func (c *Controller) DeletePhoto() error {
	fmt.Println("Input the photo's id: ")
	var id int
	_, err := fmt.Scan(&id)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}
	photo, err := c.api.GetPhoto(id)
	if err != nil {
		return fmt.Errorf("error deleting the photo: %v")
	}
	err = c.api.DeletePhoto(id)
	if err != nil {
		return fmt.Errorf("Something bad deleting the photo: %v", err)
	}
	err = c.slack.SendNotification(fmt.Sprintf("Someone has deleted a photo"), photo)
	if err != nil {
		fmt.Println("could not been posible sent notification to slack: ", err)
	}
	fmt.Println("PHOTO deleted")
	return nil
}

// capturePhoto helper to get data for create a new Photo
func (c *Controller) capturePhoto() (*entity.Photo, error) {
	scanner := bufio.NewScanner(os.Stdin)
	newPhoto := entity.Photo{}
	fmt.Println("Input the name for the photo: ")
	scanner.Scan()
	newPhoto.Title = scanner.Text()
	fmt.Println("Input the album id for the photo: ")
	_, err := fmt.Scan(&newPhoto.AlbumID)
	if err != nil {
		return nil, fmt.Errorf("Error getting album id for the photo: %v", err)
	}
	// validate if the albumID exist (in the API just exist 100 albums)
	if newPhoto.AlbumID > 100 || newPhoto.AlbumID < 1 {
		return nil, fmt.Errorf("Invalid album id, the album does not exist")
	}
	fmt.Println("Input the url for the photo: ")
	scanner.Scan()
	newPhoto.URL = scanner.Text()
	fmt.Println("Input the thumbnail url for the photo: ")
	scanner.Scan()
	newPhoto.ThumbnailURL = scanner.Text()
	return &newPhoto, nil
}
