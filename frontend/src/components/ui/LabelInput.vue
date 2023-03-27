<template>
  <div class="ipt-wrapper">
    <span class="pre">{{ label }}</span>
    <input
      class="ipt"
      v-model="modelValue"
      @input="$emit('update:modelValue', modelValue)"
      :placeholder="defaultHolder"
    />
  </div>
</template>

<script>
export default {
  props: {
    label: {
      type: String,
      required: true,
    },
    modelValue: {
      type: String,
      required: true,
    },
    holder: {
      type: String,
      default: "",
    },
  },
  emits: {
    "update:modelValue": null,
  },
  setup(props) {
    let defaultHolder = props.holder;

    if (defaultHolder.length == 0) defaultHolder = "请输入" + props.label;

    return { defaultHolder };
  },
};
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

.ipt {
  flex: 1;
  padding: 8px;
  border: 1px solid #bfcbd9;
  border-radius: 4px;
  font-size: 18px;
  outline: none;
}
</style>