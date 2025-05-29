<template>
	<view class="container">
		<form @submit="submitCustomer">
			<!-- 基本信息 -->
			<view class="card">
				<view class="card-title">基本信息</view>
				
				<view class="form-group">
					<text class="label">客户姓名 *</text>
					<input class="form-input" v-model="formData.name" placeholder="请输入客户姓名" />
				</view>
				
				<view class="form-group">
					<text class="label">联系电话</text>
					<input class="form-input" v-model="formData.phone" placeholder="请输入联系电话" />
				</view>
			</view>
			
			<!-- 地址信息 -->
			<view class="card">
				<view class="card-title">
					<text>地址信息</text>
					<button class="add-address-btn" @click="addAddress" type="button">添加地址</button>
				</view>
				
				<view class="address-list">
					<view class="address-item" v-for="(address, index) in formData.addresses" :key="index">
						<view class="address-header">
							<input class="address-label-input" v-model="address.label" placeholder="地址标签（如：家、公司）" />
							<button class="remove-address-btn" @click="removeAddress(index)" type="button">删除</button>
						</view>
						<textarea class="address-input" v-model="address.address" placeholder="请输入详细地址"></textarea>
					</view>
				</view>
				
				<view v-if="formData.addresses.length === 0" class="empty-addresses">
					<text class="empty-text">暂无地址信息</text>
				</view>
			</view>
			
			<!-- 提交按钮 -->
			<view class="submit-section">
				<button class="submit-btn" @click="submitCustomer">{{ isEdit ? '更新客户' : '创建客户' }}</button>
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
				customerId: null,
				formData: {
					name: '',
					phone: '',
					addresses: []
				}
			}
		},
		onLoad(options) {
			if (options.id) {
				this.isEdit = true
				this.customerId = options.id
				this.loadCustomer()
			} else {
				// 默认添加一个地址
				this.addAddress()
			}
		},
		methods: {
			async loadCustomer() {
				try {
					uni.showLoading({ title: '加载中...' })
					const res = await api.customers.get(this.customerId)
					if (res.code === 0) {
						const customerData = res.data.customer || res.data
						this.formData = {
							name: customerData.name || '',
							phone: customerData.phone || '',
							addresses: res.data.addresses || []
						}
						
						// 如果没有地址，添加一个空地址
						if (this.formData.addresses.length === 0) {
							this.addAddress()
						}
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
			
			addAddress() {
				this.formData.addresses.push({
					label: '',
					address: ''
				})
			},
			
			removeAddress(index) {
				this.formData.addresses.splice(index, 1)
			},
			
			validateForm() {
				if (!this.formData.name.trim()) {
					uni.showToast({ title: '请输入客户姓名', icon: 'none' })
					return false
				}
				
				// 验证地址信息
				for (let i = 0; i < this.formData.addresses.length; i++) {
					const address = this.formData.addresses[i]
					if (!address.address.trim()) {
						uni.showToast({ title: `请填写第${i+1}个地址的详细信息`, icon: 'none' })
						return false
					}
				}
				
				return true
			},
			
			async submitCustomer() {
				if (!this.validateForm()) {
					return
				}
				
				try {
					uni.showLoading({ title: this.isEdit ? '更新中...' : '创建中...' })
					
					// 准备提交数据
					const submitData = {
						name: this.formData.name.trim(),
						phone: this.formData.phone.trim(),
						addresses: this.formData.addresses.filter(addr => addr.address.trim()).map(addr => ({
							label: addr.label.trim() || '默认地址',
							address: addr.address.trim()
						}))
					}
					
					let res
					if (this.isEdit) {
						res = await api.customers.update(this.customerId, submitData)
					} else {
						res = await api.customers.create(submitData)
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
					console.error('提交客户失败:', error)
					uni.showToast({
						title: error.message || '操作失败',
						icon: 'none'
					})
				} finally {
					uni.hideLoading()
				}
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
	
	.add-address-btn {
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
	
	.address-list {
		margin-top: 20rpx;
	}
	
	.address-item {
		border: 2rpx solid #e0e0e0;
		border-radius: 10rpx;
		padding: 20rpx;
		margin-bottom: 20rpx;
		background-color: #fafafa;
	}
	
	.address-header {
		display: flex;
		gap: 15rpx;
		align-items: center;
		margin-bottom: 20rpx;
	}
	
	.address-label-input {
		flex: 1;
		padding: 20rpx;
		border: 2rpx solid #e0e0e0;
		border-radius: 8rpx;
		background-color: #fff;
		font-size: 28rpx;
	}
	
	.remove-address-btn {
		background-color: #ff3b30;
		color: white;
		padding: 20rpx 30rpx;
		border-radius: 8rpx;
		font-size: 24rpx;
		border: none;
		white-space: nowrap;
	}
	
	.address-input {
		width: 100%;
		padding: 20rpx;
		border: 2rpx solid #e0e0e0;
		border-radius: 8rpx;
		background-color: #fff;
		font-size: 28rpx;
		min-height: 120rpx;
	}
	
	.empty-addresses {
		text-align: center;
		padding: 60rpx 0;
		color: #999;
	}
	
	.empty-text {
		font-size: 28rpx;
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