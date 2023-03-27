<template>
  <div>
    <LabelInput label="旧密码" v-model="previous"></LabelInput>
    <LabelInput label="新密码" v-model="current"></LabelInput>
    <LabelButton label="最终确认" @click="submit"></LabelButton>
  </div>
</template>

<script>
import { ref } from "vue"
import axios from "axios"

import LabelInput from "/src/components/ui/LabelInput.vue"
import LabelButton from "/src/components/ui/LabelButton.vue"

export default {
  components: { LabelInput, LabelButton },
  setup() {
    const previous = ref("")
    const current  = ref("")

    const submit = () => {
      import("/src/backend").then(({ site }) => {
        axios.post(site + "/identity/dashboard/update/security", {
          id:       parseInt(localStorage.getItem("ID")),
          token:    localStorage.getItem("Identity"),
          previous: previous.value,
          current:  current.value,
        })
      })
    }

    return { current, previous, submit }
  }
}
</script>
