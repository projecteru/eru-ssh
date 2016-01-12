FROM ubuntu
MAINTAINER CMGS <ilskdw@gmail.com>

RUN apt-get update && apt-get install -y ssh
RUN mkdir /var/run/sshd
RUN echo root:eRuSyStEm | chpasswd
RUN sed -i 's/PermitRootLogin without-password/PermitRootLogin yes/' /etc/ssh/sshd_config
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd
ENV NOTVISIBLE "in users profile"
RUN echo "export VISIBLE=now" >> /etc/profile
RUN mkdir -p /root/.ssh
ADD mimiron.pub /root/.ssh/authorized_keys
RUN chmod 600 /root/.ssh/authorized_keys
