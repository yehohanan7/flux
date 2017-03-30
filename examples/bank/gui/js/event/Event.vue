<template>
<div class="row">
      <div class="eight columns">
        <label><u>{{heading}}</u></label>
        <ul>
          <li>Event 1</li>
          <li>Event 2</li>
          <li>Event 3</li>
        </ul>
      </div>
    </div>
</template>

<script>
export default {
  name: 'event',
  data () {
    return {
      heading: 'Events'
    }
  }
}
</script>