package helper

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

const (
	TYPE_UNIX = "unix"
	TYPE_TCP  = "tcp"
)

var UnsupportedProtocolError = errors.New("This protocol is not supported")

//Livestatus fetches data from livestatus.
type Livestatus struct {
	address        string
	connectionType string
}

func NewLivestatus(address string) (*Livestatus, error) {
	splitted := strings.SplitN(address, ":", 2)
	if splitted[0] == TYPE_UNIX || splitted[0] == TYPE_TCP {
		return &Livestatus{address: splitted[1], connectionType: splitted[0]}, nil
	} else {
		return nil, UnsupportedProtocolError
	}

}

//Queries livestatus and returns an list of list outer list are lines inner elements within the line.
func (l Livestatus) Query(query string) (*[][]string, error) {
	var conn net.Conn
	switch l.connectionType {
	case "tcp":
		conn, _ = net.Dial(TYPE_TCP, l.address)
	case "file":
		conn, _ = net.Dial(TYPE_UNIX, l.address)
	}
	defer conn.Close()
	fmt.Fprintf(conn, query)
	reader := bufio.NewReader(conn)

	result := [][]string{}

	length := 1
	for length > 0 {
		message, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		length = len(message)
		if length > 0 {
			csvReader := csv.NewReader(strings.NewReader(string(message)))
			csvReader.Comma = ';'
			csvReader.LazyQuotes = true
			records, err := csvReader.Read()
			if err != nil {
				return nil, err
			}
			result = append(result, records)
		}
	}
	return &result, nil
}
