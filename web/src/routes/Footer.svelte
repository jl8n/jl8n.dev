<script lang="ts">
    import { page } from "$app/stores";
    import { onMount } from "svelte";
    import { PUBLIC_BASE_URL } from "$env/static/public";
    import bars from "$lib/images/bars.svg";

    let eventSource: EventSource;

    let track = "";
    let album = "";
    let artist = "";
    let art = "";
    let isDataLoaded = false;
    let isArtLoaded = false;

    function checkIfLoading() {
        if (track && artist && album) {
            isDataLoaded = true;
        }
        if (art) {
            isArtLoaded = true;
        }
    }

    function fetchData() {
        eventSource = new EventSource(`${PUBLIC_BASE_URL}/nowplaying`);

        eventSource.addEventListener("Track", function (event) {
            track = event.data;
            checkIfLoading();
        });

        eventSource.addEventListener("Album", function (event) {
            album = event.data;
        });

        eventSource.addEventListener("Artist", function (event) {
            artist = event.data;
        });

        eventSource.addEventListener("AlbumArt", function (event) {
            art = `${PUBLIC_BASE_URL}${event.data}`;
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

                {#if isDataLoaded}
                    {#if isArtLoaded}
                        <div class="album-art">
                            <img src={art} alt="Welcome" />
                        </div>
                    {:else}
                        <div style="width: 50px; height: 50px;" />
                    {/if}

                    <div class="track-info">
                        <div class="trac foo2">{track}</div>
                        <div class="foo2">{artist} &bull; {album}</div>
                    </div>
                {:else}
                    <img src={bars} alt="Bars" width="25" height="25" />
                {/if}
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
        position: fixed;
        bottom: 0px;
        height: 55px;
        width: 100%;
        padding: 0px 25px;
        display: grid;
        grid-template-columns: 1fr;
        // border-top: 1px solid white;
        align-items: center;
        background-color: black;
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

    .track-info {
        display: flex;
        flex-direction: column;
        font-family: "Roboto";
        font-size: 1.1em;
        line-height: 20px;
    }

    .album-art {
        display: flex;
    }

    .album-art > img {
        height: 50px;
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
