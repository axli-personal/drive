<template>
  <TopNav :valid="valid"></TopNav>
  <div class="tiny-padding">
    <router-view v-slot="{ Component }">
      <transition name="fade" mode="out-in" appear>
        <component :is="Component"/>
      </transition>
    </router-view>
  </div>
</template>

<script>
import { ref } from "vue";

import TopNav from "/src/components/TopNav.vue";

export default {
  components: { TopNav },
  setup() {
    const valid = ref(localStorage.getItem("ID") != null);

    return { valid };
  },
};
</script>

<style scoped>
.tiny-padding {
  padding: 25px;
}
</style>
