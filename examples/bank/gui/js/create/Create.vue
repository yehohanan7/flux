<template>
  <div>
  <header></header>
  <div class="row" id="account">
    <div class="four columns">
      <input class="button-primary" type="button" v-bind:value="text" v-on:click="create">
    </div>
  </div>
  </div>
</template>

<script>
export default {
  name: 'create',
  data () {
    return {
      text: 'Create Account'
    }
  },
  methods:{
    create: function(){

      this.$router.push('account');
    }
  }
}
</script>