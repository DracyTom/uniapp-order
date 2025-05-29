<template>
  <view class="customer-card" @click="onCardClick">
    <view class="customer-header">
      <text class="customer-name">{{ customer.name }}</text>
      <text class="customer-phone">{{ customer.phone }}</text>
    </view>
    <view class="customer-addresses">
      <view v-for="(address, index) in customer.addresses" :key="index" class="address-item">
        <text class="address-text">{{ address.address }}</text>
        <view v-if="address.is_default" class="default-tag">默认</view>
      </view>
    </view>
  </view>
</template>

<script>
export default {
  name: 'CustomerCard',
  props: {
    customer: {
      type: Object,
      required: true
    }
  },
  methods: {
    onCardClick() {
      this.$emit('click', this.customer);
    }
  }
}
</script>

<style scoped>
.customer-card {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  margin-bottom: 16px;
  padding: 16px;
}

.customer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.customer-name {
  font-size: 16px;
  font-weight: bold;
  color: #333;
}

.customer-phone {
  font-size: 14px;
  color: #666;
}

.customer-addresses {
  border-top: 1px solid #f0f0f0;
  padding-top: 12px;
}

.address-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.address-text {
  font-size: 14px;
  color: #666;
  flex: 1;
}

.default-tag {
  background: #1890ff;
  color: #fff;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 10px;
}
</style>