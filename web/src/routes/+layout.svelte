<script lang="ts">
  import { onMount } from "svelte";
  import { page } from "$app/stores";
  import profile from "$lib/images/pfp.webp";

  let promise = fetch("http://localhost:3000/");
  let album: string;
  let artist: string;
  let track: string;

  onMount(async () => {
    const res = await promise;
    const resData = await res.json();
    album = resData.Album;
    artist = resData.Artist;
    track = resData.Track;
  });
</script>

<div class="grid">
  <header>
    <h1>Josh Layton</h1>
  </header>

  <main>
    <slot />
  </main>

  <footer>
    <nav>
      <a href="/" aria-current={$page.url.pathname === "/" ? "page" : undefined}
        >Home</a
      >
      <a
        href="/rooms"
        aria-current={$page.url.pathname === "/data" ? "page" : undefined}
        >Projects</a
      >
      <a
        href="/rooms"
        aria-current={$page.url.pathname === "/music" ? "page" : undefined}
        >Music</a
      >
      <a
        href="/rooms"
        aria-current={$page.url.pathname === "/data" ? "page" : undefined}
        >Data</a
      >

    </nav>

    <div class="now-playing">
      <div class="foo">
        <div class="album-art">
          <img src={profile} alt="Welcome" />
        </div>
        <div class="track-info">
          {#await promise}
            <p>Loading...</p>
          {:then data}
            <p>{track}</p>
            <p>{artist} &middot; {album}</p>
          {:catch error}
            <p>Something went wrong</p>
          {/await}
        </div>
      </div>

    </div>

    <ul class="links">
      <li><a href="https://github.com/jl8n">github</a></li>
      <li><a href="mailto:josh.l8n@gmail.com">email</a></li>
    </ul>
  </footer>
</div>

<style>
  @font-face {
    font-family: "Space Mono";
    src: url("$lib/fonts/space-mono-v12-latin-regular.woff2") format("woff");
  }

  .grid {
    display: grid;
    grid-template-rows: 35px 1fr 75px;
    background-color: black;
    min-height: 100vh;
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
    justify-content: center;
    align-items: center;
  }

  footer {
    display: grid;
    grid-template-columns: auto 3fr auto;
    border-top: 1px solid white;
    align-items: center;
    margin: 0px 25px;
    position: sticky;
    bottom: 0px;
    font-family: "Space Mono";
  }

  footer > nav {
    display: flex;
    gap: 3vw;
    text-transform: uppercase;
  }

  footer > nav > a {
    color: white;
  }

  footer > nav > a[aria-current="page"] {
    color: rgb(192, 21, 21);
  }

  footer > .now-playing {
    display: flex;
    flex-direction: column;
    align-items: center;
    font-family: "Space Mono", sans-serif;
    font-size: 0.8em;
    font-weight: 400;
  }

  .track-info > p {
    margin: 9px 10px;
  }

  .foo {
    display: flex;
    justify-content: center;
    align-items: center;
  }
  
  .album-art > img {
    height: 50px;
    background-color: pink;
  }

  
  
  footer > .links {
    display: flex;
    gap: 25px;
    list-style-type: none;
  }

  nav > a {
    text-decoration: none;
  }

  h1 {
    margin: 0px;
  }
</style>
