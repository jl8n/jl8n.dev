<script lang="ts">
    import profile from "$lib/images/pfp.webp";
    import TabContainer from "../components/TabContainer.svelte";

    let currentSelection: string = "home";
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

    <div>
        <article>
            <TabContainer />
        </article>


    </div>
</section>

<style lang="scss">
    h1,
    h3 {
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
        margin-bottom: 50px;

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
            font-size: 7em;
        }
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
