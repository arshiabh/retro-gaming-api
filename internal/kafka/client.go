package kafka

type Client struct {
	Brokers []string
}

func NewClient(brokers []string) *Client {
	return &Client{
		Brokers: brokers,
	}
}
