<template>
  <div id="top-nav">
    <div id="left">
      <!-- Home -->
      <div id="top-nav-home">
        <router-link to="/">云端硬盘</router-link>
      </div>
      <!-- Old
      <div id="top-nav-article" class="list-unit" @mouseenter="expend" @mouseleave="collapse">
        <img class="icon-unit" src="/icon/top-nav/article.svg" alt="article options entry">
        <div class="list-box list-center-position" style="display: none">
          <div class="list-center-flex">
            <router-link to="/article/search">文章搜索</router-link>
            <router-link v-if="$props.valid" to="/article/management">内容管理</router-link>
          </div>
        </div>
      </div>
      -->
    </div>

    <div id="right">
      <!-- Identity -->
      <div id="top-nav-identity" class="list-unit" @mouseenter="expend" @mouseleave="collapse">
        <img class="icon-unit" src="/icon/top-nav/user.svg" alt="identity options entry">
        <div class="list-box list-center-position" style="display: none">
          <div class="list-center-flex">
            <router-link to="/register">注册</router-link>
            <router-link to="/login">登录</router-link>
          </div>
          <!-- Old
          <div class="list-center-flex">
            <router-link to="/identity/dashboard">控制面板</router-link>
            <router-link to="/identity/logout">立即注销</router-link>
          </div>
          -->
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  setup() {
    const expend = (event) => {
      event.target.children.item(1).style.display = "block";
    }
    const collapse = (event) => {
      event.target.children.item(1).style.display = "none";
    }
    return { expend, collapse };
  }
}
</script>

<style scoped>
#top-nav {
  height: 50px;
  display: flex;
  padding: 0 25px;
  background: #fbfbfd;
  flex-flow: row nowrap;
  border-bottom: 1px solid #dfdfdf;
}

#top-nav-home {
  padding-right: 15px;
}

#left {
  flex: 1;
  display: flex;
  flex-direction: row;
}

#right {
  flex: 0;
  display: flex;
  flex-direction: row;
}

.list-unit {
  position: relative;
  z-index: 100;
}

.icon-unit {
  margin: 10px 15px;
  height: 30px;
}

.list-box {
  padding: 0 15px;
  border: 1px solid #dfdfdf;
  border-radius: 0 0 8px 8px;
  background: #fbfbfd;
}

#top-nav-identity .list-box {
  margin-left: -20px;
}

.list-center-position {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translate(-50%);
}

.list-center-flex {
  display: flex;
  align-items: center;
  flex-direction: column;
}

/* Below are link styles */

a {
  color: #47474b;
  font-size: 25px;
  width: max-content;
  font-weight: 700;
  line-height: 50px;
  text-align: center;
  text-decoration: none;
}

#top-nav-home a {
  color: #188a92;
}

.link-same {
  color: #5455ac;
}
</style>
