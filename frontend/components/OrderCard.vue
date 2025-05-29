<template>
  <view class="order-card" @click="onCardClick">
    <view class="order-header">
      <text class="order-id">订单 #{{ order.id }}</text>
      <view class="order-status" :class="{ 'paid': order.payment_status === 1 }">
        {{ order.payment_status === 1 ? '已付款' : '未付款' }}
      </view>
    </view>
    
    <view class="order-customer">
      <text class="customer-name">{{ order.customer_name }}</text>
      <text class="delivery-time">{{ formatDate(order.delivery_time) }}</text>
    </view>
    
    <view class="order-address">
      <text class="address-text">{{ order.delivery_address }}</text>
    </view>
    
    <view class="order-items">
      <view v-for="item in order.items" :key="item.id" class="order-item">
        <text class="item-name">{{ item.product_name }}</text>
        <text class="item-quantity">x{{ item.quantity }}</text>
        <text class="item-price">¥{{ (item.price * item.quantity).toFixed(2) }}</text>
      </view>
    </view>
    
    <view class="order-footer">
      <text class="total-price">总计: ¥{{ order.total_amount }}</text>
    </view>
    
    <view v-if="order.notes" class="order-notes">
      <text class="notes-text">备注: {{ order.notes }}</text>
    </view>
  </view>
</template>

<script>
export default {
  name: 'OrderCard',
  props: {
    order: {
      type: Object,
      required: true
    }
  },
  methods: {
    onCardClick() {
      this.$emit('click', this.order);
    },
    formatDate(dateStr) {
      if (!dateStr) return '';
      const date = new Date(dateStr);
      return `${date.getMonth() + 1}/${date.getDate()} ${date.getHours()}:${date.getMinutes().toString().padStart(2, '0')}`;
    }
  }
}
</script>

<style scoped>
.order-card {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  margin-bottom: 16px;
  padding: 16px;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.order-id {
  font-size: 16px;
  font-weight: bold;
  color: #333;
}

.order-status {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  background: #fff2e8;
  color: #fa8c16;
}

.order-status.paid {
  background: #e8f5e8;
  color: #52c41a;
}

.order-customer {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.customer-name {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.delivery-time {
  font-size: 12px;
  color: #666;
}

.order-address {
  margin-bottom: 12px;
}

.address-text {
  font-size: 14px;
  color: #666;
}

.order-items {
  border-top: 1px solid #f0f0f0;
  padding-top: 12px;
  margin-bottom: 12px;
}

.order-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.item-name {
  font-size: 14px;
  color: #333;
  flex: 1;
}

.item-quantity {
  font-size: 12px;
  color: #666;
  margin: 0 8px;
}

.item-price {
  font-size: 14px;
  color: #ff6b35;
  font-weight: 500;
}

.order-footer {
  border-top: 1px solid #f0f0f0;
  padding-top: 12px;
  text-align: right;
}

.total-price {
  font-size: 16px;
  color: #ff6b35;
  font-weight: bold;
}

.order-notes {
  margin-top: 8px;
  padding: 8px;
  background: #f9f9f9;
  border-radius: 4px;
}

.notes-text {
  font-size: 12px;
  color: #666;
}
</style>