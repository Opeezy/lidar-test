package utils

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

type Serial struct {
	Name     string
	Baud     int
	Size     byte
	StopBits serial.StopBits
	Timeout  time.Duration

	Config     *serial.Config // Serial config pointer
	Port       *serial.Port   // Serial port pointer
	ConfigMade bool           // Valid config status
	PortOpen   bool           // Port connection status
}

type SerialWorker interface {
	OpenPortConnection()
	MakeConfig()
	Listen()
	Read()
}

func (s *Serial) OpenPortConnection() error {
	if s.ConfigMade {
		port, err := serial.OpenPort(s.Config)
		s.Port = port
		if err != nil {
			return err
		} else {
			log.Printf("Successfully opened port %v", s.Name)
			s.PortOpen = true
			return nil
		}
	} else {
		err := errors.New("no config made")
		return err
	}
}

func (s *Serial) MakeConfig() {
	NewConfig := &serial.Config{
		Name:        s.Name,          // Port name (eg. "COM5")
		Baud:        s.Baud,          // Baudrate
		Size:        s.Size,          // Data size (usually 8 bytes)
		StopBits:    s.StopBits,      // Stop bits (usually 1)
		ReadTimeout: time.Second * 5} // Timeout duration
	s.Config = NewConfig
	s.Name = NewConfig.Name
	s.ConfigMade = true
	log.Println("Config generated")
}

func (s *Serial) Listen(duration time.Duration) {
	if !s.ConfigMade {
		log.Fatalln("No config made")
	} else if !s.PortOpen {
		log.Fatalln("No port connection")
	} else {
		log.Println("Reading serial data")
		TimeNow := time.Now() // Get the current time

		data := make(chan []byte) // Channel to pass our packets through

		for {
			// Check for duration and read serial data until it is reached
			if time.Since(TimeNow) >= duration {
				log.Fatalln("5 seconds passed")
				break
			}
			go s.ReadSerial(data)
			go s.WriteToFile(data)
		}
	}
}

func (s *Serial) ReadSerial(c chan []byte) {
	buf := make([]byte, 47) // 47 byte buffer (the size of our lidar packets come out to 47 bytes)

	i, _ := s.Port.Read(buf) // Read incoming serial data into buffer
	_ = i

	// Check first two bytes as they are always static. 0x54 and 0x2C indicate a complete packet
	if buf[0] == 0x54 && buf[1] == 0x2C {
		c <- buf // Send our packet to channel
	}
}

func (s *Serial) WriteToFile(c chan []byte) {
	packet := <-c
	fmt.Println(packet)
}

func NewSerial(name string, baud int, size byte, stopbits serial.StopBits, timeout time.Duration) *Serial {
	serial := &Serial{
		Name:       name,
		Baud:       baud,
		Size:       size,
		StopBits:   stopbits,
		Timeout:    timeout,
		ConfigMade: false,
		PortOpen:   false,
	}

	serial.MakeConfig()
	return serial
}
