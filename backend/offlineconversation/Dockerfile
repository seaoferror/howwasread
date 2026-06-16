FROM amazoncorretto:25-alpine-jdk as stage1
WORKDIR /app
COPY gradle gradle
COPY src src
COPY build.gradle .
COPY settings.gradle .
COPY gradlew .
RUN chmod +x gradlew
RUN ./gradlew bootJar

FROM gcr.io/distroless/java25-debian13
WORKDIR /app
COPY --from=stage1 /app/build/libs/*.jar app.jar

ENTRYPOINT ["java", "-jar", "app.jar"]