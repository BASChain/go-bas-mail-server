package bmtpserver

import (
	"github.com/realbmail/go-bas-mail-server/wallet"
	"github.com/realbmail/go-bmail-protocol/translayer"
	"log"
	"net"
	"sync"
)

type BMTPFunc func(*TcpSession) error

var (
	bmtpserverInst     BMTPServerIntf
	bmtpserverInstLock sync.Mutex
)

type BMTPServerConf struct {
	ListenPort     int
	SupportFunc    map[int]BMTPFunc
	supportVersion []uint16
	Session        map[string]*TcpSession
	listener       *net.TCPListener
	quit           chan interface{}
	wg             sync.WaitGroup
	timeout        int
	wallet         wallet.ServerWalletIntf
}

type BMTPServerIntf interface {
	StartTCPServer() error
	StopTCPServer()
	VersionInSrv(version int) bool
	SupportVersion() []uint16
}

func NewServer2(listenport int) BMTPServerIntf {
	server := &BMTPServerConf{}

	server.ListenPort = listenport
	server.quit = make(chan interface{})
	server.SupportFunc = make(map[int]BMTPFunc)
	server.SupportFunc[int(translayer.BMAILVER1)] = HandleMsgV1
	server.Session = make(map[string]*TcpSession)
	server.quit = make(chan interface{}, 1)
	server.timeout = 300 //second
	server.wallet = wallet.GetServerWallet()

	return server
}

func GetBMTPServer() BMTPServerIntf {
	if bmtpserverInst == nil {
		bmtpserverInstLock.Lock()
		bmtpserverInstLock.Unlock()
		if bmtpserverInst == nil {
			bmtpserverInst = NewServer2(translayer.BMTP_PORT)
		}
	}

	return bmtpserverInst
}

func (s *BMTPServerConf) VersionInSrv(version int) bool {
	if _, ok := s.SupportFunc[version]; !ok {
		return false
	} else {
		return true
	}
}

func (s *BMTPServerConf) SupportVersion() []uint16 {
	if s.supportVersion == nil {
		i := 0
		for k, _ := range s.SupportFunc {
			s.supportVersion = append(s.supportVersion, uint16(k))
			i++
		}
	}

	return s.supportVersion
}

func (s *BMTPServerConf) StartTCPServer() error {
	laddr := &net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: s.ListenPort}

	l, err := net.ListenTCP("tcp4", laddr)
	if err != nil {
		return err
	}

	s.listener = l
	//defer l.Close()

	s.listener = l
	s.wg.Add(1)
	go s.serve()

	s.wg.Wait()

	return nil
}

func (s *BMTPServerConf) serve() {
	defer s.wg.Done()

	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				log.Println("accept error", err)
			}
		} else {
			s.wg.Add(1)
			go func() {
				s.handleConnect(conn)
				s.wg.Done()
			}()
		}
	}

}

func (s *BMTPServerConf) handleConnect(conn *net.TCPConn) {
	raddrstr := conn.RemoteAddr().String()

	defer func(raddr string) {
		conn.Close()
		delete(s.Session, raddr)
	}(raddrstr)

	ac := &TcpSession{}
	ac.conn = conn
	ac.server = s

	s.Session[raddrstr] = ac

	//conn.SetDeadline(time.Now().Add(time.Duration(s.timeout) * time.Second))

	if err := ac.Negotiation(); err != nil {
		log.Println(err)
		return
	}
	for {
		//conn.SetDeadline(time.Now().Add(time.Duration(s.timeout) * time.Second))
		select {
		case <-s.quit:
			return
		default:
			if err := ac.Handle(ac); err != nil {
				log.Println(err)
				return
			}
			s.timeout = 600 //second
		}
	}
}

func (s *BMTPServerConf) StopTCPServer() {
	for _, c := range s.Session {
		c.conn.Close()
	}
	s.listener.Close()
	close(s.quit)

	return
}
