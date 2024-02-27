FROM maven:3.8.4-openjdk-17-slim AS build

WORKDIR /app

COPY pom.xml .

RUN mvn dependency:go-offline

COPY src ./src

RUN mvn package -DskipTests

FROM openjdk:17-jdk-slim

ARG JAR_FILE=/app/target/*.jar

COPY --from=build ${JAR_FILE} app.jar

ENTRYPOINT ["java","-jar","/app.jar"]