FROM openjdk:8-jre
COPY target/gauge-controller-0.0.1-SNAPSHOT.jar /opt/gauge-controller-0.0.1-SNAPSHOT.jar
COPY gauges.json /opt/gauges.json
CMD ["java","-jar","/opt/gauge-controller-0.0.1-SNAPSHOT.jar", "--spring.config.location=/opt/conf/gauge-controller.properties"]
