FROM alpine

RUN apk --update upgrade && \
    apk --update add ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

ARG APPNAME
ENV PARENT_DIR="/go/src/github.com/ivanturianytsia/"${APPNAME}
ENV EXE="${PARENT_DIR}/${APPNAME}"

RUN mkdir -p ${PARENT_DIR}
RUN mkdir -p ${PARENT_DIR}/client/dist/

WORKDIR ${PARENT_DIR}

ADD ./bin/${APPNAME}_alpine $EXE
ADD ./client/dist ${PARENT_DIR}/client/dist

ENTRYPOINT $EXE
