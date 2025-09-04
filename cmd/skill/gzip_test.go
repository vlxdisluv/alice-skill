package main

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/require"
	"github.com/vlxdisluv/alice-skill/internal/store"
	"github.com/vlxdisluv/alice-skill/internal/store/mock"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGzipCompression(t *testing.T) {
	// создадим конроллер моков и экземпляр мок-хранилища
	ctrl := gomock.NewController(t)
	s := mock.NewMockStore(ctrl)

	// определим, какой результат будем получать от «хранилища»
	messages := []store.Message{
		{
			Sender:  "411419e5-f5be-4cdb-83aa-2ca2b6648353",
			Time:    time.Now(),
			Payload: "Hello!",
		},
	}

	// установим условие: при любом вызове метода ListMessages возвращать массив messages без ошибки
	s.EXPECT().
		ListMessages(gomock.Any(), gomock.Any()).
		Return(messages, nil)

	// создадим экземпляр приложения и передадим ему «хранилище»
	appInstance := newApp(s)

	handler := http.HandlerFunc(gzipMiddleware(appInstance.webhook))

	srv := httptest.NewServer(handler)
	defer srv.Close()

	requestBody := `{
        "request": {
            "type": "SimpleUtterance",
            "command": "sudo do something"
        },
        "version": "1.0"
    }`

	// ожидаемое содержимое тела ответа при успешном запросе
	successBody := `{
        "response": {
            "text": "Для вас нет новых сообщений."
        },
        "version": "1.0"
    }`

	t.Run("sends_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		r := httptest.NewRequest("POST", srv.URL, buf)
		r.RequestURI = ""
		r.Header.Set("Content-Encoding", "gzip")
		r.Header.Set("Accept-Encoding", "")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.JSONEq(t, successBody, string(b))
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		buf := bytes.NewBufferString(requestBody)
		r := httptest.NewRequest("POST", srv.URL, buf)
		r.RequestURI = ""
		r.Header.Set("Accept-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()

		zr, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)

		b, err := io.ReadAll(zr)
		require.NoError(t, err)

		require.JSONEq(t, successBody, string(b))
	})
}
