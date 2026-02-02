package transport_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deadelus/go-clean-app/v2/transport"
	"github.com/stretchr/testify/assert"
)

func TestSendJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	payload := map[string]string{"foo": "bar"}
	transport.SendJSON(rec, http.StatusCreated, payload)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	var got map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, payload, got)
}

func TestSendError(t *testing.T) {
	rec := httptest.NewRecorder()
	msg := "something went wrong"
	transport.SendError(rec, http.StatusBadRequest, msg)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var resp transport.Response
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, msg, resp.Error)
}

func TestSendSuccess(t *testing.T) {
	rec := httptest.NewRecorder()
	data := map[string]string{"foo": "bar"}
	transport.SendSuccess(rec, http.StatusOK, data)
	assert.Equal(t, http.StatusOK, rec.Code)
	var resp transport.Response
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, resp.Data)
}
