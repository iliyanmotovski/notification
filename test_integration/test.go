package test_integration

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/iliyanmotovski/notification/api/notification-api/rest/router"
	"github.com/iliyanmotovski/notification/pkg/service/clock"
	"github.com/iliyanmotovski/notification/pkg/service/config"
	"github.com/iliyanmotovski/notification/pkg/service/orm"
)

var ormService orm.Engine
var registryService orm.RegistryService

func createTestEngine() {
	// make orm instance only once
	if ormService == nil {
		configService, err := config.NewConfigService("test", "../config/config_test.yaml")
		if err != nil {
			log.Fatal(err)
		}

		registryService, _, err = orm.NewORMRegistryService(configService)
		if err != nil {
			log.Fatal(err)
		}

		ormService = registryService.GetORMService()
		ormService.ExecuteAlters()
	}

	ormService.TruncateTables()
	ormService.GetCacheService().Clear()
	ormService.GetCacheService(orm.StreamsPool).Clear()
}

func sendHTTPRequest(clockService clock.IClock, method, pathAndQuery string, body interface{}, response interface{}) {
	ginEngine := gin.New()
	router.Init(registryService, clockService)(ginEngine)

	ts := httptest.NewServer(ginEngine)
	defer ts.Close()

	marshaled, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	r, err := http.NewRequest(method, ts.URL+pathAndQuery, bytes.NewBuffer(marshaled))
	if err != nil {
		log.Fatal(err)
	}

	w, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Fatal(err)
	}

	defer w.Body.Close()

	respBody, err := io.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}

	if w.StatusCode != http.StatusOK {
		log.Fatalf("got http status %d : %s", w.StatusCode, string(respBody))
	}

	if response != nil {
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			log.Fatal(err)
		}
	}
}
