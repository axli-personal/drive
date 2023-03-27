<template>
  <div>
    <p id="welcome">欢迎来到控制面板, <strong>{{ account }}</strong>.</p>

    <div id="content">
      <ul id="menu">
        <li id="menuHeader">控制面板</li>
        <li
            v-for="(opt, index) in opts"
            :key="index"
            :class="currentOpt === opt.component ? 'blue-background' : 'gray-hover'"
            @click="changeOpt(opt)"
        >
          {{ opt.title }}
        </li>
      </ul>

      <component :is="currentOpt" class="option"></component>
    </div>
  </div>
</template>

<script>
import "/src/css/color.scss";

import {ref} from "vue";

import axios from "axios";
import {site} from "/src/backend";

import Update from "/src/components/dashboard/Update.vue";
import Security from "/src/components/dashboard/Security.vue";

export default {
  components: {Update, Security},
  setup() {
    const account = ref("loading...");

    // Get account from server.
    axios
    .get(site + "/identity/dashboard/fetch/info", {
      id: parseInt(localStorage.getItem("ID")),
      token: localStorage.getItem("Identity"),
    })
    .then(({data}) => {
      if (data.success) account.value = data.account;
    });

    const opts = [
      {title: "修改账户", component: "Update"},
      {title: "修改密码", component: "Security"},
    ];

    const currentOpt = ref(opts[0].component);

    const changeOpt = (opt) => {
      currentOpt.value = opt.component;
    };

    return {account, opts, currentOpt, changeOpt};
  },
};
</script>

<style scoped>
#welcome {
  padding: 10px;
  margin-bottom: 20px;
  font-size: 18px;
  background: rgba(128, 128, 128, 0.5);
}

#content {
  display: flex;
  box-shadow: 0 0 0 1px #eee;
  background-color: #ffffff;
  border-radius: 5px;
  flex-direction: row;
}

#menu {
  display: flex;
  flex-direction: column;
  font-size: 18px;
}

#menu li {
  padding: 18px;
  border-bottom: 1px solid #e5e9ef;
  list-style: none;
  cursor: pointer;
}

#menu #menuHeader {
  color: #99a2aa;
  font-weight: bold;
  cursor: default;
}

.option {
  flex: 1;
  padding: 0 30px;
  border-left: 1px solid #e5e9ef;
  min-height: 250px;
}
</style>
