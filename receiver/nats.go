/*
 * skogul, nats receiver/subscriber
 *
 * Author(s):
 *  - Niklas Holmstedt <n.holmstedt@gmail.com>
 *
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2.1 of the License, or (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public
 * License along with this library; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA
 * 02110-1301  USA
 */

package receiver

import (
	"fmt"
	"os"
	"github.com/nats-io/nats.go"
	"github.com/telenornms/skogul"
	"crypto/tls"
	"crypto/x509"
)

var natsLog = skogul.Logger("receiver", "nats")
/*
Nats.io simple pub/sub receiver with:

- Authentication: Username/Password, TLS
- Authorization: Username/Password, UserCredentials/JWT
- Queue: Load balancing for multiple receivers in the same Queue.

*/
type Nats struct {
	Handler		skogul.HandlerRef `doc:"Handler used to parse, transform and send data."`
	Servers		string		  `doc:"Comma separated list of nats URLs"`
	Subject		string		  `doc:"Subject to subscribe to messages on"`
	Queue		string		  `doc:"Worker queue to distribute messages on"`
	Name		string		  `doc:"Client name"`
	Username	string		  `doc:"Client username"`
	Password	string		  `doc:"Client password"`
	TLSClientKey    string		  `doc:"TLS client key file path"`
	TLSClientCert   string		  `doc:"TLS client cert file path"`
	TLSCACert	string		  `doc:"CA cert file path"`
	UserCreds       string		  `doc:"Nats credentials file path"`
	NKeyFile        string		  `doc:"Nats nkey file path"`
	Insecure	bool		  `doc:"TLS InsecureSkipVerify"`
	o		*[]nats.Option
	nc		*nats.Conn
}

// Verify configuration
func (n *Nats) Verify() error {
	if n.Handler.Name == "" {
		return skogul.MissingArgument("Handler")
	}
	if n.Servers == "" {
		return skogul.MissingArgument("Address")
	}
        if n.Subject == "" {
                return skogul.MissingArgument("Subject")
        }

	return nil
}

//Unsure on how "Self contained" the receiver should be.
func getCertPool(path string) (*x509.CertPool, error) {
	// this means "use system default"
	if path == "" {
		return nil, nil
	}
	cp := x509.NewCertPool()
	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open custom root CA: %w", err)
	}
	defer func() {
		fd.Close()
	}()
	bytes := make([]byte, 1024000)
	n, err := fd.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to read custom root CA: %w", err)
	}
	ok := cp.AppendCertsFromPEM(bytes[:n])
	if !ok {
		return nil, fmt.Errorf("unable to append certificate to root CA pool")
	}
	return cp, nil
}

func (n *Nats) Start() error {
	if n.Name == "" {
		n.Name = "skogul"
	}
	n.o = &[]nats.Option{nats.Name(n.Name)}

	if n.Servers == "" {
		n.Servers = nats.DefaultURL
	}

	//User Credentials
	if n.UserCreds != "" && n.NKeyFile != "" {
		//Cred file contains nkey.
		natsLog.Fatal("Please configure usercreds or nkeyfile.")
	}
	if n.UserCreds != "" {
		*n.o = append(*n.o, nats.UserCredentials(n.UserCreds))
	}

	//Plain text passwords
	if n.Username != "" && n.Password != "" {
		if n.TLSClientKey != "" {
			natsLog.Warnf("Using plain text password over a non encrypted transport!")
		}
		*n.o = append(*n.o, nats.UserInfo(n.Username, n.Password))
	}
	//TLS authentication, Note: Fix selfsigned certificates.
	if n.TLSClientKey != "" && n.TLSClientCert != "" {
		cert, err := tls.LoadX509KeyPair(n.TLSClientCert, n.TLSClientKey)
		if err != nil {
			natsLog.Fatalf("error parsing X509 certificate/key pair: %v", err)
			return err
		}

		cp, err := getCertPool(n.TLSCACert)
                if err != nil {
                        natsLog.Fatalf("Failed to initialize root CA pool")
			return err
                }

		config := &tls.Config{
			InsecureSkipVerify:	n.Insecure,
			Certificates:		[]tls.Certificate{cert},
			RootCAs:		cp,
		}
		*n.o = append(*n.o, nats.Secure(config))
	}

	//NKey auth
	if n.NKeyFile != "" {
		opt, err := nats.NkeyOptionFromSeed(n.NKeyFile)
		if err != nil {
			natsLog.Fatal(err)
		}
		*n.o = append(*n.o, opt)
	}

	var err error
	n.nc, err = nats.Connect(n.Servers, *n.o...)
	cb := func(msg *nats.Msg) {
		if err:= n.Handler.H.Handle(msg.Data); err != nil {
			natsLog.WithError(err).Warn("Unable to handle Nats message")
		}
		return
	}
	if err != nil {
		natsLog.Errorf("Encountered an error while connecting to Nats: %w", err)
	}

	if n.Queue == "" {
		n.nc.QueueSubscribe(n.Subject, n.Queue, cb)
	} else {
		n.nc.Subscribe(n.Subject, cb)
	}
	return n.nc.LastError()
}
