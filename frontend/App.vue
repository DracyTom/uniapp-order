<template>
	<view id="app">
		<!-- 应用入口 -->
	</view>
</template>

<script>
	export default {
		onLaunch: function() {
			console.log('App Launch')
		},
		onShow: function() {
			console.log('App Show')
		},
		onHide: function() {
			console.log('App Hide')
		}
	}
</script>

<style>
	/*每个页面公共css */
	page {
		background-color: #f8f8f8;
	}
	
	.container {
		padding: 20rpx;
	}
	
	.card {
		background-color: #ffffff;
		border-radius: 10rpx;
		padding: 20rpx;
		margin-bottom: 20rpx;
		box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.1);
	}
	
	.btn-primary {
		background-color: #007aff;
		color: #ffffff;
		border-radius: 10rpx;
		padding: 20rpx;
		text-align: center;
		margin: 10rpx 0;
	}
	
	.btn-success {
		background-color: #4cd964;
		color: #ffffff;
		border-radius: 10rpx;
		padding: 20rpx;
		text-align: center;
		margin: 10rpx 0;
	}
	
	.text-center {
		text-align: center;
	}
	
	.text-muted {
		color: #999999;
	}
	
	.flex {
		display: flex;
	}
	
	.flex-between {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}
	
	.flex-center {
		display: flex;
		justify-content: center;
		align-items: center;
	}
</style>