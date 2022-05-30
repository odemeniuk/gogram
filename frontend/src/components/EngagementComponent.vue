<template>
  <div class="engagement-form">
    <input type="text" v-model="username" >
    <button @click="calculate()">Calculate</button>
    <h1 v-show="result != ''">User engagement is: {{result}}%</h1>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  name: 'EngagementComponent',
  data() {
    return {
      username: '',
      result: '',
    }
  },
  methods: {
      calculate: function() {
        const headers = {
    'Content-Type': 'application/json',
  }
        axios.post('http://localhost:8081/calculate', {"username": this.username},{
    headers: headers
  }).then(response => {
          console.log(response)
          this.result = response.data.engagement
        }).catch(error => {
          console.log(error.response)
        })
      }
    }
  
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
