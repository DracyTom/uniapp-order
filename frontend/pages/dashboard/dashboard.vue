<template>
	<view class="container">
		<view class="stats-grid">
			<view class="stat-card" v-for="stat in stats" :key="stat.title">
				<view class="stat-number">{{ stat.value }}</view>
				<view class="stat-title">{{ stat.title }}</view>
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
				]
			}
		},
		onLoad() {
			this.loadDashboard()
		},
		methods: {
			async loadDashboard() {
				try {
					const res = await api.dashboard.stats()
					if (res.code === 0) {
						this.dashboardData = res.data
						this.updateStats()
					}
				} catch (error) {
					console.error('加载驾驶仓数据失败:', error)
					uni.showToast({
						title: '加载失败',
						icon: 'none'
					})
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
			}
		}
	}
</script>

<style scoped>
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
	
	.card-title {
		font-size: 32rpx;
		font-weight: bold;
		margin-bottom: 20rpx;
		color: #333;
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
</style>