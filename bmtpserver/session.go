package bmtpserver

import (
	"net"

	"io"
	"log"
	"errors"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"github.com/BASChain/go-bmail-protocol/bmp"
	"strconv"
	"crypto/rand"
	"github.com/BASChain/go-bas-mail-server/protocol"
)


type TcpSession struct {
	sn []byte
	pubkey []byte
	conn *net.TCPConn
	Handle BMTPFunc
	size int
	server *BMTPServerConf
	version int
	buf []byte
	bmtl *translayer.BMTransLayer
	rbody protocol.MsgBody
	wbody bmp.EnvelopeMsg
}


func newSn() []byte  {
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

func (ts *TcpSession)Negotiation() error  {
	if err:=ts.readHead();err!=nil{
		return err
	}

	support:=ts.server.VersionInSrv(int(ts.bmtl.GetVersion()))

	ack:=&bmp.HELOACK{}
	if support == false{
		ack.ErrCode = 1
		ack.SupportVersion = ts.server.SupportVersion()
	}else{
		ack.ErrCode = 0
		ts.sn = newSn()
		copy(ack.SN[:],ts.sn)
		ack.SrvBca = ts.server.wallet.BCAddress()
	}

	ts.wbody = ack

	if err:=ts.WriteMsg();err!=nil{
		return err
	}

	if support{
		ts.Handle = ts.server.SupportFunc[int(ts.bmtl.GetVersion())]
		return nil
	}else{
		return errors.New("client version is not support, ver: "+strconv.Itoa(int(ts.bmtl.GetVersion())))
	}

}

func (ts *TcpSession)WriteMsg() error {
	bmtl:=translayer.NewBMTL(ts.wbody.MsgType())

	buf,err:=ts.wbody.GetBytes()
	if err!=nil{
		return err
	}

	bmtl.SetDataLen(uint32(len(buf)))

	data,_:=bmtl.Pack()

	var n int
	n,err = ts.conn.Write(data)
	if err!=nil || n != len(data){
		return errors.New("Write "+strconv.Itoa(int(ts.wbody.MsgType()))+" message head Failed")
	}

	n,err = ts.conn.Write(buf)
	if err!=nil || n != len(data){
		return errors.New("Write "+strconv.Itoa(int(ts.wbody.MsgType()))+" message body Failed")
	}

	return nil
}



func (ts *TcpSession)readNext() (error)  {
	buf:=make([]byte,ts.size)
	for{
		n,err := ts.conn.Read(buf)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			} else if err != io.EOF {
				log.Println("read error", err)
				return err
			}
		}
		if n == 0 {
			return errors.New("no data to read")
		}

		ts.buf = buf

		return nil
	}
}

func (ts *TcpSession)deriveHead() (*translayer.BMTransLayer,error) {
	if len(ts.buf) != translayer.BMHeadSize(){
		return nil,errors.New("data is not a header buffer")
	}

	bmtl:=&translayer.BMTransLayer{}

	_,err:=bmtl.UnPack(ts.buf)
	if err!=nil{
		return nil,err
	}

	return bmtl,nil

}

func (ts *TcpSession)deriveBody() error {
	ts.rbody = protocol.MsgGrid[ts.bmtl.GetMsgType()]
	return ts.rbody.UnPack(ts.buf)
}

func (ts *TcpSession)readHead() error  {
	ts.size = translayer.BMHeadSize()

	if err := ts.readNext(); err!=nil{
		return err
	}

	if head,err:=ts.deriveHead();err!=nil{
		return err
	}else{
		ts.bmtl = head
		ts.size = int(head.GetDataLen())
		return nil
	}
}

func (ts *TcpSession)readBody() error  {
	if ts.bmtl == nil || ts.size == 0{
		return errors.New("please read message header first")
	}

	if err:=ts.readNext();err!=nil{
		return err
	}else{
		if err:=ts.deriveBody();err!=nil{
			return err
		}
	}

	return nil
}


func HandleMsgV1(ts *TcpSession) error  {
	if err:=ts.readHead();err!=nil{
		return err
	}

	if err:=ts.readBody();err!=nil{
		return err
	}

	//if err:=ts.rbody.Save2DB();err!=nil{
	//	return err
	//}

	if resp,err:=ts.rbody.Response();err!=nil{
		ts.wbody = resp
		ts.WriteMsg()
		return err
	}else{
		ts.wbody = resp
		if err:=ts.WriteMsg();err!=nil{
			return err
		}
	}

	return nil

}

