// 通用工具函数

// 格式化日期
export function formatDate(date, format = 'YYYY-MM-DD') {
  if (!date) return '';
  
  const d = new Date(date);
  const year = d.getFullYear();
  const month = (d.getMonth() + 1).toString().padStart(2, '0');
  const day = d.getDate().toString().padStart(2, '0');
  const hour = d.getHours().toString().padStart(2, '0');
  const minute = d.getMinutes().toString().padStart(2, '0');
  
  switch (format) {
    case 'YYYY-MM-DD':
      return `${year}-${month}-${day}`;
    case 'YYYY-MM-DD HH:mm':
      return `${year}-${month}-${day} ${hour}:${minute}`;
    case 'MM-DD':
      return `${month}-${day}`;
    case 'HH:mm':
      return `${hour}:${minute}`;
    default:
      return `${year}-${month}-${day}`;
  }
}

// 格式化价格
export function formatPrice(price) {
  if (typeof price !== 'number') {
    price = parseFloat(price) || 0;
  }
  return price.toFixed(2);
}

// 显示提示信息
export function showToast(title, icon = 'none', duration = 2000) {
  uni.showToast({
    title,
    icon,
    duration
  });
}

// 显示确认对话框
export function showConfirm(content, title = '提示') {
  return new Promise((resolve) => {
    uni.showModal({
      title,
      content,
      success: (res) => {
        resolve(res.confirm);
      }
    });
  });
}

// 显示加载中
export function showLoading(title = '加载中...') {
  uni.showLoading({
    title,
    mask: true
  });
}

// 隐藏加载中
export function hideLoading() {
  uni.hideLoading();
}

// 获取农历信息
export function getLunarInfo(date) {
  // 这里可以集成农历转换库
  // 暂时返回模拟数据
  const lunarMonths = ['正月', '二月', '三月', '四月', '五月', '六月', 
                      '七月', '八月', '九月', '十月', '冬月', '腊月'];
  const lunarDays = ['初一', '初二', '初三', '初四', '初五', '初六', '初七', '初八', '初九', '初十',
                    '十一', '十二', '十三', '十四', '十五', '十六', '十七', '十八', '十九', '二十',
                    '廿一', '廿二', '廿三', '廿四', '廿五', '廿六', '廿七', '廿八', '廿九', '三十'];
  
  const d = new Date(date);
  const month = d.getMonth();
  const day = d.getDate();
  
  return {
    month: lunarMonths[month % 12],
    day: lunarDays[(day - 1) % 30]
  };
}

// 获取节假日信息
export function getHolidayInfo(date) {
  const holidays = {
    '01-01': '元旦',
    '02-14': '情人节',
    '03-08': '妇女节',
    '05-01': '劳动节',
    '06-01': '儿童节',
    '10-01': '国庆节',
    '12-25': '圣诞节'
  };
  
  const d = new Date(date);
  const key = `${(d.getMonth() + 1).toString().padStart(2, '0')}-${d.getDate().toString().padStart(2, '0')}`;
  return holidays[key] || null;
}

// 防抖函数
export function debounce(func, wait) {
  let timeout;
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
}

// 节流函数
export function throttle(func, limit) {
  let inThrottle;
  return function() {
    const args = arguments;
    const context = this;
    if (!inThrottle) {
      func.apply(context, args);
      inThrottle = true;
      setTimeout(() => inThrottle = false, limit);
    }
  };
}

// 深拷贝
export function deepClone(obj) {
  if (obj === null || typeof obj !== 'object') return obj;
  if (obj instanceof Date) return new Date(obj.getTime());
  if (obj instanceof Array) return obj.map(item => deepClone(item));
  if (typeof obj === 'object') {
    const clonedObj = {};
    for (const key in obj) {
      if (obj.hasOwnProperty(key)) {
        clonedObj[key] = deepClone(obj[key]);
      }
    }
    return clonedObj;
  }
}

// 生成唯一ID
export function generateId() {
  return Date.now().toString(36) + Math.random().toString(36).substr(2);
}

// 验证手机号
export function validatePhone(phone) {
  const phoneRegex = /^1[3-9]\d{9}$/;
  return phoneRegex.test(phone);
}

// 验证邮箱
export function validateEmail(email) {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}