package natsstream

import (
     "github.com/roadrunner-server/config/v2"
     "github.com/roadrunner-server/http/v2"
     "github.com/roadrunner-server/rpc/v2"

     "github.com/roadrunner-server/errors"
	 "github.com/nats-io/nats.go"
)

const PluginName = "nats-stream"
var once sync.Once


type CoreStreamInfo struct {
	Name     string
	Info     *nats.StreamInfo
	Subjects []string
}

type CoreStreamManager struct {
	ServiceBus
	StreamConnection nats.JetStreamContext
	StreamsInfo      []CoreStreamInfo
	SystemStreams    []map[string]interface{}
}

type CoreEventPublisher struct {
	Config        map[string]interface{}
	StreamManager CoreStreamManager
}


type CoreConfiguration struct {
	EventPublisher  map[string]interface{}
}

type Config struct{
     Address string `mapstructure:"address"`
}

type Plugin struct {
     cfg *Config
	 eventPublisher *CoreEventPublisher
}

type SuccessResponse struct {
	Success bool  `json:"success"`
	Error   error `json:"error"`
}


var coreEventPublisher *CoreEventPublisher
// You can also initialize some defaults values for config keys
func (cfg *Config) InitDefaults() {
     if cfg.Address == "" {
      cfg.Address = "tcp://localhost:8088"
    }
}

func (s *Plugin) Init(r *rpc.Plugin, h *http.Plugin, cfg config.Configurer) error {
 const op = errors.Op("nats_stream_plugin_init") // error operation name
 if !cfg.Has(PluginName) {
  return errors.E(op, errors.Disabled)
 }

 // unmarshall 
 err := cfg.UnmarshalKey(PluginName, &s.cfg)
 if err != nil {
  // Error will stop execution
  return errors.E(op, err)
 }

 // Init defaults
 s.cfg.InitDefaults()

 config := CoreConfiguration{}
 eventPubStreams := []map[string]interface{}{}
 eventPubStreams = append(eventPubStreams, map[string]interface{}{
	 "name": "ml-events", "subjects": []string{"Duel:Created", "Duel:Joined"}})

 config.EventPublisher = map[string]interface{}{"natsUser": "client@localhost", "name": "rpc:service:data:" + uuid.New().String(), "natsServers": []string{"ws://67.227.250.186:4555", "ws://67.227.250.186:4556", "ws://67.227.250.186:4557"},
	 "streams": eventPubStreams}
 coreEventPublisher = NewEventPublisher(config)
 var wg sync.WaitGroup
 wg.Add(1)
 coreEventPublisher.Connect(&wg)
 wg.Wait()
 coreEventPublisher.RegisterStreams()
 s.eventPublisher = coreEventPublisher
 
 return nil
}

func (s *Plugin) Serve() chan error {
	const op = errors.Op("nats_stream_plugin_serve")
	   errCh := make(chan error, 1)
   
	   err := s.DoSomeWork()
	   err != nil {
		errCh <- errors.E(op, err)
		return errCh
	   }
   
	   return nil
   }
   
   func (s *Plugin) Stop() error {
	   return s.stopServing()
   }
   
   func (s *Plugin) DoSomeWork() error {
	return nil
   }

   
func NewEventPublisher(config map[string]interface{}) *CoreEventPublisher {
	once.Do(func() {
		publisher := new(CoreEventPublisher)
		publisher.Config = config
		publisher.StreamManager = estream.NewStreamPublisherManager(config)
		coreEventPublisher = publisher
	})
	return coreEventPublisher
}

func (publisher *CoreEventPublisher) Connect(wg *sync.WaitGroup) SuccessResponse {
	defer wg.Done()
	return publisher.StreamManager.Connect(wg, false)
}

func (publisher *CoreEventSubscriber) RegisterStreams() (bool, error) {
	streams, valid := publisher.Config["streams"].([]map[string]interface{})
	errs := []string{}
	fmt.Printf("valid config for publisher.Config %v", publisher.Config)
	fmt.Printf("\nvalid config for subscription %v", streams)
	fmt.Printf("\nvalid config for valid %v", valid)

	if !valid {
		errs = append(errs, fmt.Sprintf("%v : invalid configuration: subscriptions must be defined.", "configuration"))
	} else {

		lo.ForEach(streams, func(stream map[string]interface{}, _ int) {
			subjects, valid := stream["subjects"].([]string)
			fmt.Printf("valid config for streams['subjects] %v", stream["subjects"])
			fmt.Printf("\nvalid config for subjects %v", subjects)
			fmt.Printf("\nvalid config for valid %v", valid)
			name := fmt.Sprintf("%v", stream["name"])

			subjects = stream["subjects"].([]string)
			added, err := publisher.StreamManager.CreateStream(name, &nats.StreamConfig{Name: name, Subjects: subjects})

			fmt.Printf("added response: %v\n", added)
			fmt.Printf("added err: %v\n", err)

		})
	}

	if len(errs) > 0 {
		errsString := strings.Join(errs, "\n")
		return false, fmt.Errorf("the following error occurred: \n%v", errsString)
	}
	return true, nil
}

func (publisher *CoreEventPublisher) PublishEvent(subject string, data []byte) SuccessResponse {
	pres, err := publisher.StreamManager.StreamConnection.Publish(subject, data)
	fmt.Printf("publish response: %v\n", pres)
	fmt.Printf("publish err: %v\n", err)
	return SuccessResponse{Success: err != nil, Error: err}
}