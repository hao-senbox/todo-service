package uploader

import (
	"context"
	"encoding/json"
	"fmt"
	"todo-service/pkg/constants"
	"todo-service/pkg/consul"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/consul/api"
)

type Avatar struct {
	Url string `json:"url"`
}

type ImageKey struct {
	Key string `json:"key"`
}

type ImageService interface {
	GetImageKey(ctx context.Context, key string) (*Avatar, error)
	DeleteImageKey(ctx context.Context, key string) error
}

type imageService struct {
	client *callAPI
}

type callAPI struct {
	client       consul.ServiceDiscovery
	clientServer *api.CatalogService
}

var (
	imageServiceStr = "go-main-service"
)

func NewImageService(client *api.Client) ImageService {
	mainServiceAPI := NewServiceAPI(client, imageServiceStr)
	return &imageService{
		client: mainServiceAPI,
	}
}

func NewServiceAPI(client *api.Client, serviceName string) *callAPI {
	sd, err := consul.NewServiceDiscovery(client, serviceName)
	if err != nil {
		fmt.Printf("Error creating service discovery: %v\n", err)
		return nil
	}

	var service *api.CatalogService

	for i := 0; i < 10; i++ {
		service, err = sd.DiscoverService()
		if err == nil && service != nil {
			break
		}
		fmt.Printf("Waiting for service %s... retry %d/10\n", serviceName, i+1)
		time.Sleep(3 * time.Second)
	}

	if service == nil {
		fmt.Printf("Service %s not found after retries, continuing anyway...\n", serviceName)
	}

	if os.Getenv("LOCAL_TEST") == "true" {
		fmt.Println("Running in LOCAL_TEST mode — overriding service address to localhost")
		service.ServiceAddress = "localhost"
	}

	return &callAPI{
		client:       sd,
		clientServer: service,
	}
}

func (s *imageService) GetImageKey(ctx context.Context, key string) (*Avatar, error) {

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	image, err := s.client.getImageKey(key, token)

	if err != nil {
		if sc, ok := image["status_code"].(float64); ok && int(sc) == 500 {
			return nil, nil
		}
		return nil, err
	}

	innerData, ok := image["data"].(string)
	if !ok || innerData == "" {
		return nil, nil
	}

	return &Avatar{
		Url: innerData,
	}, nil

}

func (s *imageService) DeleteImageKey(ctx context.Context, key string) error {

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok {
		return fmt.Errorf("token not found in context")
	}

	err := s.client.deleleImage(key, token)

	if err != nil {
		return err
	}

	return nil

}

func (c *callAPI) deleleImage(key string, token string) error {

	endpoint := "/v1/images/delete"

	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	body := map[string]string{
		"key": key,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling body: %v", err)
	}

	_, err = c.client.CallAPI(c.clientServer, endpoint, http.MethodPost, jsonBody, header)
	if err != nil {
		fmt.Printf("Error calling API: %v\n", err)
		return err
	}

	return nil
}

func (c *callAPI) getImageKey(key string, token string) (map[string]interface{}, error) {

	endpoint := "/v1/images"

	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	body := map[string]string{
		"key":  key,
		"mode": "public",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshalling body: %v", err)
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodPost, jsonBody, header)
	if err != nil {
		fmt.Printf("Error calling API: %v\n", err)
		return nil, err
	}

	var imageData interface{}

	err = json.Unmarshal([]byte(res), &imageData)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return nil, err
	}

	myMap := imageData.(map[string]interface{})

	return myMap, nil
}
