<template>
  <div class="page">
    <el-form class="form-default" label-width="80px" label-position="left">
      <h1>注册</h1>
      <el-form-item label="账号">
        <el-input v-model="account" placeholder="请输入账号"/>
      </el-form-item>
      <el-form-item label="昵称">
        <el-input v-model="username" placeholder="请输入昵称"/>
      </el-form-item>
      <el-form-item label="密码">
        <el-input v-model="password" placeholder="请输入密码"/>
      </el-form-item>
      <el-form-item label="介绍">
        <el-input v-model="introduction" type="textarea" maxlength="50" autosize/>
      </el-form-item>
      <el-form-item label="确认">
        <el-button @click="submit">注册</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { ElForm, ElFormItem, ElInput, ElButton } from "element-plus";
import { ElMessage } from "element-plus";
import { useRouter } from "vue-router";

import { userService } from "/src/backend";

const router = useRouter();

const account = ref("");
const password = ref("");
const username = ref("");
const introduction = ref("");

const submit = () => {
  userService.post(
    "/register",
    {
      account: account.value,
      password: password.value,
      username: username.value,
      introduction: introduction.value,
    }
  ).then(() => {
    router.push("/login");
  }).catch(() => {
    ElMessage({ type: "error", message: "注册失败" });
  })
}
</script>
