package mock

import (
	"github.com/circleyu/go-zendesk/zendesk"
)

var _ zendesk.API = (*Client)(nil)
