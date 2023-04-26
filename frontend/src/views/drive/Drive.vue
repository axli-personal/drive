<template>
  <div id="drive">
    <div id="header">
      <el-icon>
        <HomeFilled/>
      </el-icon>
      <el-breadcrumb id="path-nav">
        <el-breadcrumb-item to="/drive/my-drive">云端硬盘</el-breadcrumb-item>
      </el-breadcrumb>
      <el-icon @click="displayDrawer">
        <CirclePlus/>
      </el-icon>
      <CreateDrawer v-model="isDisplayDrawer" parent="Drive"></CreateDrawer>
    </div>

    <ul id="menu">
      <li v-for="folder in folders">
        <el-icon>
          <Folder/>
        </el-icon>
        <span @click="viewFolder(folder.id)">{{ folder.name }}</span>
        <el-icon @click="removeFolder(folder.id)">
          <Delete/>
        </el-icon>
      </li>
      <li v-for="file in files">
        <el-icon>
          <Document/>
        </el-icon>
        <span @click="viewFile(file.id, file.name)">{{ file.name }}</span>
        <el-icon @click="removeFile(file.id)">
          <Delete/>
        </el-icon>
      </li>
    </ul>
  </div>
</template>

<script>
import { ref } from "vue";
import { useRouter } from "vue-router";

import { ElMessage } from "element-plus";
import { HomeFilled, CirclePlus, Folder, Document, Delete } from "@element-plus/icons-vue";

import CreateDrawer from "/src/components/drive/CreateDrawer.vue";

import { driveService } from "/src/backend";
import { getFileViewType } from "./util";

export default {
  components: { CreateDrawer, HomeFilled, CirclePlus, Folder, Document, Delete },
  setup() {
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

    const viewFolder = (folderId) => {
      router.push(`/drive/folders/${folderId}`);
    }

    const viewFile = (fileId, fileName) => {
      router.push(`/drive/files/${getFileViewType(fileName)}/${fileId}`);
    }

    const removeFolder = async (folderId) => {
      try {
        const { data } = await driveService.post(`/folders/remove/${folderId}`);
        window.location.reload();
      } catch (e) {
        ElMessage({ type: "error", message: "移除失败" });
      }
    }

    const removeFile = async (fileId) => {
      try {
        const { data } = await driveService.post(`/files/remove/${fileId}`);
        window.location.reload();
      } catch (e) {
        ElMessage({ type: "error", message: "移除失败" });
      }
    }

    const isDisplayDrawer = ref(false);

    const displayDrawer = () => {
      isDisplayDrawer.value = true;
    }

    return {
      folders,
      files,
      isDisplayDrawer,
      displayDrawer,
      viewFolder,
      viewFile,
      removeFile,
      removeFolder
    };
  },
};
</script>

<style scoped>
#drive {
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

#header #path-nav,
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
