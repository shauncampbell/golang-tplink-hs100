package connector_test

import (
	"net"
	"time"

	c "github.com/jaedle/golang-tplink-hs100/internal/connector"
	"github.com/jaedle/golang-tplink-hs100/internal/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Connector", func() {

	const aRequest = `{"expected": "command"}}`
	const aResponse = `{"response": "expected"}}`

	It("sends command", func() {
		l := startServer()
		defer l.Close()
		request := make(chan []byte)
		go handleRequest(l, request, crypto.EncryptWithHeader(aResponse))

		_, err := localHostDevice().SendCommand(command(aRequest))
		requestPayload := awaitRequest(request)

		Expect(err).NotTo(HaveOccurred())
		Expect(requestPayload).To(Equal(crypto.EncryptWithHeader(aRequest)))
	})

	It("sends command and receives response", func() {
		l := startServer()
		defer l.Close()
		request := make(chan []byte)
		go handleRequest(l, request, crypto.EncryptWithHeader(aResponse))

		resp, err := localHostDevice().SendCommand(command(aRequest))
		awaitRequest(request)

		Expect(err).NotTo(HaveOccurred())
		Expect(resp).To(Equal(aResponse))
	})

	It("fails if cannot connect", func() {
		dev := localHostDevice()
		response, err := dev.SendCommand(command(""))
		Expect(err).To(HaveOccurred())
		Expect(response).To(BeEmpty())
	})
})

func awaitRequest(requestChan chan []byte) []byte {
	const t = 1 * time.Second

	var request []byte
	select {
	case request = <-requestChan:
		break
	case <-time.After(t):
		Fail("received no return value")
	}

	return request
}

func startServer() net.Listener {
	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		Fail("could not start test server")
	}
	return l
}

func handleRequest(l net.Listener, request chan []byte, response []byte) {
	conn, err := l.Accept()
	if err != nil {
		Fail("can not start server")
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		Fail("can not read request")
	}
	received := buf[:n]
	_, _ = conn.Write(response)
	_ = conn.Close()
	request <- received
}

func command(cmd string) c.Command {
	return c.NewCommand(cmd)
}

func localHostDevice() c.Device {
	return c.NewDevice("127.0.0.1")
}
