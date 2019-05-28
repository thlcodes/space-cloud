package functions

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/client"
)

func (m *Module) initWorkers(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go m.worker()
	}
}

func (m *Module) worker() {
	for msg := range m.channel {
		req := new(model.FunctionsPayload)
		err := json.Unmarshal(msg.Data, req)
		if err != nil {
			log.Println("Functions Error -", err)
			m.publishErrorResponse(msg.Reply, err)
			continue
		}

		t, p := m.services.Load(req.Service)
		if !p {
			err := errors.New("No service available")
			log.Println("Functions Error -", err)
			m.publishErrorResponse(msg.Reply, err)
		}

		service := t.(*servicesStub)
		c := service.clients[rand.Intn(len(service.clients))]
		m.requestService(c, req, msg.Reply)
	}
}

func (m *Module) removeStaleRequests() {
	ticker := time.NewTicker(2 * time.Minute)

	for range ticker.C {
		m.pendingRequests.Range(func(key interface{}, value interface{}) bool {
			req := value.(*pendingRequest)

			// Remove the request if its more than 30 seconds old
			diff := time.Now().Sub(req.reqTime)
			if diff.Seconds() > 30 {
				m.pendingRequests.Delete(key)
			}

			return true
		})
	}
}

func (m *Module) requestService(client client.Client, req *model.FunctionsPayload, reply string) {
	// Generate a unique id for request
	id := uuid.NewV1().String()

	// Add request to the map of pending requests
	m.pendingRequests.Store(id, &pendingRequest{reply: reply, reqTime: time.Now()})

	// Send the request to the service
	client.Write(&model.Message{
		Type: utils.TypeServiceRequest,
		Data: req,
		ID:   id,
	})
}

func (m *Module) publishErrorResponse(subject string, err error) {
	res := &model.FunctionsPayload{Params: map[string]interface{}{"error": err}}
	data, _ := json.Marshal(res)
	m.nc.Publish(subject, data)
}
