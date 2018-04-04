package cn.cityworks

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
import org.springframework.cloud.netflix.eureka.server.EnableEurekaServer

@SpringBootApplication
@EnableEurekaServer
class Launch

fun main(args: Array<String>) {
    runApplication<Launch>(*args)
}
