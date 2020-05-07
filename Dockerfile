FROM node:14-alpine

ADD scripts/install.sh /scripts/
ADD default.config.js /semantic-release/
ADD release/linux/amd64/semantic-release-plugin /bin/

RUN sh /scripts/install.sh

CMD ["/bin/semantic-release-plugin"]
