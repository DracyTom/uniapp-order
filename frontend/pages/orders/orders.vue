<template>
	<view class="container">
		<!-- 搜索和筛选 -->
		<view class="search-bar">
			<input class="search-input" placeholder="搜索订单..." v-model="searchText" @input="onSearch" />
			<picker @change="onStatusChange" :value="statusIndex" :range="statusOptions">
				<view class="filter-btn">{{ statusOptions[statusIndex] }}</view>
			</picker>
		</view>
		
		<!-- 订单列表 -->
		<view class="order-list">
			<view class="order-item card" v-for="order in filteredOrders" :key="order.id" @click="viewOrder(order)">
				<view class="order-header flex-between">
					<text class="order-number">订单 #{{ order.id }}</text>
					<text class="order-status" :class="getStatusClass(order.payment_status)">
						{{ getStatusText(order.payment_status) }}
					</text>
				</view>
				
				<view class="order-info">
					<view class="info-row">
						<text class="label">客户：</text>
						<text class="value">{{ order.customer_name }}</text>
					</view>
					<view class="info-row">
						<text class="label">收货人：</text>
						<text class="value">{{ order.recipient }}</text>
					</view>
					<view class="info-row">
						<text class="label">送货时间：</text>
						<text class="value">{{ formatDate(order.delivery_time) }}</text>
					</view>
					<view class="info-row">
						<text class="label">送货地址：</text>
						<text class="value">{{ order.delivery_addr }}</text>
					</view>
				</view>
				
				<view class="order-footer flex-between">
					<text class="order-amount">¥{{ order.total_amount }}</text>
					<text class="order-date">{{ formatDateTime(order.created_at) }}</text>
				</view>
			</view>
		</view>
		
		<!-- 空状态 -->
		<view v-if="filteredOrders.length === 0" class="empty-state">
			<text class="empty-text">暂无订单数据</text>
		</view>
		
		<!-- 添加按钮 -->
		<view class="fab" @click="addOrder">
			<text class="fab-icon">+</text>
		</view>
		
		<!-- 订单详情弹窗 -->
		<uni-popup ref="orderPopup" type="bottom">
			<view class="popup-content">
				<view class="popup-header">
					<text class="popup-title">订单详情</text>
					<text class="close-btn" @click="closePopup">×</text>
				</view>
				
				<view v-if="selectedOrder" class="order-detail">
					<view class="detail-section">
						<text class="section-title">基本信息</text>
						<view class="detail-row">
							<text class="detail-label">订单号：</text>
							<text class="detail-value">#{{ selectedOrder.id }}</text>
						</view>
						<view class="detail-row">
							<text class="detail-label">客户：</text>
							<text class="detail-value">{{ selectedOrder.customer_name }}</text>
						</view>
						<view class="detail-row">
							<text class="detail-label">收货人：</text>
							<text class="detail-value">{{ selectedOrder.recipient }}</text>
						</view>
						<view class="detail-row">
							<text class="detail-label">联系电话：</text>
							<text class="detail-value">{{ selectedOrder.phone || '未填写' }}</text>
						</view>
					</view>
					
					<view class="detail-section">
						<text class="section-title">配送信息</text>
						<view class="detail-row">
							<text class="detail-label">送货时间：</text>
							<text class="detail-value">{{ formatDate(selectedOrder.delivery_time) }}</text>
						</view>
						<view class="detail-row">
							<text class="detail-label">送货地址：</text>
							<text class="detail-value">{{ selectedOrder.delivery_addr }}</text>
						</view>
						<view class="detail-row">
							<text class="detail-label">特殊要求：</text>
							<text class="detail-value">{{ selectedOrder.special_req || '无' }}</text>
						</view>
					</view>
					
					<view class="detail-section">
						<text class="section-title">商品清单</text>
						<view class="product-item" v-for="item in selectedOrder.order_items" :key="item.id">
							<view class="product-info">
								<text class="product-name">{{ item.product_name || '商品' + item.product_id }}</text>
								<text class="product-spec">数量: {{ item.quantity }} × ¥{{ item.price }}</text>
							</view>
							<text class="product-amount">¥{{ item.amount }}</text>
						</view>
					</view>
					
					<view class="detail-section">
						<text class="section-title">费用信息</text>
						<view class="detail-row">
							<text class="detail-label">订单总额：</text>
							<text class="detail-value total-amount">¥{{ selectedOrder.total_amount }}</text>
						</view>
						<view class="detail-row">
							<text class="detail-label">支付状态：</text>
							<text class="detail-value" :class="getStatusClass(selectedOrder.payment_status)">
								{{ getStatusText(selectedOrder.payment_status) }}
							</text>
						</view>
					</view>
					
					<view v-if="selectedOrder.remarks" class="detail-section">
						<text class="section-title">备注</text>
						<text class="detail-value">{{ selectedOrder.remarks }}</text>
					</view>
				</view>
				
				<view class="popup-actions">
					<button class="action-btn edit-btn" @click="editOrder">编辑</button>
					<button class="action-btn delete-btn" @click="deleteOrder">删除</button>
				</view>
			</view>
		</uni-popup>
	</view>
