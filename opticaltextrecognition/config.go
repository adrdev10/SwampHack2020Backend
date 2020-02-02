package opticaltextrecognition

//ConfigurationFile represents our config file
type ConfigurationFile struct {
	ResultTopic     string
	ResultBucket    string
	TranlateTopic   string
	SentimentalData string
	Languages       []string
}

//NewConfig creates a new configurable file
func NewConfig(rt, rb, tp, sd string, lang []string) *ConfigurationFile {
	return &ConfigurationFile{
		ResultTopic:     rt,
		ResultBucket:    rb,
		TranlateTopic:   tp,
		SentimentalData: sd,
		Languages:       lang,
	}
}
