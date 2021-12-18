package client

import "fmt"

func (c *Client) handlePayload(payload []byte) {
	fmt.Println(payload)
}
