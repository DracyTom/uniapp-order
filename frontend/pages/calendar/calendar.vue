<template>
	<view class="container">
		<!-- 日历头部 -->
		<view class="calendar-header">
			<view class="header-controls">
				<button class="nav-btn" @click="prevMonth">‹</button>
				<text class="current-month">{{ currentYear }}年{{ currentMonth }}月</text>
				<button class="nav-btn" @click="nextMonth">›</button>
			</view>
			
			<view class="view-switcher">
				<button 
					class="view-btn" 
					:class="{ active: currentView === 'month' }" 
					@click="switchView('month')"
				>月</button>
				<button 
					class="view-btn" 
					:class="{ active: currentView === 'week' }" 
					@click="switchView('week')"
				>周</button>
				<button 
					class="view-btn" 
					:class="{ active: currentView === 'year' }" 
					@click="switchView('year')"
				>年</button>
			</view>
		</view>
		
		<!-- 月视图 -->
		<view v-if="currentView === 'month'" class="month-view">
			<!-- 星期标题 -->
			<view class="week-header">
				<text class="week-day" v-for="day in weekDays" :key="day">{{ day }}</text>
			</view>
			
			<!-- 日期网格 -->
			<view class="date-grid">
				<view 
					class="date-cell" 
					v-for="date in monthDates" 
					:key="date.key"
					:class="{ 
						'other-month': !date.isCurrentMonth,
						'today': date.isToday,
						'has-orders': date.orderCount > 0
					}"
					@click="selectDate(date)"
				>
					<view class="date-number">{{ date.day }}</view>
					<view class="lunar-date">{{ date.lunar }}</view>
					<view v-if="date.holiday" class="holiday">{{ date.holiday }}</view>
					<view v-if="date.orderCount > 0" class="order-count">{{ date.orderCount }}单</view>
				</view>
			</view>
		</view>
		
		<!-- 周视图 -->
		<view v-if="currentView === 'week'" class="week-view">
			<view class="week-header">
				<text class="week-day" v-for="day in weekDays" :key="day">{{ day }}</text>
			</view>
			
			<view class="week-dates">
				<view 
					class="week-date-cell" 
					v-for="date in weekDates" 
					:key="date.key"
					:class="{ 
						'today': date.isToday,
						'has-orders': date.orderCount > 0
					}"
					@click="selectDate(date)"
				>
					<view class="date-number">{{ date.day }}</view>
					<view class="lunar-date">{{ date.lunar }}</view>
					<view v-if="date.holiday" class="holiday">{{ date.holiday }}</view>
					<view v-if="date.orderCount > 0" class="order-detail">
						<text class="order-count">{{ date.orderCount }}单</text>
						<view class="order-products" v-if="date.products">
							<text class="product-item" v-for="product in date.products" :key="product.name">
								{{ product.name }}×{{ product.quantity }}
							</text>
						</view>
					</view>
				</view>
			</view>
		</view>
		
		<!-- 年视图 -->
		<view v-if="currentView === 'year'" class="year-view">
			<view class="year-grid">
				<view class="month-card" v-for="month in yearMonths" :key="month.month" @click="selectMonth(month)">
					<view class="month-title">{{ month.month }}月</view>
					<view class="month-stats">
						<text class="stat-item">{{ month.orderCount }}单</text>
						<text class="stat-item">¥{{ month.totalAmount }}</text>
					</view>
				</view>
			</view>
		</view>
		
		<!-- 选中日期详情弹窗 -->
		<uni-popup ref="datePopup" type="bottom">
			<view class="popup-content">
				<view class="popup-header">
					<text class="popup-title">{{ selectedDateInfo.dateStr }} 订单详情</text>
					<text class="close-btn" @click="closeDatePopup">×</text>
				</view>
				
				<view v-if="selectedDateInfo.orders && selectedDateInfo.orders.length > 0" class="date-orders">
					<view class="order-summary">
						<text class="summary-text">共 {{ selectedDateInfo.orders.length }} 个订单，总金额 ¥{{ selectedDateInfo.totalAmount }}</text>
					</view>
					
					<view class="order-list">
						<view class="order-item" v-for="order in selectedDateInfo.orders" :key="order.id">
							<view class="order-header">
								<text class="order-id">#{{ order.id }}</text>
								<text class="order-amount">¥{{ order.total_amount }}</text>
							</view>
							<view class="order-info">
								<text class="customer-name">{{ order.customer_name }}</text>
								<text class="recipient">{{ order.recipient }}</text>
							</view>
						</view>
					</view>
				</view>
				
				<view v-else class="no-orders">
					<text class="no-orders-text">当天暂无订单</text>
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
				currentView: 'month', // month, week, year
				currentYear: new Date().getFullYear(),
				currentMonth: new Date().getMonth() + 1,
				currentWeek: 0,
				weekDays: ['日', '一', '二', '三', '四', '五', '六'],
				monthDates: [],
				weekDates: [],
				yearMonths: [],
				calendarData: {},
				selectedDateInfo: {
					dateStr: '',
					orders: [],
					totalAmount: 0
				}
			}
		},
		onLoad() {
			this.loadCalendarData()
		},
		methods: {
			async loadCalendarData() {
				try {
					uni.showLoading({ title: '加载中...' })
					const res = await api.orders.calendar(this.currentYear, this.currentMonth)
					if (res.code === 0) {
						this.calendarData = res.data || {}
						this.generateCalendar()
					}
				} catch (error) {
					console.error('加载日历数据失败:', error)
					uni.showToast({
						title: '加载失败',
						icon: 'none'
					})
				} finally {
					uni.hideLoading()
				}
			},
			
			generateCalendar() {
				if (this.currentView === 'month') {
					this.generateMonthView()
				} else if (this.currentView === 'week') {
					this.generateWeekView()
				} else if (this.currentView === 'year') {
					this.generateYearView()
				}
			},
			
			generateMonthView() {
				const year = this.currentYear
				const month = this.currentMonth
				const firstDay = new Date(year, month - 1, 1)
				const lastDay = new Date(year, month, 0)
				const startDate = new Date(firstDay)
				startDate.setDate(startDate.getDate() - firstDay.getDay())
				
				const dates = []
				const today = new Date()
				
				for (let i = 0; i < 42; i++) {
					const date = new Date(startDate)
					date.setDate(startDate.getDate() + i)
					
					const dateStr = this.formatDate(date)
					const dayData = this.calendarData[dateStr] || {}
					
					dates.push({
						key: dateStr,
						day: date.getDate(),
						date: date,
						dateStr: dateStr,
						isCurrentMonth: date.getMonth() === month - 1,
						isToday: this.isSameDay(date, today),
						lunar: this.getLunarDate(date),
						holiday: this.getHoliday(date),
						orderCount: dayData.order_count || 0,
						orders: dayData.orders || [],
						totalAmount: dayData.total_amount || 0,
						products: dayData.products || []
					})
				}
				
				this.monthDates = dates
			},
			
			generateWeekView() {
				const today = new Date()
				const startOfWeek = new Date(today)
				startOfWeek.setDate(today.getDate() - today.getDay())
				
				const dates = []
				
				for (let i = 0; i < 7; i++) {
					const date = new Date(startOfWeek)
					date.setDate(startOfWeek.getDate() + i)
					
					const dateStr = this.formatDate(date)
					const dayData = this.calendarData[dateStr] || {}
					
					dates.push({
						key: dateStr,
						day: date.getDate(),
						date: date,
						dateStr: dateStr,
						isToday: this.isSameDay(date, today),
						lunar: this.getLunarDate(date),
						holiday: this.getHoliday(date),
						orderCount: dayData.order_count || 0,
						orders: dayData.orders || [],
						totalAmount: dayData.total_amount || 0,
						products: dayData.products || []
					})
				}
				
				this.weekDates = dates
			},
			
			generateYearView() {
				const months = []
				
				for (let month = 1; month <= 12; month++) {
					const monthData = this.getMonthData(this.currentYear, month)
					months.push({
						month: month,
						orderCount: monthData.orderCount,
						totalAmount: monthData.totalAmount
					})
				}
				
				this.yearMonths = months
			},
			
			getMonthData(year, month) {
				let orderCount = 0
				let totalAmount = 0
				
				Object.keys(this.calendarData).forEach(dateStr => {
					const date = new Date(dateStr)
					if (date.getFullYear() === year && date.getMonth() + 1 === month) {
						const dayData = this.calendarData[dateStr]
						orderCount += dayData.order_count || 0
						totalAmount += dayData.total_amount || 0
					}
				})
				
				return { orderCount, totalAmount: totalAmount.toFixed(2) }
			},
			
			switchView(view) {
				this.currentView = view
				this.generateCalendar()
			},
			
			prevMonth() {
				if (this.currentMonth === 1) {
					this.currentYear--
					this.currentMonth = 12
				} else {
					this.currentMonth--
				}
				this.loadCalendarData()
			},
			
			nextMonth() {
				if (this.currentMonth === 12) {
					this.currentYear++
					this.currentMonth = 1
				} else {
					this.currentMonth++
				}
				this.loadCalendarData()
			},
			
			selectDate(date) {
				this.selectedDateInfo = {
					dateStr: this.formatDateStr(date.date),
					orders: date.orders || [],
					totalAmount: date.totalAmount || 0
				}
				this.$refs.datePopup.open()
			},
			
			selectMonth(month) {
				this.currentMonth = month.month
				this.switchView('month')
			},
			
			closeDatePopup() {
				this.$refs.datePopup.close()
			},
			
			formatDate(date) {
				return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
			},
			
			formatDateStr(date) {
				return `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日`
			},
			
			isSameDay(date1, date2) {
				return date1.getFullYear() === date2.getFullYear() &&
					   date1.getMonth() === date2.getMonth() &&
					   date1.getDate() === date2.getDate()
			},
			
			getLunarDate(date) {
				// 简化的农历转换，实际项目中应使用专业的农历库
				const lunarMonths = ['正', '二', '三', '四', '五', '六', '七', '八', '九', '十', '冬', '腊']
				const lunarDays = ['初一', '初二', '初三', '初四', '初五', '初六', '初七', '初八', '初九', '初十',
								  '十一', '十二', '十三', '十四', '十五', '十六', '十七', '十八', '十九', '二十',
								  '廿一', '廿二', '廿三', '廿四', '廿五', '廿六', '廿七', '廿八', '廿九', '三十']
				
				// 这里使用简单的模拟数据，实际应该调用农历转换API
				const day = date.getDate()
				const month = date.getMonth()
				return lunarDays[(day - 1) % 30] || '初一'
			},
			
			getHoliday(date) {
				// 简化的节假日判断
				const month = date.getMonth() + 1
				const day = date.getDate()
				
				const holidays = {
					'1-1': '元旦',
					'2-14': '情人节',
					'3-8': '妇女节',
					'5-1': '劳动节',
					'6-1': '儿童节',
					'10-1': '国庆节',
					'12-25': '圣诞节'
				}
				
				return holidays[`${month}-${day}`] || ''
			}
		}
	}
