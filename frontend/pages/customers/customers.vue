<template>
	<view class="container">
		<!-- 搜索栏 -->
		<view class="search-bar">
			<input class="search-input" placeholder="搜索客户..." v-model="searchText" @input="onSearch" />
		</view>
		
		<!-- 客户列表 -->
		<view class="customer-list">
			<view class="customer-card card" v-for="customer in filteredCustomers" :key="customer.id" @click="viewCustomer(customer)">
				<view class="customer-header">
					<view class="customer-avatar">
						<text class="avatar-text">{{ getAvatarText(customer.name) }}</text>
					</view>
					<view class="customer-info">
						<text class="customer-name">{{ customer.name }}</text>
						<text class="customer-phone">{{ customer.phone || '未填写电话' }}</text>
					</view>
				</view>
				
				<view class="customer-addresses">
					<view class="address-item" v-for="(address, index) in customer.addresses" :key="index">
						<view class="address-label">{{ address.label || '地址' + (index + 1) }}</view>
						<view class="address-text">{{ address.address }}</view>
					</view>
					<view v-if="!customer.addresses || customer.addresses.length === 0" class="no-address">
						<text class="no-address-text">暂无地址信息</text>
					</view>
				</view>
				
				<view class="customer-footer">
					<text class="create-time">{{ formatDateTime(customer.created_at) }}</text>
				</view>
			</view>
		</view>
		
		<!-- 空状态 -->
		<view v-if="filteredCustomers.length === 0" class="empty-state">
			<text class="empty-text">暂无客户数据</text>
		</view>
		
		<!-- 添加按钮 -->
		<view class="fab" @click="addCustomer">
			<text class="fab-icon">+</text>
		</view>
		
		<!-- 客户详情弹窗 -->
		<uni-popup ref="customerPopup" type="bottom">
			<view class="popup-content">
				<view class="popup-header">
					<text class="popup-title">客户详情</text>
					<text class="close-btn" @click="closePopup">×</text>
				</view>
				
				<view v-if="selectedCustomer" class="customer-detail">
					<view class="detail-header">
						<view class="detail-avatar">
							<text class="avatar-text">{{ getAvatarText(selectedCustomer.name) }}</text>
						</view>
						<view class="detail-info">
							<text class="detail-name">{{ selectedCustomer.name }}</text>
							<text class="detail-phone">{{ selectedCustomer.phone || '未填写电话' }}</text>
						</view>
					</view>
					
					<view class="detail-section">
						<text class="section-title">基本信息</text>
						<view class="detail-row">
							<text class="detail-label">客户姓名：</text>
							<text class="detail-value">{{ selectedCustomer.name }}</text>
						</view>
						<view class="detail-row">
							<text class="detail-label">联系电话：</text>
							<text class="detail-value">{{ selectedCustomer.phone || '未填写' }}</text>
						</view>
						<view class="detail-row">
							<text class="detail-label">创建时间：</text>
							<text class="detail-value">{{ formatDateTime(selectedCustomer.created_at) }}</text>
						</view>
					</view>
					
					<view class="detail-section">
						<text class="section-title">地址信息</text>
						<view v-if="selectedCustomer.addresses && selectedCustomer.addresses.length > 0">
							<view class="address-detail" v-for="(address, index) in selectedCustomer.addresses" :key="index">
								<view class="address-header">
									<text class="address-label-detail">{{ address.label || '地址' + (index + 1) }}</text>
								</view>
								<text class="address-content">{{ address.address }}</text>
							</view>
						</view>
						<view v-else class="no-address-detail">
							<text class="no-address-text">暂无地址信息</text>
						</view>
					</view>
				</view>
				
				<view class="popup-actions">
					<button class="action-btn edit-btn" @click="editCustomer">编辑</button>
					<button class="action-btn delete-btn" @click="deleteCustomer">删除</button>
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
				customers: [],
				filteredCustomers: [],
				searchText: '',
				selectedCustomer: null
			}
		},
		onLoad() {
			this.loadCustomers()
		},
		onShow() {
			this.loadCustomers()
		},
		methods: {
			async loadCustomers() {
				try {
					uni.showLoading({ title: '加载中...' })
					const res = await api.customers.list()
					if (res.code === 0) {
						this.customers = Array.isArray(res.data) ? res.data : [res.data].filter(Boolean)
						this.filterCustomers()
					}
				} catch (error) {
					console.error('加载客户失败:', error)
					uni.showToast({
						title: '加载失败',
						icon: 'none'
					})
				} finally {
					uni.hideLoading()
				}
			},
			
			filterCustomers() {
				let filtered = this.customers
				
				// 按搜索文本筛选
				if (this.searchText) {
					const text = this.searchText.toLowerCase()
					filtered = filtered.filter(customer => {
						const customerData = customer.customer || customer
						return customerData.name?.toLowerCase().includes(text) ||
							   customerData.phone?.toLowerCase().includes(text) ||
							   (customer.addresses && customer.addresses.some(addr => 
								   addr.address?.toLowerCase().includes(text)
							   ))
					})
				}
				
				this.filteredCustomers = filtered
			},
			
			onSearch() {
				this.filterCustomers()
			},
			
			viewCustomer(customer) {
				this.selectedCustomer = customer.customer || customer
				this.selectedCustomer.addresses = customer.addresses || []
				this.$refs.customerPopup.open()
			},
			
			closePopup() {
				this.$refs.customerPopup.close()
			},
			
			addCustomer() {
				uni.navigateTo({
					url: '/pages/customers/customer-form'
				})
			},
			
			editCustomer() {
				this.closePopup()
				uni.navigateTo({
					url: `/pages/customers/customer-form?id=${this.selectedCustomer.id}`
				})
			},
			
			async deleteCustomer() {
				const that = this
				uni.showModal({
					title: '确认删除',
					content: '确定要删除这个客户吗？',
					success: async function(res) {
						if (res.confirm) {
							try {
								await api.customers.delete(that.selectedCustomer.id)
								uni.showToast({
									title: '删除成功',
									icon: 'success'
								})
								that.closePopup()
								that.loadCustomers()
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
			
			getAvatarText(name) {
				if (!name) return '?'
				return name.charAt(0).toUpperCase()
			},
			
			formatDateTime(dateStr) {
				if (!dateStr) return ''
				const date = new Date(dateStr)
				return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
			}
		}
	}
</script>

<style scoped>
	.search-bar {
		margin-bottom: 30rpx;
	}
	
	.search-input {
		width: 100%;
		padding: 20rpx;
		border: 2rpx solid #e0e0e0;
		border-radius: 10rpx;
		background-color: #fff;
	}
	
	.customer-list {
		margin-bottom: 100rpx;
	}
	
	.customer-card {
		margin-bottom: 20rpx;
		cursor: pointer;
	}
	
	.customer-header {
		display: flex;
		align-items: center;
		margin-bottom: 20rpx;
	}
	
	.customer-avatar {
		width: 80rpx;
		height: 80rpx;
		border-radius: 40rpx;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		display: flex;
		align-items: center;
		justify-content: center;
		margin-right: 20rpx;
	}
	
	.avatar-text {
		color: white;
		font-size: 32rpx;
		font-weight: bold;
	}
	
	.customer-info {
		flex: 1;
	}
	
	.customer-name {
		display: block;
		font-size: 32rpx;
		font-weight: bold;
		color: #333;
		margin-bottom: 8rpx;
	}
	
	.customer-phone {
		display: block;
		font-size: 26rpx;
		color: #666;
	}
	
	.customer-addresses {
		margin-bottom: 20rpx;
	}
	
	.address-item {
		margin-bottom: 15rpx;
		padding: 15rpx;
		background-color: #f8f9fa;
		border-radius: 8rpx;
	}
	
	.address-label {
		font-size: 24rpx;
		color: #007aff;
		margin-bottom: 8rpx;
		font-weight: bold;
	}
	
	.address-text {
		font-size: 26rpx;
		color: #333;
		line-height: 1.4;
	}
	
	.no-address {
		text-align: center;
		padding: 30rpx 0;
	}
	
	.no-address-text {
		color: #999;
		font-size: 26rpx;
	}
	
	.customer-footer {
		border-top: 2rpx solid #f0f0f0;
		padding-top: 15rpx;
		text-align: right;
	}
	
	.create-time {
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
	
	.customer-detail {
		padding: 30rpx;
	}
	
	.detail-header {
		display: flex;
		align-items: center;
		margin-bottom: 40rpx;
		padding-bottom: 30rpx;
		border-bottom: 2rpx solid #f0f0f0;
	}
	
	.detail-avatar {
		width: 120rpx;
		height: 120rpx;
		border-radius: 60rpx;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		display: flex;
		align-items: center;
		justify-content: center;
		margin-right: 30rpx;
	}
	
	.detail-info {
		flex: 1;
	}
	
	.detail-name {
		display: block;
		font-size: 36rpx;
		font-weight: bold;
		color: #333;
		margin-bottom: 12rpx;
	}
	
	.detail-phone {
		display: block;
		font-size: 28rpx;
		color: #666;
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
	
	.address-detail {
		margin-bottom: 20rpx;
		padding: 20rpx;
		background-color: #f8f9fa;
		border-radius: 10rpx;
	}
	
	.address-header {
		margin-bottom: 12rpx;
	}
	
	.address-label-detail {
		font-size: 26rpx;
		color: #007aff;
		font-weight: bold;
	}
	
	.address-content {
		font-size: 28rpx;
		color: #333;
		line-height: 1.5;
	}
	
	.no-address-detail {
		text-align: center;
		padding: 40rpx 0;
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