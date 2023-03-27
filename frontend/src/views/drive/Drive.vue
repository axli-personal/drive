<template>
  <div id="drive">
    <div id="header">
      <span>我的云端硬盘</span>
      <el-upload :http-request="upload" :show-file-list="false">
        <el-icon :size="20">
          <Upload/>
        </el-icon>
      </el-upload>
    </div>

    <ul id="menu">
      <li v-for="folder in folders" @click="changeFolder(folder.id)">
        <img src="/icon/folder.png" alt="folder">
        {{ folder.name }}
      </li>
      <li v-for="file in files" @click="viewFile(file.id, file.name)">
        <img src="/icon/file.png" alt="file">
        {{ file.name }}
      </li>
    </ul>
  </div>
</template>

<script>
import { ref } from "vue";
import { useRoute, useRouter } from "vue-router";

import { ElUpload, ElIcon, ElMessage } from "element-plus";
import { Upload } from "@element-plus/icons-vue";

import { driveService, storageService } from "/src/backend";

export default {
  components: { ElUpload, ElIcon, Upload },
  setup() {
    const route = useRoute();
    const router = useRouter();

    const folders = ref([]);
    const files = ref([]);

    driveService.get("/drive").then(({ data }) => {
      folders.value = data.children.folders;
      files.value = data.children.files;
    }).catch((err) => {
      const data = err.response.data;

      if (data.code) {
        if (data.code === "Unauthenticated") {
          router.push("/login");
          return;
        }
        if (data.code === "NotCreateDrive") {
          router.push("/drive/plan");
          return;
        }
      }

      ElMessage({ type: "error", message: "获取硬盘失败" });
    });

    const changeFolder = (folderId) => {
      router.push(`/drive/folders/${folderId}`);
    }

    const viewFile = (fileId, fileName) => {
      let viewType = "binary";
      const dotPosition = fileName.lastIndexOf(".");
      if (dotPosition != -1) {
        switch (fileName.substring(dotPosition + 1)) {
          case "txt":
            viewType = "text"
            break;
          case "md":
            viewType = "markdown"
            break;
        }
      }
      router.push(`/drive/files/${viewType}/${fileId}`);
    }

    const upload = (options) => {
      const file = options.file;
      const form = new FormData();
      form.append('file', file);
      form.append('parent', "Drive");
      storageService.post(
        "/upload",
        form,
        {
          headers: { "content-type": "multipart/form-data" },
        }
      ).then(() => {
        window.location.reload();
      }).catch(() => {
        ElMessage({ type: "error", message: "上传失败" });
      })
    }

    return { folders, files, changeFolder, viewFile, upload };
  },
};
</script>

<style scoped>
#drive {
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

#menu {
  display: flex;
  flex-direction: column;
  margin: 0;
  padding: 0;
  font-size: 18px;
}

#menu li {
  padding: 18px;
  border-bottom: 1px solid #e5e9ef;
  list-style: none;
  cursor: pointer;
}

#menu img {
  height: 16px;
}

#menu a {
  color: inherit;
  text-decoration: none;
}
</style>
