<script lang="ts">
  import Footer from "./Footer.svelte";
</script>

<div class="grid">
  <header class="mobile-only" />
  <main>
    <div class="footer-margin">
      <slot />
    </div>
  </main>

  <Footer />
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