</script>

<style scoped>
	.calendar-header {
		margin-bottom: 30rpx;
	}
	
	.header-controls {
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 20rpx;
	}
	
	.nav-btn {
		background-color: #007aff;
		color: white;
		border: none;
		border-radius: 50%;
		width: 60rpx;
		height: 60rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 32rpx;
		margin: 0 30rpx;
	}
	
	.current-month {
		font-size: 36rpx;
		font-weight: bold;
		color: #333;
		min-width: 200rpx;
		text-align: center;
	}
	
	.view-switcher {
		display: flex;
		justify-content: center;
		gap: 20rpx;
	}
	
	.view-btn {
		padding: 16rpx 32rpx;
		border: 2rpx solid #007aff;
		border-radius: 25rpx;
		background-color: white;
		color: #007aff;
		font-size: 28rpx;
	}
	
	.view-btn.active {
		background-color: #007aff;
		color: white;
	}
	
	.week-header {
		display: grid;
		grid-template-columns: repeat(7, 1fr);
		gap: 2rpx;
		margin-bottom: 20rpx;
	}
	
	.week-day {
		text-align: center;
		padding: 20rpx 0;
		font-size: 28rpx;
		color: #666;
		font-weight: bold;
	}
	
	.date-grid {
		display: grid;
		grid-template-columns: repeat(7, 1fr);
		gap: 2rpx;
	}
	
	.date-cell {
		aspect-ratio: 1;
		background-color: white;
		border-radius: 8rpx;
		padding: 10rpx;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: flex-start;
		cursor: pointer;
		position: relative;
	}
	
	.date-cell.other-month {
		opacity: 0.3;
	}
	
	.date-cell.today {
		background-color: #007aff;
		color: white;
	}
	
	.date-cell.has-orders {
		border: 2rpx solid #4caf50;
	}
	
	.date-number {
		font-size: 28rpx;
		font-weight: bold;
		margin-bottom: 4rpx;
	}
	
	.lunar-date {
		font-size: 20rpx;
		color: #666;
		margin-bottom: 4rpx;
	}
	
	.date-cell.today .lunar-date {
		color: rgba(255, 255, 255, 0.8);
	}
	
	.holiday {
		font-size: 18rpx;
		color: #ff3b30;
		margin-bottom: 4rpx;
	}
	
	.order-count {
		font-size: 18rpx;
		background-color: #4caf50;
		color: white;
		padding: 2rpx 8rpx;
		border-radius: 10rpx;
	}
	
	.week-dates {
		display: grid;
		grid-template-columns: repeat(7, 1fr);
		gap: 20rpx;
	}
	
	.week-date-cell {
		background-color: white;
		border-radius: 15rpx;
		padding: 20rpx;
		text-align: center;
		cursor: pointer;
		box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.1);
	}
	
	.week-date-cell.today {
		background-color: #007aff;
		color: white;
	}
	
	.week-date-cell.has-orders {
		border: 2rpx solid #4caf50;
	}
	
	.order-detail {
		margin-top: 10rpx;
	}
	
	.order-products {
		margin-top: 8rpx;
	}
	
	.product-item {
		display: block;
		font-size: 20rpx;
		color: #666;
		margin-bottom: 4rpx;
	}
	
	.week-date-cell.today .product-item {
		color: rgba(255, 255, 255, 0.8);
	}
	
	.year-grid {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: 20rpx;
	}
	
	.month-card {
		background-color: white;
		border-radius: 15rpx;
		padding: 30rpx;
		text-align: center;
		cursor: pointer;
		box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.1);
	}
	
	.month-title {
		font-size: 32rpx;
		font-weight: bold;
		color: #333;
		margin-bottom: 20rpx;
	}
	
	.month-stats {
		display: flex;
		flex-direction: column;
		gap: 10rpx;
	}
	
	.stat-item {
		font-size: 24rpx;
		color: #666;
	}
	
	.popup-content {
		background-color: white;
		border-radius: 20rpx 20rpx 0 0;
		max-height: 70vh;
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
		font-size: 32rpx;
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
	
	.date-orders {
		padding: 30rpx;
	}
	
	.order-summary {
		background-color: #f8f9fa;
		padding: 20rpx;
		border-radius: 10rpx;
		margin-bottom: 30rpx;
		text-align: center;
	}
	
	.summary-text {
		font-size: 28rpx;
		color: #333;
	}
	
	.order-list {
		display: flex;
		flex-direction: column;
		gap: 20rpx;
	}
	
	.order-item {
		background-color: #f8f9fa;
		padding: 20rpx;
		border-radius: 10rpx;
	}
	
	.order-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 10rpx;
	}
	
	.order-id {
		font-size: 28rpx;
		font-weight: bold;
		color: #333;
	}
	
	.order-amount {
		font-size: 28rpx;
		color: #e91e63;
		font-weight: bold;
	}
	
	.order-info {
		display: flex;
		gap: 20rpx;
	}
	
	.customer-name, .recipient {
		font-size: 26rpx;
		color: #666;
	}
	
	.no-orders {
		padding: 60rpx;
		text-align: center;
	}
	
	.no-orders-text {
		font-size: 28rpx;
		color: #999;
	}
</style>