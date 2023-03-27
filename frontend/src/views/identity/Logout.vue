<template>
  <div>正在跳转...</div>
</template>

<script>
import { inject } from 'vue';
import { useRouter } from "vue-router";

export default {
  setup() {
    let router = useRouter();

    const logout = inject("logout");

    localStorage.removeItem("ID");
    localStorage.removeItem("Identity");

    logout();

    router.push("/home");
  }
}
</script>
