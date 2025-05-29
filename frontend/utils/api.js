// API配置
const BASE_URL = 'http://localhost:12000/api'

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
	products: {
		list: () => request('/products'),
		get: (id) => request(`/products/${id}`),
		create: (data) => request('/products', { method: 'POST', data }),
		update: (id, data) => request(`/products/${id}`, { method: 'PUT', data }),
		delete: (id) => request(`/products/${id}`, { method: 'DELETE' }),
		categories: () => request('/products/categories')
	},
	
	// 客户管理
	customers: {
		list: () => request('/customers'),
		get: (id) => request(`/customers/${id}`),
		create: (data) => request('/customers', { method: 'POST', data }),
		update: (id, data) => request(`/customers/${id}`, { method: 'PUT', data }),
		delete: (id) => request(`/customers/${id}`, { method: 'DELETE' })
	},
	
	// 订单管理
	orders: {
		list: () => request('/orders'),
		get: (id) => request(`/orders/${id}`),
		create: (data) => request('/orders', { method: 'POST', data }),
		update: (id, data) => request(`/orders/${id}`, { method: 'PUT', data }),
		delete: (id) => request(`/orders/${id}`, { method: 'DELETE' }),
		calendar: (year, month) => request(`/orders/calendar?year=${year}&month=${month}`)
	},
	
	// 驾驶仓
	dashboard: {
		stats: () => request('/dashboard')
	}
}

export default api