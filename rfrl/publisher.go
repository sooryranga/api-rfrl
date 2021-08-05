package rfrl

import "os"

var (
	JavascriptTopic string = os.Getenv("JAVASCRIPT_TOPIC")
	PythonTopic     string = os.Getenv("PYTHON_TOPIC")
	GoLangTopic     string = os.Getenv("GO_LANG_TOPIC")
)

type Publisher interface {
	CreateTopic(topicName string) error
	Subscribe(topicName string, abort chan bool) (chan []byte, error)
	Publish(topicName string, data []byte) error
}
