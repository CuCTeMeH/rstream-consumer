package config

type ConsumerConfig struct {
	GroupSize                   int64  `mapstructure:"group_size"`
	ConsumerIDsListName         string `mapstructure:"consumer_ids_list_name"`
	PublishedMessagesStreamName string `mapstructure:"published_messages_stream_name"`
	ConsumedMessagesSetName     string `mapstructure:"consumed_messages_set_name"`
}
