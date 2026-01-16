# TLS Certificates for Nginx (Development / Self-Signed)

In production, TLS certificates are typically issued by a trusted Certificate Authority (for example, via Let’s Encrypt).
For local development and demo environments, a self-signed certificate can be used.

## Generate a self-signed certificate (Linux/macOS)
From the repository root:

```bash
mkdir -p deploy/shagram/nginx/certs

openssl req -x509 -newkey rsa:4096 \
  -keyout deploy/shagram/nginx/certs/key.pem \
  -out deploy/shagram/nginx/certs/cert.pem \
  -sha256 -days 365 -nodes \
  -subj "/C=RU/ST=Moscow/L=Korolyov/O=Shagram/CN=localhost"
```

## Verify
```bash
openssl x509 -in deploy/shagram/nginx/certs/cert.pem -noout -text | head
```

## Usage
The Nginx configuration expects:
- `deploy/shagram/nginx/certs/cert.pem`
- `deploy/shagram/nginx/certs/key.pem`

Start the deployment:

```bash
cd deploy/shagram
docker compose up -d
```
