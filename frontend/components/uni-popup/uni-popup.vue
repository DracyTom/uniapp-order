<template>
	<view v-if="showPopup" class="uni-popup" :class="[type]" @click="maskClick">
		<view class="uni-popup-content" @click.stop>
			<slot></slot>
		</view>
	</view>
</template>

<script>
	export default {
		name: 'uni-popup',
		props: {
			type: {
				type: String,
				default: 'center'
			},
			maskClick: {
				type: Boolean,
				default: true
			}
		},
		data() {
			return {
				showPopup: false
			}
		},
		methods: {
			open() {
				this.showPopup = true
			},
			close() {
				this.showPopup = false
			},
			maskClick() {
				if (this.maskClick) {
					this.close()
				}
			}
		}
	}
</script>

<style scoped>
	.uni-popup {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: rgba(0, 0, 0, 0.5);
		z-index: 999;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.uni-popup.bottom {
		align-items: flex-end;
		justify-content: center;
	}
	
	.uni-popup.top {
		align-items: flex-start;
		justify-content: center;
	}
	
	.uni-popup-content {
		background-color: white;
		border-radius: 10rpx;
		max-width: 90%;
		max-height: 90%;
		overflow: hidden;
	}
	
	.uni-popup.bottom .uni-popup-content {
		width: 100%;
		border-radius: 20rpx 20rpx 0 0;
	}
	
	.uni-popup.top .uni-popup-content {
		width: 100%;
		border-radius: 0 0 20rpx 20rpx;
	}
</style>