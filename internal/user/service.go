package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"todo-service/pkg/constants"
	"todo-service/pkg/consul"

	"github.com/hashicorp/consul/api"
)

type UserService interface {
	GetUserInfor(ctx context.Context, userID string) (*UserInfor, error)
	GetStudentInfor(ctx context.Context, studentID string) (*UserInfor, error)
	GetTeacherInfor(ctx context.Context, studentID string) (*UserInfor, error)
	GetStaffInfor(ctx context.Context, studentID string) (*UserInfor, error)
	GetListTeacherInfor(ctx context.Context, userID string) ([]*UserInfor, error)
	GetListStaffInfor(ctx context.Context, userID string) ([]*UserInfor, error)
	GetTeacherInforByOrg(ctx context.Context, teacherID, orgID string) (*UserInfor, error)
	GetStaffInforByOrg(ctx context.Context, staffID, orgID string) (*UserInfor, error)
}

type userService struct {
	client *callAPI
}

type callAPI struct {
	client       consul.ServiceDiscovery
	clientServer *api.CatalogService
}

var (
	mainService = "go-main-service"
)

func NewUserService(client *api.Client) UserService {
	mainServiceAPI := NewServiceAPI(client, mainService)
	return &userService{
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

func safeGetString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func safeGetMapString(v interface{}) (map[string]interface{}, bool) {
	if v == nil {
		return nil, false
	}
	m, ok := v.(map[string]interface{})
	return m, ok
}

func parseAvatarSafely(data map[string]interface{}) Avatar {
	var avatar Avatar

	if rawAvatar, exists := data["avatar"]; exists && rawAvatar != nil {
		if avatarMap, ok := safeGetMapString(rawAvatar); ok {
			avatar = Avatar{
				ImageID:  uint64(castToInt64(avatarMap["image_id"])),
				ImageKey: safeGetString(avatarMap["image_key"]),
				ImageUrl: safeGetString(avatarMap["image_url"]),
				Index:    int(castToInt64(avatarMap["index"])),
				IsMain:   castToBool(avatarMap["is_main"]),
			}
		}
	}

	return avatar
}

func parseUserInforSafely(data map[string]interface{}) (*UserInfor, error) {
	if data == nil {
		return nil, nil
	}

	innerData, ok := safeGetMapString(data["data"])
	if !ok || innerData == nil {
		return nil, nil
	}

	avatar := parseAvatarSafely(innerData)

	return &UserInfor{
		UserID:   safeGetString(innerData["id"]),
		UserName: safeGetString(innerData["name"]),
		Avartar:  avatar,
	}, nil
}

func (u *userService) GetUserInfor(ctx context.Context, userID string) (*UserInfor, error) {
	if u.client == nil {
		return nil, fmt.Errorf("client is not initialized")
	}

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := u.client.getUserInfor(userID, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return parseUserInforSafely(data)
}

func (u *userService) GetStudentInfor(ctx context.Context, studentID string) (*UserInfor, error) {
	if u.client == nil {
		return nil, fmt.Errorf("client is not initialized")
	}

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := u.client.getStudentInfor(studentID, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get student info: %w", err)
	}

	return parseUserInforSafely(data)
}

func (u *userService) GetTeacherInfor(ctx context.Context, teacherID string) (*UserInfor, error) {
	if u.client == nil {
		return nil, fmt.Errorf("client is not initialized")
	}

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := u.client.getTeacherInfor(teacherID, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher info: %w", err)
	}

	return parseUserInforSafely(data)
}

func (u *userService) GetStaffInfor(ctx context.Context, staffID string) (*UserInfor, error) {
	if u.client == nil {
		return nil, fmt.Errorf("client is not initialized")
	}

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := u.client.getStaffInfor(staffID, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get staff info: %w", err)
	}

	return parseUserInforSafely(data)
}

func (u *userService) GetListTeacherInfor(ctx context.Context, userID string) ([]*UserInfor, error) {
	if u.client == nil {
		return nil, fmt.Errorf("client is not initialized")
	}

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := u.client.getListTeacherInfor(userID, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher list: %w", err)
	}

	return parseListUserInforSafely(data)
}

func (u *userService) GetListStaffInfor(ctx context.Context, userID string) ([]*UserInfor, error) {
	if u.client == nil {
		return nil, fmt.Errorf("client is not initialized")
	}

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := u.client.getListStaffTeacherInfor(userID, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get staff list: %w", err)
	}

	return parseListUserInforSafely(data)
}

func (u *userService) GetTeacherInforByOrg(ctx context.Context, teacherID, orgID string) (*UserInfor, error) {
	if u.client == nil {
		return nil, fmt.Errorf("client is not initialized")
	}

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := u.client.getTeacherInforByOrg(teacherID, orgID, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher info: %w", err)
	}

	return parseUserInforSafely(data)
}

func (u *userService) GetStaffInforByOrg(ctx context.Context, stafffID, orgID string) (*UserInfor, error) {
	if u.client == nil {
		return nil, fmt.Errorf("client is not initialized")
	}

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := u.client.getStaffInforByOrg(stafffID, orgID, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get staff info: %w", err)
	}

	return parseUserInforSafely(data)
}

func parseListUserInforSafely(data map[string]interface{}) ([]*UserInfor, error) {
	if data == nil {
		return nil, nil
	}

	rawData, ok := data["data"]
	if !ok || rawData == nil {
		return nil, nil
	}

	list, ok := rawData.([]interface{})
	if !ok {
		return nil, nil // Return nil instead of error for non-array data
	}

	var result []*UserInfor
	for _, item := range list {
		if item == nil {
			continue
		}

		if itemMap, ok := safeGetMapString(item); ok {
			avatar := parseAvatarSafely(itemMap)

			ui := &UserInfor{
				UserID:   safeGetString(itemMap["id"]),
				UserName: safeGetString(itemMap["name"]),
				Avartar:  avatar,
			}
			result = append(result, ui)
		}
	}

	return result, nil
}

func (c *callAPI) getUserInfor(userID string, token string) (map[string]interface{}, error) {
	if c == nil || c.client == nil || c.clientServer == nil {
		return nil, fmt.Errorf("client is not properly initialized")
	}

	endpoint := fmt.Sprintf("/v1/gateway/users/%s", userID)
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	if res == "" {
		return nil, nil
	}

	var userData map[string]interface{}
	if err := json.Unmarshal([]byte(res), &userData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return userData, nil
}

func (c *callAPI) getStudentInfor(studentID string, token string) (map[string]interface{}, error) {
	if c == nil || c.client == nil || c.clientServer == nil {
		return nil, fmt.Errorf("client is not properly initialized")
	}

	endpoint := fmt.Sprintf("/v1/gateway/students/%s", studentID)
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	if res == "" {
		return nil, nil
	}

	var userData interface{}
	if err := json.Unmarshal([]byte(res), &userData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if userData == nil {
		return nil, nil
	}

	if myMap, ok := userData.(map[string]interface{}); ok {
		return myMap, nil
	}

	return nil, fmt.Errorf("unexpected response format")
}

func (c *callAPI) getTeacherInfor(teacherID string, token string) (map[string]interface{}, error) {
	if c == nil || c.client == nil || c.clientServer == nil {
		return nil, fmt.Errorf("client is not properly initialized")
	}

	endpoint := fmt.Sprintf("/v1/gateway/teachers/%s", teacherID)
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	if res == "" {
		return nil, nil
	}

	var userData interface{}
	if err := json.Unmarshal([]byte(res), &userData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if userData == nil {
		return nil, nil
	}

	if myMap, ok := userData.(map[string]interface{}); ok {
		return myMap, nil
	}

	return nil, fmt.Errorf("unexpected response format")
}

func (c *callAPI) getStaffInfor(staffID string, token string) (map[string]interface{}, error) {
	if c == nil || c.client == nil || c.clientServer == nil {
		return nil, fmt.Errorf("client is not properly initialized")
	}

	endpoint := fmt.Sprintf("/v1/gateway/staffs/%s", staffID)
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	if res == "" {
		return nil, nil
	}

	var userData interface{}
	if err := json.Unmarshal([]byte(res), &userData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if userData == nil {
		return nil, nil
	}

	if myMap, ok := userData.(map[string]interface{}); ok {
		return myMap, nil
	}

	return nil, fmt.Errorf("unexpected response format")
}

func (c *callAPI) getListTeacherInfor(userID string, token string) (map[string]interface{}, error) {
	if c == nil || c.client == nil || c.clientServer == nil {
		return nil, fmt.Errorf("client is not properly initialized")
	}

	endpoint := fmt.Sprintf("/v1/gateway/teachers/get-by-user/%s", userID)
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	if res == "" {
		return nil, nil
	}

	var userData interface{}
	if err := json.Unmarshal([]byte(res), &userData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if userData == nil {
		return nil, nil
	}

	if myMap, ok := userData.(map[string]interface{}); ok {
		return myMap, nil
	}

	return nil, fmt.Errorf("unexpected response format")
}

func (c *callAPI) getListStaffTeacherInfor(userID string, token string) (map[string]interface{}, error) {
	if c == nil || c.client == nil || c.clientServer == nil {
		return nil, fmt.Errorf("client is not properly initialized")
	}

	endpoint := fmt.Sprintf("/v1/gateway/staffs/get-by-user/%s", userID)
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	if res == "" {
		return nil, nil
	}

	var userData interface{}
	if err := json.Unmarshal([]byte(res), &userData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if userData == nil {
		return nil, nil
	}

	if myMap, ok := userData.(map[string]interface{}); ok {
		return myMap, nil
	}

	return nil, fmt.Errorf("unexpected response format")
}

func (c *callAPI) getTeacherInforByOrg(teacherID, orgID string, token string) (map[string]interface{}, error) {
	if c == nil || c.client == nil || c.clientServer == nil {
		return nil, fmt.Errorf("client is not properly initialized")
	}

	endpoint := fmt.Sprintf("/v1/gateway/teachers/organization/%s/user/%s", orgID, teacherID)
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	if res == "" {
		return nil, nil
	}

	var userData interface{}
	if err := json.Unmarshal([]byte(res), &userData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if userData == nil {
		return nil, nil
	}

	if myMap, ok := userData.(map[string]interface{}); ok {
		return myMap, nil
	}

	return nil, fmt.Errorf("unexpected response format")
}

func (c *callAPI) getStaffInforByOrg(staffID, orgID string, token string) (map[string]interface{}, error) {
	if c == nil || c.client == nil || c.clientServer == nil {
		return nil, fmt.Errorf("client is not properly initialized")
	}

	endpoint := fmt.Sprintf("/v1/gateway/staffs/organization/%s/user/%s", orgID, staffID)
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}
	fmt.Println(res)
	if res == "" {
		return nil, nil
	}

	var userData interface{}
	if err := json.Unmarshal([]byte(res), &userData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if userData == nil {
		return nil, nil
	}

	if myMap, ok := userData.(map[string]interface{}); ok {
		return myMap, nil
	}

	return nil, fmt.Errorf("unexpected response format")
}

func castToInt64(v interface{}) int64 {
	if v == nil {
		return 0
	}

	switch val := v.(type) {
	case float64:
		return int64(val)
	case int:
		return int64(val)
	case int64:
		return val
	case int32:
		return int64(val)
	case string:
		// Try to parse string as number, return 0 if failed
		return 0
	default:
		return 0
	}
}

func castToBool(v interface{}) bool {
	if v == nil {
		return false
	}

	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val == "true" || val == "1"
	case int, int64, int32:
		return val != 0
	case float64:
		return val != 0.0
	default:
		return false
	}
}
