<template>
  <div class="page">
    <el-form label-width="120px" label-position="left">
      <h1>登陆</h1>
      <el-form-item label="账号">
        <el-input v-model="account" placeholder="请输入账号"/>
      </el-form-item>
      <el-form-item label="密码">
        <el-input v-model="password" placeholder="请输入密码"/>
      </el-form-item>
      <el-form-item label="确认">
        <el-button @click="submit">登陆</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { ElForm, ElFormItem, ElInput, ElButton, ElMessage } from "element-plus";
import { useRouter } from "vue-router";

import { userService } from "/src/backend";

const router = useRouter();

const account = ref("");
const password = ref("");

const submit = () => {
  userService.post(
    "/login",
    {
      account: account.value,
      password: password.value,
    }
  ).then(() => {
    router.push("/");
  }).catch(() => {
    ElMessage({ type: "error", message: "登陆失败" });
  })
}
</script>
