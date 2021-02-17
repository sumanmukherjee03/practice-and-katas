FROM amazoncorretto:11

RUN mkdir -p /heapdump
ENV LANG=C.UTF-8
ENV JAVA_TOOL_OPTIONS="-XX:+UseContainerSupport -XX:MaxRAMPercentage=50.0 -XX:+PrintCommandLineFlags -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=/heapdump"
ENV JAVA_OPTS="-XshowSettings:vm -showversion"
ADD target/*.jar app.jar
RUN sh -c 'touch /app.jar'
ENTRYPOINT ["sh", "-c", "java $JAVA_OPTS -cp /app.jar com.testprobe.app.App"]
