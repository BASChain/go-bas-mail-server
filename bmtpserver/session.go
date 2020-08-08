package bmtpserver

import (
	"errors"
	"github.com/realbmail/go-bas-mail-server/protocol"
	"github.com/realbmail/go-bas-mail-server/tools"
	"github.com/realbmail/go-bmail-protocol/bmp"
	"github.com/realbmail/go-bmail-protocol/translayer"
	"log"
	"net"
	"strconv"
)

type TcpSession struct {
	sn      []byte
	pubkey  []byte
	conn    *net.TCPConn
	Handle  BMTPFunc
	size    int
	server  *BMTPServerConf
	version int
	buf     []byte
	bmtl    *translayer.BMTransLayer
	rbody   protocol.RBody
	wbody   protocol.WBody
}

func (ts *TcpSession) Negotiation() error {
	if err := ts.readHead(); err != nil {
		return err
	}

	support := ts.server.VersionInSrv(int(ts.bmtl.GetVersion()))

	ack := &bmp.HELOACK{}

	if support == false {
		ack.ErrCode = 1
		ack.SupportVersion = ts.server.SupportVersion()
	} else {
		ack.ErrCode = 0
		ts.sn = tools.NewSn(tools.SerialNumberLength)
		copy(ack.SN[:], ts.sn)
		ack.SrvBca = ts.server.wallet.BCAddress()
	}

	ts.wbody = ack

	if err := ts.WriteMsg(); err != nil {
		return err
	}

	if support {
		//fmt.Println("version", ts.bmtl.GetVersion())
		ts.Handle = ts.server.SupportFunc[int(ts.bmtl.GetVersion())]
		return nil
	} else {
		return errors.New("client version is not support, ver: " + strconv.Itoa(int(ts.bmtl.GetVersion())))
	}

}

func (ts *TcpSession) WriteMsg() error {
	bmtl := translayer.NewBMTL(ts.wbody.MsgType())

	buf, err := ts.wbody.GetBytes()
	if err != nil {
		return err
	}

	bmtl.SetDataLen(uint32(len(buf)))

	log.Println(len(buf), bmtl.GetDataLen())

	data, _ := bmtl.Pack()

	var n int
	n, err = ts.conn.Write(data)
	if err != nil || n != len(data) {
		return errors.New("Write " + strconv.Itoa(int(ts.wbody.MsgType())) + " message head Failed")
	}

	n, err = ts.conn.Write(buf)
	if err != nil || n != len(buf) {
		return errors.New("Write " + strconv.Itoa(int(ts.wbody.MsgType())) + " message body Failed")
	}

	log.Println("write buf:", string(buf))

	return nil
}

func (ts *TcpSession) readNext() error {
	buf := make([]byte, ts.size)
	total := 0
	for {
		n, err := ts.conn.Read(buf[total:])
		if err != nil && total < ts.size {
			//if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			//	continue
			//} else if err != io.EOF {
			//	log.Println("read error", err)
			//	return err
			//}
			return err
		}
		total += n
		if total >= ts.size {
			break
		}

	}
	ts.buf = buf

	return nil
}

func (ts *TcpSession) deriveHead() (*translayer.BMTransLayer, error) {
	if len(ts.buf) != translayer.BMHeadSize() {
		return nil, errors.New("data is not a header buffer")
	}

	bmtl := &translayer.BMTransLayer{}

	_, err := bmtl.UnPack(ts.buf)
	if err != nil {
		return nil, err
	}

	return bmtl, nil

}

func (ts *TcpSession) deriveBody() error {
	ts.rbody = protocol.MsgGrid[ts.bmtl.GetMsgType()]
	log.Println(string(ts.buf))
	return ts.rbody.UnPack(ts.buf)
}

func (ts *TcpSession) readHead() error {
	ts.size = translayer.BMHeadSize()

	if err := ts.readNext(); err != nil {
		return err
	}

	if head, err := ts.deriveHead(); err != nil {
		return err
	} else {
		ts.bmtl = head
		ts.size = int(head.GetDataLen())
		return nil
	}
}

func (ts *TcpSession) readBody() error {
	if ts.bmtl == nil || ts.size == 0 {
		return errors.New("please read message header first")
	}

	if err := ts.readNext(); err != nil {
		return err
	} else {
		if err := ts.deriveBody(); err != nil {
			return err
		}
	}

	return nil
}

func HandleMsgV1(ts *TcpSession) error {
	if err := ts.readHead(); err != nil {
		return err
	}

	if err := ts.readBody(); err != nil {
		return err
	}
	ts.rbody.SetCurrentSn(ts.sn)

	if !ts.rbody.Verify() {
		return errors.New("error")
	}

	if resp, err := ts.rbody.Response(); err != nil {
		log.Println(err)
		ts.wbody = resp
		err = ts.WriteMsg()
		log.Println(err)
		return err
	} else {
		ts.wbody = resp
		if err := ts.WriteMsg(); err != nil {
			log.Println(err)
			return err
		}
	}
	if err := ts.rbody.Dispatch(); err != nil {
		//todo...
	}

	return nil

}
