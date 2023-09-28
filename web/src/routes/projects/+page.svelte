<script lang="ts">
    import { onMount } from "svelte";
    let eventSource: EventSource;
    let data = "";
    let track = "";
    let album = "";
    let artist = "";

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
    }

    onMount(() => {
        return () => {
            if (eventSource) {
                eventSource.close();
            }
        };
    });
</script>

<div class="flex">
    <button on:click={fetchData}> Fetch Data </button>
    <div>Latest Track: {track}</div>
    <div>Latest Album: {album}</div>
    <div>Latest Artist: {artist}</div>
</div>

<style>
    .flex {
        display: flex;
        flex-direction: column;
        padding-top: 3em;
    }
</style>
