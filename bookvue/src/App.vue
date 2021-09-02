<template>
  <AuthInfo :authinfo="authinfo" />
  <h1>Book List</h1>
  <BookList :books="books" />
</template>

<script>
import AuthInfo from './components/AuthInfo.vue'
import BookList from './components/BookList.vue'
import axios from "axios"
export default {
  components: {
    AuthInfo,
    BookList,
  },
  data() {
    return {
      books: [],
      authinfo: {
        FullName: "",
        Email: "",
      },
      loading: false,
      error: null,
    }
  },
  methods: {    
    async fetchBooks() {
      try {
        this.error = null
        this.loading = true
        const url = `http://localhost:9001/api/books`
        axios.defaults.withCredentials = true;
        axios.defaults.headers.post['Content-Type'] = 'application/json';
        const response = await axios.get(url)

        if (response.status != 200) {
          console.log(response)
        } else {
          this.books = response.data.books
          this.authinfo = response.data.authinfo
        }
        
      } catch (err) {       
        console.log(err)
      }
      this.loading = false
    },
  },
  mounted() {
    this.fetchBooks()
  },

}
</script>

<style>
#app {
  font-family: "SF Mono", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: left;
  color: #2c3e50;
  margin-top: 20px;
}
</style>
