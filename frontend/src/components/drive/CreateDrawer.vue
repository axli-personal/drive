<template>
  <el-drawer :with-header="false">
    <el-menu default-active="1" mode="horizontal">
      <el-menu-item index="1" @click="displayMenuContent('create-folder')">新建文件夹</el-menu-item>
      <el-menu-item index="2" @click="displayMenuContent('create-file')">上传文件</el-menu-item>
    </el-menu>
    <div id="menu-content">
      <el-form v-if="displayKey === 'create-folder'" label-width="120px" label-position="left">
        <el-form-item label="目录名">
          <el-input v-model="folderName" placeholder="请输入目录名"/>
        </el-form-item>
        <el-form-item label="确认">
          <el-button @click="createFolder">创建</el-button>
        </el-form-item>
      </el-form>
      <el-form v-if="displayKey === 'create-file'" label-width="120px" label-position="left">
        <el-form-item label="文件名">
          <el-input disabled placeholder="选择文件后显示"/>
        </el-form-item>
        <el-form-item label="上传文件">
          <el-upload style="width: 100%" drag :http-request="uploadFile" :show-file-list="false">
            <el-icon class="el-icon--upload">
              <UploadFilled/>
            </el-icon>
            <div class="el-upload__text">拖拽或点击上传</div>
          </el-upload>
        </el-form-item>
      </el-form>
    </div>
  </el-drawer>
</template>

<script>
import { ref } from "vue";
import { ElMessage } from "element-plus";
import { UploadFilled } from "@element-plus/icons-vue";

import { driveService, storageService } from "/src/backend";

export default {
  components: { UploadFilled },
  props: {
    parent: {
      type: String,
      default: "Drive"
    }
  },
  setup(props) {
    const displayKey = ref("create-folder");

    const displayMenuContent = (key) => {
      displayKey.value = key;
    }

    const folderName = ref("");

    const createFolder = () => {
      driveService.post("/folders/create", {
        parent: props.parent,
        folderName: folderName.value,
      }).then(() => {
        location.reload();
      }).catch(() => {
        ElMessage({ type: "error", message: "创建失败" });
      })
    }

    const uploadFile = (options) => {
      const form = new FormData();
      form.append('parent', props.parent);
      form.append('file', options.file);
      storageService.post("/upload", form, {
        headers: { "content-type": "multipart/form-data" }
      }).then(() => {
        location.reload();
      }).catch(() => {
        ElMessage({ type: "error", message: "上传失败" });
      })
    }

    return { displayKey, displayMenuContent, folderName, createFolder, uploadFile };
  },
};
</script>

<style scoped>
#menu-content {
  padding: 20px;
}
</style>
