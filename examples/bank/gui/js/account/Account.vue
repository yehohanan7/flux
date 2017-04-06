<template>
  <div class="row" id="account">
        <div class="four columns">
          <div class="row">
            <label><u>{{heading}}</u></label>
          </div>
          <div class="row">
            <strong>$</strong><h3 v-model="balance">{{balance}}</h3>
          </div>
        </div>
        <div class="four columns">
          <form>
            <div class="row">
              <div class="six columns">
                <label for="amount"><u>To Account</u></label>
                <input class="u-full-width" type="number" placeholder="Enter number here" id="amount" v-model="amount" >
              </div>

            </div>
            <input class="button-primary" type="button" value="Add" v-on:click="credit">
            <input class="button-primary" type="button" value="Remove" v-on:click="debit">
          </form>
        </div>
      </div>
</template>

<script>
import bankApi from '../api/BankApi.js';
export default {
  name: 'account',
  data () {
    return {
      heading: 'Balance',
      amount: 0,
      balance: 0
    }
  },
  methods:{
    debit: function(){
      this.balance -= amount.valueAsNumber;
      bankApi.debit({amount: amount.valueAsNumber});
      this.balance = bankApi.currentBalance();

    },
    credit: function(){
      this.balance += amount.valueAsNumber;
      bankApi.credit({amount: amount.valueAsNumber});
    },
  }
}
</script>
