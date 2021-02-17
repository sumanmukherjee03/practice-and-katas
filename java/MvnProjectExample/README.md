### Development

Creating a barebones simple maven project
```
mvn archetype:generate -DgroupId=com.testprobe.app -DartifactId=testprobe-app -DarchetypeArtifactId=maven-archetype-quickstart -DarchetypeVersion=1.4 -DinteractiveMode=false
cd testprobe-app
mvn package
java -cp target/testprobe-app-1.0-SNAPSHOT.jar com.testprobe.app.App
```
