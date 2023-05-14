<script lang="ts">
    import { page } from "$app/stores";
    import profile from "$lib/images/pfp.webp";
    import type { NowPlaying } from "./+layout";

    export let data: NowPlaying;
    const { Artist, Album, Track, Art, ArtAvailable } = data;
</script>

<footer>
    <nav>
        <a
            href="/"
            aria-current={$page.url.pathname === "/" ? "page" : undefined}
            >Home</a
        >
        <a
            href="/projects"
            aria-current={$page.url.pathname === "/projects"
                ? "page"
                : undefined}>Projects</a
        >
        <a
            href="/listens"
            aria-current={$page.url.pathname === "/listens"
                ? "page"
                : undefined}>Listens</a
        >
        <a
            href="/rooms"
            aria-current={$page.url.pathname === "/data" ? "page" : undefined}
            >Data</a
        >
    </nav>

    <div class="now-playing">
        {#if Track}
            <div class="foo">
                {#if ArtAvailable}
                    <div class="album-art">
                        <img src={Art} alt="Welcome" />
                    </div>
                {/if}
                <div class="track-info">
                    <div class="track">{Track}</div>
                    <div>{Artist} &bull; {Album}</div>
                </div>
            </div>
        {:else}
            <div>Josh is currently not listening to music</div>
        {/if}
    </div>

    <ul class="links">
        <li><a href="https://github.com/jl8n">github</a></li>
        <li><a href="mailto:josh.l8n@gmail.com">email</a></li>
    </ul>
</footer>

<style>
    footer {
        display: grid;
        grid-template-columns: auto 3fr auto;
        border-top: 1px solid white;
        align-items: center;
        margin: 0px 25px;
        position: sticky;
        bottom: 0px;
        background-color: black;
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
        font-size: 0.8em;
        font-weight: 400;
    }

    .foo {
        display: flex;
        gap: 15px;
        align-items: center;
    }

    .track {
        font-weight: 600;
        font-size: 1.1em;
    }

    .track-info {
        display: flex;
        flex-direction: column;
        font-family: "Roboto";
        font-size: 1.1em;
        line-height: 20px;
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
</style>
