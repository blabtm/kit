FROM gradle:8.14-jdk21-corretto AS build

WORKDIR /app

COPY gradle.properties .
COPY settings.gradle .
COPY gradle ./gradle
COPY buildSrc ./buildSrc
COPY libs/build.gradle ./libs/build.gradle
COPY libs/java ./libs/java
COPY svc/build.gradle ./svc/build.gradle
COPY --exclude=sim svc/phy/beam.em-es ./svc/phy/beam.em-es
COPY svc/phy/ccd ./svc/phy/ccd

RUN gradle --no-daemon :svc:phy:beam.em-es:shadowJar

FROM flink:2.1-java21
COPY --from=build /app/svc/phy/beam.em-es/build/libs/*.jar /opt/flink/usrlib/