</template>

<script>
	import { api } from '../../utils/api.js'
	
	export default {
		data() {
			return {
				orders: [],
				filteredOrders: [],
				searchText: '',
				statusIndex: 0,
				statusOptions: ['全部', '未付款', '已付款'],
				selectedOrder: null
			}
		},
		onLoad() {
			this.loadOrders()
		},
		onShow() {
			this.loadOrders()
		},
		methods: {
			async loadOrders() {
				try {
					uni.showLoading({ title: '加载中...' })
					const res = await api.orders.list()
					if (res.code === 0) {
						this.orders = Array.isArray(res.data) ? res.data : [res.data].filter(Boolean)
						this.filterOrders()
					}
				} catch (error) {
					console.error('加载订单失败:', error)
					uni.showToast({
						title: '加载失败',
						icon: 'none'
					})
				} finally {
					uni.hideLoading()
				}
			},
			
			filterOrders() {
				let filtered = this.orders
				
				// 按状态筛选
				if (this.statusIndex === 1) {
					filtered = filtered.filter(order => order.payment_status === 'unpaid')
				} else if (this.statusIndex === 2) {
					filtered = filtered.filter(order => order.payment_status === 'paid')
				}
				
				// 按搜索文本筛选
				if (this.searchText) {
					const text = this.searchText.toLowerCase()
					filtered = filtered.filter(order => 
						order.customer_name?.toLowerCase().includes(text) ||
						order.recipient?.toLowerCase().includes(text) ||
						order.id.toString().includes(text)
					)
				}
				
				this.filteredOrders = filtered
			},
			
			onSearch() {
				this.filterOrders()
			},
			
			onStatusChange(e) {
				this.statusIndex = e.detail.value
				this.filterOrders()
			},
			
			viewOrder(order) {
				this.selectedOrder = order
				this.$refs.orderPopup.open()
			},
			
			closePopup() {
				this.$refs.orderPopup.close()
			},
			
			addOrder() {
				uni.navigateTo({
					url: '/pages/orders/order-form'
				})
			},
			
			editOrder() {
				this.closePopup()
				uni.navigateTo({
					url: `/pages/orders/order-form?id=${this.selectedOrder.id}`
				})
			},
			
			async deleteOrder() {
				const that = this
				uni.showModal({
					title: '确认删除',
					content: '确定要删除这个订单吗？',
					success: async function(res) {
						if (res.confirm) {
							try {
								await api.orders.delete(that.selectedOrder.id)
								uni.showToast({
									title: '删除成功',
									icon: 'success'
								})
								that.closePopup()
								that.loadOrders()
							} catch (error) {
								uni.showToast({
									title: '删除失败',
									icon: 'none'
								})
							}
						}
					}
				})
			},
			
			getStatusClass(status) {
				return status === 'paid' ? 'status-paid' : 'status-unpaid'
			},
			
			getStatusText(status) {
				return status === 'paid' ? '已付款' : '未付款'
			},
			
			formatDate(dateStr) {
				if (!dateStr) return ''
				const date = new Date(dateStr)
				return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
			},
			
			formatDateTime(dateStr) {
				if (!dateStr) return ''
				const date = new Date(dateStr)
				return `${date.getMonth() + 1}/${date.getDate()} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
			}
		}
	}
</script>

<style scoped>
	.search-bar {
		display: flex;
		gap: 20rpx;
		margin-bottom: 30rpx;
	}
	
	.search-input {
		flex: 1;
		padding: 20rpx;
		border: 2rpx solid #e0e0e0;
		border-radius: 10rpx;
		background-color: #fff;
	}
	
	.filter-btn {
		padding: 20rpx 30rpx;
		background-color: #007aff;
		color: white;
		border-radius: 10rpx;
		text-align: center;
		min-width: 120rpx;
	}
	
	.order-list {
		margin-bottom: 100rpx;
	}
	
	.order-item {
		margin-bottom: 20rpx;
		cursor: pointer;
	}
	
	.order-header {
		margin-bottom: 20rpx;
	}
	
	.order-number {
		font-size: 32rpx;
		font-weight: bold;
		color: #333;
	}
	
	.order-status {
		padding: 8rpx 16rpx;
		border-radius: 20rpx;
		font-size: 24rpx;
	}
	
	.status-paid {
		background-color: #e8f5e8;
		color: #4caf50;
	}
	
	.status-unpaid {
		background-color: #fff3e0;
		color: #ff9800;
	}
	
	.order-info {
		margin-bottom: 20rpx;
	}
	
	.info-row {
		display: flex;
		margin-bottom: 10rpx;
	}
	
	.label {
		color: #666;
		width: 140rpx;
		font-size: 28rpx;
	}
	
	.value {
		flex: 1;
		color: #333;
		font-size: 28rpx;
	}
	
	.order-footer {
		border-top: 2rpx solid #f0f0f0;
		padding-top: 20rpx;
	}
	
	.order-amount {
		font-size: 36rpx;
		font-weight: bold;
		color: #e91e63;
	}
	
	.order-date {
		color: #999;
		font-size: 24rpx;
	}
	
	.empty-state {
		text-align: center;
		padding: 100rpx 0;
	}
	
	.empty-text {
		color: #999;
		font-size: 28rpx;
	}
	
	.fab {
		position: fixed;
		right: 30rpx;
		bottom: 30rpx;
		width: 120rpx;
		height: 120rpx;
		background-color: #007aff;
		border-radius: 60rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		box-shadow: 0 4rpx 20rpx rgba(0, 122, 255, 0.3);
		z-index: 999;
	}
	
	.fab-icon {
		color: white;
		font-size: 48rpx;
		font-weight: bold;
	}
	
	.popup-content {
		background-color: white;
		border-radius: 20rpx 20rpx 0 0;
		max-height: 80vh;
		overflow-y: auto;
	}
	
	.popup-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 30rpx;
		border-bottom: 2rpx solid #f0f0f0;
	}
	
	.popup-title {
		font-size: 36rpx;
		font-weight: bold;
		color: #333;
	}
	
	.close-btn {
		font-size: 48rpx;
		color: #999;
		width: 60rpx;
		height: 60rpx;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.order-detail {
		padding: 30rpx;
	}
	
	.detail-section {
		margin-bottom: 40rpx;
	}
	
	.section-title {
		font-size: 32rpx;
		font-weight: bold;
		color: #333;
		margin-bottom: 20rpx;
		display: block;
	}
	
	.detail-row {
		display: flex;
		margin-bottom: 15rpx;
	}
	
	.detail-label {
		color: #666;
		width: 160rpx;
		font-size: 28rpx;
	}
	
	.detail-value {
		flex: 1;
		color: #333;
		font-size: 28rpx;
	}
	
	.total-amount {
		color: #e91e63;
		font-weight: bold;
		font-size: 32rpx;
	}
	
	.product-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 20rpx 0;
		border-bottom: 2rpx solid #f8f8f8;
	}
	
	.product-info {
		flex: 1;
	}
	
	.product-name {
		display: block;
		font-size: 28rpx;
		color: #333;
		margin-bottom: 8rpx;
	}
	
	.product-spec {
		display: block;
		font-size: 24rpx;
		color: #666;
	}
	
	.product-amount {
		font-size: 28rpx;
		color: #e91e63;
		font-weight: bold;
	}
	
	.popup-actions {
		display: flex;
		gap: 20rpx;
		padding: 30rpx;
		border-top: 2rpx solid #f0f0f0;
	}
	
	.action-btn {
		flex: 1;
		padding: 24rpx;
		border-radius: 10rpx;
		text-align: center;
		font-size: 28rpx;
		border: none;
	}
	
	.edit-btn {
		background-color: #007aff;
		color: white;
	}
	
	.delete-btn {
		background-color: #ff3b30;
		color: white;
	}
</style>