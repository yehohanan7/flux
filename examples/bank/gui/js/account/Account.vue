<template>
  <div class="row" id="account">
        <div class="four columns">
          <div class="row">
            <label><u>{{heading}}</u></label>
          </div>
          <div class="row">
            <strong>$</strong><h3>200</h3>
          </div>
        </div>
        <div class="four columns">
          <form>
            <div class="row">
              <div class="six columns">
                <label for="amount"><u>To Account</u></label>
                <input class="u-full-width" type="number" placeholder="Enter number here" id="amount">
              </div>

            </div>
            <input class="button-primary" type="submit" value="Add">
            <input class="button-primary" type="submit" value="Remove">
          </form>
        </div>
      </div>
</template>

<script>
export default {
  name: 'account',
  data () {
    return {
      heading: 'Balance'
    }
  }
}
</script>
