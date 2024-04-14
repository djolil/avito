package tests

import (
	"avito/internal/dto"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBanner(t *testing.T) {
	var token string

	t.Run("admin login", func(t *testing.T) {
		req := dto.UserLoginRequest{
			Email:    "baiden@mail.ru",
			Password: "123",
		}
		data, err := json.Marshal(req)
		assert.Nil(t, err, "Ошибка при маршалинге тела запроса")
		buf := bytes.NewBuffer(data)

		res, err := http.Post("http://0.0.0.0:8080/user/login", "application/json", buf)
		assert.Nil(t, err, "Ошибка при авторизации")

		assert.Equal(t, res.StatusCode, http.StatusOK)
		tokenBytes, err := io.ReadAll(res.Body)
		assert.Nil(t, err, "Ошибка при получении токена")

		response := struct {
			Data string `json:"data"`
		}{}
		err = json.Unmarshal(tokenBytes, &response)
		assert.Nil(t, err, "Ошибка при анмаршалинге токена")
		token = response.Data
	})

	var bannerID int

	t.Run("create banner", func(t *testing.T) {
		req := dto.BannerCreateRequest{
			TagIDs:    []uint32{1, 2},
			FeatureID: 4,
			Content: dto.Banner{
				Name: "test",
			},
			IsActive: true,
		}
		data, err := json.Marshal(req)
		assert.Nil(t, err, "Ошибка при маршалинге тела запроса")

		request, err := http.NewRequest("POST", "http://0.0.0.0:8080/banner", bytes.NewBuffer(data))
		assert.Nil(t, err, "Ошибка при создании запроса")

		t.Log("FFFFFFFFFFFFFFFF " + token)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("token", "Bearer "+token)

		client := &http.Client{
			Timeout: time.Second,
		}
		res, err := client.Do(request)
		assert.Nil(t, err, "Ошибка выполнении запроса")
		assert.Equal(t, res.StatusCode, http.StatusCreated)

		bannerBytes, err := io.ReadAll(res.Body)
		assert.Nil(t, err, "Ошибка при получении тела ответа")

		var b dto.BannerCreateResponse
		err = json.Unmarshal(bannerBytes, &b)
		assert.Nil(t, err, "Ошибка при анмаршалинге баннера")

		bannerID = int(b.BannerID)
	})

	t.Run("get banner", func(t *testing.T) {
		request, err := http.NewRequest("GET", "http://0.0.0.0:8080/user_banner?tag_id=1&feature_id=4&use_last_revision=false", nil)
		assert.Nil(t, err, "Ошибка при создании запроса")

		request.Header.Add("token", "Bearer "+token)

		client := &http.Client{
			Timeout: time.Second,
		}
		res, err := client.Do(request)
		assert.Nil(t, err, "Ошибка выполнении запроса")
		assert.Equal(t, res.StatusCode, http.StatusOK)

		bannerBytes, err := io.ReadAll(res.Body)
		assert.Nil(t, err, "Ошибка при получении тела ответа")

		var b dto.BannerResponse
		err = json.Unmarshal(bannerBytes, &b)
		assert.Nil(t, err, "Ошибка при анмаршалинге баннера")

		assert.Equal(t, b.Name, "test")
	})

	t.Run("update banner", func(t *testing.T) {
		req := dto.BannerUpdateRequest{
			TagIDs:    []uint32{1, 2},
			FeatureID: 4,
			Content: dto.Banner{
				Name: "test2",
			},
			IsActive: false,
		}
		data, err := json.Marshal(req)
		assert.Nil(t, err, "Ошибка при маршалинге тела запроса")

		request, err := http.NewRequest("PATCH", "http://0.0.0.0:8080/banner/"+strconv.Itoa(bannerID), bytes.NewBuffer(data))
		assert.Nil(t, err, "Ошибка при создании запроса")

		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("token", "Bearer "+token)

		client := &http.Client{
			Timeout: time.Second,
		}
		res, err := client.Do(request)
		assert.Nil(t, err, "Ошибка выполнении запроса")
		assert.Equal(t, res.StatusCode, http.StatusOK)
	})

	t.Run("get banner", func(t *testing.T) {
		request, err := http.NewRequest("GET", "http://0.0.0.0:8080/user_banner?tag_id=1&feature_id=4&use_last_revision=true", nil)
		assert.Nil(t, err, "Ошибка при создании запроса")

		request.Header.Add("token", "Bearer "+token)

		client := &http.Client{
			Timeout: time.Second,
		}
		res, err := client.Do(request)
		assert.Nil(t, err, "Ошибка выполнении запроса")
		assert.Equal(t, res.StatusCode, http.StatusOK)

		bannerBytes, err := io.ReadAll(res.Body)
		assert.Nil(t, err, "Ошибка при получении тела ответа")

		var b dto.BannerResponse
		err = json.Unmarshal(bannerBytes, &b)
		assert.Nil(t, err, "Ошибка при анмаршалинге баннера")

		assert.Equal(t, b.Name, "test2")
	})
}
