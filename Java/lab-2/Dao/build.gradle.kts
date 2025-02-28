plugins {
    id("java")
}

group = "com.stt0504"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

dependencies {
    testImplementation(platform("org.junit:junit-bom:5.9.1"))
    testImplementation("org.junit.jupiter:junit-jupiter")
    implementation("org.hibernate:hibernate-core:6.0.0.Final")
    compileOnly("org.projectlombok:lombok:1.18.20")
    annotationProcessor("org.projectlombok:lombok:1.18.20")
    implementation(project(":Entity"))
    implementation("org.postgresql:postgresql:42.2.24")
}

tasks.test {
    useJUnitPlatform()
}
