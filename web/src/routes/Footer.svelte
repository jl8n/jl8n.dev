<script lang="ts">
    import { page } from "$app/stores";
    import { onMount } from "svelte";
    let eventSource: EventSource;

    let track = "";
    let album = "";
    let artist = "";
    let art = "";

    function fetchData() {
        eventSource = new EventSource("http://localhost:3000/nowplaying");

        eventSource.addEventListener("Track", function (event) {
            track = event.data;
        });

        eventSource.addEventListener("Album", function (event) {
            album = event.data;
        });

        eventSource.addEventListener("Artist", function (event) {
            artist = event.data;
        });

        eventSource.addEventListener("AlbumArt", function (event) {
            art = "http://localhost:3000" + event.data;
        });
    }

    onMount(() => {
        fetchData();
        return () => {
            if (eventSource) {
                eventSource.close();
            }
        };
    });
</script>

<footer>
    <nav class="desktop-only">
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
        <!-- <a
            href="/rooms"
            aria-current={$page.url.pathname === "/data" ? "page" : undefined}
            >Data</a
        > -->
    </nav>

    <div class="now-playing">
        {#if track}
            <div class="foo">
                <!-- {#if ArtAvailable}
                    <div class="album-art">
                        <img src={Art} alt="Welcome" />
                    </div>
                {/if} -->
                <div class="album-art">
                    <img src={art} alt="Welcome" />
                </div>
                <div class="track-info">
                    <div class="trac foo2">{track}</div>
                    <div class="foo2">{artist} &bull; {album}</div>
                </div>
            </div>
        {:else}
            <div>Josh is currently not listening to music</div>
        {/if}
    </div>

    <div class="desktop-only">
        <ul class="links">
            <li><a href="https://github.com/jl8n">github</a></li>
            <li><a href="mailto:josh.l8n@gmail.com">email</a></li>
        </ul>
    </div>
</footer>

<style lang="scss">
    footer {
        display: grid;
        grid-template-columns: 1fr;
        border-top: 1px solid white;
        align-items: center;
        padding: 0px 25px;
        position: fixed;
        bottom: 0px;
        background-color: black;
        width: 100%;
        box-sizing: border-box;

        @media (min-width: 768px) {
            grid-template-columns: auto 3fr auto;
        }
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

    .foo2 {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        max-width: 65vw;

        @media (min-width: 768px) {
            max-width: 30vw;
        }
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

    .links {
        display: flex;
        gap: 25px;
        list-style-type: none;
    }

    nav > a {
        text-decoration: none;
    }

    .desktop-only {
        display: none;

        @media (min-width: 768px) {
            display: block;
        }
    }
</style>
