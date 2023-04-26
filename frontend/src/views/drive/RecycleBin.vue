<template>
  <div id="recycle-bin">
    <div id="header">
      <el-icon>
        <HomeFilled/>
      </el-icon>
      <span>回收站</span>
    </div>

    <ul id="menu">
      <li v-for="folder in folders">
        <el-icon>
          <Folder/>
        </el-icon>
        <span>{{ folder.name }}</span>
        <el-icon @click="restoreFolder(folder.id)">
          <RefreshLeft/>
        </el-icon>
      </li>
      <li v-for="file in files">
        <el-icon>
          <Document/>
        </el-icon>
        <span>{{ file.name }}</span>
        <el-icon @click="restoreFile(file.id)">
          <RefreshLeft/>
        </el-icon>
      </li>
    </ul>
  </div>
</template>

<script>
import { ref } from "vue";
import { useRouter } from "vue-router";

import { ElUpload, ElIcon, ElMessage } from "element-plus";
import { HomeFilled, Folder, Document, RefreshLeft } from "@element-plus/icons-vue";

import { driveService } from "/src/backend";

export default {
  components: { ElUpload, ElIcon, HomeFilled, Folder, Document, RefreshLeft },
  setup() {
    const router = useRouter();

    const folders = ref([]);
    const files = ref([]);

    const restoreFolder = (folderId) => {
      driveService.post(
        `/folders/restore/${folderId}`
      ).then(() => {
        window.location.reload();
      }).catch(() => {
        ElMessage({ type: "error", message: "恢复失败" });
      })
    }

    const restoreFile = (fileId) => {
      driveService.post(
        `/files/restore/${fileId}`
      ).then(() => {
        window.location.reload();
      }).catch(() => {
        ElMessage({ type: "error", message: "恢复失败" });
      })
    }

    driveService.get("/recycle-bin").then(({ data }) => {
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

      ElMessage({ type: "error", message: "获取回收站失败" });
    });

    return { folders, files, restoreFolder, restoreFile };
  },
};
</script>

<style scoped>
#recycle-bin {
  font-size: 18px;
  border-radius: 5px;
  background-color: #ffffff;
  box-shadow: 0 0 0 1px #eee;
}

#header,
#menu li {
  display: flex;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e5e9ef;
}

#header span,
#menu li span {
  flex: 1;
  padding: 0 15px;
}

#menu {
  display: flex;
  flex-direction: column;
  margin: 0;
  padding: 0;
}

#menu li {
  list-style: none;
  cursor: pointer;
}
</style>
