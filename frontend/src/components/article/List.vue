<template>
  <ul class="list">
    <li v-for="head in heads">
      <h4><a :href="'/drive/content/' + head.id">{{ head.title }}</a></h4>
      <p>{{ head.abstract }}</p>
      <span class="tag">作者: {{ head.member.nickname }}</span>
      <span class="tag">浏览量: {{ head.counts }}</span>
    </li>
  </ul>
</template>

<script>
export default {
  props: {
    heads: {
      type: Array,
      required: true
    }
  }
}
</script>

<style lang="scss" scoped>
.list {
  margin: 0;
  padding: 0;
  color: #555666;

  li {
    padding: 10px 0;
    border-bottom: 1px solid #dfdfdf;
    list-style: none;
  }

  h4 {
    margin: 0 0 10px 0;
    color: #222226;
    font-size: 18px;
  }

  a {
    color: inherit;
    text-decoration: none;
    border-bottom: 3px solid #4FC08D;
  }

  p {
    min-height: 25px;
    margin: 0 0 10px 0;
  }

  .tag {
    margin-right: 25px;
  }
}
</style>