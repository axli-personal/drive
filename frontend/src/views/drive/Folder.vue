<template>
  <div id="folder">
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
      <el-icon @click="displayDrawer">
        <CirclePlus/>
      </el-icon>
      <CreateDrawer v-model="isDisplayDrawer" :parent="id"></CreateDrawer>
    </div>

    <ul id="menu">
      <li v-if="parent">
        <el-icon>
          <Folder/>
        </el-icon>
        <span @click="backToParent">..</span>
      </li>
      <li v-for="folder in folders">
        <el-icon>
          <Folder/>
        </el-icon>
        <span @click="viewFolder(folder.id)">{{ folder.name }}</span>
        <el-icon @click="viewFolder(folder.id)">
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
import { useRoute, useRouter, onBeforeRouteUpdate } from "vue-router";

import { ElMessage } from "element-plus";
import { HomeFilled, CirclePlus, Folder, Document, Delete } from "@element-plus/icons-vue";

import CreateDrawer from "/src/components/drive/CreateDrawer.vue";

import { driveService } from "/src/backend";
import { getFileViewType } from "./util";

export default {
  components: { CreateDrawer, HomeFilled, CirclePlus, Folder, Document, Delete },
  setup() {
    const route = useRoute();
    const router = useRouter();

    const id = ref("");
    const name = ref("");
    const parent = ref("");
    const folders = ref([]);
    const files = ref([]);
    const path = ref([]);

    const initFolder = async (folderId) => {
      try {
        id.value = folderId;

        const { data: folderData } = await driveService.get(`/folders/${folderId}`);
        name.value = folderData.name;
        parent.value = folderData.parent;
        folders.value = folderData.children.folders;
        files.value = folderData.children.files;

        const { data: pathData } = await driveService.get(`/path/${parent.value}`);
        path.value = pathData.Folders ? pathData.Folders.reverse() : [];
      } catch (e) {
        ElMessage({ type: "error", message: "获取目录失败" });
      }
    }

    initFolder(route.params.folderId);

    onBeforeRouteUpdate((to) => {
      initFolder(to.params.folderId);
    })

    const backToParent = () => {
      if (parent.value === "Drive") {
        router.push("/drive/my-drive");
      } else {
        router.push(`/drive/folders/${parent.value}`);
      }
    }

    const viewFolder = (folderId) => {
      router.push(`/drive/folders/${folderId}`);
    }

    const viewFile = (fileId, fileName) => {
      router.push(`/drive/files/${getFileViewType(fileName)}/${fileId}`);
    }

    const isDisplayDrawer = ref(false);

    const displayDrawer = () => {
      isDisplayDrawer.value = true;
    }

    return {
      id,
      name,
      parent,
      files,
      folders,
      path,
      backToParent,
      viewFolder,
      viewFile,
      isDisplayDrawer,
      displayDrawer
    };
  },
};
</script>

<style scoped>
#folder {
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
