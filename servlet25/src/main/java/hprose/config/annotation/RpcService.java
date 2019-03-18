package hprose.config.annotation;

import org.springframework.stereotype.Component;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

/**
 * RPC 服务注解（标注在服务实现类上）
 *参考 https://gitee.com/huangyong/rpc.git
 * @author huangyong
 * @since 1.0.0
 */
@Target({ElementType.TYPE})
@Retention(RetentionPolicy.RUNTIME)
@Component
public @interface RpcService {

    /**
     * 服务接口类
     */
    Class<?> value();

    /**
     * 服务版本号
     */
    String version() default "0.0.1";

}
