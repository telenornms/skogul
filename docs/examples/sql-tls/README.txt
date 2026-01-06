SQL TLS Examples
================

Configuration options for SQL sender/receiver TLS:

  CAFile:   CA certificate for server verification
  CertFile: Client certificate for mutual TLS (requires KeyFile)
  KeyFile:  Client private key for mutual TLS (requires CertFile)
  Insecure: Skip server certificate verification (testing only)

Valid combinations:
  1. No TLS fields           - Plain connection
  2. CAFile only             - Verify server certificate
  3. CAFile + Cert + Key     - Full mutual TLS (recommended)
  4. Cert + Key + Insecure   - Client auth without server verification
  5. Insecure only           - TLS without verification

See tls-example.json for MySQL and PostgreSQL configuration examples.
