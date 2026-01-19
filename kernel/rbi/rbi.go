package rbi

import (
	"fmt"
	"kernel/tools"
	"net"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

type Rbi struct {
	open bool

	Conn *net.UDPConn
}

var rbi *Rbi

func init() {
	rbi = new(Rbi)
}

func SetOpen(o bool) {
	rbi.open = o
}

func IsOpen() bool {
	return rbi.open
}

func Init(url, port string) error {
	if !IsOpen() {
		return nil
	}

	//address := fmt.Sprintf("%s:%s", Rbi_Server_Url, Rbi_Server_Port)
	//if test {
	//	address = fmt.Sprintf("%s:%s", Rbi_Server_Url, Rbi_Server_Port_Test)
	//}
	address := fmt.Sprintf("%s:%s", url, port)
	serverAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Error("ResolveUDPAddr failed", zap.Error(err))
		return err
	}

	rbi.Conn, err = net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		log.Error("DialUDP failed", zap.Error(err))
		return err
	}

	return nil
}

func GetRbi() *Rbi {
	return rbi
}

func Write(data string) error {
	//dataBytes, err := encodeGob(data)
	//if err != nil {
	//	log.Error("udp encode gob err", zap.Error(err))
	//	return err
	//}
	if data == "" {
		return nil
	}
	_, err := GetRbi().Conn.Write([]byte(data))
	if err != nil {
		log.Error("udp write err", zap.Error(err))
		return err
	}

	return nil
}

func RbiWrite(data IRbi) {
	if !IsOpen() {
		return
	}

	go tools.GoSafe("rbi "+data.Name(), func() {
		if err := Write(StructToPipeString(data)); err != nil {
			log.Error(data.Name()+" err", zap.Error(err))
		}
	})
}
