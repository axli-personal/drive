<template>
  <div class="page" id="page-search">
    <section id="section-search">
      <el-input @change="search" v-model="keyword" placeholder="Please Input the title keyword to search">
        <template #prepend>搜索文章</template>
      </el-input>

      <el-divider />

      <div class="result">
        <div class="title">搜索结果</div>
        <List :heads="heads"></List>
      </div>
    </section>
  </div>
</template>

<script>
import axios from "axios";

import { ref }       from "vue";
import { useRouter } from "vue-router";
import { site }      from "/src/backend";

import { ElMessage, ElInput, ElDivider } from "element-plus";
import List from "/src/components/article/List.vue";

export default {
  components: { List, ElInput, ElDivider },
  setup() {
    const router  = useRouter();
    const keyword = ref("");
    const heads   = ref([]);

    const search = () => {
      axios
      .post(site + "/article/search/title", {
        keyword: keyword.value,
      })
      .then(({ data }) => {
        if (data.success) {
          heads.value = data.Heads;
        } else {
          ElMessage({ type: "error", message: data.detail });
        }
      })
    }

    return { keyword, heads, search }
  }
}
</script>

<style>
#section-search .el-input-group__prepend,
#section-search input {
  height: 70px;
  color: #161616;
  font-size: 20px;
}

.title {
  font-size: 24px;
  padding-bottom: 10px;
  border-bottom: 1px solid #dfdfdf;
}

.result {
  min-height: 150px;
  border: 1px solid #dfdfdf;
  padding: 20px;
  background-color: #fbfbfd;
}
</style>
