package cn.cityworks

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication

@SpringBootApplication
class Launch

fun main(args: Array<String>) {
    runApplication<Launch>(*args)
}
