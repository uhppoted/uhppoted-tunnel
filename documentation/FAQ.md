# FAQ

1. How do I reduce the number of log messages?

   Set the log level to _warn_ or _error_:
   - in `console` mode, set the --log-level <debug|info|warn|error> command line argument, e.g.
```
   uhppoted-tunnel --console --log-level warn ...
```

   - in `service` mode, edit the service configuration to set the appropriate level with the `--log-level <debug|info|warn|error>` command line argument

2. In `service` mode how do redirect the log messages to _/dev/null_ ?

   Really not recommend because if anything goes wrong there will nothing whatsoever to work 
   with. But if you're dead set on making things difficult - edit the service command line to
   include the `--console` command line argument. All messages will then go to _stdout_ and 
   you can redirect them using the standard _redirect_ or _pipe_ operators.

3. How do I generate CA, server and client keys for use with the TLS host and client tunnels?

   There are many, many guides out there that describe how to do it properly, in varying degrees 
   of complexity from dead simple to insanely complicated. But for a quick and dirty set of keys
   for testing:

   _CA_
   ```
   openssl req -x509 -days 1825 -newkey rsa:4096 -nodes -keyout CA.key -out CA.cert
   ```
   
      _server_
   ```
   openssl req -new -newkey rsa:2048 -nodes -keyout server.key -out server.csr \
           -subj   '/C=US/ST=CA/L=San Diego/O=uhppoted/CN=127.0.0.1' \
           -addext "subjectAltName=IP:127.0.0.1"

   openssl x509 -req -in server.csr -days 365 -CA CA.cert -CAkey CA.key -CAcreateserial -out server.cert \
           -extfile <(printf "subjectAltName=IP:127.0.0.1")

   rm server.csr

   ```
   
      _client_
   ```
   openssl req -new -newkey rsa:2048 -nodes -keyout client.key -out client.csr \
           -subj   '/C=US/ST=CA/L=San Diego/O=uhppoted/CN=127.0.0.1' \
           -addext "subjectAltName=IP:127.0.0.1"

   openssl x509 -req -in client.csr -days 365 -CA CA.cert -CAkey CA.key -CAcreateserial -out client.cert \
           -extfile <(printf "subjectAltName=IP:127.0.0.1")

   rm client.csr
   ```

   References:
   - [OpenSSL CA](https://openssl-ca.readthedocs.io/en/latest/index.html)
   - [OpenSSL Essentials: Working with SSL Certificates, Private Keys and CSRs](https://www.digitalocean.comcommunity/tutorials/openssl-essentials-working-with-ssl-certificates-private-keys-and-csrs)
   - [How do you sign a Certificate Signing Request with your Certification Authority?](https://stackoverflow.com/questions/21297139/how-do-you-sign-a-certificate-signing-request-with-your-certification-authority)
