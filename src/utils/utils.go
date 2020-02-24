package utils

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/f0o/turbo-pancake/common"
)

//GetSTDIN : Returns contents of STDIN as string
func GetSTDIN() string {
	reader := bufio.NewReader(os.Stdin)
	var output []rune
	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}
	return string(output)
}

//UUIDv4 - Create Pseudo-Random UUID Version 4 (RFC4122)
func UUIDv4() (string, error) {
	var uuid [16]byte
	_, err := io.ReadFull(rand.Reader, uuid[:])
	if err != nil {
		return "", err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	var str [36]byte
	hex.Encode(str[:], uuid[:4])
	str[8] = '-'
	hex.Encode(str[9:13], uuid[4:6])
	str[13] = '-'
	hex.Encode(str[14:18], uuid[6:8])
	str[18] = '-'
	hex.Encode(str[19:23], uuid[8:10])
	str[23] = '-'
	hex.Encode(str[24:], uuid[10:])
	return string(str[:]), nil
}

//Debug : Print Debug Line
func Debug(prefix string, format string, a ...interface{}) {
	if *common.DebugFlag == true {
		log := fmt.Sprintf(format, a...)
		common.Logger.Printf("[%s] %s", prefix, log)
	}
}

//NewAndPublishEvent : Creates a New EventSpec and Publishes it
func NewAndPublishEvent(t string, v string, p interface{}) {
	e := NewEvent()
	e.Type = t
	e.Verb = v
	e.Payload = p
	PublishEvent(e)
}

//NewEvent : Create New EventSpec
func NewEvent() common.EventSpec {
	event := common.EventSpec{}
	event.ID, _ = UUIDv4()
	return event
}

//PublishEvent : Naive EventSpec Publisher
func PublishEvent(event common.EventSpec) error {
	select {
	case common.EventQueue <- event:
		return nil
	default:
		return errors.New("Could not Publish Event")
	}
}
