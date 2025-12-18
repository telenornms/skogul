SQL TLS Configuration Examples
==============================

This directory contains example configurations for SQL senders and
receivers with TLS/certificate authentication.

Configuration Options
---------------------

CAFile:   Path to CA certificate for server verification
CertFile: Path to client certificate for mutual TLS authentication
KeyFile:  Path to client private key for mutual TLS authentication
Insecure: Skip certificate verification (for testing only)

Valid Combinations
------------------

1. No TLS fields                      - Plain connection (specify password in connstr)
2. CAFile only                        - Verify server certificate with custom CA
3. CAFile + CertFile + KeyFile        - Full mutual TLS (recommended)
4. CertFile + KeyFile + Insecure:true - Client cert auth without server verification
5. Insecure: true                     - TLS without certificate verification

Examples
--------

mysql.json    - MySQL with mutual TLS (client cert + server verification)
postgres.json - PostgreSQL with mutual TLS (client cert + server verification)

Each example includes both a SQL receiver (reads from database) and
a SQL sender (writes to database).

Notes
-----

MySQL:
- TLS is configured via mysql.RegisterTLSConfig() internally
- Connection string uses DSN format: user@tcp(host:port)/database
- Add ?parseTime=true to parse TIME/DATETIME columns

PostgreSQL:
- TLS options are appended to the connection string automatically
- sslmode is set based on configuration:
  - verify-full when CAFile is provided
  - require when Insecure is set
- You can alternatively specify SSL options directly in connstr

Both drivers:
- CertFile and KeyFile must be provided together or not at all
- Client certs require either CAFile (for server verification) or Insecure: true
- For password auth, include credentials in the connection string
