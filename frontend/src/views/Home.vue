<template>
  <div class="recommend">
    <div class="title">推荐文章</div>

    <List :heads="heads"></List>

    <div v-if="nextPage !== -1" class="more" @click="loadMore"></div>
    <div v-else class="end" @click="backToStart"></div>
  </div>
</template>

<style lang="scss" scoped>
.recommend {
  border: 1px solid #dfdfdf;
  padding: 20px;
  background-color: #fbfbfd;

  .title {
    font-size: 24px;
    padding-bottom: 10px;
    border-bottom: 1px solid #dfdfdf;
  }

  .more, .end {
    margin-top: 20px;
    height: 40px;
    cursor: pointer;
  }

  .more {
    background-image: url("/icon/arrow/bottom.png");
  }

  .end {
    background-image: url("/icon/arrow/top.png");
  }
}
</style>

<script>
import axios from "axios";
import { ref } from "vue";
import { site } from "/src/backend";
import List from "/src/components/article/List.vue";

export default {
  components: { List },
  setup() {
    const heads = ref([]);

    let nextPage = ref(0);

    const loadMore = () => {
      if (nextPage.value === -1) return;

      axios
      .get(site + "/article/recommend/" + nextPage.value)
      .then(({data}) => {
        if (data == null) {
          nextPage.value = -1;
          return;
        }

        for (let i = 0; i < data.length; i++) {
          heads.value.push(data[i]);
        }
        nextPage.value++;
      })
      .catch((err) => {
        nextPage.value = -1;
      })
    }

    const backToStart = () => {
      scrollTo(0, 0);
    }

    loadMore();

    return { heads, nextPage, loadMore, backToStart };
  }
}
</script>