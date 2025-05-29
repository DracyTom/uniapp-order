// API配置
const BASE_URL = 'http://localhost:12001/api'

// 请求封装
const request = (url, options = {}) => {
	return new Promise((resolve, reject) => {
		uni.request({
			url: BASE_URL + url,
			method: options.method || 'GET',
			data: options.data || {},
			header: {
				'Content-Type': 'application/json',
				...options.header
			},
			success: (res) => {
				if (res.statusCode === 200) {
					resolve(res.data)
				} else {
					reject(res)
				}
			},
			fail: (err) => {
				reject(err)
			}
		})
	})
}

// API接口
export const api = {
	// 商品管理
	product: {
		list: (params) => request('/products', { method: 'GET', data: params }),
		get: (id) => request(`/products/${id}`),
		create: (data) => request('/products', { method: 'POST', data }),
		update: (id, data) => request(`/products/${id}`, { method: 'PUT', data }),
		delete: (id) => request(`/products/${id}`, { method: 'DELETE' }),
		categories: () => request('/products/categories')
	},
	
	// 兼容旧的API调用
	products: {
		list: (params) => request('/products', { method: 'GET', data: params }),
		get: (id) => request(`/products/${id}`),
		create: (data) => request('/products', { method: 'POST', data }),
		update: (id, data) => request(`/products/${id}`, { method: 'PUT', data }),
		delete: (id) => request(`/products/${id}`, { method: 'DELETE' }),
		categories: () => request('/products/categories')
	},
	
	// 客户管理
	customer: {
		list: (params) => request('/customers', { method: 'GET', data: params }),
		get: (id) => request(`/customers/${id}`),
		create: (data) => request('/customers', { method: 'POST', data }),
		update: (id, data) => request(`/customers/${id}`, { method: 'PUT', data }),
		delete: (id) => request(`/customers/${id}`, { method: 'DELETE' })
	},
	
	// 兼容旧的API调用
	customers: {
		list: (params) => request('/customers', { method: 'GET', data: params }),
		get: (id) => request(`/customers/${id}`),
		create: (data) => request('/customers', { method: 'POST', data }),
		update: (id, data) => request(`/customers/${id}`, { method: 'PUT', data }),
		delete: (id) => request(`/customers/${id}`, { method: 'DELETE' })
	},
	
	// 订单管理
	order: {
		list: (params) => request('/orders', { method: 'GET', data: params }),
		get: (id) => request(`/orders/${id}`),
		create: (data) => request('/orders', { method: 'POST', data }),
		update: (id, data) => request(`/orders/${id}`, { method: 'PUT', data }),
		delete: (id) => request(`/orders/${id}`, { method: 'DELETE' }),
		calendar: (year, month) => request(`/orders/calendar?year=${year}&month=${month}`)
	},
	
	// 兼容旧的API调用
	orders: {
		list: (params) => request('/orders', { method: 'GET', data: params }),
		get: (id) => request(`/orders/${id}`),
		create: (data) => request('/orders', { method: 'POST', data }),
		update: (id, data) => request(`/orders/${id}`, { method: 'PUT', data }),
		delete: (id) => request(`/orders/${id}`, { method: 'DELETE' }),
		calendar: (year, month) => request(`/orders/calendar?year=${year}&month=${month}`)
	},
	
	// 驾驶仓
	dashboard: {
		stats: () => request('/dashboard/stats'),
		hotProducts: () => request('/dashboard/hot-products')
	}
}

export default api