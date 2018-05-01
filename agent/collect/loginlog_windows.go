// +build windows
// 迭代和设计过程：https://mp.weixin.qq.com/s/rHDJ2tQWEaZLikMt5bgCsw

package collect

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"syscall"
	"time"
	"yulong-hids/agent/common"

	"golang.org/x/sys/windows"

	"github.com/elastic/beats/winlogbeat/sys"
	win "github.com/elastic/beats/winlogbeat/sys/wineventlog"
)

const (
	// renderBufferSize is the size in bytes of the buffer used to render events.
	renderBufferSize   = 1 << 14
	winodwsEvtxFile    = "C:\\Windows\\System32\\winevt\\Logs\\Security.evtx"
	winodwsEvtxFilex32 = "C:\\Windows\\Sysnative\\winevt\\Logs\\Security.evtx"
)

var localAddress = []string{"-", "127.0.0.1", "::1"}

type Record struct {
	sys.Event
	API string // The event log API type used to read the record.
	XML string // XML representation of the event.
}

// EventLog is an interface to a Windows Event Log.
type EventLog interface {
	// Open the event log. recordNumber is the last successfully read event log
	// record number. Read will resume from recordNumber + 1. To start reading
	// from the first event specify a recordNumber of 0.
	Open(recordNumber uint64) error

	// Read records from the event log.
	Read() ([]Record, error)

	// Close the event log. It should not be re-opened after closing.
	Close() error
}

// query contains parameters used to customize the event log data that is
// queried from the log.
type query struct {
	IgnoreOlder time.Duration // Ignore records older than this period of time.
	EventID     string        // White-list and black-list of events.
	Level       string        // Severity level.
	Provider    []string      // Provider (source name).
}

// Validate that winEventLog implements the EventLog interface.
var _ EventLog = &winEventLog{}
var first = true

// winEventLog implements the EventLog interface for reading from the Windows
// Event Log API.
type winEventLog struct {
	// config       winEventLogConfig
	query        string
	channelName  string        // Name of the channel from which to read.
	subscription win.EvtHandle // Handle to the subscription.
	maxRead      int           // Maximum number returned in one Read.
	lastRead     uint64        // Record number of the last read event.

	render    func(event win.EvtHandle, out io.Writer) error // Function for rendering the event to XML.
	renderBuf []byte                                         // Buffer used for rendering event.
	outputBuf *sys.ByteBuffer                                // Buffer for receiving XML
	// cache     *messageFilesCache                             // Cached mapping of source name to event message file handles.

	logPrefix string // String to prefix on log messages.
}

func (l *winEventLog) Open(recordNumber uint64) error {
	bookmark, err := win.CreateBookmark(l.channelName, recordNumber)
	if err != nil {
		return err
	}
	defer win.Close(bookmark)

	// Using a pull subscription to receive events. See:
	// https://msdn.microsoft.com/en-us/library/windows/desktop/aa385771(v=vs.85).aspx#pull
	signalEvent, err := windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		return nil
	}

	subscriptionHandle, err := win.Subscribe(
		0, // Session - nil for localhost
		signalEvent,
		"",       // Channel - empty b/c channel is in the query
		l.query,  // Query - nil means all events
		bookmark, // Bookmark - for resuming from a specific event
		win.EvtSubscribeStartAfterBookmark)
	if err != nil {
		return err
	}

	l.subscription = subscriptionHandle
	return nil
}

func (l *winEventLog) Read() ([]Record, error) {
	handles, _, err := l.eventHandles(l.maxRead)
	if err != nil || len(handles) == 0 {
		return nil, err
	}
	defer func() {
		for _, h := range handles {
			win.Close(h)
		}
	}()

	var records []Record
	for _, h := range handles {
		l.outputBuf.Reset()
		err := l.render(h, l.outputBuf)
		if bufErr, ok := err.(sys.InsufficientBufferError); ok {
			l.renderBuf = make([]byte, bufErr.RequiredSize)
			l.outputBuf.Reset()
			err = l.render(h, l.outputBuf)
		}
		if err != nil && l.outputBuf.Len() == 0 {
			continue
		}

		r, err := l.buildRecordFromXML(l.outputBuf.Bytes(), err)
		if err != nil {
			continue
		}
		records = append(records, r)
		l.lastRead = r.RecordID
	}
	return records, nil
}

func (l *winEventLog) Close() error {
	return win.Close(l.subscription)
}

