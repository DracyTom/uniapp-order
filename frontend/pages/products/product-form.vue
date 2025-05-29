<template>
	<view class="container">
		<form @submit="submitProduct">
			<!-- 基本信息 -->
			<view class="card">
				<view class="card-title">基本信息</view>
				
				<view class="form-group">
					<text class="label">商品名称 *</text>
					<input class="form-input" v-model="formData.name" placeholder="请输入商品名称" />
				</view>
				
				<view class="form-group">
					<text class="label">商品分类 *</text>
					<input class="form-input" v-model="formData.category" placeholder="请输入商品分类" />
				</view>
				
				<view class="form-group">
					<text class="label">商品价格 *</text>
					<input class="form-input" type="digit" v-model="formData.price" placeholder="请输入商品价格" />
				</view>
				
				<view class="form-group">
					<text class="label">上架状态</text>
					<switch :checked="formData.status" @change="onStatusChange" color="#007aff" />
				</view>
			</view>
			
			<!-- 商品图片 -->
			<view class="card">
				<view class="card-title">商品图片</view>
				
				<view class="image-upload">
					<view v-if="formData.image" class="image-preview">
						<image :src="formData.image" mode="aspectFill"></image>
						<view class="image-actions">
							<button class="image-btn" @click="chooseImage" type="button">更换</button>
							<button class="image-btn delete" @click="removeImage" type="button">删除</button>
						</view>
					</view>
					<view v-else class="image-placeholder" @click="chooseImage">
						<text class="upload-icon">📷</text>
						<text class="upload-text">点击上传图片</text>
					</view>
				</view>
			</view>
			
			<!-- 自定义属性 -->
			<view class="card">
				<view class="card-title">
					<text>自定义属性</text>
					<button class="add-attr-btn" @click="addAttribute" type="button">添加属性</button>
				</view>
				
				<view class="attr-list">
					<view class="attr-item" v-for="(attr, index) in attributes" :key="index">
						<view class="attr-inputs">
							<input class="attr-key" v-model="attr.key" placeholder="属性名" />
							<input class="attr-value" v-model="attr.value" placeholder="属性值" />
							<button class="remove-attr-btn" @click="removeAttribute(index)" type="button">删除</button>
						</view>
					</view>
				</view>
				
				<view v-if="attributes.length === 0" class="empty-attrs">
					<text class="empty-text">暂无自定义属性</text>
				</view>
			</view>
			
			<!-- 提交按钮 -->
			<view class="submit-section">
				<button class="submit-btn" @click="submitProduct">{{ isEdit ? '更新商品' : '创建商品' }}</button>
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
				productId: null,
				formData: {
					name: '',
					category: '',
					price: '',
					status: true,
					image: '',
					custom_attrs: {}
				},
				attributes: []
			}
		},
		onLoad(options) {
			if (options.id) {
				this.isEdit = true
				this.productId = options.id
				this.loadProduct()
			}
		},
		methods: {
			async loadProduct() {
				try {
					uni.showLoading({ title: '加载中...' })
					const res = await api.products.get(this.productId)
					if (res.code === 0) {
						this.formData = {
							name: res.data.name || '',
							category: res.data.category || '',
							price: res.data.price?.toString() || '',
							status: res.data.status !== false,
							image: res.data.image || '',
							custom_attrs: res.data.custom_attrs || {}
						}
						
						// 转换自定义属性为数组格式
						this.attributes = Object.entries(this.formData.custom_attrs).map(([key, value]) => ({
							key,
							value
						}))
					}
				} catch (error) {
					console.error('加载商品失败:', error)
					uni.showToast({
						title: '加载失败',
						icon: 'none'
					})
				} finally {
					uni.hideLoading()
				}
			},
			
			onStatusChange(e) {
				this.formData.status = e.detail.value
			},
			
			chooseImage() {
				const that = this
				uni.chooseImage({
					count: 1,
					sizeType: ['compressed'],
					sourceType: ['album', 'camera'],
					success: function(res) {
						// 这里应该上传图片到服务器，暂时使用本地路径
						that.formData.image = res.tempFilePaths[0]
						
						// 实际项目中应该调用上传接口
						// that.uploadImage(res.tempFilePaths[0])
					},
					fail: function(error) {
						console.error('选择图片失败:', error)
						uni.showToast({
							title: '选择图片失败',
							icon: 'none'
						})
					}
				})
			},
			
			removeImage() {
				this.formData.image = ''
			},
			
			addAttribute() {
				this.attributes.push({
					key: '',
					value: ''
				})
			},
			
			removeAttribute(index) {
				this.attributes.splice(index, 1)
			},
			
			validateForm() {
				if (!this.formData.name.trim()) {
					uni.showToast({ title: '请输入商品名称', icon: 'none' })
					return false
				}
				if (!this.formData.category.trim()) {
					uni.showToast({ title: '请输入商品分类', icon: 'none' })
					return false
				}
				if (!this.formData.price || isNaN(parseFloat(this.formData.price))) {
					uni.showToast({ title: '请输入有效的商品价格', icon: 'none' })
					return false
				}
				
				// 验证自定义属性
				for (let i = 0; i < this.attributes.length; i++) {
					const attr = this.attributes[i]
					if (!attr.key.trim() || !attr.value.trim()) {
						uni.showToast({ title: `请完善第${i+1}个自定义属性`, icon: 'none' })
						return false
					}
				}
				
				return true
			},
			
			async submitProduct() {
				if (!this.validateForm()) {
					return
				}
				
				try {
					uni.showLoading({ title: this.isEdit ? '更新中...' : '创建中...' })
					
					// 准备提交数据
					const submitData = {
						name: this.formData.name.trim(),
						category: this.formData.category.trim(),
						price: parseFloat(this.formData.price),
						status: this.formData.status,
						image: this.formData.image,
						custom_attrs: {}
					}
					
					// 转换自定义属性为对象格式
					this.attributes.forEach(attr => {
						if (attr.key.trim() && attr.value.trim()) {
							submitData.custom_attrs[attr.key.trim()] = attr.value.trim()
						}
					})
					
					let res
					if (this.isEdit) {
						res = await api.products.update(this.productId, submitData)
					} else {
						res = await api.products.create(submitData)
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
					console.error('提交商品失败:', error)
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
	
	.add-attr-btn {
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
	
	.image-upload {
		margin-top: 20rpx;
	}
	
	.image-preview {
		position: relative;
		width: 200rpx;
		height: 200rpx;
		border-radius: 15rpx;
		overflow: hidden;
		margin: 0 auto;
	}
	
	.image-preview image {
		width: 100%;
		height: 100%;
	}
	
	.image-actions {
		position: absolute;
		bottom: 0;
		left: 0;
		right: 0;
		background: rgba(0, 0, 0, 0.7);
		display: flex;
		justify-content: center;
		gap: 20rpx;
		padding: 15rpx;
	}
	
	.image-btn {
		background-color: rgba(255, 255, 255, 0.9);
		color: #333;
		padding: 8rpx 16rpx;
		border-radius: 6rpx;
		font-size: 24rpx;
		border: none;
	}
	
	.image-btn.delete {
		background-color: rgba(255, 59, 48, 0.9);
		color: white;
	}
	
	.image-placeholder {
		width: 200rpx;
		height: 200rpx;
		border: 4rpx dashed #e0e0e0;
		border-radius: 15rpx;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		margin: 0 auto;
		cursor: pointer;
		background-color: #fafafa;
	}
	
	.upload-icon {
		font-size: 48rpx;
		margin-bottom: 10rpx;
	}
	
	.upload-text {
		font-size: 24rpx;
		color: #666;
	}
	
	.attr-list {
		margin-top: 20rpx;
	}
	
	.attr-item {
		margin-bottom: 20rpx;
	}
	
	.attr-inputs {
		display: flex;
		gap: 15rpx;
		align-items: center;
	}
	
	.attr-key, .attr-value {
		flex: 1;
		padding: 20rpx;
		border: 2rpx solid #e0e0e0;
		border-radius: 8rpx;
		font-size: 28rpx;
		background-color: #fff;
	}
	
	.remove-attr-btn {
		background-color: #ff3b30;
		color: white;
		padding: 20rpx 24rpx;
		border-radius: 8rpx;
		font-size: 24rpx;
		border: none;
		white-space: nowrap;
	}
	
	.empty-attrs {
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