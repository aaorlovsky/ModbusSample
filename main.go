package main

import (
	"fmt"
	"time"

	"github.com/simonvetter/modbus"
)

func main() {
    var client  *modbus.ModbusClient
    var err      error

    // for a TCP endpoint
    // (see examples/tls_client.go for TLS usage and options)
    client, err = modbus.NewClient(&modbus.ClientConfiguration{
        URL:      "tcp://hostname-or-ip-address:502",
        Timeout:  1 * time.Second,
    })
    // note: use udp:// for modbus TCP over UDP

    // for an RTU (serial) device/bus
    client, err = modbus.NewClient(&modbus.ClientConfiguration{
        URL:      "rtu:///dev/ttyUSB0",
        Speed:    19200,                   // default
        DataBits: 8,                       // default, optional
        Parity:   modbus.PARITY_NONE,      // default, optional
        StopBits: 2,                       // default if no parity, optional
        Timeout:  300 * time.Millisecond,
    })

    // for an RTU over TCP device/bus (remote serial port or
    // simple TCP-to-serial bridge)
    client, err = modbus.NewClient(&modbus.ClientConfiguration{
        URL:      "rtuovertcp://hostname-or-ip-address:502",
        Speed:    19200,                   // serial link speed
        Timeout:  1 * time.Second,
    })
    // note: use rtuoverudp:// for modbus RTU over UDP

    if err != nil {
        // error out if client creation failed
    }

    // now that the client is created and configured, attempt to connect
    err = client.Open()
    if err != nil {
        // error out if we failed to connect/open the device
        // note: multiple Open() attempts can be made on the same client until
        // the connection succeeds (i.e. err == nil), calling the constructor again
        // is unnecessary.
        // likewise, a client can be opened and closed as many times as needed.
    }

    // read a single 16-bit holding register at address 100
    var reg16   uint16
    reg16, err  = client.ReadRegister(100, modbus.HOLDING_REGISTER)
    if err != nil {
      // error out
    } else {
      // use value
      fmt.Printf("value: %v", reg16)        // as unsigned integer
      fmt.Printf("value: %v", int16(reg16)) // as signed integer
    }

    // read 4 consecutive 16-bit input registers starting at address 100
    var reg16s  []uint16
    reg16s, err = client.ReadRegisters(100, 4, modbus.INPUT_REGISTER)

	if err != nil {
		// error out
	  } else {
		// use value
		fmt.Printf("value: %v", reg16s)        // as unsigned integer
		//fmt.Printf("value: %v", uint16(reg16s)) // as signed integer
	  }

    // read the same 4 consecutive 16-bit input registers as 2 32-bit integers
    //var reg32s  []uint32
    //reg32s, err = client.ReadUint32s(100, 2, modbus.INPUT_REGISTER)

    // read the same 4 consecutive 16-bit registers as a single 64-bit integer
    //var reg64   uint64
    //reg64, err  = client.ReadUint64(100, modbus.INPUT_REGISTER)

    // read the same 4 consecutive 16-bit registers as a slice of bytes
    //var regBs   []byte
    //regBs, err  = client.ReadBytes(100, 8, modbus.INPUT_REGISTER)

    // by default, 16-bit integers are decoded as big-endian and 32/64-bit values as
    // big-endian with the high word first.
    // change the byte/word ordering of subsequent requests to little endian, with
    // the low word first (note that the second argument only affects 32/64-bit values)
    client.SetEncoding(modbus.LITTLE_ENDIAN, modbus.LOW_WORD_FIRST)

    // read the same 4 consecutive 16-bit input registers as 2 32-bit floats
    //var fl32s   []float32
    //fl32s, err  = client.ReadFloat32s(100, 2, modbus.INPUT_REGISTER)

    // write -200 to 16-bit (holding) register 100, as a signed integer
    var s int16 = -200
    err         = client.WriteRegister(100, uint16(s))

    // Switch to unit ID (a.k.a. slave ID) #4
    client.SetUnitId(4)

    // write 3 floats to registers 100 to 105
    err         = client.WriteFloat32s(100, []float32{
        3.14,
        1.1,
        -783.22,
    })

    // write 0x0102030405060708 to 16-bit (holding) registers 10 through 13
    // (8 bytes i.e. 4 consecutive modbus registers)
    err         = client.WriteBytes(10, []byte{
        0x01, 0x02, 0x03, 0x04,
        0x05, 0x06, 0x07, 0x08,
    })

    // close the TCP connection/serial port
    client.Close()
}