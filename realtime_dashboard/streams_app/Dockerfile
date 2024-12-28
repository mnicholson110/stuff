FROM maven:3.8.1-openjdk-11-slim AS build
WORKDIR /app

COPY pom.xml .
COPY src ./src

RUN mvn clean package -Dskiptests

FROM openjdk:11-jre-slim-buster
WORKDIR /app

COPY --from=build /app/target/store-aggregation-app.jar /app/app.jar
CMD ["java", "-jar", "/app/app.jar"]
