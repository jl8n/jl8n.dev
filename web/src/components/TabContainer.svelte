<script lang="ts">
    let activeTab = "Languages";
    const tabs = ["Languages", "Experience", "About"];

    interface Language {
        name: string;
        frameworks: string[];
    }

    let languages = [
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
        {
            name: "go",
            frameworks: [],
        },
    ];

    function onLanguageClick(language: Language) {
        console.log("fire", language);
    }

    function setActiveTab(tab: string) {
        activeTab = tab;
    }
</script>

<div class="tabs">
    {#each tabs as tab}
        <button
            on:click={() => setActiveTab(tab)}
            class:active={tab === activeTab}
            aria-controls={tab}
        >
            {tab}
        </button>
    {/each}
</div>

<div class="content">
    {#if activeTab === "Languages"}
        {#each languages as language (language)}
            <li class="language">
                <button on:click={() => onLanguageClick(language)}>
                    {language.name}
                </button>
            </li>
        {/each}

        <article>
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
    {:else if activeTab === "Experience"}
        <p>Experience content here</p>
    {:else if activeTab === "About"}
        <p>About content here</p>
    {/if}
</div>

<style lang="scss">
    button {
        border: none;
        background: transparent;
        color: #ffffff;
        font-size: 1em;
        cursor: pointer;
    }

    .tabs {
        display: flex;
        justify-content: space-around;
    }
    .tabs button {
        flex: 1; /* Each button takes an equal fraction of the space */
        outline: none;
        border-bottom: 2px solid rgba(250, 250, 250, 0.137);
        color: #ffffff;
        font-size: 1.25em;
        text-transform: uppercase;
        font-family: Roboto;
        padding: 25px 0px 15px;
        margin-bottom: 35px;
        /* Add additional styling for buttons here */
    }

    .tabs button.active {
        border-bottom: 2px solid rgb(255, 255, 255); /* Set the white bottom border for the active button */
    }

    .language {
        list-style: none;
        cursor: pointer;
        font-size: 1.25em;
        color: #ffffff;
        padding: 7px 0px;

        &:hover {
            background-color: #ffffff17;
        }
    }
</style>
