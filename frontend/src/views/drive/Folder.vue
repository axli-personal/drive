<template>
  <div id="folder">
    <div id="header">
      <span>{{ name }}</span>
      <el-upload :http-request="upload" :show-file-list="false">
        <el-icon :size="20">
          <Upload/>
        </el-icon>
      </el-upload>
    </div>

    <ul id="menu">
      <li v-if="parent" @click="backToParent">
        <img src="/icon/folder.png" alt="folder">
        ..
      </li>
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

    const name = ref("");
    const parent = ref("");
    const folders = ref([]);
    const files = ref([]);

    const changeFolder = (folderId) => {
      router.push(`/drive/folders/${folderId}`);
      driveService.get(
        `/folders/${folderId}`
      ).then(({ data }) => {
        name.value = data.name;
        parent.value = data.parent;
        folders.value = data.children.folders;
        files.value = data.children.files;
      }).catch(() => {
        ElMessage({ type: "error", message: "获取目录失败" });
      });
    }

    const backToParent = () => {
      if (parent.value === "Drive") {
        router.push("/drive/my-drive");
      } else {
        changeFolder(parent.value);
      }
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
      form.append('parent', route.params["folderId"]);
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

    changeFolder(route.params["folderId"]);

    return { name, parent, files, folders, backToParent, changeFolder, viewFile, upload };
  },
};
</script>

<style scoped>
#folder {
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
