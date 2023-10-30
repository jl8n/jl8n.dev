<script lang="ts">
    import { onMount } from "svelte";
    import { PUBLIC_BASE_URL } from "$env/static/public";

    // frequency of size is the chance an image will be that size
    // 5 `small` out of 8 total means there's a 5/8 chance of small
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
    let bricks: { size: string; url: string }[] = [];

    function handleImageLoad(event: Event) {
        const loadingImage = event.target as HTMLImageElement;
        loadingImage.classList.add("loaded");
        loadingImage.classList.remove("loading");
    }

    function shuffleBricks() {
        bricks = [...bricks.sort(() => Math.random() - 0.5)];
    }

    onMount(async () => {
        // randomize order of sizes
        const sizeArray = Array.from(
            { length: 100 },
            () => sizes[Math.floor(Math.random() * sizes.length)]
        );

        console.log(sizeArray);

        const res = await fetch(`${PUBLIC_BASE_URL}/art`);
        if (res.ok) {
            const resData = await res.json();
            const fileURLs = resData.files;

            // Combine sizes and URLs into bricks array
            bricks = sizeArray.map((size, i) => ({
                size,
                url: fileURLs[i % fileURLs.length], // Loop over URLs if there are more bricks than URLs
            }));

            // Shuffle bricks array
            for (let i = bricks.length - 1; i > 0; i--) {
                const j = Math.floor(Math.random() * (i + 1));
                [bricks[i], bricks[j]] = [bricks[j], bricks[i]];
            }
        }
    });
</script>

<div class="grid">
    <div class="wall">
        <button on:click={shuffleBricks} />
        {#each bricks as brick, index (index)}
            <div class={`brick ${brick.size}`}>
                <img
                    src={brick.url}
                    alt="Album"
                    class="loading"
                    on:load={handleImageLoad}
                />
            </div>
        {/each}
    </div>
</div>

<style lang="scss">
    .grid {
        height: 100%;
        width: 100%;
    }

    img {
        z-index: 0;
    }

    img.loading {
        opacity: 0;
    }

    .loaded {
        opacity: 0;
        animation-name: fadeIn;
        animation-duration: 0.5s;
        animation-timing-function: ease-in-out;
        animation-fill-mode: forwards;
    }

    @keyframes fadeIn {
        from {
            opacity: 0;
        }
        to {
            opacity: 1;
        }
    }

    .wall {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
        grid-auto-rows: 1fr;
        grid-auto-flow: dense;
        grid-gap: 3px;
        transform: translateY(-50px);
    }

    .brick {
        position: relative;
        width: 100%;

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
            transition: transform 0.3s ease;
            z-index: 1;
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
