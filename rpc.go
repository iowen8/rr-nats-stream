package natsstream

import (
	"encoding/json"
    "fmt"

)

type rpc struct {
	srv *Plugin
	log logger.Logger
}




func (s *rpc) PublishEvent(subject string, data map[string]interface{}, output *SuccessResponse) error {
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
   
   // eventService returns associated rpc service.
   func (p *Plugin) RPC() interface{} {
	return &rpc{srv: p, log: p.log}
   }