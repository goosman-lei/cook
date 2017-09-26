package mserver

import (
	"bytes"
	cook_io "cook/io"
	cook_log "cook/log"
	cook_util "cook/util"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"sync"
	"testing"
	"time"
)

type TestMessage struct {
	text     string
	seqnum   uint32
	clientNo uint32
}

var message_status_table *cook_util.CMap

var message_status_table_lock sync.Mutex

var gt *testing.T

func update_message_status(clientNo int, seqnum, status int) {
	message_status_table_lock.Lock()

	clientNoStr := strconv.FormatInt(int64(clientNo), 10)
	seqnumStr := strconv.FormatInt(int64(seqnum), 10)
	clientMessageMap := message_status_table.MustGet(clientNoStr).(*cook_util.CMap)

	oStatus, exists := clientMessageMap.Get(seqnumStr)
	if !exists {
		gt.Logf("Wrong seqnum: %d.%d", clientNo, seqnum)
		gt.Fail()
	}

	if oStatus.(int) != status-1 {
		gt.Logf("Wrong status: %d => %d. seqnum = %d.%d", oStatus.(int), status, clientNo, seqnum)
		gt.Fail()
	}

	clientMessageMap.Set(seqnumStr, status)
	message_status_table_lock.Unlock()
}

func TestCase(t *testing.T) {
	bw := new(bytes.Buffer)
	cook_log.Ldebug.SetOutput(bw)
	cook_log.Linfo.SetOutput(bw)
	cook_log.Lwarn.SetOutput(bw)
	cook_log.Lfatal.SetOutput(bw)
	go func() {
		log.Println(http.ListenAndServe("10.10.200.12:6060", nil))
	}()

	nClient := 10
	gt = t

	message_status_table = cook_util.NewCMap()
	for i := 0; i < nClient; i++ {
		cmap := cook_util.NewCMap()
		cmap.Copy(map[string]interface{}{
			"0": 0,
			"1": 0,
		})
		iStr := strconv.FormatInt(int64(i), 10)
		message_status_table.Set(iStr, cmap)
	}

	serverDone := make(chan bool)
	clientDone := make(chan bool, nClient)
	go CaseServer(serverDone, clientDone, nClient, t)
	for i := 0; i < nClient; i++ {
		go CaseClient(clientDone, uint32(i), t)
	}

	<-serverDone
}

func CaseServer(serverDone chan bool, clientDone chan bool, nClient int, t *testing.T) {
	cfg := MSvrConf{
		SN:          "TestServer",
		ListenAddrs: []string{":8801", ":8802", ":8803", ":8804"},

		ClientMsgQueueSize: 64,
	}
	ops := MSvrOps{
		ServerStartup:  ops_server_startup,
		ServerShutdown: ops_server_shutdown,

		ClientOnline:  ops_client_online,
		ClientOffline: ops_client_offline,

		SendMsg: ops_send_msg,
		RecvMsg: ops_recv_msg,
	}

	s := NewMServer(cfg, ops)

	s.Startup()

	// waiting test client done
	for i := 0; i < nClient; i++ {
		<-clientDone
	}

	s.Shutdown("normal")

	serverDone <- true
}

