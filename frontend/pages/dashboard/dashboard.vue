<template>
	<view class="container">
		<!-- 顶部统计卡片 -->
		<view class="stats-grid">
			<view class="stat-card" v-for="stat in stats" :key="stat.title">
				<view class="stat-number">{{ stat.value }}</view>
				<view class="stat-title">{{ stat.title }}</view>
			</view>
		</view>
		
		<!-- 快速操作 -->
		<view class="card">
			<view class="card-title">快速操作</view>
			<view class="quick-actions">
				<view class="action-item" @click="navigateTo('/pages/orders/add')">
					<view class="action-icon">📝</view>
					<text class="action-text">新建订单</text>
				</view>
				<view class="action-item" @click="navigateTo('/pages/products/add')">
					<view class="action-icon">📦</view>
					<text class="action-text">添加商品</text>
				</view>
				<view class="action-item" @click="navigateTo('/pages/customers/add')">
					<view class="action-icon">👤</view>
					<text class="action-text">添加客户</text>
				</view>
				<view class="action-item" @click="navigateTo('/pages/calendar/calendar')">
					<view class="action-icon">📅</view>
					<text class="action-text">查看日历</text>
				</view>
			</view>
		</view>
		
		<view class="card">
			<view class="card-title">今日概况</view>
			<view class="today-stats">
				<view class="today-item">
					<text class="label">订单数量</text>
					<text class="value">{{ dashboardData.today?.orders || 0 }}</text>
				</view>
				<view class="today-item">
					<text class="label">销售金额</text>
					<text class="value">¥{{ dashboardData.today?.total_amount || 0 }}</text>
				</view>
			</view>
		</view>
		
		<view class="card">
			<view class="card-title">本月统计</view>
			<view class="month-stats">
				<view class="month-item">
					<text class="label">订单总数</text>
					<text class="value">{{ dashboardData.month?.orders || 0 }}</text>
				</view>
				<view class="month-item">
					<text class="label">销售总额</text>
					<text class="value">¥{{ dashboardData.month?.total_amount || 0 }}</text>
				</view>
			</view>
		</view>
		
		<!-- 今日订单列表 -->
		<view class="card">
			<view class="card-header">
				<view class="card-title">今日订单</view>
				<text class="view-all" @click="navigateTo('/pages/orders/orders')">查看全部</text>
			</view>
			<view v-if="todayOrders.length === 0" class="empty-state">
				<text class="empty-text">今日暂无订单</text>
			</view>
			<view v-else class="order-list">
				<view v-for="order in todayOrders" :key="order.id" class="order-item" @click="viewOrderDetail(order)">
					<view class="order-header">
						<text class="order-id">#{{ order.id }}</text>
						<view class="order-status" :class="{ 'paid': order.payment_status === 1 }">
							{{ order.payment_status === 1 ? '已付款' : '未付款' }}
						</view>
					</view>
					<text class="customer-name">{{ order.customer_name }}</text>
					<text class="order-amount">¥{{ order.total_amount }}</text>
				</view>
			</view>
		</view>
		
		<!-- 热销商品 -->
		<view class="card">
			<view class="card-header">
				<view class="card-title">热销商品</view>
				<text class="view-all" @click="navigateTo('/pages/products/products')">查看全部</text>
			</view>
			<view v-if="hotProducts.length === 0" class="empty-state">
				<text class="empty-text">暂无数据</text>
			</view>
			<view v-else class="product-list">
				<view v-for="product in hotProducts" :key="product.id" class="product-item">
					<image class="product-image" :src="product.image || '/static/default-product.png'" mode="aspectFill"></image>
					<view class="product-info">
						<text class="product-name">{{ product.name }}</text>
						<text class="product-sales">销量: {{ product.sales || 0 }}</text>
					</view>
					<text class="product-price">¥{{ product.price }}</text>
				</view>
			</view>
		</view>
	</view>
</template>

<script>
	import { api } from '../../utils/api.js'
	
	export default {
		data() {
			return {
				dashboardData: {},
				stats: [
					{ title: '今日订单', value: 0 },
					{ title: '本月订单', value: 0 },
					{ title: '客户总数', value: 0 },
					{ title: '商品总数', value: 0 }
				],
				todayOrders: [],
				hotProducts: []
			}
		},
		onLoad() {
			this.loadDashboard()
		},
		onShow() {
			// 每次显示页面时刷新数据
			this.loadDashboard()
		},
		methods: {
			async loadDashboard() {
				try {
					uni.showLoading({ title: '加载中...' });
					
					// 并行加载所有数据
					const [statsRes, ordersRes, productsRes] = await Promise.all([
						api.dashboard.stats(),
						api.order.list({ date: new Date().toISOString().split('T')[0] }),
						api.dashboard.hotProducts()
					]);

					if (statsRes.code === 0) {
						this.dashboardData = statsRes.data;
						this.updateStats();
					}
					
					if (ordersRes.code === 0) {
						this.todayOrders = ordersRes.data.slice(0, 5); // 只显示前5个
					}
					
					if (productsRes.code === 0) {
						this.hotProducts = productsRes.data.slice(0, 5); // 只显示前5个
					}
					
				} catch (error) {
					console.error('加载驾驶仓数据失败:', error)
					uni.showToast({
						title: '加载失败',
						icon: 'none'
					})
				} finally {
					uni.hideLoading();
				}
			},
			updateStats() {
				const data = this.dashboardData
				this.stats = [
					{ title: '今日订单', value: data.today?.orders || 0 },
					{ title: '本月订单', value: data.month?.orders || 0 },
					{ title: '客户总数', value: data.customer_stats?.total_customers || 0 },
					{ title: '商品总数', value: data.product_stats?.total_products || 0 }
				]
			},
			navigateTo(url) {
				uni.navigateTo({ url });
			},
			viewOrderDetail(order) {
				uni.navigateTo({
					url: `/pages/orders/detail?id=${order.id}`
				});
			}
		}
	}
