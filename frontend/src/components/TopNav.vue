<template>
  <div id="top-nav">
    <div id="left">
      <router-link id="home" to="/">云端硬盘</router-link>
      <router-link to="/drive/my-recycle-bin">回收站</router-link>
    </div>
    <div id="right">
      <router-link to="/register">注册</router-link>
      <router-link to="/login">登录</router-link>
    </div>
  </div>
</template>

<script>
export default {}
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

#left {
  flex: 1;
}

#left a {
  margin-right: 25px;
}

#right a {
  margin-left: 25px;
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

#home {
  color: #188a92;
}

.link-same {
  color: #5455ac;
}
</style>
