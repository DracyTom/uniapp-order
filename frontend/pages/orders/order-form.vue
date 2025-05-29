<template>
	<view class="container">
		<form @submit="submitOrder">
			<!-- 客户信息 -->
			<view class="card">
				<view class="card-title">客户信息</view>
				
				<view class="form-group">
					<text class="label">选择客户 *</text>
					<picker @change="onCustomerChange" :value="customerIndex" :range="customerOptions" range-key="name">
						<view class="picker-input">
							{{ customerOptions[customerIndex]?.name || '请选择客户' }}
						</view>
					</picker>
				</view>
				
				<view class="form-group">
					<text class="label">收货人 *</text>
					<input class="form-input" v-model="formData.recipient" placeholder="请输入收货人姓名" />
				</view>
				
				<view class="form-group">
					<text class="label">联系电话</text>
					<input class="form-input" v-model="formData.phone" placeholder="请输入联系电话" />
				</view>
			</view>
			
			<!-- 配送信息 -->
			<view class="card">
				<view class="card-title">配送信息</view>
				
				<view class="form-group">
					<text class="label">送货时间 *</text>
					<picker mode="date" @change="onDateChange" :value="deliveryDate">
						<view class="picker-input">{{ deliveryDate || '请选择送货日期' }}</view>
					</picker>
				</view>
				
				<view class="form-group">
					<text class="label">送货地址 *</text>
					<picker @change="onAddressChange" :value="addressIndex" :range="addressOptions" range-key="address">
						<view class="picker-input">
							{{ addressOptions[addressIndex]?.address || '请选择送货地址' }}
						</view>
					</picker>
				</view>
				
				<view class="form-group">
					<text class="label">特殊要求</text>
					<textarea class="form-textarea" v-model="formData.special_req" placeholder="请输入特殊要求（可选）"></textarea>
				</view>
			</view>
			
			<!-- 商品信息 -->
			<view class="card">
				<view class="card-title">
					<text>商品清单</text>
					<button class="add-product-btn" @click="addProduct" type="button">添加商品</button>
				</view>
				
				<view class="product-list">
					<view class="product-item" v-for="(item, index) in formData.order_items" :key="index">
						<view class="product-header">
							<picker @change="(e) => onProductChange(e, index)" :value="item.productIndex" :range="productOptions" range-key="name">
								<view class="product-picker">
									{{ productOptions[item.productIndex]?.name || '请选择商品' }}
								</view>
							</picker>
							<button class="remove-btn" @click="removeProduct(index)" type="button">删除</button>
						</view>
						
						<view class="product-details">
							<view class="quantity-input">
								<text class="input-label">数量:</text>
								<input type="number" v-model="item.quantity" @input="calculateAmount(index)" placeholder="0" />
							</view>
							<view class="price-display">
								<text class="input-label">单价:</text>
								<text class="price-text">¥{{ item.price || 0 }}</text>
							</view>
							<view class="amount-display">
								<text class="input-label">小计:</text>
								<text class="amount-text">¥{{ item.amount || 0 }}</text>
							</view>
						</view>
					</view>
				</view>
				
				<view v-if="formData.order_items.length === 0" class="empty-products">
					<text class="empty-text">请添加商品</text>
				</view>
			</view>
			
			<!-- 订单信息 -->
			<view class="card">
				<view class="card-title">订单信息</view>
				
				<view class="form-group">
					<text class="label">支付状态</text>
					<picker @change="onPaymentStatusChange" :value="paymentStatusIndex" :range="paymentStatusOptions">
						<view class="picker-input">{{ paymentStatusOptions[paymentStatusIndex] }}</view>
					</picker>
				</view>
				
				<view class="form-group">
					<text class="label">备注</text>
					<textarea class="form-textarea" v-model="formData.remarks" placeholder="请输入备注信息（可选）"></textarea>
				</view>
				
				<view class="total-section">
					<view class="total-row">
						<text class="total-label">订单总额:</text>
						<text class="total-amount">¥{{ totalAmount }}</text>
					</view>
				</view>
			</view>
			
			<!-- 提交按钮 -->
			<view class="submit-section">
				<button class="submit-btn" @click="submitOrder">{{ isEdit ? '更新订单' : '创建订单' }}</button>
			</view>
		</form>
	</view>
