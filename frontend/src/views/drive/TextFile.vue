<template>
  <div id="file">
    <div id="header">
      <el-icon>
        <HomeFilled/>
      </el-icon>
      <el-breadcrumb id="path-nav">
        <el-breadcrumb-item to="/drive/my-drive">云端硬盘</el-breadcrumb-item>
        <el-breadcrumb-item v-for="folder in path" :to="`/drive/folders/${folder.Id}`">
          {{ folder.Name }}
        </el-breadcrumb-item>
        <el-breadcrumb-item>{{ name }}</el-breadcrumb-item>
      </el-breadcrumb>
      <el-icon @click="download">
        <Download/>
      </el-icon>
    </div>
    <div id="description">
      <el-descriptions border>
        <el-descriptions-item label="文件名">{{ name }}</el-descriptions-item>
        <el-descriptions-item label="文件大小">{{ size }}字节</el-descriptions-item>
        <el-descriptions-item label="下载次数">{{ downloadCounts }}</el-descriptions-item>
      </el-descriptions>
    </div>
    <div id="view">
      <pre>{{ content }}</pre>
    </div>
  </div>
</template>

<script>
import { ref } from "vue";
import { useRoute } from "vue-router";

import { ElMessage } from "element-plus";
import { HomeFilled, Download } from "@element-plus/icons-vue";

import { driveService, storageService } from "/src/backend";
import marked from "/src/marked";

export default {
  components: { HomeFilled, Download },
  setup() {
    const route = useRoute();

    const name = ref("");
    const parent = ref("");
    const size = ref(0);
    const downloadCounts = ref(0);
    const path = ref([]);
    const content = ref("");

    const initFile = async (fileId) => {
      try {
        const { data: fileData } = await driveService.get(`/files/${fileId}`);
        name.value = fileData.name;
        parent.value = fileData.parent;
        size.value = fileData.bytes;
        downloadCounts.value = fileData.downloadCounts;

        const { data: pathData } = await driveService.get(`/path/${parent.value}`);
        path.value = pathData.Folders ? pathData.Folders.reverse() : [];

        const { data: textData } = await storageService.get(`/download/${fileId}`);
        content.value = textData;
      } catch (e) {
        ElMessage({ type: "error", message: "获取文件失败" });
      }
    }

    initFile(route.params.fileId);

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

    return { name, size, downloadCounts, path, content, download };
  }
}
</script>

<style scoped>
#file {
  font-size: 18px;
  border-radius: 5px;
  background-color: #ffffff;
  box-shadow: 0 0 0 1px #eee;
}

#header {
  display: flex;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e5e9ef;
}

#header #path-nav {
  flex: 1;
  padding: 0 15px;
}

#view {
  padding: 20px;
}
</style>
