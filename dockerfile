FROM golang:1.19 as builder

WORKDIR /app

COPY dockerfiles .

CMD ["make", "build"]

FROM alpine:3.16

WORKDIR /app

COPY --from=builder /app/output/newsdumper /app/newsdumper

ENV OUTOUT_PATH="./"
ENV RUN_STR="@every 1h"
ENV CONFIG_PATH="./config/config.yaml"

RUN ["bash", "-c ", "mkdir -p $DUMP_PATH && ./newsdumper --output-dir $DUMP_PATH \
--cron-str $RUN_STR \
--config $CONFIG_PATH"]