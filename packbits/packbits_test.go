package packbits

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestEncode(t *testing.T) {
	assertEncodedEqual([]byte{0x00}, []byte{0x00, 0x00}, t)
	assertEncodedEqual([]byte{0x00, 0x01, 0x02}, []byte{0x02, 0x00, 0x01, 0x02}, t)
	assertEncodedEqual([]byte{0x00, 0x00}, []byte{0xff, 0x00}, t)
	assertEncodedEqual([]byte{0x00, 0x00, 0x00}, []byte{0xfe, 0x00}, t)
	assertEncodedEqual([]byte{0x02, 0x01, 0x00, 0x00, 0x00, 0x03, 0x04}, []byte{0x01, 0x02, 0x01, 0xfe, 0x00, 0x01, 0x03, 0x04}, t)
}

func TestRecoverOriginalBytes(t *testing.T) {
	assertRecoverd([]byte{0x00}, t)
	assertRecoverd([]byte{0x00, 0x01, 0x02}, t)
	assertRecoverd([]byte{0x00, 0x00}, t)
	assertRecoverd([]byte{0x00, 0x00, 0x00}, t)
	const randomIteration = 1000
	const randomLength = 100
	for i := 0; i < randomIteration; i++ {
		max := rand.Intn(256)
		data := randomBytes(randomLength, byte(max))
		assertRecoverd(data, t)
	}
}

func randomBytes(n int, max byte) []byte {
	bytes := make([]byte, n)
	for i, _ := range bytes {
		bytes[i] = byte(rand.Intn(int(max) + 1))
	}
	return bytes
}

func assertEncodedEqual(value []byte, expected []byte, t *testing.T) {
	encodedValue, _ := Encode(value)
	if bytes.Compare(encodedValue, expected) != 0 {
		t.Fatalf("expected %v but %v", expected, encodedValue)
	}
}

func assertRecoverd(value []byte, t *testing.T) {
	encodedValue, _ := Encode(value)
	decodedValue, _ := Decode(encodedValue)
	expected := value
	if bytes.Compare(decodedValue, expected) != 0 {
		t.Fatalf("expected %v but %v\n(value:%v, encoded:%v)", expected, decodedValue, value, encodedValue)
	}
}
