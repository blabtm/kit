FROM gradle:jdk21-corretto AS build

WORKDIR /app
COPY --exclude=svc \
    --exclude=platform/service-manager \
    --exclude=platform/grafana . .

COPY platform/nifi platform/nifi
COPY libs/build.gradle libs/build.gradle
COPY libs/java/build.gradle libs/java/build.gradle
COPY libs/java/vcas libs/java/vcas

RUN gradle -p platform/nifi nar --no-daemon

FROM apache/nifi:latest
COPY --from=build /app/platform/nifi/nifi-*-bundle/nifi-*-nar/build/libs/*.nar /opt/nifi/nifi-current/lib/

# Docker PostgreSQL driver
RUN curl -o /opt/nifi/nifi-current/lib/postgresql-42.7.8.jar https://jdbc.postgresql.org/download/postgresql-42.7.8.jar
