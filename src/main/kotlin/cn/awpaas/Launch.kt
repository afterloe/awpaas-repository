package cn.cityworks

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
import org.springframework.cloud.netflix.eureka.server.EnableEurekaServer
//import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity
import org.springframework.web.reactive.config.EnableWebFlux

@SpringBootApplication
@EnableEurekaServer
@EnableWebFlux
//@EnableWebSecurity
class Launch

fun main(args: Array<String>) {
    runApplication<Launch>(*args)
}
