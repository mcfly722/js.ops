package event

import (
	context "context"
	"crypto/tls"
	"crypto/x509"
	fmt "fmt"
	"log"
	"net"
	"syscall"

	"github.com/dop251/goja"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type powershellTask struct {
	callback   goja.Callable
	finished   bool
	scriptBody string
	response   string
}

func (current *powershellTask) HasFinished(vm *goja.Runtime) bool {
	if current.finished {
		current.callback(nil, vm.ToValue(current.response))
	}
	return current.finished
}

var (
	modcrypt32                          = syscall.NewLazyDLL("crypt32.dll")
	procCryptSignMessage                = modcrypt32.NewProc("CryptSignMessage")
	procCertDuplicateCertificateContext = modcrypt32.NewProc("CertDuplicateCertificateContext")
)

type certificate struct {
	certContext uintptr
	*x509.Certificate
}

type custom struct {
	credentials.TransportCredentials
}

func (c *custom) serverHandshake(conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	conn, authInfo, err := c.TransportCredentials.ServerHandshake(conn)
	tlsInfo := authInfo.(credentials.TLSInfo)
	name := tlsInfo.State.PeerCertificates[0].Subject.CommonName
	fmt.Printf("%s\n", name)
	return conn, authInfo, err
}

func verifyPeerCertificate(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	if len(rawCerts) != 1 {
		return fmt.Errorf("Expected 1 certificate, but got %d", len(rawCerts))
	}

	cert, err := x509.ParseCertificate(rawCerts[0])

	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("%s%s", "verifyPeerCertificate SUBJECT=", cert.Subject))
	return nil
}

// NewPowershellTask cunstructor
func NewPowershellTask(call goja.Callable, scriptBody string) Task {

	log.Println("GRPC version ", grpc.Version)

	cert, err := tls.LoadX509KeyPair("..\\testServer.pem", "..\\testServer.pem")
	if err != nil {
		log.Fatal(fmt.Sprintf("%s.%s", "LoadX509KeyPair", err))
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify:    true,
		Certificates:          []tls.Certificate{cert},
		VerifyPeerCertificate: verifyPeerCertificate,
	}

	conn, err := grpc.Dial("localhost:5001", grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))

	if err != nil {
		log.Fatal(fmt.Sprintf("%s.%s", "Dial", err))
	}
	log.Println(fmt.Sprintf("connection state=%s", conn.GetState()))

	client := NewGreeterClient(conn)

	res, err := client.SayHello(context.Background(), &HelloRequest{Name: "test Hello Request"})
	if err != nil {
		log.Fatal(fmt.Sprintf("%s.%s", "SayHello", err))
	}

	task := &powershellTask{
		callback:   call,
		finished:   true,
		scriptBody: scriptBody,
		response:   res.GetMessage(),
	}

	return task
}
