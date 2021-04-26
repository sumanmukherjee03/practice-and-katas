## TLS

This README is not specifically about kubernetes. It covers the generic topic of TLS.

Symetric encryption
  When the key used to encrypt the data from the sender is also passed along to the receiver
  Because it is the same key that is also used by the receiver to decrypt the encrypted data

Asymetric encryption
  You can think of the key as a private key used to unlock a lock. The lock can be thought of as a publicly available lock
  that anyone can access and use to lock something. The public key in asymetric encryption is analogous to the lock here.

  Asymetric encryption with ssh :
    A user can generate a public key (or what we are calling as public lock) and a private key with `ssh-keygen -t rsa -b 4096 -C client@example.com`.
    And the server admin can place the public key (or public lock) of the user on the server door in ~/.ssh/authorized_keys file.
    This lock can only be opened with the users private key.
    So, user generates lock and key. Gives the lock to the server admnin and server admin puts it on the servers door.

  Asymetric encryption with web traffic :
    With web traffic, in symetric encryption the key that encrypts the data is sent over the wire to the server from the user
    so that the server can decrypt the data. This is a concern because a hacker can sniff the key along with the encrypted data.
    But if the key can be shared from the client to the server with asymetric encryption then the client/server
    can continue talking with symetric encryption after that point onwards.

    Opposite to ssh where the client generates the key/lock pair and gives the lock to the server admin so that
    the server admin can place it on the servers entrypoint, the server instead generates the key/lock pair.
    On the server use `openssl genrsa -out server.key 1024` to generate the private key and `openssl rsa -in server.key -pubout > serverpub.pem` to generate the public lock.

    When the client first accesses the server, they get the public lock from the server.
    The client then encrypts the symetric key with the servers public lock and sends it back to the server.
    The server decrypts it with it's private key and gets the symetric key of the client.
    Then the client can continue talking to the server with symetric encryption with the servers.

    To verify the identity of the server, when the server sends it's public lock it sends it in a certificate.
    The certificate contains information about who it is issued to - ie name, alternate names etc, by whom it was issued, location of the server etc.
    The certificate's validity is verified by the browser based on the name/alernate names on the cert and by the name of the issuer.
    A certificate issued by a well-known certificate authority can be easily verified by the client's browser.

    A server can get it's certificate signed by a CA by generating a Certificate Signing Request (CSR) with `openssl req -new -key server.key -out server.csr -subj "/C=US/ST=CA/O=Example,Inc/CN=example.com"`
    Once the certificate authority validates the certificate it sends you back a signed certificate.

    To validate the certificate authority itself, the browsers themselves have the well known CA authorities public lock
    baked into them. The certificate from the server was signed using the CA's private key.
    The client browser can use the public lock of the CA to validate if the certificate from the server was signed by the CA itself as it says.
        This is a bit of a difference from the public lock/private key analogy. In reality, the public lock is also a key.
        You can encrypt data with one key and decrypt with the other one. So, you can encrypt with the private key
        And decrypt with the public key and vice versa. It all depends on what you are encrypting and what you are
        decrypting with. Only in case of the certs, are the certs encrypted by the CA's private key and validated
        with the CA's public key. But usually data is always encrypted with the public key and decrypted with the private key.
    If it turns out to be valid, the client continues with encypting the symetric key with the servers public lock.
    You can add the public lock of your private CA to the trusted public keys in your browser.

    The above is how a client validates a server. But the server cant know for sure that the client is who they say they are.
    This is where a client cert comes into picture. This is not something that's usually done though.
    However, the server can request a certificate from the client during the initial trust building.
    So, the client must generate a CSR and get a signed cert from a CA. This client cert is sent over to the server.
    The server already has the public keys necessary to verify the certificate of the client.

    Certs with public key are named `.crt` or `.pem`.
    Private keys are usually named as `.key` or `-key.pem`.
