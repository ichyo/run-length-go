package packbits

import (
	"bytes"
	"log"
)

const (
	literalMaxLen   = 128
	contiguasMaxLen = 256 - literalMaxLen + 1
)

func Encode(b []byte) ([]byte, error) {
	var res bytes.Buffer
	var lit bytes.Buffer
	last_idx := 0

	flushLitralBuffer := func() {
		if lit.Len() != 0 {
			if lit.Len() > literalMaxLen {
				log.Fatalf("lit.Len() is over %d.", literalMaxLen)
			}
			res.WriteByte(byte(lit.Len() - 1))
			lit.WriteTo(&res)
			lit.Reset()
		}
	}
	appendSingleValue := func(x byte) {
		lit.WriteByte(x)
		if lit.Len() == literalMaxLen {
			flushLitralBuffer()
		}
	}
	appendContiguasValue := func(x byte, count int) {
		flushLitralBuffer()
		if count > contiguasMaxLen {
			log.Fatalf("count is over %d.", contiguasMaxLen)
		}
		res.WriteByte(byte(255 - count + 2))
		res.WriteByte(x)
	}
	appendValue := func(x byte, count int) {
		if count == 1 {
			appendSingleValue(x)
		} else {
			appendContiguasValue(x, count)
		}
	}

	for i, x := range b {
		if l := i - last_idx; x != b[last_idx] || l == contiguasMaxLen {
			appendValue(b[last_idx], l)
			last_idx = i
		}
	}

	appendValue(b[last_idx], len(b)-last_idx)
	flushLitralBuffer()

	return res.Bytes(), nil
}

func Decode(b []byte) ([]byte, error) {
	var res bytes.Buffer
	for i := 0; i < len(b); i++ {
		if b[i]+1 <= literalMaxLen {
			cnt := int(b[i] + 1)
			for j := 0; j < cnt; j++ {
				res.WriteByte(b[i+j+1])
			}
			i += cnt
		} else {
			cnt := 257 - int(b[i])
			for j := 0; j < cnt; j++ {
				res.WriteByte(b[i+1])
			}
			i += 1
		}
	}
	return res.Bytes(), nil
}
