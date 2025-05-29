<template>
  <view class="product-card" @click="onCardClick">
    <image class="product-image" :src="product.image || '/static/default-product.png'" mode="aspectFill"></image>
    <view class="product-info">
      <text class="product-name">{{ product.name }}</text>
      <text class="product-category">{{ product.category }}</text>
      <view class="product-price-row">
        <text class="product-price">¥{{ product.price }}</text>
        <view class="product-status" :class="{ 'active': product.status === 1 }">
          {{ product.status === 1 ? '上架' : '下架' }}
        </view>
      </view>
    </view>
  </view>
</template>

<script>
export default {
  name: 'ProductCard',
  props: {
    product: {
      type: Object,
      required: true
    }
  },
  methods: {
    onCardClick() {
      this.$emit('click', this.product);
    }
  }
}
</script>

<style scoped>
.product-card {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  margin-bottom: 16px;
  overflow: hidden;
}

.product-image {
  width: 100%;
  height: 200px;
}

.product-info {
  padding: 12px;
}

.product-name {
  font-size: 16px;
  font-weight: bold;
  color: #333;
  display: block;
  margin-bottom: 4px;
}

.product-category {
  font-size: 12px;
  color: #999;
  display: block;
  margin-bottom: 8px;
}

.product-price-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.product-price {
  font-size: 18px;
  color: #ff6b35;
  font-weight: bold;
}

.product-status {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  background: #f5f5f5;
  color: #999;
}

.product-status.active {
  background: #e8f5e8;
  color: #52c41a;
}
</style>