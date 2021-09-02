<template>
  <p>{{ authinfo.FullName }} - {{ authinfo.Email }}</p>
</template>

<script>
export default {
  props: {
    authinfo: {
      type: Object,
      required: true,
    },
  }
}
</script>