package tutorme

const (
	JavascriptTopic string = "javascript_topic"
	PythonTopic     string = "python_topic"
	GoLangTopic     string = "go_lang_topic"
)

type Publisher interface {
	CreateTopic(topicName string) error
	Subscribe(topicName string, abort chan bool) (chan []byte, error)
	Publish(topicName string, data []byte) error
}
