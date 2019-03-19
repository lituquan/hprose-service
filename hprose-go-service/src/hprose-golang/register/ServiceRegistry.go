package register

/**
 * 服务注册接口
 *
 * @author huangyong
 * @since 1.0.0
 */
type ServiceRegistry interface  {

    /**
     * 注册服务名称与服务地址
     *
     * @param serviceName    服务名称
     * @param serviceAddress 服务地址
     */
     Register( serviceName,  serviceAddress string)
}