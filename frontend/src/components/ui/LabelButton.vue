<template>
  <div class="ipt-wrapper">
    <span class="pre">{{ label }}</span>
    <button @click="$emit('click')" class="btn">{{ msg }}</button>
  </div>
</template>

<script>
export default {
  props: {
    label: {
      type: String,
      default: "最终确认"
    },
    msg: {
      type: String,
      default: "更新"
    }
  },
  emits: {
    click: null
  }
}
</script>

<style scoped>
.ipt-wrapper {
  display: flex;
  height: 40px;
  flex-direction: row;
  margin: 20px 0;
}

.pre {
  min-width: 160px;
  line-height: 40px;
  font-size: 22px;
  font-weight: bold;
}

.btn {
  flex: 1;
  padding: 8px;
  border: 1px solid #bfcbd9;
  border-radius: 4px;
  font-size: 18px;
  cursor: pointer;
}

.btn:hover {
  background: rgb(59, 182, 223);
}
</style>