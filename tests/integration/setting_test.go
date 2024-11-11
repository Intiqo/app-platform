package integration

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gofrs/uuid/v5"

	"github.com/Intiqo/app-platform/internal/domain"
	"github.com/Intiqo/app-platform/tests/helper"
)

func TestFilterSettingsByCriteria(t *testing.T) {
	t.Run("filter all settings", func(t *testing.T) {
		// Setup the tests
		tApi, e, teardownSuite := helper.SetupSuite(t)
		defer teardownSuite(t)

		// Create and send a request
		reqBody := domain.FilterSettingsByCriteriaInput{}
		rec, err := helper.SendRequest(e, tApi.SettingHandler.Filter, http.MethodPost, "/setting/filter", nil, nil, reqBody)
		if err != nil {
			t.Fatalf("Error sending request: %v", err)
		}

		// Check the status code
		codeWanted := http.StatusOK
		codeGot := rec.Code
		if codeWanted != codeGot {
			t.Fatalf("Wanted status code %v, got %v", codeWanted, codeGot)
		}

		// Parse & verify the response
		var resp domain.PaginationResponse
		helper.ParseResponse(t, rec, &resp)
		tWanted := int64(3)
		tGot := resp.Total
		if tWanted != tGot {
			t.Fatalf("Wanted %v settings, got %v", tWanted, tGot)
		}

		// Parse & verify the data
		var entityData []domain.Setting
		helper.ParseEntityData(t, resp.Data, &entityData)
		keyWanted := "app.name"
		keyFound := false
		for _, v := range entityData {
			if v.Key == keyWanted {
				keyFound = true
				break
			}
		}
		if !keyFound {
			t.Fatalf("Wanted setting key %v, couldn't find it", keyWanted)
		}
	})

	t.Run("settings not found", func(t *testing.T) {
		// Setup the tests
		tApi, e, teardownSuite := helper.SetupSuite(t)
		defer teardownSuite(t)

		// Create and send a request
		reqBody := domain.FilterSettingsByCriteriaInput{
			Keys: []string{"unknown.key"},
		}
		rec, err := helper.SendRequest(e, tApi.SettingHandler.Filter, http.MethodPost, "/setting/filter", nil, nil, reqBody)
		if err != nil {
			t.Fatalf("Error sending request: %v", err)
		}

		// Check the status code
		codeWanted := http.StatusOK
		codeGot := rec.Code
		if codeWanted != codeGot {
			t.Fatalf("Wanted status code %v, got %v", codeWanted, codeGot)
		}

		// Parse & verify the response
		var resp domain.PaginationResponse
		helper.ParseResponse(t, rec, &resp)
		tWanted := int64(0)
		tGot := resp.Total
		if tWanted != tGot {
			t.Fatalf("Wanted %v settings, got %v", tWanted, tGot)
		}

		// Parse & verify the data
		var entityData []domain.Setting
		helper.ParseEntityData(t, resp.Data, &entityData)
		lenWanted := 0
		lenGot := len(entityData)
		if lenWanted != lenGot {
			t.Fatalf("Wanted %v settings, got %v", lenWanted, lenGot)
		}
	})

	t.Run("should search settings by key", func(t *testing.T) {
		// Setup the tests
		tApi, e, teardownSuite := helper.SetupSuite(t)
		defer teardownSuite(t)

		// Create and send a request
		reqBody := domain.FilterSettingsByCriteriaInput{
			Keys: []string{"app.name"},
		}
		rec, err := helper.SendRequest(e, tApi.SettingHandler.Filter, http.MethodPost, "/setting/filter", nil, nil, reqBody)
		if err != nil {
			t.Fatalf("Error sending request: %v", err)
		}

		// Check the status code
		codeWanted := http.StatusOK
		codeGot := rec.Code
		if codeWanted != codeGot {
			t.Fatalf("Wanted status code %v, got %v", codeWanted, codeGot)
		}

		// Parse & verify the response
		var resp domain.PaginationResponse
		helper.ParseResponse(t, rec, &resp)
		tWanted := int64(1)
		tGot := resp.Total
		if tWanted != tGot {
			t.Fatalf("Wanted %v settings, got %v", tWanted, tGot)
		}

		// Parse & verify the data
		var entityData []domain.Setting
		helper.ParseEntityData(t, resp.Data, &entityData)
		keyWanted := "app.name"
		keyFound := false
		for _, v := range entityData {
			if v.Key == keyWanted {
				keyFound = true
				break
			}
		}
		if !keyFound {
			t.Fatalf("Wanted setting key %v, couldn't find it", keyWanted)
		}
	})
}

