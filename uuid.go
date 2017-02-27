package gouuid

import (
	"encoding/binary"
	"encoding/hex"
	"sync"
	"time"
)

const (
	//this is greogian epoch 1582-10-15 00:00
	epochStart = 122192928000000000
)

type UUID [16]byte

var (
	clockseq uint16
	node     []byte
	mu       sync.Mutex
	muOnce   sync.Once
	lastTime uint64
	timeNow  uint64
)

func initClockSeq() {
	b, err := rng(2)
	if err != nil {
		panic(err)
	}
	clockseq = binary.BigEndian.Uint16(b)
}

func createNode() {
	b, err := rng(6)
	if err != nil {
		panic(err)
	}
	node = []byte{
		b[0] | 0x01,
		b[1],
		b[2],
		b[3],
		b[4],
		b[5],
	}
}

func getTimeIntervals() uint64 {
	return epochStart + uint64((time.Now().UnixNano() / 100))
}

func getStorage() (uint16, uint64, []byte) {
	muOnce.Do(initialize)

	mu.Lock()
	defer mu.Unlock()

	timeNow := getTimeIntervals()

	if timeNow <= lastTime {
		clockseq = clockseq + 1&0x3fff
	}
	lastTime = timeNow
	return clockseq, timeNow, node
}

func initialize() {
	initClockSeq()
	createNode()
}

func NewV1() *UUID {
	u := UUID{}
	clockseq, timeNow, node := getStorage()
	//Set the time_low field equal to the least significant 32 bits
	//(bits zero through 31) of the timestamp in the same order of
	//significance.
	binary.BigEndian.PutUint32(u[0:], uint32(timeNow))
	//Set the time_mid field equal to bits 32 through 47 from the
	//timestamp in the same order of significance.
	binary.BigEndian.PutUint16(u[4:], uint16(timeNow>>32))
	//Set the 12 least significant bits (bits zero through 11) of the
	//time_hi_and_version field equal to bits 48 through 59 from the
	//timestamp in the same order of significance.
	//Set the four most significant bits (bits 12 through 15) of the
	//time_hi_and_version field to the 4-bit version number
	//corresponding to the UUID version being created, as shown in the
	//table above.
	setTimeHiAndVersion(&u, 1, timeNow)

	//Set the clock_seq_low field to the eight least significant bits
	//(bits zero through 7) of the clock sequence in the same order of
	//significance.
	//this is setting both 8 and 9th octet with clockseq
	binary.BigEndian.PutUint16(u[8:], clockseq)
	//Set the 6 least significant bits (bits zero through 5) of the
	//clock_seq_hi_and_reserved field to the 6 most significant bits
	//(bits 8 through 13) of the clock sequence in the same order of
	//significance.
	setVariant(&u)
	//Set the node field to the 48-bit IEEE address in the same order of
	//significance as the address.
	copy(u[10:], node)

	return &u
}

//Format func return xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func (u *UUID) Format() string {
	//UUID = time-low "-" time-mid "-"
	//time-high-and-version "-"
	//clock-seq-and-reserved
	//clock-seq-low "-" node
	b := make([]byte, 36)
	//time_low
	hex.Encode(b[0:8], u[0:4])
	b[8] = '-'
	//time_mid
	hex.Encode(b[9:13], u[4:6])
	b[13] = '-'
	//time_hi_and_version
	hex.Encode(b[14:18], u[6:8])
	//clock_seq_hi_and_reserved
	b[18] = '-'
	hex.Encode(b[19:23], u[8:10])
	b[23] = '-'
	hex.Encode(b[24:], u[10:])
	return string(b)
}

func setVariant(u *UUID) {
	//set 6 lsb to zero
	//x := u[8] & 0x3f
	u[8] = (u[8] & 0x3f) | 0x80
}

func setTimeHiAndVersion(u *UUID, v int, timeNow uint64) {
	binary.BigEndian.PutUint16(u[6:], uint16(timeNow>>48))
	setVersion(u, 1)
}
func setVersion(u *UUID, v byte) {
	//clear the 4 msb
	u[6] = (u[6] & 0x0f) | (v << 4)
	// 8 bit e.g. 1 i.e 0001 <<4 becomes 00010000 2 i.e. 0010 << 4 becomes 00100000
	v8 := v << 4
	//now OR to set version
	u[6] = u[6] | v8
}

func NewV4() *UUID {
	u := UUID{}
	r, _ := rng(16)
	//clockseq, timeNow, node := getStorage()
	// if time != 0 {
	// 	timeNow = time
	// }
	copy(u[0:], r)
	setVersion(&u, 4)
	setVariant(&u)
	return &u
}