func (l *winEventLog) eventHandles(maxRead int) ([]win.EvtHandle, int, error) {
	handles, err := win.EventHandles(l.subscription, maxRead)
	switch err {
	case nil:
		return handles, maxRead, nil
	case win.ERROR_NO_MORE_ITEMS:
		return nil, maxRead, nil
	case win.RPC_S_INVALID_BOUND:
		if err := l.Close(); err != nil {
			return nil, 0, err
		}
		if err := l.Open(l.lastRead); err != nil {
			return nil, 0, err
		}
		return l.eventHandles(maxRead / 2)
	default:
		return nil, 0, err
	}
}

func (l *winEventLog) buildRecordFromXML(x []byte, recoveredErr error) (Record, error) {
	e, err := sys.UnmarshalEventXML(x)
	if err != nil {
		return Record{}, fmt.Errorf("Failed to unmarshal XML='%s'. %v", x, err)
	}

	sys.PopulateAccount(&e.User)

	if e.RenderErrorCode != 0 {
		// Convert the render error code to an error message that can be
		// included in the "message_error" field.
		e.RenderErr = syscall.Errno(e.RenderErrorCode).Error()
	} else if recoveredErr != nil {
		e.RenderErr = recoveredErr.Error()
	}

	if e.Level == "" {
		// Fallback on LevelRaw if the Level is not set in the RenderingInfo.
		e.Level = win.EventLevel(e.LevelRaw).String()
	}

	r := Record{
		Event: e,
	}
	return r, nil
}

func newWinEventLog(eventID string) (EventLog, error) {
	var ignoreOlder time.Duration
	if first {
		ignoreOlder = time.Hour * 17520
		first = false
	} else {
		ignoreOlder = time.Second * 60
	}
	query, err := win.Query{
		Log:         "Security",
		IgnoreOlder: ignoreOlder,
		Level:       "",
		EventID:     eventID,
		Provider:    []string{},
	}.Build()
	if err != nil {
		return nil, err
	}

	l := &winEventLog{
		query:       query,
		channelName: "Security",
		maxRead:     1000,
		renderBuf:   make([]byte, renderBufferSize),
		outputBuf:   sys.NewByteBuffer(renderBufferSize),
	}

	l.render = func(event win.EvtHandle, out io.Writer) error {
		return win.RenderEvent(event, 0, l.renderBuf, nil, out)
	}
	return l, nil
}

// GetLoginLog 获取系统登录日志
func GetLoginLog() (resultData []map[string]string) {
	var loginFile string
	var timestamp int64
	if common.Config.Lasttime == "all" {
		timestamp = 615147123
	} else {
		ti, _ := time.Parse("2006-01-02T15:04:05Z07:00", common.Config.Lasttime)
		timestamp = ti.Unix()
	}
	if runtime.GOARCH == "386" {
		loginFile = winodwsEvtxFilex32
	} else {
		loginFile = winodwsEvtxFile
	}
	if _, err := os.Stat(loginFile); err != nil {
		// 不支持2003
		log.Println(err.Error())
		return
	}
	resultData = getSuccessLog(timestamp)
	resultData = append(resultData, getFailedLog(timestamp)...)
	return
}

func getSuccessLog(timestamp int64) (resultData []map[string]string) {
	l, err := newWinEventLog("4624")
	if err != nil {
		return
	}
	err = l.Open(0)
	if err != nil {
		return
	}
	reList, _ := l.Read()
	for _, rec := range reList {
		// rec.EventData.Pairs[10].Value != "5" &&
		if rec.TimeCreated.SystemTime.Local().Unix() > timestamp {
			if common.InArray(localAddress, rec.EventData.Pairs[18].Value, false) {
				continue
			}
			m := make(map[string]string)
			m["status"] = "true"
			m["username"] = rec.EventData.Pairs[5].Value
			m["remote"] = rec.EventData.Pairs[18].Value
			m["time"] = rec.TimeCreated.SystemTime.Local().Format("2006-01-02T15:04:05Z07:00")
			resultData = append(resultData, m)
		}
	}
	return
}
func getFailedLog(timestamp int64) (resultData []map[string]string) {
	l, err := newWinEventLog("4625")
	if err != nil {
		return
	}
	err = l.Open(0)
	if err != nil {
		return
	}
	reList, _ := l.Read()
	for _, rec := range reList {
		// rec.EventData.Pairs[8].Value != "5" &&
		if rec.TimeCreated.SystemTime.Local().Unix() > timestamp {
			if common.InArray(localAddress, rec.EventData.Pairs[18].Value, false) {
				continue
			}
			m := make(map[string]string)
			m["status"] = "false"
			m["username"] = rec.EventData.Pairs[5].Value
			m["remote"] = rec.EventData.Pairs[18].Value
			m["time"] = rec.TimeCreated.SystemTime.Local().Format("2006-01-02T15:04:05Z07:00")
			resultData = append(resultData, m)
		}
	}
	return
}
