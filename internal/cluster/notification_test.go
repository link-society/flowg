package cluster

import (
	"bytes"
	"context"
	"io"
	"reflect"
	"testing"

	"link-society.com/flowg/internal/storage"
	"link-society.com/flowg/internal/storage/changefeed"
)

type fakeStreamable struct {
	applied [][]changefeed.Record
	err     error
}

var _ storage.Streamable = (*fakeStreamable)(nil)

func (f *fakeStreamable) Dump(context.Context, io.Writer, uint64) (uint64, error) { return 0, nil }
func (f *fakeStreamable) Load(context.Context, io.Reader) error                   { return nil }
func (f *fakeStreamable) Merge(context.Context, io.Reader) error                  { return nil }

func (f *fakeStreamable) ApplyReplicated(_ context.Context, records []changefeed.Record) error {
	f.applied = append(f.applied, records)
	return f.err
}

func sampleRecords() []changefeed.Record {
	return []changefeed.Record{
		{Key: []byte("pipeline:foo"), Value: []byte("\x01envelope-bytes")},
		{Key: []byte("transformer:bar"), Value: []byte{0x00, 0xff, 0x10}},
	}
}

func TestWriteNotificationRoundTrip(t *testing.T) {
	original := &writeNotification{
		Namespace: changefeed.NamespaceConfig,
		Records:   sampleRecords(),
	}

	data := original.Marshal()
	if len(data) == 0 {
		t.Fatal("Marshal returned empty payload")
	}
	if data[0] != writeNotificationTag {
		t.Fatalf("expected tag byte %d, got %d", writeNotificationTag, data[0])
	}

	parsed, err := parseNotification(data)
	if err != nil {
		t.Fatalf("parseNotification: %v", err)
	}

	wn, ok := parsed.(*writeNotification)
	if !ok {
		t.Fatalf("expected *writeNotification, got %T", parsed)
	}
	if wn.Namespace != original.Namespace {
		t.Errorf("namespace: got %q want %q", wn.Namespace, original.Namespace)
	}
	if !reflect.DeepEqual(wn.Records, original.Records) {
		t.Errorf("records round-trip mismatch:\n got %#v\nwant %#v", wn.Records, original.Records)
	}
}

func TestParseNotificationRejectsInvalid(t *testing.T) {
	cases := map[string][]byte{
		"empty":        nil,
		"unknown tag":  {0xff, 0x01, 0x02},
		"invalid json": {writeNotificationTag, '{', 'n', 'o'},
	}

	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := parseNotification(data); err == nil {
				t.Fatalf("expected error for %s payload", name)
			}
		})
	}
}

func TestWriteNotificationHandleDispatch(t *testing.T) {
	fake := &fakeStreamable{}
	d := &delegate{
		storages: map[string]storage.Streamable{
			changefeed.NamespaceConfig: fake,
		},
	}

	records := sampleRecords()
	wn := &writeNotification{Namespace: changefeed.NamespaceConfig, Records: records}

	if err := wn.Handle(context.Background(), d); err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if len(fake.applied) != 1 {
		t.Fatalf("expected ApplyReplicated to be called once, got %d", len(fake.applied))
	}
	if !reflect.DeepEqual(fake.applied[0], records) {
		t.Errorf("dispatched records mismatch:\n got %#v\nwant %#v", fake.applied[0], records)
	}
}

func TestWriteNotificationHandleUnknownNamespace(t *testing.T) {
	d := &delegate{storages: map[string]storage.Streamable{}}
	wn := &writeNotification{Namespace: "does-not-exist", Records: sampleRecords()}

	if err := wn.Handle(context.Background(), d); err == nil {
		t.Fatal("expected an error for an unknown namespace")
	}
}

// TestWriteNotificationHandlesBinaryRecords guards against any future switch to
// an encoding that cannot carry arbitrary bytes (LWW envelopes are binary).
func TestWriteNotificationHandlesBinaryRecords(t *testing.T) {
	records := []changefeed.Record{
		{Key: []byte{0x00, 0x01, 0x02}, Value: bytes.Repeat([]byte{0xff}, 32)},
	}
	wn := &writeNotification{Namespace: changefeed.NamespaceAuth, Records: records}

	parsed, err := parseNotification(wn.Marshal())
	if err != nil {
		t.Fatalf("parseNotification: %v", err)
	}
	if !reflect.DeepEqual(parsed.(*writeNotification).Records, records) {
		t.Error("binary records did not survive marshal round-trip")
	}
}
