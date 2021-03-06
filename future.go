package raft

type deferError struct {
	errCh chan error
}

func (d *deferError) init() {
	d.errCh = make(chan error, 1)
}

func (d *deferError) respond(err error) {
	d.errCh <- err
	close(d.errCh)
}

func (d *deferError) error() <-chan error {
	return d.errCh
}

type Future struct {
	deferError
	respCh chan interface{}
}

func newFuture() *Future {
	f := &Future{
		respCh: make(chan interface{}, 1),
	}
	f.init()
	return f
}

func (f *Future) respond(resp interface{}, err error) {
	if err == nil {
		f.respCh <- resp
		close(f.respCh)
	} else {
		f.deferError.respond(err)
	}
}

func (f *Future) Response() (resp interface{}, err error) {
	select {
	case err = <-f.error():
		return
	case resp = <-f.respCh:
		return
	}
}
func (f *Future) AsyncResponse() (respCh <-chan interface{}, errCh <-chan error) {
	return f.respCh, f.errCh
}
