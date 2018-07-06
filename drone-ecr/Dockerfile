FROM docker:17.12.0-ce-dind
LABEL maintainer="Florian Boehmak <florian@fnbk.cc>"

RUN \
	mkdir -p /aws && \
	apk -Uuv add groff less python py-pip && \
	pip install awscli && \
	apk --purge -v del py-pip && \
	rm /var/cache/apk/*

COPY run.sh /run.sh
RUN chmod +x /run.sh

CMD ["/usr/local/bin/dockerd-entrypoint.sh", "/run.sh"]
