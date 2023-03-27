<template>
  <div>
    <LabelInput label="新账号名" v-model="newAccount"></LabelInput>
    <LabelButton label="最终确认" @click="submit"></LabelButton>
  </div>
</template>

<script>
import { ref } from "vue"
import axios from "axios"

import LabelButton from "/src/components/ui/LabelButton.vue"
import LabelInput from "/src/components/ui/LabelInput.vue"

export default {
  components: { LabelButton, LabelInput },
  setup() {
    const newAccount = ref("");

    const submit = () => {
      import("/src/backend").then(({ site }) => {
        axios.post(site + "/identity/dashboard/update/info", {
          id: parseInt(localStorage.getItem("ID")),
          token: localStorage.getItem("Identity"),
          account: newAccount.value,
        })
      })
    }

    return { newAccount, submit }
  }
}
</script>
