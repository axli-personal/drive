import "./css/global.scss";

import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import ElementPlus from 'element-plus';

const app = createApp(App);

app.use(router).use(ElementPlus).mount("#app");
