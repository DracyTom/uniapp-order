<template>
	<view class="container">
		<!-- 搜索和筛选 -->
		<view class="search-bar">
			<input class="search-input" placeholder="搜索商品..." v-model="searchText" @input="onSearch" />
			<picker @change="onCategoryChange" :value="categoryIndex" :range="categoryOptions">
				<view class="filter-btn">{{ categoryOptions[categoryIndex] }}</view>
			</picker>
		</view>
		
		<!-- 商品列表 -->
		<view class="product-grid">
			<view class="product-card" v-for="product in filteredProducts" :key="product.id" @click="viewProduct(product)">
				<view class="product-image">
					<image v-if="product.image" :src="product.image" mode="aspectFill"></image>
					<view v-else class="placeholder-image">
						<text class="placeholder-text">{{ product.name.charAt(0) }}</text>
					</view>
					<view class="status-badge" :class="product.status ? 'status-active' : 'status-inactive'">
						{{ product.status ? '上架' : '下架' }}
					</view>
				</view>
				
				<view class="product-info">
					<text class="product-name">{{ product.name }}</text>
					<text class="product-category">{{ product.category }}</text>
					<view class="product-price">
						<text class="price-symbol">¥</text>
						<text class="price-value">{{ product.price }}</text>
					</view>
					
					<!-- 自定义属性 -->
					<view v-if="product.custom_attrs" class="custom-attrs">
						<view class="attr-item" v-for="(value, key) in product.custom_attrs" :key="key">
							<text class="attr-label">{{ key }}:</text>
							<text class="attr-value">{{ value }}</text>
						</view>
					</view>
				</view>
			</view>
		</view>
		
		<!-- 空状态 -->
		<view v-if="filteredProducts.length === 0" class="empty-state">
			<text class="empty-text">暂无商品数据</text>
		</view>
		
		<!-- 添加按钮 -->
		<view class="fab" @click="addProduct">
			<text class="fab-icon">+</text>
		</view>
		
		<!-- 商品详情弹窗 -->
		<uni-popup ref="productPopup" type="bottom">
			<view class="popup-content">
				<view class="popup-header">
					<text class="popup-title">商品详情</text>
					<text class="close-btn" @click="closePopup">×</text>
				</view>
				
				<view v-if="selectedProduct" class="product-detail">
					<view class="detail-image">
						<image v-if="selectedProduct.image" :src="selectedProduct.image" mode="aspectFill"></image>
						<view v-else class="detail-placeholder">
							<text class="placeholder-text">{{ selectedProduct.name.charAt(0) }}</text>
						</view>
					</view>
					
					<view class="detail-info">
						<view class="detail-row">
							<text class="detail-label">商品名称：</text>
							<text class="detail-value">{{ selectedProduct.name }}</text>
						</view>
						<view class="detail-row">
							<text class="detail-label">商品分类：</text>
							<text class="detail-value">{{ selectedProduct.category }}</text>
						</view>
						<view class="detail-row">
							<text class="detail-label">商品价格：</text>
							<text class="detail-value price">¥{{ selectedProduct.price }}</text>
						</view>
						<view class="detail-row">
							<text class="detail-label">上架状态：</text>
							<text class="detail-value" :class="selectedProduct.status ? 'status-active' : 'status-inactive'">
								{{ selectedProduct.status ? '已上架' : '已下架' }}
							</text>
						</view>
						
						<!-- 自定义属性 -->
						<view v-if="selectedProduct.custom_attrs" class="custom-attrs-detail">
							<text class="detail-label">自定义属性：</text>
							<view class="attrs-list">
								<view class="attr-detail" v-for="(value, key) in selectedProduct.custom_attrs" :key="key">
									<text class="attr-key">{{ key }}:</text>
									<text class="attr-val">{{ value }}</text>
								</view>
							</view>
						</view>
						
						<view class="detail-row">
							<text class="detail-label">创建时间：</text>
							<text class="detail-value">{{ formatDateTime(selectedProduct.created_at) }}</text>
						</view>
					</view>
				</view>
				
				<view class="popup-actions">
					<button class="action-btn edit-btn" @click="editProduct">编辑</button>
					<button class="action-btn delete-btn" @click="deleteProduct">删除</button>
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
				products: [],
				filteredProducts: [],
				searchText: '',
				categoryIndex: 0,
				categoryOptions: ['全部分类'],
				selectedProduct: null
			}
		},
		onLoad() {
			this.loadProducts()
		},
		onShow() {
			this.loadProducts()
		},
		methods: {
			async loadProducts() {
				try {
					uni.showLoading({ title: '加载中...' })
					const res = await api.products.list()
					if (res.code === 0) {
						this.products = Array.isArray(res.data) ? res.data : [res.data].filter(Boolean)
						this.updateCategories()
						this.filterProducts()
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
			
			updateCategories() {
				const categories = new Set(['全部分类'])
				this.products.forEach(product => {
					if (product.category) {
						categories.add(product.category)
					}
				})
				this.categoryOptions = Array.from(categories)
			},
			
			filterProducts() {
				let filtered = this.products
				
				// 按分类筛选
				if (this.categoryIndex > 0) {
					const selectedCategory = this.categoryOptions[this.categoryIndex]
					filtered = filtered.filter(product => product.category === selectedCategory)
				}
				
				// 按搜索文本筛选
				if (this.searchText) {
					const text = this.searchText.toLowerCase()
					filtered = filtered.filter(product => 
						product.name?.toLowerCase().includes(text) ||
						product.category?.toLowerCase().includes(text)
					)
				}
				
				this.filteredProducts = filtered
			},
			
			onSearch() {
				this.filterProducts()
			},
			
			onCategoryChange(e) {
				this.categoryIndex = e.detail.value
				this.filterProducts()
			},
			
			viewProduct(product) {
				this.selectedProduct = product
				this.$refs.productPopup.open()
			},
			
			closePopup() {
				this.$refs.productPopup.close()
			},
			
			addProduct() {
				uni.navigateTo({
					url: '/pages/products/product-form'
				})
			},
			
			editProduct() {
				this.closePopup()
				uni.navigateTo({
					url: `/pages/products/product-form?id=${this.selectedProduct.id}`
				})
			},
			
			async deleteProduct() {
				const that = this
				uni.showModal({
					title: '确认删除',
					content: '确定要删除这个商品吗？',
					success: async function(res) {
						if (res.confirm) {
							try {
								await api.products.delete(that.selectedProduct.id)
								uni.showToast({
									title: '删除成功',
									icon: 'success'
								})
								that.closePopup()
								that.loadProducts()
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
	
	.product-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 20rpx;
		margin-bottom: 100rpx;
	}
	
	.product-card {
		background-color: white;
		border-radius: 15rpx;
		overflow: hidden;
		box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.1);
		cursor: pointer;
	}
	
	.product-image {
		position: relative;
		height: 200rpx;
		background-color: #f5f5f5;
	}
	
	.product-image image {
		width: 100%;
		height: 100%;
	}
	
	.placeholder-image {
		width: 100%;
		height: 100%;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.placeholder-text {
		color: white;
		font-size: 48rpx;
		font-weight: bold;
	}
	
	.status-badge {
		position: absolute;
		top: 10rpx;
		right: 10rpx;
		padding: 6rpx 12rpx;
		border-radius: 20rpx;
		font-size: 20rpx;
		color: white;
	}
	
	.status-active {
		background-color: #4caf50;
	}
	
	.status-inactive {
		background-color: #ff9800;
	}
	
	.product-info {
		padding: 20rpx;
	}
	
	.product-name {
		display: block;
		font-size: 28rpx;
		font-weight: bold;
		color: #333;
		margin-bottom: 8rpx;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	
	.product-category {
		display: block;
		font-size: 24rpx;
		color: #666;
		margin-bottom: 12rpx;
	}
	
	.product-price {
		display: flex;
		align-items: baseline;
		margin-bottom: 12rpx;
	}
	
	.price-symbol {
		font-size: 24rpx;
		color: #e91e63;
	}
	
	.price-value {
		font-size: 32rpx;
		font-weight: bold;
		color: #e91e63;
	}
	
	.custom-attrs {
		margin-top: 12rpx;
	}
	
	.attr-item {
		display: flex;
		margin-bottom: 6rpx;
	}
	
	.attr-label {
		font-size: 20rpx;
		color: #666;
		margin-right: 8rpx;
	}
	
	.attr-value {
		font-size: 20rpx;
		color: #333;
		flex: 1;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
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
	
	.product-detail {
		padding: 30rpx;
	}
	
	.detail-image {
		width: 200rpx;
		height: 200rpx;
		margin: 0 auto 30rpx;
		border-radius: 15rpx;
		overflow: hidden;
		background-color: #f5f5f5;
	}
	
	.detail-image image {
		width: 100%;
		height: 100%;
	}
	
	.detail-placeholder {
		width: 100%;
		height: 100%;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.detail-info {
		margin-top: 30rpx;
	}
	
	.detail-row {
		display: flex;
		margin-bottom: 20rpx;
		align-items: flex-start;
	}
	
	.detail-label {
		color: #666;
		width: 160rpx;
		font-size: 28rpx;
		flex-shrink: 0;
	}
	
	.detail-value {
		flex: 1;
		color: #333;
		font-size: 28rpx;
	}
	
	.detail-value.price {
		color: #e91e63;
		font-weight: bold;
		font-size: 32rpx;
	}
	
	.custom-attrs-detail {
		display: flex;
		margin-bottom: 20rpx;
		align-items: flex-start;
	}
	
	.attrs-list {
		flex: 1;
	}
	
	.attr-detail {
		display: flex;
		margin-bottom: 10rpx;
	}
	
	.attr-key {
		color: #666;
		font-size: 26rpx;
		margin-right: 12rpx;
		min-width: 120rpx;
	}
	
	.attr-val {
		color: #333;
		font-size: 26rpx;
		flex: 1;
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