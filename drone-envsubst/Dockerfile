FROM alpine
LABEL maintainer="Florian Boehmak <florian@fnbk.cc>"

RUN \
  apk add --update ca-certificates && \
  apk add -t deps curl && \
  apk add bash && \
  apk add --update libintl && \
  apk add --virtual build_deps gettext &&  \
  cp /usr/bin/envsubst /usr/local/bin/envsubst && \
  apk del build_deps

COPY run.sh /run.sh
RUN chmod +x /run.sh

CMD ["/run.sh"]