<div class="grid">

  <header>
    <h1>Josh Layton</h1>

  </header>

  <main>
    <slot />
  </main>

  <footer>
    <nav>
      <a href="/">Home</a>
      <a href="/rooms">Rooms</a>
    </nav>

    <div class="now-playing">
      Now playing
    </div>

    <div class="links">
      github
    </div>
  </footer>

</div>


<style>
  .grid {
    display: grid;
    grid-template-rows: 35px 1fr 75px;
    background-color: black;
    min-height: 100vh;
  }

  header {
    display: flex;
    justify-content: center;
    background-color: aqua;
  }

  main {

  }

  footer {
    display: grid;
    grid-template-columns: auto 3fr auto;
    border-top: 1px solid white;
    align-items: center;
    margin: 0px 25px;
    position: sticky;
    bottom: 0px;
  }

  footer > nav {
    display: flex;
    text-transform: uppercase;
  }

  footer > nav > a {
    margin-right: 30px;
    color: white;
  }

  footer > .now-playing {
    display: flex;
    justify-content: center;
    font-family: 'Space Mono', sans-serif;
    font-size: 1.1vw;
    font-weight: 400;
  }

  nav > a {
    text-decoration: none;
  }

  h1 {
    margin: 0px;
  }
</style>