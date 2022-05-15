FROM migrate/migrate:latest
RUN apk update && apk add bash
COPY wait-for-it.sh .
RUN chmod +x wait-for-it.sh