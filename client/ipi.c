
#include <stdio.h>
#include <stdlib.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <netinet/udp.h>
#include <sys/types.h>
#include <netdb.h>
#include <string.h>

#define RECV_BUF_SIZE 80
#define DEFAULT_SERVER "plebis.net"

int main(int argc, char **argv)
{
    int udp_socket, result;
    ssize_t received;
    struct addrinfo addr_hints;
    struct addrinfo *addr_result;
    char recv_buf[RECV_BUF_SIZE];
    char *server_ip = DEFAULT_SERVER;

    memset(&addr_hints, 0, sizeof(struct addrinfo));
    addr_hints.ai_family = AF_INET;
    addr_hints.ai_socktype = SOCK_DGRAM;

    if (argc == 2) {
        server_ip = argv[1];
    }
    result = getaddrinfo(server_ip, "7777", &addr_hints, &addr_result);
    if (result != 0) {
        fprintf(stderr, "getaddrinfo: %s\n", gai_strerror(result));
        exit(EXIT_FAILURE);
    }

    udp_socket = socket(AF_INET, SOCK_DGRAM, 0);
    if (udp_socket == -1) {
        perror("socket");
        exit(EXIT_FAILURE);
    }

    result = sendto(
                 udp_socket, "", 0, 0, addr_result->ai_addr, addr_result->ai_addrlen);
    if (result == -1) {
        perror("sendto");
        exit(EXIT_FAILURE);
    }

    received = recv(udp_socket, recv_buf, RECV_BUF_SIZE, 0);
    if (received == -1) {
        perror("recvfrom");
        exit(EXIT_FAILURE);
    }

    recv_buf[received] = '\0';
    printf("%s", recv_buf);	// Response includes \n
    return EXIT_SUCCESS;
}
