<script lang="ts">
  import Footer from "./Footer.svelte";

  /* types */
</script>

<div class="grid">
  <header>
    <h1>F</h1>
  </header>

  <main>
    <slot />
  </main>

  <Footer />
</div>

<style>
  @font-face {
    font-family: "Space Mono";
    src: url("$lib/fonts/space-mono-v12-latin-regular.woff2") format("woff");
  }

  @font-face {
    font-family: "Roboto";
    src: url("$lib/fonts/roboto-v30-latin-regular.woff2") format("woff");
  }

  .grid {
    display: grid;
    grid-template-rows: 35px 1fr 75px;
    background-color: rgb(37, 16, 16);
    height: 100vh;
  }

  header {
    display: flex;
    justify-content: center;
    position: sticky;
    top: 0;
    background-color: aqua;
  }

  main {
    display: flex;
  }

  h1 {
    margin: 0px;
  }

  .footer {
    background-color: aquamarine;
  }
</style>
