import ssl
import socket

def check_ssl(hostname):
    context = ssl.create_default_context()
    with socket.create_connection((hostname, 443)) as sock:
        with context.wrap_socket(sock, server_hostname=hostname) as ssock:
            return ssock.version()

if __name__ == "__main__":
    import sys
    target = sys.argv[1]
    print(check_ssl(target))
