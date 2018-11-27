from docker
run apk add --no-cache iptables supervisor vim tmux
run ln -f -s vim /usr/bin/vi
volume /var/lib/docker
copy supervisord.conf /etc/
entrypoint ["/usr/bin/supervisord", "-n"]
copy .vimrc /root/
copy tmux.conf /root/.tmux.conf
env DOCKER_BUILDKIT=1