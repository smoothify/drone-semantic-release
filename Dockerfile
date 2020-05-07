FROM node:14-alpine

ADD scripts/install.sh /scripts/
RUN sh /scripts/install.sh

ADD default.config.js /semantic-release/
ADD release/linux/amd64/semantic-release-plugin /bin/
ENV NODE_PATH=/usr/local/lib/node_modules

CMD ["/bin/semantic-release-plugin"]