func CaseClient(clientDone chan bool, clientNo uint32, t *testing.T) {
	time.Sleep(1e7)

	addrs := []string{"127.0.0.1:8801", "127.0.0.1:8802", "127.0.0.1:8803", "127.0.0.1:8804"}
	addr := addrs[rand.Intn(len(addrs))]
	cook_log.Infof("client [%d] connect server addr is: %s", int(clientNo), addr)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Logf("client [%d] connect to server failed", clientNo)
		t.Fail()
		return
	}

	sendDone := make(chan bool)
	recvDone := make(chan bool)

	clientNoStr := strconv.FormatInt(int64(clientNo), 10)
	clientMessageMap := message_status_table.MustGet(clientNoStr).(*cook_util.CMap)
	message_ids := make([]int, 0, clientMessageMap.Len())
	clientMessageMap.Iterate(func(k string, v interface{}) {
		clientNoInt, _ := strconv.ParseInt(k, 10, 64)
		message_ids = append(message_ids, int(clientNoInt))
	})

	go func() {
		for _, id := range message_ids {
			m := &TestMessage{
				text:     "client send to server",
				seqnum:   uint32(id),
				clientNo: clientNo,
			}

			text := "client send to server"
			buffer := make([]byte, len(m.text)+12)
			binary.BigEndian.PutUint32(buffer[0:4], uint32(id))
			binary.BigEndian.PutUint32(buffer[4:8], uint32(clientNo))
			binary.BigEndian.PutUint32(buffer[8:12], uint32(len(text)))
			copy(buffer[12:], []byte(text))

			update_message_status(int(clientNo), id, 1)
			if _, err := conn.Write(buffer); err != nil {
				t.Logf("send msg to server failed: %s", err)
				t.Fail()
			}
			cook_log.Debugf("client [%d] send msg: seqnum = %d.%d text = %s status = 1", clientNo, clientNo, id, text)
		}
		cook_log.Debugf("client [%d] send all message done", clientNo)
		sendDone <- true
	}()

	go func() {
		for i := 0; i < len(message_ids); i++ {
			headBuffer := make([]byte, 12)
			_, err := conn.Read(headBuffer)
			if err != nil {
				t.Logf("recv msg header from server failed: %s", err)
				t.Fail()
			}

			seqNum := binary.BigEndian.Uint32(headBuffer[0:4])
			clientNo := binary.BigEndian.Uint32(headBuffer[4:8])
			textLen := int(binary.BigEndian.Uint32(headBuffer[8:12]))

			textBuffer := make([]byte, textLen)
			_, err = conn.Read(textBuffer)
			if err != nil {
				t.Logf("recv msg text from server failed: %s", err)
				t.Fail()
			}
			m := &TestMessage{seqnum: seqNum, text: string(textBuffer), clientNo: clientNo}

			update_message_status(int(clientNo), int(m.seqnum), 4)
			cook_log.Debugf("client [%d] recv msg: seqnum = %d.%d text = %s. status = 5", clientNo, m.clientNo, m.seqnum, m.text)
		}
		cook_log.Debugf("client [%d] recv all message done", clientNo)
		<-sendDone
		recvDone <- true
	}()
	<-recvDone

	conn.Close()

	cook_log.Debugf("test client [%d] done", clientNo)
	clientDone <- true
}

func ops_server_startup(s *MServer) error {
	s.Debugf("ops_server_startup")
	return nil
}

func ops_server_shutdown(s *MServer) error {
	s.Debugf("ops_server_shutdown")
	return nil
}

func ops_client_online(c *MClient) error {
	c.Debugf("ops_client_online")
	return nil
}

func ops_client_offline(c *MClient) error {
	c.Debugf("ops_client_offline")
	return nil
}

func ops_send_msg(c *MClient, msg interface{}) error {
	m, ok := msg.(*TestMessage)
	if !ok {
		return fmt.Errorf("error message")
	}

	c.Debugf("send message: seqnum = %d.%d text = %s. status = 4", m.clientNo, m.seqnum, m.text)
	update_message_status(int(m.clientNo), int(m.seqnum), 3)

	buffer := make([]byte, len(m.text)+12)
	binary.BigEndian.PutUint32(buffer[0:4], m.seqnum)
	binary.BigEndian.PutUint32(buffer[4:8], m.clientNo)
	binary.BigEndian.PutUint32(buffer[8:12], uint32(len(m.text)))
	copy(buffer[12:], []byte(m.text))

	return cook_io.WriteN(c.Conn, buffer)
}

func ops_recv_msg(c *MClient) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	headBuffer, err := cook_io.ReadN_Net(c.Conn, 12, time.Millisecond*10)
	if err != nil {
		if err != io.EOF {
			c.Warnf("read message header error: %s", err)
		}
		return err
	}

	seqNum := binary.BigEndian.Uint32(headBuffer[0:4])
	clientNo := binary.BigEndian.Uint32(headBuffer[4:8])
	textLen := int(binary.BigEndian.Uint32(headBuffer[8:12]))

	text, err := cook_io.ReadN_Net(c.Conn, textLen, time.Millisecond*10)
	if err != nil {
		if err != io.EOF {
			c.Warnf("read message body error: %s", err)
		}
		return err
	}

	c.Debugf("recv message: seqnum = %d.%d, text = %s. status = 2", int(clientNo), int(seqNum), text)
	update_message_status(int(clientNo), int(seqNum), 2)

	c.Send(&TestMessage{
		text:     "server reply to client",
		seqnum:   seqNum,
		clientNo: clientNo,
	})

	return nil
}