</script>

<style scoped>
	.container {
		padding: 20rpx;
		background: #f5f5f5;
		min-height: 100vh;
	}
	
	.stats-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 20rpx;
		margin-bottom: 30rpx;
	}
	
	.stat-card {
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		color: white;
		padding: 30rpx;
		border-radius: 15rpx;
		text-align: center;
	}
	
	.stat-number {
		font-size: 48rpx;
		font-weight: bold;
		margin-bottom: 10rpx;
	}
	
	.stat-title {
		font-size: 28rpx;
		opacity: 0.9;
	}
	
	.card {
		background: #fff;
		border-radius: 15rpx;
		padding: 30rpx;
		margin-bottom: 30rpx;
		box-shadow: 0 4rpx 20rpx rgba(0,0,0,0.1);
	}
	
	.card-title {
		font-size: 32rpx;
		font-weight: bold;
		margin-bottom: 20rpx;
		color: #333;
	}
	
	.card-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 20rpx;
	}
	
	.view-all {
		font-size: 28rpx;
		color: #1890ff;
	}
	
	.quick-actions {
		display: grid;
		grid-template-columns: 1fr 1fr 1fr 1fr;
		gap: 20rpx;
	}
	
	.action-item {
		text-align: center;
		padding: 20rpx;
		border-radius: 10rpx;
		background: #f8f9fa;
	}
	
	.action-icon {
		font-size: 40rpx;
		margin-bottom: 10rpx;
	}
	
	.action-text {
		font-size: 24rpx;
		color: #666;
		display: block;
	}
	
	.today-stats, .month-stats {
		display: flex;
		justify-content: space-between;
	}
	
	.today-item, .month-item {
		flex: 1;
		text-align: center;
		padding: 20rpx;
		background-color: #f8f9fa;
		border-radius: 10rpx;
		margin: 0 10rpx;
	}
	
	.label {
		display: block;
		font-size: 28rpx;
		color: #666;
		margin-bottom: 10rpx;
	}
	
	.value {
		display: block;
		font-size: 36rpx;
		font-weight: bold;
		color: #333;
	}
	
	.empty-state {
		text-align: center;
		padding: 60rpx 0;
	}
	
	.empty-text {
		font-size: 28rpx;
		color: #999;
	}
	
	.order-list {
		/* 订单列表样式 */
	}
	
	.order-item {
		padding: 20rpx;
		border-bottom: 1rpx solid #f0f0f0;
		display: flex;
		flex-direction: column;
	}
	
	.order-item:last-child {
		border-bottom: none;
	}
	
	.order-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 10rpx;
	}
	
	.order-id {
		font-size: 28rpx;
		font-weight: bold;
		color: #333;
	}
	
	.order-status {
		padding: 4rpx 12rpx;
		border-radius: 8rpx;
		font-size: 24rpx;
		background: #fff2e8;
		color: #fa8c16;
	}
	
	.order-status.paid {
		background: #e8f5e8;
		color: #52c41a;
	}
	
	.customer-name {
		font-size: 28rpx;
		color: #666;
		margin-bottom: 8rpx;
	}
	
	.order-amount {
		font-size: 32rpx;
		color: #ff6b35;
		font-weight: bold;
	}
	
	.product-list {
		/* 商品列表样式 */
	}
	
	.product-item {
		display: flex;
		align-items: center;
		padding: 20rpx 0;
		border-bottom: 1rpx solid #f0f0f0;
	}
	
	.product-item:last-child {
		border-bottom: none;
	}
	
	.product-image {
		width: 80rpx;
		height: 80rpx;
		border-radius: 8rpx;
		margin-right: 20rpx;
	}
	
	.product-info {
		flex: 1;
	}
	
	.product-name {
		font-size: 28rpx;
		color: #333;
		display: block;
		margin-bottom: 8rpx;
	}
	
	.product-sales {
		font-size: 24rpx;
		color: #666;
	}
	
	.product-price {
		font-size: 32rpx;
		color: #ff6b35;
		font-weight: bold;
	}
</style>