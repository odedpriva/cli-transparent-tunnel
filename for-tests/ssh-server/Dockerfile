FROM ghcr.io/linuxserver/openssh-server
# Enable tunnelling or else we cannot connect to MongoDB via SSH:
# - Fixes observed error: "Could not establish connection from local ('127.0.0.1', 14988) to remote ('localhost', 27017) side of the tunnel: open new channel ssh error: ChannelException(1, 'Administratively prohibited')"
# - Ref: https://github.com/linuxserver/docker-openssh-server/issues/43
# - Ref: https://unix.stackexchange.com/questions/14160/ssh-tunneling-error-channel-1-open-failed-administratively-prohibited-open
RUN sed -i 's@#PermitTunnel no@PermitTunnel yes@' /etc/ssh/sshd_config
RUN sed -i 's@AllowTcpForwarding no@AllowTcpForwarding yes@' /etc/ssh/sshd_config