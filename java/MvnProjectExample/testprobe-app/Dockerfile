FROM maven:3.6.0-jdk-11-slim AS build
COPY src /home/app/src
COPY pom.xml /home/app
RUN mvn -f /home/app/pom.xml clean compile assembly:single

FROM amazoncorretto:11
COPY --from=build /home/app/target/*.jar app.jar
RUN mkdir -p /heapdump
ENV JAVA_TOOL_OPTIONS="-XX:+UseContainerSupport -XX:MaxRAMPercentage=60.0 -XX:+PrintCommandLineFlags -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=/heapdump"
# ENV JAVA_TOOL_OPTIONS="-XX:+PrintCommandLineFlags -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=/heapdump"
ENV JAVA_OPTS="-XshowSettings:vm -showversion"
RUN sh -c 'touch /app.jar'

EXPOSE 28500

ENTRYPOINT ["sh", "-c", "java $JAVA_OPTS -cp /app.jar com.testprobe.app.App"]