</template>

<script>
	import { api } from '../../utils/api.js'
	
	export default {
		data() {
			return {
				isEdit: false,
				orderId: null,
				formData: {
					customer_id: '',
					recipient: '',
					phone: '',
					delivery_time: '',
					delivery_addr: '',
					special_req: '',
					payment_status: 'unpaid',
					remarks: '',
					order_items: []
				},
				customers: [],
				customerOptions: [],
				customerIndex: 0,
				addresses: [],
				addressOptions: [],
				addressIndex: 0,
				products: [],
				productOptions: [],
				deliveryDate: '',
				paymentStatusIndex: 0,
				paymentStatusOptions: ['未付款', '已付款']
			}
		},
		computed: {
			totalAmount() {
				return this.formData.order_items.reduce((total, item) => {
					return total + (parseFloat(item.amount) || 0)
				}, 0).toFixed(2)
			}
		},
		onLoad(options) {
			if (options.id) {
				this.isEdit = true
				this.orderId = options.id
			}
			this.loadData()
		},
		methods: {
			async loadData() {
				try {
					uni.showLoading({ title: '加载中...' })
					
					// 加载客户列表
					const customersRes = await api.customers.list()
					if (customersRes.code === 0) {
						this.customers = Array.isArray(customersRes.data) ? customersRes.data : [customersRes.data].filter(Boolean)
						this.customerOptions = this.customers.map(customer => ({
							id: customer.customer?.id || customer.id,
							name: customer.customer?.name || customer.name
						}))
					}
					
					// 加载商品列表
					const productsRes = await api.products.list()
					if (productsRes.code === 0) {
						this.products = Array.isArray(productsRes.data) ? productsRes.data : [productsRes.data].filter(Boolean)
						this.productOptions = this.products.map(product => ({
							id: product.id,
							name: product.name,
							price: product.price
						}))
					}
					
					// 如果是编辑模式，加载订单数据
					if (this.isEdit) {
						await this.loadOrder()
					}
					
				} catch (error) {
					console.error('加载数据失败:', error)
					uni.showToast({
						title: '加载失败',
						icon: 'none'
					})
				} finally {
					uni.hideLoading()
				}
			},
			
			async loadOrder() {
				try {
					const res = await api.orders.get(this.orderId)
					if (res.code === 0) {
						const order = res.data.order || res.data
						this.formData = {
							customer_id: order.customer_id,
							recipient: order.recipient,
							phone: order.phone || '',
							delivery_time: order.delivery_time,
							delivery_addr: order.delivery_addr,
							special_req: order.special_req || '',
							payment_status: order.payment_status,
							remarks: order.remarks || '',
							order_items: res.data.order_items || []
						}
						
						// 设置选择器的值
						this.deliveryDate = this.formatDate(order.delivery_time)
						this.paymentStatusIndex = order.payment_status === 'paid' ? 1 : 0
						
						// 设置客户选择器
						const customerIndex = this.customerOptions.findIndex(c => c.id === order.customer_id)
						if (customerIndex >= 0) {
							this.customerIndex = customerIndex
							await this.loadCustomerAddresses(order.customer_id)
						}
					}
				} catch (error) {
					console.error('加载订单失败:', error)
				}
			},
			
			async onCustomerChange(e) {
				this.customerIndex = e.detail.value
				const customer = this.customerOptions[this.customerIndex]
				if (customer) {
					this.formData.customer_id = customer.id
					await this.loadCustomerAddresses(customer.id)
				}
			},
			
			async loadCustomerAddresses(customerId) {
				try {
					const customer = this.customers.find(c => (c.customer?.id || c.id) === customerId)
					if (customer && customer.addresses) {
						this.addresses = customer.addresses
						this.addressOptions = customer.addresses.map(addr => ({
							id: addr.id,
							address: addr.address,
							label: addr.label
						}))
						
						// 默认选择第一个地址
						if (this.addressOptions.length > 0) {
							this.addressIndex = 0
							this.formData.delivery_addr = this.addressOptions[0].address
						}
					}
				} catch (error) {
					console.error('加载地址失败:', error)
				}
			},
			
			onAddressChange(e) {
				this.addressIndex = e.detail.value
				const address = this.addressOptions[this.addressIndex]
				if (address) {
					this.formData.delivery_addr = address.address
				}
			},
			
			onDateChange(e) {
				this.deliveryDate = e.detail.value
				this.formData.delivery_time = e.detail.value + 'T00:00:00Z'
			},
			
			onPaymentStatusChange(e) {
				this.paymentStatusIndex = e.detail.value
				this.formData.payment_status = e.detail.value === 1 ? 'paid' : 'unpaid'
			},
			
			addProduct() {
				this.formData.order_items.push({
					product_id: '',
					productIndex: 0,
					quantity: 1,
					price: 0,
					amount: 0
				})
			},
			
			removeProduct(index) {
				this.formData.order_items.splice(index, 1)
			},
			
			onProductChange(e, index) {
				const productIndex = e.detail.value
				const product = this.productOptions[productIndex]
				if (product) {
					this.formData.order_items[index].productIndex = productIndex
					this.formData.order_items[index].product_id = product.id
					this.formData.order_items[index].price = product.price
					this.calculateAmount(index)
				}
			},
			
			calculateAmount(index) {
				const item = this.formData.order_items[index]
				const quantity = parseFloat(item.quantity) || 0
				const price = parseFloat(item.price) || 0
				item.amount = (quantity * price).toFixed(2)
			},
			
			validateForm() {
				if (!this.formData.customer_id) {
					uni.showToast({ title: '请选择客户', icon: 'none' })
					return false
				}
				if (!this.formData.recipient) {
					uni.showToast({ title: '请输入收货人', icon: 'none' })
					return false
				}
				if (!this.formData.delivery_time) {
					uni.showToast({ title: '请选择送货时间', icon: 'none' })
					return false
				}
				if (!this.formData.delivery_addr) {
					uni.showToast({ title: '请选择送货地址', icon: 'none' })
					return false
				}
				if (this.formData.order_items.length === 0) {
					uni.showToast({ title: '请添加商品', icon: 'none' })
					return false
				}
				
				// 验证商品信息
				for (let i = 0; i < this.formData.order_items.length; i++) {
					const item = this.formData.order_items[i]
					if (!item.product_id) {
						uni.showToast({ title: `请选择第${i+1}个商品`, icon: 'none' })
						return false
					}
					if (!item.quantity || item.quantity <= 0) {
						uni.showToast({ title: `请输入第${i+1}个商品的数量`, icon: 'none' })
						return false
					}
				}
				
				return true
			},
			
			async submitOrder() {
				if (!this.validateForm()) {
					return
				}
				
				try {
					uni.showLoading({ title: this.isEdit ? '更新中...' : '创建中...' })
					
					// 准备提交数据
					const submitData = {
						...this.formData,
						total_amount: parseFloat(this.totalAmount)
					}
					
					// 清理order_items中的临时字段
					submitData.order_items = submitData.order_items.map(item => ({
						product_id: item.product_id,
						quantity: parseInt(item.quantity),
						price: parseFloat(item.price),
						amount: parseFloat(item.amount)
					}))
					
					let res
					if (this.isEdit) {
						res = await api.orders.update(this.orderId, submitData)
					} else {
						res = await api.orders.create(submitData)
					}
					
					if (res.code === 0) {
						uni.showToast({
							title: this.isEdit ? '更新成功' : '创建成功',
							icon: 'success'
						})
						
						setTimeout(() => {
							uni.navigateBack()
						}, 1500)
					} else {
						throw new Error(res.message || '操作失败')
					}
					
				} catch (error) {
					console.error('提交订单失败:', error)
					uni.showToast({
						title: error.message || '操作失败',
						icon: 'none'
					})
				} finally {
					uni.hideLoading()
				}
			},
			
			formatDate(dateStr) {
				if (!dateStr) return ''
				const date = new Date(dateStr)
				return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
			}
		}
	}
