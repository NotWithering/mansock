package mansock

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/NotWithering/argo"
)

var config = map[string]any{
	"IP":    "127.0.0.1",
	"PORT":  "23",
	"PROTO": "tcp",
}

// errors
const (
	errorCustom         string = "error: %s\n"
	errorNEA            string = "error: not enough arguments"
	errorSyntax         string = "error: syntax error"
	errorUnknownCommand string = "error: unknown command"
)

func main() {
	var socket net.Conn

	buf := bytes.NewBuffer([]byte{})

	var exit bool
	for !exit {
		fmt.Print("MANSOCK> ")

		reader := bufio.NewReader(os.Stdin)
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		command = strings.TrimSpace(command)

		args, incomplete := argo.Parse(command)
		if incomplete {
			fmt.Println(errorSyntax)
			continue
		}

		switch args[0] {
		case "h", "help":
			fmt.Println("h, help           : show this help menu       : h")
			fmt.Println("s, set            : set a configuration value : s PORT string \"127.0.0.1\"")
			fmt.Println("l, list           : list configurations       : l")
			fmt.Println("c, connect        : connect to socket         : c")
			fmt.Println("b, buffer         : see the buffer data       : b")
			fmt.Println("wb, writebuffer   : write to the buffer       : wb uint8 255")
			fmt.Println("cb, clearbuffer   : clear the buffer          : cb")
			fmt.Println("db, deflatebuffer : deflate the buffer        : db")
			fmt.Println("ws, writesocket   : write to the socket       : ws")
			fmt.Println("cs, closesocket   : close the socket          : cs")
			fmt.Println("e, exit           : exits the program         : e")
		case "s", "set": // SET key type value
			if len(args) < 4 {
				fmt.Println(errorNEA)
				continue
			}
			value, err := convert(args[3], args[2])
			if err != nil {
				fmt.Printf(errorCustom, err)
				continue
			}
			config[args[1]] = value
		case "l", "list":
			for k, v := range config {
				fmt.Printf("%s: %v (%s)\n", k, v, reflect.TypeOf(v).Name())
			}
		case "c", "connect":
			socket, err = net.Dial(fmt.Sprint(config["PROTO"]), fmt.Sprint(config["IP"])+":"+fmt.Sprint(config["PORT"]))
			if err != nil {
				fmt.Printf(errorCustom, err)
				continue
			}
		case "cs", "closesocket":
			if err := socket.Close(); err != nil {
				fmt.Printf(errorCustom, err)
				continue
			}
		case "b", "buffer":
			fmt.Println(buf.Bytes())
		case "wb", "writebuffer": // WRITEBUFFER type value
			if len(args) < 3 {
				fmt.Println(errorNEA)
				continue
			}

			value, err := convert(args[2], args[1])
			if err != nil {
				fmt.Printf(errorCustom, err)
				continue
			}

			switch val := value.(type) {
			case string:
				var fixedSizeString = make([]byte, len(val))
				copy(fixedSizeString[:], []byte(fmt.Sprint(value)))
				value = fixedSizeString
			}

			if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
				fmt.Printf(errorCustom, err)
				continue
			}
		case "cb", "clearbuffer":
			buf.Reset()
		case "db", "deflatebuffer":
			input := buf.Bytes()
			deflated := new(bytes.Buffer)

			writer := zlib.NewWriter(deflated)

			_, err := writer.Write(input)
			if err != nil {
				fmt.Printf(errorCustom, err)
				continue
			}

			buf.Reset()
			buf.Write(deflated.Bytes())
		case "ws", "writesocket":
			if socket == nil {
				fmt.Println("socket not connected")
				continue
			}
			if _, err := socket.Write(buf.Bytes()); err != nil {
				fmt.Printf(errorCustom, err)
				continue
			}
		case "e", "exit":
			exit = true
		case "":
		default:
			fmt.Println("unknown command")
		}
	}
	if socket != nil {
		if err := socket.Close(); err != nil {
			fmt.Printf(errorCustom, err)
		}
	}
}

func convert(value string, targetType string) (any, error) {
	switch targetType {
	case "int", "i":
		result, err := strconv.Atoi(value)
		return result, err
	case "int8", "i8":
		result, err := strconv.ParseInt(value, 10, 8)
		return int8(result), err
	case "int16", "i16":
		result, err := strconv.ParseInt(value, 10, 16)
		return int16(result), err
	case "int32", "i32":
		result, err := strconv.ParseInt(value, 10, 32)
		return int32(result), err
	case "int64", "i64":
		result, err := strconv.ParseInt(value, 10, 64)
		return result, err
	case "uint", "u":
		result, err := strconv.ParseUint(value, 10, 0)
		return result, err
	case "uint8", "u8":
		result, err := strconv.ParseUint(value, 10, 8)
		return uint8(result), err
	case "uint16", "u16":
		result, err := strconv.ParseUint(value, 10, 16)
		return uint16(result), err
	case "uint32", "u32":
		result, err := strconv.ParseUint(value, 10, 32)
		return uint32(result), err
	case "uint64", "u64":
		result, err := strconv.ParseUint(value, 10, 64)
		return result, err
	case "float32", "f32":
		result, err := strconv.ParseFloat(value, 32)
		return float32(result), err
	case "float64", "f64":
		result, err := strconv.ParseFloat(value, 64)
		return result, err
	case "bool", "b":
		result, err := strconv.ParseBool(value)
		return result, err
	case "string", "s":
		return value, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", targetType)
	}
}
