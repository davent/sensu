package sensu

import "testing"

func TestLoadConfigDir(t *testing.T) {
	config, err := LoadConfigDir("example/")
	if err != nil {
		t.Errorf("LoadConfigDir Fail: %s", err)
	}

	if config.RabbitMQ.Host != "10.0.0.1" {
		t.Errorf("LoadConfigDir conf.d/rabbitmq.json not loaded correctly. config.RabbitMQ.Host should be '10.0.0.1' but is actually '%s'", config.RabbitMQ.Host)
	}
}
