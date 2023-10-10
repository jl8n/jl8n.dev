<script>
    import { slide } from "svelte/transition";
    export let show = false;
</script>

{#if show}
    <div class="drawer" transition:slide={{ duration: 300 }}>
        <button on:click={() => (show = false)}>Close</button>
        <!-- Your drawer content goes here -->
    </div>
{/if}

<style>
    .drawer {
        position: fixed;
        top: 0;
        bottom: 0;
        width: 100%;
        background: white;
        z-index: 1;
        transition: transform 0.3s ease-out;
    }
</style>
