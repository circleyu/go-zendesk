package mock

import (
	"go-zendesk/zendesk"
)

var _ zendesk.API = (*Client)(nil)
