<script lang="ts">
  import Footer from "./Footer.svelte";
  import Drawer from "../components/LeftDrawer.svelte";
  import Test from "./Test.svelte";
</script>

<div class="grid">
  <header class="mobile-only">
    <Drawer />
  </header>

  <main>
    <div class="footer-margin">
      <Test />
      <slot />
    </div>
  </main>
  <div class="footer">
    <Footer />
  </div>
</div>

<style lang="scss">
  @font-face {
    font-family: "Space Mono";
    src: url("$lib/fonts/space-mono-v12-latin-regular.woff2") format("woff");
  }

  @font-face {
    font-family: "Roboto";
    src: url("$lib/fonts/roboto-v30-latin-regular.woff2") format("woff");
  }

  header {
    position: fixed;
    top: 0;
    width: 100%;
    height: 42px;
    background-color: antiquewhite;
  }

  main {
    padding-top: 42px;

    @media (min-width: 768px) {
      padding-top: 0px;
    }
  }

  .grid {
    display: grid;
    grid-template-rows: 1fr 75px;
  }

  .footer-margin {
    margin-bottom: 50px;
  }

  .mobile-only {
    display: block;

    @media (min-width: 768px) {
      display: none;
    }
  }

  // .menu {
  //   font-size: 48px;
  //   color: #ffffff;
  //   outline: none;
  //   border: none;
  //   background-color: transparent;
  //   cursor: pointer;
  // }
</style>
