import { createApp } from 'vue'
import App from './App.vue'
import '@/assets/icon/iconfont.css'

const app = createApp(App)

// element-plus
import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css';
// 注册所有图标
import * as ElementPlusIconsVue from '@element-plus/icons-vue';
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component);
}
app.use(ElementPlus)

import { createPinia } from 'pinia';
const pinia = createPinia()
app.use(pinia)

import router from './router';
app.use(router)

app.mount('#app')