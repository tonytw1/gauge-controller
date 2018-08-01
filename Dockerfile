FROM openjdk:10-jre
COPY target/gauge-controller-0.0.1-SNAPSHOT.jar /opt/gauge-controller-0.0.1-SNAPSHOT.jar
COPY gauges.json /opt/gauges.json
CMD ["java","-XshowSettings:vm", "-XX:+PrintCommandLineFlags", "-jar","/opt/gauge-controller-0.0.1-SNAPSHOT.jar", "--spring.config.location=/opt/gauge-controller/conf/gauge-controller.properties"]
