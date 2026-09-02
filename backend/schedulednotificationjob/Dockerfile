# Stage 1: Build the fat JAR
FROM eclipse-temurin:21-jdk AS builder

WORKDIR /app
COPY gradlew settings.gradle build.gradle ./
COPY gradle/ gradle/
RUN chmod +x gradlew
# Cache dependencies
RUN ./gradlew dependencies --no-daemon || true

COPY src/ src/
RUN ./gradlew shadowJar --no-daemon

# Stage 2: Add the JAR to the official Flink image
FROM flink:2.2.1-java21

COPY --from=builder /app/build/libs/schedulednotificationjob-1.0-SNAPSHOT.jar /opt/flink/usrlib/schedulednotificationjob.jar
