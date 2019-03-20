package register

/**
 * 服务注册接口
 *
 * @author lituquan
 * @since 1.0.0
 */
type ServiceRegistry interface {

	/**
	 * 注册服务名称与服务地址
	 *
	 * @param serviceName    服务名称
	 * @param serviceAddress 服务地址
	 */
	Register(serviceName, serviceAddress string)

	/**
	 * 根据服务名称查找服务地址
	 *
	 * @param serviceName 服务名称
	 * @return 服务地址
	 */
	Discover(name string) (string, error)
}
