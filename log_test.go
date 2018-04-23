package log

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"testing"

	"github.com/frrist/testify/assert"
	tracer "github.com/ipfs/go-log/tracer"
	writer "github.com/ipfs/go-log/writer"
)

func TestSingleEvent(t *testing.T) {
	assert := assert.New(t)

	// Set up a pipe to use as backend and log stream
	lgs, lgb := io.Pipe()
	// event logs will be written to lgb
	// event logs will be read from lgs
	writer.WriterGroup.AddWriter(lgb)

	// create a logger
	lgr := Logger("test")

	// create a root context
	ctx := context.Background()

	// start an event
	ctx = lgr.Start(ctx, "event1")

	// finish the event
	lgr.Finish(ctx)

	// decode the log event
	var ls tracer.LoggableSpan
	evtDecoder := json.NewDecoder(lgs)
	evtDecoder.Decode(&ls)

	// event name and system should be
	assert.Equal("event1", ls.Operation)
	assert.Equal("test", ls.Tags["system"])
	// greater than zero should work for now
	assert.NotZero(ls.Duration)
	assert.NotZero(ls.Start)
}

func TestSingleEventWithErr(t *testing.T) {
	assert := assert.New(t)

	// Set up a pipe to use as backend and log stream
	lgs, lgb := io.Pipe()
	// event logs will be written to lgb
	// event logs will be read from lgs
	writer.WriterGroup.AddWriter(lgb)

	// create a logger
	lgr := Logger("test")

	// create a root context
	ctx := context.Background()

	// start an event
	ctx = lgr.Start(ctx, "event1")

	// finish the event
	lgr.FinishWithErr(ctx, errors.New("rawer im an error"))

	// decode the log event
	var ls tracer.LoggableSpan
	evtDecoder := json.NewDecoder(lgs)
	evtDecoder.Decode(&ls)

	// event name and system should be
	assert.Equal("event1", ls.Operation)
	assert.Equal("test", ls.Tags["system"])
	assert.Equal(true, ls.Tags["error"])
	assert.Contains(ls.Logs[0].Field[0].Value, "rawer im an error")
	// greater than zero should work for now
	assert.NotZero(ls.Duration)
	assert.NotZero(ls.Start)
}

func TestEventWithTag(t *testing.T) {
	assert := assert.New(t)

	// Set up a pipe to use as backend and log stream
	lgs, lgb := io.Pipe()
	// event logs will be written to lgb
	// event logs will be read from lgs
	writer.WriterGroup.AddWriter(lgb)

	// create a logger
	lgr := Logger("test")

	// create a root context
	ctx := context.Background()

	// start an event
	ctx = lgr.Start(ctx, "event1")
	lgr.SetTag(ctx, "tk", "tv")

	// finish the event
	lgr.Finish(ctx)

	// decode the log event
	var ls tracer.LoggableSpan
	evtDecoder := json.NewDecoder(lgs)
	evtDecoder.Decode(&ls)

	// event name and system should be
	assert.Equal("event1", ls.Operation)
	assert.Equal("test", ls.Tags["system"])
	assert.Equal("tv", ls.Tags["tk"])
	// greater than zero should work for now
	assert.NotZero(ls.Duration)
	assert.NotZero(ls.Start)
}

func TestEventWithTags(t *testing.T) {
	assert := assert.New(t)

	// Set up a pipe to use as backend and log stream
	lgs, lgb := io.Pipe()
	// event logs will be written to lgb
	// event logs will be read from lgs
	writer.WriterGroup.AddWriter(lgb)

	// create a logger
	lgr := Logger("test")

	// create a root context
	ctx := context.Background()

	// start an event
	ctx = lgr.Start(ctx, "event1")
	lgr.SetTags(ctx, map[string]interface{}{
		"tk1": "tv1",
		"tk2": "tv2",
	})

	// finish the event
	lgr.Finish(ctx)

	// decode the log event
	var ls tracer.LoggableSpan
	evtDecoder := json.NewDecoder(lgs)
	evtDecoder.Decode(&ls)

	// event name and system should be
	assert.Equal("event1", ls.Operation)
	assert.Equal("test", ls.Tags["system"])
	assert.Equal("tv1", ls.Tags["tk1"])
	assert.Equal("tv2", ls.Tags["tk2"])
	// greater than zero should work for now
	assert.NotZero(ls.Duration)
	assert.NotZero(ls.Start)
}

