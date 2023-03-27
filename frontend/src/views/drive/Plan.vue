<template>
  <div>正在为您创建一块免费的硬盘...</div>
</template>

<script>
import { driveService } from "../../backend";
import { ElMessage } from "element-plus";
import { useRouter } from "vue-router";

export default {
  setup() {
    const router = useRouter();

    driveService.post("/drive/create").then(() => {
      ElMessage({ type: "success", message: "硬盘创建成功" });
      router.push("/");
    }).catch((err) => {
      console.log(err.response);
    });
  }
}
</script>

<style scoped>

</style>