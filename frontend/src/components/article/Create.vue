<template>
  <div>
    <div class="form-item">
      <label class="form-item-label" for="form-title">文章标题</label>
      <input class="form-item-input" id="form-title" v-model="title" placeholder="请输入文章标题">
    </div>

    <div class="form-item">
      <label class="form-item-label" for="form-abstract">文章摘要</label>
      <input class="form-item-input" id="form-abstract" v-model="abstract" placeholder="请输入文章摘要">
    </div>

    <div class="form-item">
      <label class="form-item-label" for="form-body">文章主体</label>
      <textarea class="form-item-textarea" id="form-body" v-model="body" @input="adjustHeight"></textarea>
    </div>

    <hr class="form-divider">

    <div class="form-item">
      <button class="form-item-button" @click="request">发表</button>
    </div>
  </div>
</template>

<script>
import axios from "axios";

import { ref } from "vue";
import { site } from "/src/backend";

import { ElMessage } from "element-plus";

export default {
  setup() {
    const title = ref("");
    const abstract = ref("");
    const body = ref("");

    const request = () => {
      axios
        .post(site + "/article/publish", {
          id: parseInt(localStorage.getItem("ID")),
          token: localStorage.getItem("Identity"),
          title: title.value,
          abstract: abstract.value,
          body: body.value,
        })
        .then(({ data }) => {
          if (data.success) {
            ElMessage({ type: "success", message: data.detail });
          } else {
            ElMessage({ type: "error", message: data.detail });
          }
        });
    };

    const adjustHeight = function(event) {
      // 16 is the top padding plus the bottom padding.
      event.target.style.height = event.target.scrollHeight - 16 + "px";
    }

    return { title, abstract, body, request, adjustHeight };
  },
};
</script>

<style scoped>
.form-item {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  margin: 20px 0;
}

.form-item-label {
  min-width: 160px;
  height: 40px;
  line-height: 40px;
  font-size: 22px;
  font-weight: 600;
}

.form-item-input {
  flex: 1;
  height: 40px;
  padding: 8px;
  border: 1px solid #bfcbd9;
  border-radius: 4px;
  font-size: 18px;
  outline: none;
}

.form-item-textarea {
  flex: 1;
  min-height: 75px;
  max-height: 300px;
  resize: none;
  padding: 8px;
  border: 1px solid #bfcbd9;
  border-radius: 4px;
  font-size: 18px;
  outline: none;
}

.form-item-button {
  width: 400px;
  padding: 8px;
  border: 1px solid #bfcbd9;
  border-radius: 4px;
  font-size: 18px;
  cursor: pointer;
}

.form-divider {
  border-color: #e5e9ef;
  margin-bottom: 40px;
}
</style>
