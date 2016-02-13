FROM alpine
MAINTAINER Alex Ungur <alexaandru@gmail.com>

RUN apk update && apk add openssh openssh-sftp-server && cd /etc/ssh && ssh-keygen -A && adduser -D test && echo "test:1234"| chpasswd && echo -e "\nPubkeyAuthentication=no\n" >> /etc/ssh/sshd_config

CMD ["/usr/sbin/sshd", "-D", "-d"]
