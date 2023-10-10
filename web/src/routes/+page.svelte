<script lang="ts">
    import profile from "$lib/images/pfp.webp";

    interface Language {
        name: string;
        frameworks: string[];
    }

    let currentSelection: string = "home";
    let languages = [
        {
            name: "home",
            frameworks: [],
        },
        {
            name: "c++",
            frameworks: [],
        },
        {
            name: "python",
            frameworks: [],
        },
        {
            name: "javascript / typescript",
            frameworks: ["node", "vue", "sveltekit"],
        },
    ];

    function onLanguageClick(language: Language) {
        currentSelection = language.name;
        console.log("fire", language);
    }
</script>

<section class="grid">
    <div class="stage">
        <!-- <img src={profile} alt="Welcome" style="width: 250px;" /> -->
        {#if currentSelection === "home"}
            <div class="foo"><h1>Hi, I'm Josh</h1></div>
        {:else if currentSelection === "c++"}
            <h1 class="foo">c++</h1>
        {:else if currentSelection === "javascript / typescript"}
            <div class="typescript">
                <h1 class="foo">javascript / typescript</h1>
                <article>
                    <ul>
                        <li>Vue</li>
                        <li>Svelte / Sveltekit</li>
                    </ul>
                </article>
            </div>
        {:else if currentSelection === "python"}
            <h1 class="foo">python</h1>
        {/if}
    </div>
    <div class="right">
        <article>
            <p>
                I'm primarily a full stack developer, but my interests span
                everything from embedded sytems to dev-ops.
            </p>
            <p>
                Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
                eiusmod tempor incididunt ut labore et dolore magna aliqua.
                Condimentum lacinia quis vel eros donec. Egestas fringilla
                phasellus faucibus scelerisque eleifend donec pretium vulputate
                sapien. Massa tincidunt dui ut ornare lectus sit amet. Rhoncus
                est pellentesque elit ullamcorper dignissim cras tincidunt
                lobortis feugiat. Malesuada proin libero nunc consequat interdum
                varius sit. Ornare suspendisse sed nisi lacus sed. Aliquam
                faucibus purus in massa tempor nec. Phasellus faucibus
                scelerisque eleifend donec pretium. At elementum eu facilisis
                sed odio morbi. Praesent elementum facilisis leo vel fringilla
                est. Duis at consectetur lorem donec massa sapien faucibus et.
                Tortor condimentum lacinia quis vel eros donec. Ultrices gravida
                dictum fusce ut placerat orci nulla pellentesque dignissim. Sed
                vulputate mi sit amet. Leo vel orci porta non pulvinar neque. Eu
                ultrices vitae auctor eu augue ut lectus arcu bibendum. Purus
                semper eget duis at tellus. Augue neque gravida in fermentum et
                sollicitudin ac orci. Pharetra sit amet aliquam id diam maecenas
                ultricies mi eget. Ac tortor vitae purus faucibus ornare.
                Interdum velit laoreet id donec. Dictum non consectetur a erat
                nam at lectus. Libero justo laoreet sit amet cursus sit amet
                dictum sit. Amet volutpat consequat mauris nunc congue nisi
                vitae suscipit tellus.
            </p>
        </article>

        <article>
            <h3>Languages</h3>
            <ul>
                {#each languages as language (language)}
                    <li>
                        <button
                            on:click={() => onLanguageClick(language)}
                            class="language"
                        >
                            {language.name}
                        </button>
                    </li>
                {/each}
            </ul>

            <h3>Technologies</h3>
            <ul>
                <li>linux</li>
                <li>dev-ops</li>
                <li>docker</li>
                <li>CI/CD</li>
                <li>git</li>
                <li>security</li>
            </ul>
        </article>
    </div>
</section>

<style lang="scss">
    h1,
    h2,
    h3,
    h5 {
        margin: 0px;
    }

    /* section {
        display: flex;
        gap: 10px;
    } */

    .grid {
        display: grid;
        grid-template-columns: 1fr;
        grid-template-rows: auto auto;

        @media (min-width: 768px) {
            grid-template-columns: 1fr 1fr;
            grid-template-rows: auto;
        }
    }

    .foo {
        text-transform: uppercase;
        font-size: 5em;
        line-height: 1.5em;
        overflow: hidden;
        text-overflow: ellipsis;

        @media (min-width: 768px) {
            font-size: 8em;
        }
    }

    .right {
        display: grid;
        grid-template-rows: 50vh 50vh;
        height: 100%;
    }

    .stage {
        height: 70vh;

        @media (min-width: 768px) {
            height: 100%;
        }
    }

    article {
        font-size: 0.9em;
        color: white;
        overflow: auto;
        text-align: justify;
        padding: 15px;
    }

    ul {
        gap: 10px;
        list-style: none;
        margin: 0px;
    }

    .language {
        background: transparent;
        border: none;
        font-size: 16px;
        font-weight: 500;
        color: rgba(255, 255, 255, 0.87); /* Default text color */
        transition: color 0.3s ease; /* Smooth color transition */
    }

    .language:hover {
        cursor: pointer;
    }

    .language:active {
        color: #3f51b5; /* Text color when clicked */
    }

    .typescript {
        width: 100%;
        height: 100%;
        background-color: #007acc;

        background: linear-gradient(to bottom left, #f0dc4e 50%, #2d79c7 50%);
    }
</style>
