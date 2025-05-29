<template>
	<view class="container">
		<view v-if="loading" class="loading">
			<text>加载中...</text>
		</view>
		
		<view v-else-if="order" class="order-detail">
			<!-- 订单头部信息 -->
			<view class="order-header">
				<view class="order-id">
					<text class="label">订单编号:</text>
					<text class="value">#{{ order.id }}</text>
				</view>
				<view class="order-status" :class="{ 'paid': order.payment_status === 1 }">
					{{ order.payment_status === 1 ? '已付款' : '未付款' }}
				</view>
			</view>
			
			<!-- 客户信息 -->
			<view class="section">
				<view class="section-title">客户信息</view>
				<view class="info-row">
					<text class="label">客户姓名:</text>
					<text class="value">{{ order.customer_name }}</text>
				</view>
				<view class="info-row">
					<text class="label">收货人:</text>
					<text class="value">{{ order.recipient_name }}</text>
				</view>
				<view class="info-row">
					<text class="label">送货地址:</text>
					<text class="value">{{ order.delivery_address }}</text>
				</view>
				<view class="info-row">
					<text class="label">送货时间:</text>
					<text class="value">{{ formatDateTime(order.delivery_time) }}</text>
				</view>
			</view>
			
			<!-- 商品清单 -->
			<view class="section">
				<view class="section-title">商品清单</view>
				<view class="items-list">
					<view 
						v-for="item in order.items" 
						:key="item.id"
						class="item-row"
					>
						<view class="item-info">
							<text class="item-name">{{ item.product_name }}</text>
							<text class="item-price">¥{{ item.unit_price }} × {{ item.quantity }}</text>
						</view>
						<text class="item-total">¥{{ (item.unit_price * item.quantity).toFixed(2) }}</text>
					</view>
				</view>
				
				<view class="total-section">
					<view class="total-row">
						<text class="total-label">合计金额:</text>
						<text class="total-amount">¥{{ order.total_amount }}</text>
					</view>
				</view>
			</view>
			
			<!-- 备注信息 -->
			<view v-if="order.notes" class="section">
				<view class="section-title">备注信息</view>
				<text class="notes-text">{{ order.notes }}</text>
			</view>
			
			<!-- 操作按钮 -->
			<view class="actions">
				<button class="action-btn edit-btn" @click="editOrder">编辑订单</button>
				<button 
					class="action-btn payment-btn" 
					:class="{ 'paid': order.payment_status === 1 }"
					@click="togglePaymentStatus"
				>
					{{ order.payment_status === 1 ? '标记未付款' : '标记已付款' }}
				</button>
				<button class="action-btn delete-btn" @click="deleteOrder">删除订单</button>
			</view>
		</view>
		
		<view v-else class="error">
			<text>订单不存在或加载失败</text>
		</view>
	</view>
</template>

<script>
	import { api } from '../../utils/api.js'
	import { formatDateTime } from '../../utils/common.js'
	
	export default {
		data() {
			return {
				orderId: null,
				order: null,
				loading: true
			}
		},
		
		onLoad(options) {
			if (options.id) {
				this.orderId = options.id;
				this.loadOrder();
			}
		},
		
		methods: {
			async loadOrder() {
				try {
					this.loading = true;
					const res = await api.order.get(this.orderId);
					if (res.code === 0) {
						this.order = res.data;
					} else {
						uni.showToast({
							title: res.message || '加载失败',
							icon: 'none'
						});
					}
				} catch (error) {
					console.error('加载订单失败:', error);
					uni.showToast({
						title: '加载失败',
						icon: 'none'
					});
				} finally {
					this.loading = false;
				}
			},
			
			editOrder() {
				uni.navigateTo({
					url: `/pages/orders/order-form?id=${this.orderId}`
				});
			},
			
			async togglePaymentStatus() {
				try {
					uni.showLoading({ title: '更新中...' });
					
					const newStatus = this.order.payment_status === 1 ? 0 : 1;
					const res = await api.order.update(this.orderId, {
						payment_status: newStatus
					});
					
					if (res.code === 0) {
						this.order.payment_status = newStatus;
						uni.showToast({
							title: '更新成功',
							icon: 'success'
						});
					} else {
						uni.showToast({
							title: res.message || '更新失败',
							icon: 'none'
						});
					}
				} catch (error) {
					console.error('更新付款状态失败:', error);
					uni.showToast({
						title: '更新失败',
						icon: 'none'
					});
				} finally {
					uni.hideLoading();
				}
			},
			
			deleteOrder() {
				uni.showModal({
					title: '确认删除',
					content: '确定要删除这个订单吗？此操作不可恢复。',
					success: async (res) => {
						if (res.confirm) {
							try {
								uni.showLoading({ title: '删除中...' });
								
								const result = await api.order.delete(this.orderId);
								if (result.code === 0) {
									uni.showToast({
										title: '删除成功',
										icon: 'success'
									});
									
									setTimeout(() => {
										uni.navigateBack();
									}, 1500);
								} else {
									uni.showToast({
										title: result.message || '删除失败',
										icon: 'none'
									});
								}
							} catch (error) {
								console.error('删除订单失败:', error);
								uni.showToast({
									title: '删除失败',
									icon: 'none'
								});
							} finally {
								uni.hideLoading();
							}
						}
					}
				});
			},
			
			formatDateTime
		}
	}