</script>

<style scoped>
	.card-title {
		font-size: 32rpx;
		font-weight: bold;
		margin-bottom: 30rpx;
		color: #333;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}
	
	.add-product-btn {
		background-color: #4caf50;
		color: white;
		padding: 12rpx 24rpx;
		border-radius: 8rpx;
		font-size: 24rpx;
		border: none;
	}
	
	.form-group {
		margin-bottom: 30rpx;
	}
	
	.label {
		display: block;
		font-size: 28rpx;
		color: #333;
		margin-bottom: 15rpx;
	}
	
	.form-input {
		width: 100%;
		padding: 24rpx;
		border: 2rpx solid #e0e0e0;
		border-radius: 10rpx;
		font-size: 28rpx;
		background-color: #fff;
	}
	
	.form-textarea {
		width: 100%;
		padding: 24rpx;
		border: 2rpx solid #e0e0e0;
		border-radius: 10rpx;
		font-size: 28rpx;
		background-color: #fff;
		min-height: 120rpx;
	}
	
	.picker-input {
		padding: 24rpx;
		border: 2rpx solid #e0e0e0;
		border-radius: 10rpx;
		background-color: #fff;
		color: #333;
		font-size: 28rpx;
	}
	
	.product-list {
		margin-top: 20rpx;
	}
	
	.product-item {
		border: 2rpx solid #e0e0e0;
		border-radius: 10rpx;
		padding: 20rpx;
		margin-bottom: 20rpx;
		background-color: #fafafa;
	}
	
	.product-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 20rpx;
	}
	
	.product-picker {
		flex: 1;
		padding: 20rpx;
		border: 2rpx solid #e0e0e0;
		border-radius: 8rpx;
		background-color: #fff;
		margin-right: 20rpx;
		font-size: 28rpx;
	}
	
	.remove-btn {
		background-color: #ff3b30;
		color: white;
		padding: 20rpx 30rpx;
		border-radius: 8rpx;
		font-size: 24rpx;
		border: none;
	}
	
	.product-details {
		display: flex;
		gap: 20rpx;
		align-items: center;
	}
	
	.quantity-input {
		flex: 1;
		display: flex;
		align-items: center;
		gap: 10rpx;
	}
	
	.quantity-input input {
		flex: 1;
		padding: 16rpx;
		border: 2rpx solid #e0e0e0;
		border-radius: 6rpx;
		text-align: center;
		font-size: 28rpx;
	}
	
	.price-display, .amount-display {
		display: flex;
		align-items: center;
		gap: 10rpx;
	}
	
	.input-label {
		font-size: 24rpx;
		color: #666;
		white-space: nowrap;
	}
	
	.price-text, .amount-text {
		font-size: 28rpx;
		color: #333;
		font-weight: bold;
	}
	
	.amount-text {
		color: #e91e63;
	}
	
	.empty-products {
		text-align: center;
		padding: 60rpx 0;
		color: #999;
	}
	
	.empty-text {
		font-size: 28rpx;
	}
	
	.total-section {
		border-top: 2rpx solid #e0e0e0;
		padding-top: 30rpx;
		margin-top: 30rpx;
	}
	
	.total-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}
	
	.total-label {
		font-size: 32rpx;
		color: #333;
		font-weight: bold;
	}
	
	.total-amount {
		font-size: 48rpx;
		color: #e91e63;
		font-weight: bold;
	}
	
	.submit-section {
		margin-top: 40rpx;
		margin-bottom: 40rpx;
	}
	
	.submit-btn {
		width: 100%;
		padding: 32rpx;
		background-color: #007aff;
		color: white;
		border-radius: 15rpx;
		font-size: 32rpx;
		font-weight: bold;
		border: none;
	}
</style>