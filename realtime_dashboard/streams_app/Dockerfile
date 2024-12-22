# --------------------------------
# 1) Build Stage
# --------------------------------
FROM maven:latest AS build
WORKDIR /app

# Copy your entire Maven project into the container
COPY pom.xml .
COPY src ./src

# Build the project (create the fat/uber JAR)
RUN mvn clean package -DskipTests

# --------------------------------
# 2) Runtime Stage
# --------------------------------
FROM amazoncorretto:17
WORKDIR /app

# Copy the built JAR from the build stage into the runtime stage
COPY --from=build /app/target/store-aggregation-app.jar /app/app.jar

# Expose a port if your app has a web/metrics endpoint (optional)
# EXPOSE 8080

# By default, just run the Java application
CMD ["java", "-jar", "/app/app.jar"]
