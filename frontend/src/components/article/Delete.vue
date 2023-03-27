<template>
  <div>
    <LabelInput label="文章ID" v-model="articleId"></LabelInput>
    <LabelButton @click="request"></LabelButton>
  </div>
</template>

<script>
import axios from "axios";

import { ref } from "vue";
import { ElMessage } from "element-plus";
import { site } from "/src/backend";

import LabelInput from "/src/components/ui/LabelInput.vue";
import LabelButton from "/src/components/ui/LabelButton.vue";

export default {
  components: { LabelInput, LabelButton },
  setup() {
    const articleId = ref("");

    const request = () => {
      axios
        .post(site + "/article/delete", {
          id: parseInt(localStorage.getItem("ID")),
          token: localStorage.getItem("Identity"),
          articleId: parseInt(articleId.value),
        })
        .then(({ data }) => {
          if (data.success) {
            ElMessage({ type: "success", message: data.detail });
            router.push("/drive/search");
          } else {
            ElMessage({ type: "error", message: data.detail });
          }
        });
    };

    return { articleId, request };
  },
};
</script>