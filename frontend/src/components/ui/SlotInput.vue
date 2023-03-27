<template>
  <!-- Provide parent position -->
  <div class="pos-wrapper">
    <!-- Tag to display status -->
    <div class="tag" :style="status ? 'background: #67C23A;' : 'background: #F56C6C;'"></div>
    <!-- Wrapper to hold input -->
    <div class="ipt-wrapper">
      <span class="label">{{ label }}</span>
      <slot></slot>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    label: {
      type: String,
      required: true,
    },
    status: {
      type: Boolean,
      default: false,
    }
  },
};
</script>

<style scoped>
.pos-wrapper {
  position: relative;
}

.ipt-wrapper {
  display: flex;
  flex-direction: row;
  align-items: center;
  margin: 20px 0;
}

.tag {
  position:  absolute;
  top:       0;
  right:     10px;

  width:  16px;
  height: 10px;
  border-radius: 0 0 3px 3px;
}

.label {
  min-width: 160px;
  height: 40px;
  line-height: 40px;
  font-size: 22px;
  font-weight: 600;
}

:slotted(input) {
  flex: 1;
  height: 40px;
  padding: 8px;
  border: 1px solid #bfcbd9;
  border-radius: 4px;
  font-size: 18px;
  outline: none;
}

:slotted(textarea) {
  flex: 1;
  min-height: 75px;
  resize: none;
  padding: 8px;
  border: 1px solid #bfcbd9;
  border-radius: 4px;
  font-size: 18px;
  outline: none;
}
</style>