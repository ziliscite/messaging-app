package mongo

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/ziliscite/messaging-app/pkg/must"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
)

func New(connection string) *mongo.Client {
	rootCAs := x509.NewCertPool()
	cert, err := ioutil.ReadFile("ca.pem") // Path to the downloaded certificate
	if err != nil {
		log.Fatalf("Failed to read root certificate: %v", err)
	}
	if ok := rootCAs.AppendCertsFromPEM(cert); !ok {
		log.Fatalf("Failed to append root certificate")
	}

	// Create a TLS configuration using the root CA
	tlsConfig := &tls.Config{
		RootCAs: rootCAs,
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	return must.Must(mongo.Connect(context.Background(), options.Client().ApplyURI(connection).SetTLSConfig(tlsConfig).SetServerAPIOptions(serverAPI)))
}
