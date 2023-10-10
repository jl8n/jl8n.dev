<script lang="ts">
    import { onMount } from "svelte";
    import album from "$lib/images/album.jpg";

    let sizes = [
        "small",
        "small",
        "small",
        "small",
        "small",
        "medium",
        "medium",
        "large",
    ];
    let bricks: string[] = [];

    onMount(() => {
        // randomize order of sizes
        bricks = Array.from(
            { length: 75 },
            () => sizes[Math.floor(Math.random() * sizes.length)]
        );
    });
</script>

<div class="grid">
    <div class="wall">
        {#each bricks as brick}
            <div class={`brick ${brick}`}>
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
        grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
        grid-auto-rows: 1fr;
        grid-auto-flow: dense;
        grid-gap: 3px; /* This creates a 3px gap */
    }

    .brick {
        position: relative;
        width: 100%;
        z-index: -1;

        &:before {
            content: "";
            display: block;
            padding-top: 100%;
        }

        img {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            object-fit: cover;
        }
    }

    .medium {
        grid-row: span 2;
        grid-column: span 2;
    }

    .large {
        grid-row: span 3;
        grid-column: span 3;
    }
</style>
