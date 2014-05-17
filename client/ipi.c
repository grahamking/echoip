/*
 * IP Inquiry. Client for echoip. Displays your remote IP address and location.
 *
 * Usage: ipi [echoip-server]
 * Default echoip-server is plebis.net
 *
 * Build it: make ipi
 *
 * https://github.com/grahamking/echoip
*/

#include <stdio.h>			// printf, fprintf, perror
#include <netdb.h>			// getaddrinfo, gai_strerror
#include <string.h>			// memset

#define RECV_BUF_SIZE 80
#define DEFAULT_SERVER "plebis.net"
#define EXIT_FAILURE 1
#define EXIT_SUCCESS 0

// Send a single empty UDP packet to the server, and display contents of
// first response packet.
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
        return EXIT_FAILURE;
    }

    udp_socket = socket(AF_INET, SOCK_DGRAM, 0);
    if (udp_socket == -1) {
        perror("socket");
        return EXIT_FAILURE;
    }

    result = sendto(
                 udp_socket, "", 0, 0, addr_result->ai_addr, addr_result->ai_addrlen);
    if (result == -1) {
        perror("sendto");
        return EXIT_FAILURE;
    }

    received = recv(udp_socket, recv_buf, RECV_BUF_SIZE, 0);
    if (received == -1) {
        perror("recv");
        return EXIT_FAILURE;
    }

    recv_buf[received] = '\0';
    printf("%s", recv_buf);	// Response includes \n
    return EXIT_SUCCESS;
}