func TestEventWithLogs(t *testing.T) {
	assert := assert.New(t)

	// Set up a pipe to use as backend and log stream
	lgs, lgb := io.Pipe()
	// event logs will be written to lgb
	// event logs will be read from lgs
	writer.WriterGroup.AddWriter(lgb)

	// create a logger
	lgr := Logger("test")

	// create a root context
	ctx := context.Background()

	// start an event
	ctx = lgr.Start(ctx, "event1")
	lgr.LogKV(ctx, "log1", "logv1", "log2", "logv2")
	lgr.LogKV(ctx, "treeLog", []string{"Pine", "Juniper", "Spruce", "Ginkgo"})

	// finish the event
	lgr.Finish(ctx)

	// decode the log event
	var ls tracer.LoggableSpan
	evtDecoder := json.NewDecoder(lgs)
	evtDecoder.Decode(&ls)

	// event name and system should be
	assert.Equal("event1", ls.Operation)
	assert.Equal("test", ls.Tags["system"])

	assert.Equal("log1", ls.Logs[0].Field[0].Key)
	assert.Equal("logv1", ls.Logs[0].Field[0].Value)
	assert.Equal("log2", ls.Logs[0].Field[1].Key)
	assert.Equal("logv2", ls.Logs[0].Field[1].Value)

	// Should be a differnt log (different timestamp)
	assert.Equal("treeLog", ls.Logs[1].Field[0].Key)
	assert.Equal("[Pine Juniper Spruce Ginkgo]", ls.Logs[1].Field[0].Value)

	// greater than zero should work for now
	assert.NotZero(ls.Duration)
	assert.NotZero(ls.Start)
}

func TestMultiEvent(t *testing.T) {
	assert := assert.New(t)

	// Set up a pipe to use as backend and log stream
	lgs, lgb := io.Pipe()
	// event logs will be written to lgb
	// event logs will be read from lgs
	writer.WriterGroup.AddWriter(lgb)
	evtDecoder := json.NewDecoder(lgs)

	// create a logger
	lgr := Logger("test")

	// create a root context
	ctx := context.Background()

	ctx = lgr.Start(ctx, "root")

	doEvent(ctx, "e1", lgr)
	doEvent(ctx, "e2", lgr)

	lgr.Finish(ctx)

	e1 := getEvent(evtDecoder)
	assert.Equal("e1", e1.Operation)
	assert.Equal("test", e1.Tags["system"])
	assert.NotZero(e1.Duration)
	assert.NotZero(e1.Start)

	// I hope your clocks work...
	e2 := getEvent(evtDecoder)
	assert.Equal("e2", e2.Operation)
	assert.Equal("test", e2.Tags["system"])
	assert.NotZero(e2.Duration)
	assert.True(e1.Start.Nanosecond() < e2.Start.Nanosecond())

	er := getEvent(evtDecoder)
	assert.Equal("root", er.Operation)
	assert.Equal("test", er.Tags["system"])
	assert.True(er.Duration.Nanoseconds() > e1.Duration.Nanoseconds()+e2.Duration.Nanoseconds())
	assert.NotZero(er.Start)

}

