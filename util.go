package dispatcher

import "fmt"

type TopicNameBuilder func(e EventHeader) string

func DefaultTopicNameBuilder(e EventHeader) string {
	return fmt.Sprintf("%s_%s_%s_%s", e.King, e.Noble, e.Knight, e.Peasant)
}
