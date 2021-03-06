package cyclon

import "github.com/Gaboose/go-pubsub/pnet"

// overflowBuffer forms a last-in last-out queue between the given channels.
// Input channel never blocks from outside. If the buffer is full,
// it'll discard the last-in value.
func overflowBuffer(n int, in <-chan pnet.Peer, out chan<- pnet.Peer) {
	var outMaybe chan<- pnet.Peer
	//buf, i and j together form a circular buffer
	i, j := 0, 0
	buf := make([]pnet.Peer, n)
	for {
		select {
		case v, ok := <-in:
			if !ok {
				close(out)
				return
			}
			if i == j {
				if outMaybe == nil {
					// buffer is empty
					outMaybe = out
				} else {
					// buffer is full
					// overwrite last value and nudge the output index
					j = (j + 1) % n
				}
			}
			buf[i] = v
			i = (i + 1) % n

		case outMaybe <- buf[j]:
			j = (j + 1) % n
			if i == j {
				outMaybe = nil
			}
		}
	}
}
