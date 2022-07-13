package natsstream

import (
	"encoding/json"
    "fmt"

)

type event struct {
	srv *Plugin
	log logger.Logger
}




func (s *event) PublishEvent(subject string, data map[string]interface{}, output *SuccessResponse) error {
	msg := json.Marshall(data)
	resp := srv.eventPublisher.PublishEvent(subject, msg)
	*output = resp
	return nil
   }

   // CollectTarget resettable service.
func (p *Plugin) CollectTarget(name endure.Named, r Informer) error {
	p.registry[name.Name()] = r
	return nil
   }
   
   // Collects declares services to be collected.
   func (p *Plugin) Collects() []interface{} {
	return []interface{}{
	 p.CollectTarget,
	}
   }
   
   // Name of the service.
   func (p *Plugin) Name() string {
	return PluginName
   }
   
   // eventService returns associated event service.
   func (p *Plugin) event() interface{} {
	return &event{srv: p, log: p.log}
   }