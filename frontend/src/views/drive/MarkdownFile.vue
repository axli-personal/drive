<template>
  <div id="file">
    <div id="header">
      <span>{{ name }}, {{ size }}字节, 已下载{{ downloadCounts }}次.</span>
      <el-icon :size="20" @click="download">
        <Download/>
      </el-icon>
    </div>
    <div id="view">
      <div v-html="content"></div>
    </div>
  </div>
</template>

<script>
import { ref } from "vue";
import { useRoute } from "vue-router";

import { ElIcon, ElMessage } from "element-plus";
import { Download } from "@element-plus/icons-vue";

import { driveService, storageService } from "/src/backend";

import marked from "/src/marked";

export default {
  components: { ElIcon, Download },
  setup() {
    const route = useRoute();

    const name = ref("");
    const size = ref(0);
    const downloadCounts = ref(0);
    const content = ref("");

    driveService.get(
      `/files/${route.params["fileId"]}`
    ).then(({ data }) => {
      name.value = data.name;
      size.value = data.bytes;
      downloadCounts.value = data.downloadCounts;
    }).catch(() => {
      ElMessage({ type: "error", message: "获取文件失败" });
    });

    storageService.get(
      `/download/${route.params["fileId"]}`,
    ).then(({ data }) => {
      content.value = marked(data);
    }).catch(() => {
      ElMessage({ type: "error", message: "下载失败" });
    });

    const download = () => {
      storageService.get(
        `/download/${route.params["fileId"]}`,
        {
          responseType: "blob",
        }
      ).then(({ data }) => {
        const fileURL = URL.createObjectURL(data);
        const link = document.createElement('a');

        link.href = fileURL;
        link.setAttribute('download', name.value);

        document.body.appendChild(link);
        link.click();

        document.body.removeChild(link);
        URL.revokeObjectURL(fileURL);
      }).catch(() => {
        ElMessage({ type: "error", message: "下载失败" });
      });
    };

    return { name, size, downloadCounts, content, download };
  }
}
</script>

<style scoped>
#file {
  border-radius: 5px;
  background-color: #ffffff;
  box-shadow: 0 0 0 1px #eee;
}

#header {
  display: flex;
  flex-direction: row;
  padding: 10px;
  border-bottom: 1px solid #e5e9ef;
  font-size: 18px;
}

#header span {
  flex: 1;
}

#view {
  padding: 10px;
}
</style>