</script>

<style scoped>
	.container {
		padding: 20rpx;
		background: #f5f5f5;
		min-height: 100vh;
	}
	
	.loading, .error {
		text-align: center;
		padding: 100rpx 0;
		color: #999;
		font-size: 28rpx;
	}
	
	.order-detail {
		/* 订单详情容器 */
	}
	
	.order-header {
		background: #fff;
		padding: 30rpx;
		border-radius: 15rpx;
		margin-bottom: 20rpx;
		display: flex;
		justify-content: space-between;
		align-items: center;
		box-shadow: 0 2rpx 10rpx rgba(0,0,0,0.1);
	}
	
	.order-id {
		/* 订单ID样式 */
	}
	
	.order-id .label {
		font-size: 26rpx;
		color: #666;
		margin-right: 10rpx;
	}
	
	.order-id .value {
		font-size: 32rpx;
		font-weight: bold;
		color: #333;
	}
	
	.order-status {
		padding: 12rpx 24rpx;
		border-radius: 20rpx;
		font-size: 26rpx;
		background: #fff2e8;
		color: #fa8c16;
		border: 2rpx solid #fa8c16;
	}
	
	.order-status.paid {
		background: #e8f5e8;
		color: #52c41a;
		border-color: #52c41a;
	}
	
	.section {
		background: #fff;
		padding: 30rpx;
		border-radius: 15rpx;
		margin-bottom: 20rpx;
		box-shadow: 0 2rpx 10rpx rgba(0,0,0,0.1);
	}
	
	.section-title {
		font-size: 32rpx;
		font-weight: bold;
		color: #333;
		margin-bottom: 20rpx;
		padding-bottom: 15rpx;
		border-bottom: 2rpx solid #f0f0f0;
	}
	
	.info-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 20rpx;
		padding: 15rpx 0;
		border-bottom: 1rpx solid #f5f5f5;
	}
	
	.info-row:last-child {
		margin-bottom: 0;
		border-bottom: none;
	}
	
	.info-row .label {
		font-size: 28rpx;
		color: #666;
		flex-shrink: 0;
		width: 160rpx;
	}
	
	.info-row .value {
		font-size: 28rpx;
		color: #333;
		flex: 1;
		text-align: right;
	}
	
	.items-list {
		/* 商品列表样式 */
	}
	
	.item-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 20rpx 0;
		border-bottom: 1rpx solid #f5f5f5;
	}
	
	.item-row:last-child {
		border-bottom: none;
	}
	
	.item-info {
		flex: 1;
	}
	
	.item-name {
		font-size: 28rpx;
		color: #333;
		font-weight: 500;
		display: block;
		margin-bottom: 8rpx;
	}
	
	.item-price {
		font-size: 24rpx;
		color: #999;
	}
	
	.item-total {
		font-size: 28rpx;
		color: #ff6b35;
		font-weight: bold;
		margin-left: 20rpx;
	}
	
	.total-section {
		margin-top: 20rpx;
		padding-top: 20rpx;
		border-top: 2rpx solid #f0f0f0;
	}
	
	.total-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}
	
	.total-label {
		font-size: 32rpx;
		font-weight: bold;
		color: #333;
	}
	
	.total-amount {
		font-size: 36rpx;
		font-weight: bold;
		color: #ff6b35;
	}
	
	.notes-text {
		font-size: 28rpx;
		color: #666;
		line-height: 1.6;
		background: #f9f9f9;
		padding: 20rpx;
		border-radius: 8rpx;
	}
	
	.actions {
		display: flex;
		gap: 20rpx;
		margin-top: 30rpx;
		padding: 0 20rpx;
	}
	
	.action-btn {
		flex: 1;
		padding: 24rpx 0;
		border-radius: 25rpx;
		font-size: 28rpx;
		font-weight: 500;
		border: none;
		color: #fff;
	}
	
	.edit-btn {
		background: #1890ff;
	}
	
	.edit-btn:active {
		background: #40a9ff;
	}
	
	.payment-btn {
		background: #52c41a;
	}
	
	.payment-btn.paid {
		background: #fa8c16;
	}
	
	.payment-btn:active {
		opacity: 0.8;
	}
	
	.delete-btn {
		background: #ff4d4f;
	}
	
	.delete-btn:active {
		background: #ff7875;
	}
</style>