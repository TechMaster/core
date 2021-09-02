<template>
  <ul id="array-rendering">
    <li v-for="(book, index) in books" :key="index" :book="book">
      {{ book.title + " : " + book.author}}
    </li>
  </ul>
</template>

<script>
export default {
  props: {
    books: {
      type: Array,
      required: true,
    },
  }
}
</script>