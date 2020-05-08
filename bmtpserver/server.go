package bmtpserver

import (
	"github.com/BASChain/go-bmail-protocol/translayer"
	"net"
	"github.com/pkg/errors"
	"github.com/BASChain/go-bmail-protocol/bmprotocol"
	"crypto/rand"
)

type BMTPServer struct {
	listenPort int
}

var gac map[string]*AcceptConn

func init()  {
	gac = make(map[string]*AcceptConn)
}

type AcceptConn struct {
	sn []byte
	pubkey []byte
	//eid translayer.EnveUniqID
	conn *net.TCPConn
}

func (ac *AcceptConn)RcvEnvelope(bmtl *translayer.BMTransLayer)  error {
	if bmtl.GetDataLen() == 0{
		return errors.New("Receive error envelope message")
	}

	buf:=make([]byte,bmtl.GetDataLen())
	n,err := ac.conn.Read(buf)
	if err!=nil || n != len(buf){
		return err
	}

	se:=&bmprotocol.SendEnvelope{}
	se.BMTransLayer = *bmtl

	_,err=se.UnPack(buf)
	if err!=nil{
		return err
	}

	//check sig
	//todo...

	//send
	rse := bmprotocol.NewRespSendEnvelope()
	rse.NewSn = NewSn()
	rse.Sn = se.Sn
	rse.EId = se.EId
	rse.ErrId = 0

	ac.sn = rse.NewSn

	//ac.eid = se.EId
	ac.pubkey = se.LPubKey


	//pack and send

	data,err:=rse.Pack()
	if err!=nil{
		return err
	}

	n,err=ac.conn.Write(data)
	if err!=nil{
		return err
	}

	return nil

}


func (ac *AcceptConn)RcvCryptEnvelope(bmtl *translayer.BMTransLayer)  error{
	if bmtl.GetDataLen() == 0{
		return errors.New("Receive error envelope message")
	}

	buf:=make([]byte,bmtl.GetDataLen())
	n,err := ac.conn.Read(buf)
	if err!=nil || n != len(buf){
		return err
	}

	se:=&bmprotocol.SendCryptEnvelope{}
	se.BMTransLayer = *bmtl

	_,err=se.UnPack(buf)
	if err!=nil{
		return err
	}

	//check sig
	//todo...

	//send
	rse := bmprotocol.NewRespSendCryptEnvelope()
	rse.NewSn = NewSn()
	rse.Sn = se.Sn
	rse.EId = se.EId
	rse.ErrId = 0

	ac.sn = rse.NewSn

	//ac.eid = se.EId
	ac.pubkey = se.LPubKey


	//pack and send

	data,err:=rse.Pack()
	if err!=nil{
		return err
	}

	n,err=ac.conn.Write(data)
	if err!=nil{
		return err
	}

	return nil
}

func NewSn() []byte  {
	sn := make([]byte, 16)

	for {
		n, _ := rand.Read(sn)
		if n != len(sn) {
			continue
		}
		break
	}

	return sn
}

func (ac *AcceptConn)RcvHelo(bmtl *translayer.BMTransLayer) error {
	if bmtl.GetDataLen() > 0{
		return errors.New("Receive error helo message")
	}

	ack:=bmprotocol.NewBMHelloACK(NewSn())
	data,err:=ack.Pack()
	if err!=nil{
		return err
	}

	ac.conn.Write(data)

	ac.sn = ack.GetSn()


	return nil
}


func (s *BMTPServer)StartTCPServer() error {
	laddr:=&net.TCPAddr{IP:net.ParseIP("0.0.0.0"),Port:s.listenPort}

	l,err:=net.ListenTCP("tcp4",laddr)
	if err!=nil{
		return err
	}

	defer l.Close()

	for{
		c, err:=l.AcceptTCP()
		if err!=nil{
			return err
		}

		go handleConnect(c)
	}

}

func handleConnect(c *net.TCPConn) error {

	raddrstr:=c.RemoteAddr().String()

	defer func(raddr string) {
		c.Close()
		delete(gac,raddr)
	}(raddrstr)

	ac:=&AcceptConn{}
	ac.conn = c

	gac[raddrstr] = ac

	for{
		buf:=make([]byte,translayer.BMHeadSize())
		n,err := ac.conn.Read(buf)
		if err!=nil || n != len(buf){
			return err
		}

		bmtl:=&translayer.BMTransLayer{}
		bmtl.UnPack(buf)

		switch bmtl.GetMsgType() {
		case translayer.SEND_ENVELOPE:
			if err:=ac.RcvEnvelope(bmtl);err!=nil{
				return err
			}
		case translayer.HELLO:
			if err := ac.RcvHelo(bmtl);err!=nil{
				return err
			}
		case translayer.SEND_CRYPT_ENVELOPE:
			if err:=ac.RcvCryptEnvelope(bmtl);err!=nil{
				return err
			}
		}


	}


}