func TestEventSerialization(t *testing.T) {
	assert := assert.New(t)

	// Set up a pipe to use as backend and log stream
	lgs, lgb := io.Pipe()
	// event logs will be written to lgb
	// event logs will be read from lgs
	writer.WriterGroup.AddWriter(lgb)
	evtDecoder := json.NewDecoder(lgs)

	// create a logger
	lgr := Logger("test")

	// start an event
	sndctx := lgr.Start(context.Background(), "send")

	// **imagine** that we are putting `bc` (byte context) into a protobuf message
	// and send the message to another peer on the network
	bc, err := lgr.SerializeContext(sndctx)
	assert.NoError(err)

	// now  **imagine** some peer getting a protobuf message and extracting
	// `bc` from the message to continue the operation
	rcvctx, err := lgr.StartFromParentState(context.Background(), "recv", bc)
	assert.NoError(err)

	// at some point the sender completes their operation
	lgr.Finish(sndctx)
	e := getEvent(evtDecoder)
	assert.Equal("send", e.Operation)
	assert.Equal("test", e.Tags["system"])
	assert.NotZero(e.Start)
	assert.NotZero(e.Start)

	// and then the receiver finishes theirs
	lgr.Finish(rcvctx)
	e = getEvent(evtDecoder)
	assert.Equal("recv", e.Operation)
	assert.Equal("test", e.Tags["system"])
	assert.NotZero(e.Start)
	assert.NotZero(e.Start)

}

func doEvent(ctx context.Context, name string, el EventLogger) context.Context {
	ctx = el.Start(ctx, name)
	defer func() {
		el.Finish(ctx)
	}()
	return ctx
}

func getEvent(ed *json.Decoder) tracer.LoggableSpan {
	// decode the log event
	var ls tracer.LoggableSpan
	ed.Decode(&ls)
	return ls
}

// DEPRECATED methods tested below
func TestEventBegin(t *testing.T) {
	assert := assert.New(t)

	// Set up a pipe to use as backend and log stream
	lgs, lgb := io.Pipe()
	// event logs will be written to lgb
	// event logs will be read from lgs
	writer.WriterGroup.AddWriter(lgb)
	evtDecoder := json.NewDecoder(lgs)

	// create a logger
	lgr := Logger("test")

	// create a root context
	ctx := context.Background()

	// start an event in progress with metadata
	eip := lgr.EventBegin(ctx, "event", LoggableMap{"key": "val"})

	// append more metadata
	eip.Append(LoggableMap{"foo": "bar"})

	// set an error
	eip.SetError(errors.New("gerrr im an error"))

	// finish the event
	eip.Done()

	// decode the log event
	var ls tracer.LoggableSpan
	evtDecoder.Decode(&ls)

	assert.Equal("event", ls.Operation)
	assert.Equal("test", ls.Tags["system"])
	assert.Contains(ls.Logs[0].Field[0].Value, "val")
	assert.Contains(ls.Logs[1].Field[0].Value, "bar")
	assert.Contains(ls.Logs[2].Field[0].Value, "gerrr im an error")
	// greater than zero should work for now
	assert.NotZero(ls.Duration)
	assert.NotZero(ls.Start)
}

func TestEventBeginWithErr(t *testing.T) {
	assert := assert.New(t)

	// Set up a pipe to use as backend and log stream
	lgs, lgb := io.Pipe()
	// event logs will be written to lgb
	// event logs will be read from lgs
	writer.WriterGroup.AddWriter(lgb)
	evtDecoder := json.NewDecoder(lgs)

	// create a logger
	lgr := Logger("test")

	// create a root context
	ctx := context.Background()

	// start an event in progress with metadata
	eip := lgr.EventBegin(ctx, "event", LoggableMap{"key": "val"})

	// append more metadata
	eip.Append(LoggableMap{"foo": "bar"})

	// finish the event with an error
	eip.DoneWithErr(errors.New("gerrr im an error"))

	// decode the log event
	var ls tracer.LoggableSpan
	evtDecoder.Decode(&ls)

	assert.Equal("event", ls.Operation)
	assert.Equal("test", ls.Tags["system"])
	assert.Contains(ls.Logs[0].Field[0].Value, "val")
	assert.Contains(ls.Logs[1].Field[0].Value, "bar")
	assert.Contains(ls.Logs[2].Field[0].Value, "gerrr im an error")
	// greater than zero should work for now
	assert.NotZero(ls.Duration)
	assert.NotZero(ls.Start)
}
