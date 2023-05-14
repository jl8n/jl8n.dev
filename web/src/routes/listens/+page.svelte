<script lang="ts">
    import album from "$lib/images/album.jpg";
</script>

<div class="grid">
    <h1>Music</h1>
    <div class="wall">
        {#each { length: 50 } as _, i}
            <div class="brick">
                <img src={album} alt="Album" />
            </div>
        {/each}
    </div>
</div>

<style lang="scss">
    .grid {
        height: 100%;
        width: 100%;
    }

    .wall {
        display: grid;
        grid-auto-flow: dense;
        grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
        justify-items: center;
        align-items: center;
        background-size: cover;
        object-fit: cover;
    }

    .brick {
        img {
            max-width: 100%;
            height: 100%;
        }

        &:nth-of-type(3n) {
            grid-row: span 2;
            grid-column: span 2;
        }

        &:nth-of-type(9n) {
            grid-row: span 3;
            grid-column: span 3;
        }
    }
</style>