func TestFindSettingByID(t *testing.T) {
	t.Run("should return setting by ID", func(t *testing.T) {
		// Setup the tests
		tApi, e, teardownSuite := helper.SetupSuite(t)
		defer teardownSuite(t)

		// Find an existing setting
		reqBody := domain.FilterSettingsByCriteriaInput{}
		rec, err := helper.SendRequest(e, tApi.SettingHandler.Filter, http.MethodPost, "/setting/filter", nil, nil, reqBody)
		if err != nil {
			t.Fatalf("Error sending request: %v", err)
		}

		// Check the status code
		codeWanted := http.StatusOK
		codeGot := rec.Code
		if codeWanted != codeGot {
			t.Fatalf("Wanted status code %v, got %v", codeWanted, codeGot)
		}

		// Parse & verify the response
		var resp domain.PaginationResponse
		helper.ParseResponse(t, rec, &resp)

		// Parse & verify the data
		var entityData []domain.Setting
		helper.ParseEntityData(t, resp.Data, &entityData)
		settingID := entityData[0].ID
		if settingID == uuid.Nil {
			t.Fatalf("Wanted valid setting ID, got %v", settingID)
		}

		// Create and send a request
		pathParams := map[string]string{}
		pathParams["id"] = settingID.String()
		rec, err = helper.SendRequest(e, tApi.SettingHandler.FindByID, http.MethodGet, "/setting/"+settingID.String(), pathParams, nil, nil)
		if err != nil {
			t.Fatalf("Error sending request: %v", err)
		}

		// Check the status code
		codeWanted = http.StatusOK
		codeGot = rec.Code
		if codeWanted != codeGot {
			t.Fatalf("Wanted status code %v, got %v", codeWanted, codeGot)
		}

		// Parse & verify the response
		var bResp domain.BaseResponse
		helper.ParseResponse(t, rec, &bResp)

		// Parse & verify the data
		var cData domain.Setting
		helper.ParseEntityData(t, bResp.Data, &cData)
		if cData.ID != settingID {
			t.Fatalf("Wanted setting ID %v, got %v", settingID, cData.ID)
		}
	})

	t.Run("should return error for invalid setting ID", func(t *testing.T) {
		// Setup the tests
		tApi, e, teardownSuite := helper.SetupSuite(t)
		defer teardownSuite(t)

		id := uuid.Must(uuid.NewV4()).String()
		pathParams := map[string]string{}
		pathParams["id"] = id

		// Create and send a request
		_, err := helper.SendRequest(e, tApi.SettingHandler.FindByID, http.MethodGet, "/setting/"+id, pathParams, nil, nil)
		if err == nil {
			t.Fatalf("Expected error, but got nothing")
		}
	})

	t.Run("should return error for an invalid UUID", func(t *testing.T) {
		// Setup the tests
		tApi, e, teardownSuite := helper.SetupSuite(t)
		defer teardownSuite(t)

		id := strings.Repeat("a", 10)
		pathParams := map[string]string{}
		pathParams["id"] = id

		// Create and send a request
		_, err := helper.SendRequest(e, tApi.SettingHandler.FindByID, http.MethodGet, "/setting/"+id, pathParams, nil, nil)
		if err == nil {
			t.Fatalf("Expected error, but got nothing")
		}
	})
}
