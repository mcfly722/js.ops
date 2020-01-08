package event

import (
	context "context"
	"crypto/tls"
	"crypto/x509"
	fmt "fmt"
	"log"
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

// NewPowershellTask cunstructor
func NewPowershellTask(call goja.Callable, scriptBody string) Task {

	fmt.Println("GRPC version ", grpc.Version)

	cert, err := tls.LoadX509KeyPair("testServer.pem", "testServer.pem")
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}

	conn, err := grpc.Dial("localhost:5001", grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))

	if err != nil {
		log.Fatal(fmt.Sprintf("%s.%s", "Dial", err))
	}

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
