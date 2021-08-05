package proxy

import (
	"testing"
)

func Test_WriteReadBuffer(t *testing.T) {
	buf := AllocateBuffer()
	buf.WriteByte(0)        // 1 bytes
	buf.WriteUint16(1)      // 2 bytes
	buf.WriteUint32(2)      // 4 bytes
	buf.WriteUint(3)        // 4 bytes
	buf.WriteUint64(4)      // 8 bytes
	buf.WriteInt16(5)       // 2 bytes
	buf.WriteInt32(6)       // 4 bytes
	buf.WriteInt(7)         // 4 bytes
	buf.WriteInt64(8)       // 8 bytes
	buf.WriteString("9")    // 1 bytes
	buf.Write([]byte("10")) // 2 bytes

	n := 1*2 + 2*3 + 4*4 + 8*2
	// check buffer length
	if buf.Len() != n {
		t.Fatalf("expect buf length %d, actual %d, may be buff bug.", n, buf.Len())
	}

	// check buffer data
	if b, err := buf.ReadByte(); b != 0 || err != nil {
		t.Fatalf("expect read byte 0, actual %d, err: %v", b, err)
	}
	if i, err := buf.ReadUint16(); i != 1 || err != nil {
		t.Fatalf("expect read uint16 1, actual %d, err: %v", i, err)
	}
	if i, err := buf.ReadUint32(); i != 2 || err != nil {
		t.Fatalf("expect read uint32 2, actual %d, err: %v", i, err)
	}
	if i, err := buf.ReadUint(); i != 3 || err != nil {
		t.Fatalf("expect read uint 3, actual %d, err: %v", i, err)
	}
	if i, err := buf.ReadUint64(); i != 4 || err != nil {
		t.Fatalf("expect read uint64 4, actual %d, err: %v", i, err)
	}
	if i, err := buf.ReadInt16(); i != 5 || err != nil {
		t.Fatalf("expect read int16 5, actual %d, err: %v", i, err)
	}
	if i, err := buf.ReadInt32(); i != 6 || err != nil {
		t.Fatalf("expect read int32 6, actual %d, err: %v", i, err)
	}
	if i, err := buf.ReadInt(); i != 7 || err != nil {
		t.Fatalf("expect read int 7, actual %d, err: %v", i, err)
	}
	if i, err := buf.ReadInt64(); i != 8 || err != nil {
		t.Fatalf("expect read int64 8, actual %d, err: %v", i, err)
	}
	if s := string(buf.Bytes()); s != "910" {
		t.Fatalf("expect read string 910, actual %s", s)
	}
}

func Test_PutGetBuffer(t *testing.T) {
	buf := AllocateBuffer()
	buf.PutByte(0, 0)    // 1 bytes
	buf.PutUint16(1, 1)  // 2 bytes
	buf.PutUint32(3, 2)  // 4 bytes
	buf.PutUint(7, 3)    // 4 bytes
	buf.PutUint64(11, 4) // 8 bytes
	buf.PutInt16(19, 5)  // 2 bytes
	buf.PutInt32(21, 6)  // 4 bytes
	buf.PutInt(25, 7)    // 4 bytes
	buf.PutInt64(29, 8)  // 8 bytes

	n := 1 + 2*2 + 4*4 + 8*2
	// check buffer length
	if buf.Len() != n {
		t.Fatalf("expect buf length %d, actual %d, may be buff bug.", n, buf.Len())
	}

	// check pos not changed
	if buf.Pos() != 0 {
		t.Fatalf("buf pos should not be moved, actual %d", buf.Pos())
	}

	// check buffer data
	if b, err := buf.GetByte(0); b != 0 || err != nil {
		t.Fatalf("expect get byte 0, actual %d, err: %v", b, err)
	}
	if i, err := buf.GetUint16(1); i != 1 || err != nil {
		t.Fatalf("expect get uint16 1, actual %d, err: %v", i, err)
	}
	if i, err := buf.GetUint32(3); i != 2 || err != nil {
		t.Fatalf("expect get uint32 2, actual %d, err: %v", i, err)
	}
	if i, err := buf.GetUint(7); i != 3 || err != nil {
		t.Fatalf("expect get uint 3, actual %d, err: %v", i, err)
	}
	if i, err := buf.GetUint64(11); i != 4 || err != nil {
		t.Fatalf("expect get uint64 4, actual %d, err: %v", i, err)
	}
	if i, err := buf.GetInt16(19); i != 5 || err != nil {
		t.Fatalf("expect get int16 5, actual %d, err: %v", i, err)
	}
	if i, err := buf.GetInt32(21); i != 6 || err != nil {
		t.Fatalf("expect get int32 6, actual %d, err: %v", i, err)
	}
	if i, err := buf.GetInt(25); i != 7 || err != nil {
		t.Fatalf("expect get int 7, actual %d, err: %v", i, err)
	}
	if i, err := buf.GetInt64(29); i != 8 || err != nil {
		t.Fatalf("expect get int64 8, actual %d, err: %v", i, err)
	}

}
